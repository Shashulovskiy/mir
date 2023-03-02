package main

import (
	"context"
	"crypto"
	"fmt"
	"github.com/filecoin-project/mir"
	"github.com/filecoin-project/mir/pkg/brb"
	"github.com/filecoin-project/mir/pkg/brbct"
	"github.com/filecoin-project/mir/pkg/brbdxr"
	mirCrypto "github.com/filecoin-project/mir/pkg/crypto"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/merkletree"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/net/grpc"
	"github.com/filecoin-project/mir/pkg/pb/brbctpb"
	"github.com/filecoin-project/mir/pkg/pb/brbdxrpb"
	"github.com/filecoin-project/mir/pkg/pb/brbpb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	t "github.com/filecoin-project/mir/pkg/types"
	grpctools "github.com/filecoin-project/mir/pkg/util/grpc"
	"gopkg.in/alecthomas/kingpin.v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"time"
)

const (

	// Base port number for the nodes to listen to messages from each other.
	// The nodes will listen on ports starting from nodeBasePort through nodeBasePort+3.
	nodeBasePort = 10000

	// The number of nodes participating in the chat.
	nodeNumber = 4

	// The index of the leader node of BCB.
	leaderNode = 0
)

// parsedArgs represents parsed command-line parameters passed to the program.
type parsedArgs struct {

	// ID of this node.
	// The package github.com/hyperledger-labs/mir/pkg/types defines this and other types used by the library.
	OwnID t.NodeID

	// If set, print debug output to stdout.
	Verbose bool

	// If set, print trace output to stdout.
	Trace bool
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = os.Getenv("HOME") + "/.kube/config"
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	// Create a Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// Get the IP addresses of all pods in the deployment
	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), v1.ListOptions{LabelSelector: "app=myapp"})
	if err != nil {
		panic(err)
	}
	ips := make([]string, len(pods.Items))
	for i, pod := range pods.Items {
		ips[i] = pod.Status.PodIP
	}

	// Pass the IP addresses to each replica
	for _, ip := range ips {
		fmt.Printf("IP address: %s\n", ip)
	}

	println("Node starting..")
	args := parseArgs(os.Args)
	// Initialize logger that will be used throughout the code to print log messages.
	var logger logging.Logger
	if args.Trace {
		logger = logging.ConsoleTraceLogger // Print trace-level info.
	} else if args.Verbose {
		logger = logging.ConsoleDebugLogger // Print debug-level info in verbose mode.
	} else {
		logger = logging.ConsoleWarnLogger // Only print errors and warnings by default.
	}

	// IDs of nodes that are part of the system.
	// This example uses a static configuration of nodeNumber nodes.
	nodeIDs := make([]t.NodeID, nodeNumber)
	for i := 0; i < nodeNumber; i++ {
		nodeIDs[i] = t.NewNodeIDFromInt(i)
	}

	nodeAddrs := make(map[t.NodeID]t.NodeAddress)
	for i := range nodeIDs {
		nodeAddrs[t.NewNodeIDFromInt(i)] = t.NodeAddress(grpctools.NewDummyMultiaddr(i + nodeBasePort))
	}

	transportModule, err := grpc.NewTransport(args.OwnID, nodeAddrs[args.OwnID], logger)
	if err != nil {
		return fmt.Errorf("failed to get network transport %w", err)
	}
	if err := transportModule.Start(); err != nil {
		return fmt.Errorf("could not start network transport: %w", err)
	}
	transportModule.Connect(nodeAddrs)

	merkle := merkletree.NewVerifier()
	hasher := mirCrypto.NewHasher(crypto.SHA1)

	brbBrachaModule := brb.NewModule(
		&brb.ModuleConfig{
			Self:     "brb",
			Consumer: "control",
			Net:      "net",
			Crypto:   "crypto",
		},
		&brb.ModuleParams{
			InstanceUID: []byte("testing instance"),
			AllNodes:    nodeIDs,
			Leader:      nodeIDs[leaderNode],
		},
		args.OwnID,
	)

	brbCtModule, err := brbct.NewModule(
		&brbct.ModuleConfig{
			Self:                "brbct",
			Consumer:            "control",
			Net:                 "net",
			Crypto:              "crypto",
			Hasher:              "hasher",
			MerkleProofVerifier: "merkle",
		},
		&brbct.ModuleParams{
			InstanceUID: []byte("testing instance"),
			AllNodes:    nodeIDs,
			Leader:      nodeIDs[leaderNode],
		},
		args.OwnID,
	)

	brbDxrModule, err := brbdxr.NewModule(
		&brbdxr.ModuleConfig{
			Self:     "brbdxr",
			Consumer: "control",
			Net:      "net",
			Crypto:   "crypto",
			Hasher:   "hasher",
			//MerkleProofVerifier: "merkle",
		},
		&brbdxr.ModuleParams{
			InstanceUID: []byte("testing instance"),
			AllNodes:    nodeIDs,
			Leader:      nodeIDs[leaderNode],
		},
		args.OwnID,
	)

	if err != nil {
		return nil
	}

	counter := 0

	control := newControlModule(
		/*isLeader=*/ args.OwnID == nodeIDs[leaderNode],
		func(id int64, message *[]byte, algorithm string) *events.EventList {
			if algorithm == "brbdxr" {
				return events.ListOf(&eventpb.Event{
					DestModule: "brbdxr",
					Type: &eventpb.Event_Brbdxr{
						Brbdxr: &brbdxrpb.Event{
							Type: &brbdxrpb.Event_Request{
								Request: &brbdxrpb.BroadcastRequest{
									Id:   id,
									Data: *message,
								},
							},
						},
					},
				})
			} else if algorithm == "brbct" {
				return events.ListOf(&eventpb.Event{
					DestModule: "brbct",
					Type: &eventpb.Event_Brbct{
						Brbct: &brbctpb.Event{
							Type: &brbctpb.Event_Request{
								Request: &brbctpb.BroadcastRequest{
									Id:   id,
									Data: *message,
								},
							},
						},
					},
				})
			} else if algorithm == "brb" {
				return events.ListOf(&eventpb.Event{
					DestModule: "brb",
					Type: &eventpb.Event_Brb{
						Brb: &brbpb.Event{
							Type: &brbpb.Event_Request{
								Request: &brbpb.BroadcastRequest{
									Id:   id,
									Data: *message,
								},
							},
						},
					},
				})
			} else {
				panic("Unknown algorithm")
			}
		},
		func(bytes []byte) {
			counter++
		},
	)

	go func() {
		time.Sleep(time.Minute)
		println(counter)
	}()

	m := map[t.ModuleID]modules.Module{
		"net":     transportModule,
		"crypto":  mirCrypto.New(&mirCrypto.DummyCrypto{DummySig: []byte{0}}),
		"brb":     brbBrachaModule,
		"brbct":   brbCtModule,
		"brbdxr":  brbDxrModule,
		"control": control,
		"hasher":  hasher,
		"merkle":  merkle,
	}

	// create a Mir node
	node, err := mir.NewNode("client", mir.DefaultNodeConfig().WithLogger(logger), m, nil, nil)
	if err != nil {
		return fmt.Errorf("error creating a Mir node: %w", err)
	}

	// run the node
	err = node.Run(context.Background())
	if err != nil {
		return fmt.Errorf("error running node: %w", err)
	}

	return nil
}

func parseArgs(args []string) *parsedArgs {
	app := kingpin.New("brb-channel", "BRB Chanel for continuously sending messages")
	verbose := app.Flag("verbose", "Verbose mode.").Short('v').Bool()
	trace := app.Flag("trace", "Very verbose mode.").Bool()
	ownID := app.Arg("id", "ID of this node").Required().String()

	if _, err := app.Parse(args[1:]); err != nil { // Skip args[0], which is the name of the program, not an argument.
		app.FatalUsage("could not parse arguments: %v\n", err)
	}

	return &parsedArgs{
		OwnID:   t.NodeID(*ownID),
		Verbose: *verbose,
		Trace:   *trace,
	}
}
