package main

import (
	"bufio"
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
	"golang.org/x/exp/slices"
	"gopkg.in/alecthomas/kingpin.v2"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// parsedArgs represents parsed command-line parameters passed to the program.
type parsedArgs struct {

	// File with tests
	testFile string

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
	temp, err := os.CreateTemp("", "logs")
	if err != nil {
		return err
	}
	l := log.New(temp, "INFO\t", log.Ldate|log.Ltime)
	l.Println("Node starting")

	config, err := rest.InClusterConfig()
	if err != nil {
		l.Println(err.Error())
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		l.Println(err.Error())
		panic(err.Error())
	}

	// get the list of nodes
	nodes, err := clientset.CoreV1().Pods("default").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			fmt.Println("Nodes not found")
		} else {
			panic(err.Error())
		}
	}
	thisIp := os.Getenv("MY_POD_IP")
	ips := make([]string, 0)

	fmt.Printf("Self IP: %s\n", thisIp)

	var ownID t.NodeID

	// print the IP addresses of all nodes
	for _, node := range nodes.Items {
		ips = append(ips, node.Status.PodIP)
		fmt.Printf("Node IP: %s\n", node.Status.PodIP)
	}

	slices.Sort(ips)

	for i, ip := range ips {
		if ip == thisIp {
			ownID = t.NewNodeIDFromInt(i)
		}
	}

	fmt.Printf("Own ID: %s\n", ownID)

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

	tests := parseTests(args.testFile)

	nodeNumber := len(ips)
	nodeBasePort := 10000
	leaderNode := 0

	// IDs of nodes that are part of the system.
	// This example uses a static configuration of nodeNumber nodes.
	nodeIDs := make([]t.NodeID, nodeNumber)
	for i := 0; i < nodeNumber; i++ {
		nodeIDs[i] = t.NewNodeIDFromInt(i)
	}

	nodeAddrs := make(map[t.NodeID]t.NodeAddress)
	for i := range nodeIDs {
		nodeAddrs[t.NewNodeIDFromInt(i)] = t.NodeAddress(grpctools.NewMultiaddr(ips[i], nodeBasePort))
	}

	transportModule, err := grpc.NewTransport(ownID, nodeAddrs[ownID], logger)
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
		ownID,
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
		ownID,
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
		ownID,
	)

	if err != nil {
		return nil
	}

	counter := 0

	control := newControlModule(
		/*isLeader=*/ ownID == nodeIDs[leaderNode],
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
				panic("Unknown algorithm " + algorithm)
			}
		},
		func(bytes []byte) {
			counter++
		},
		tests,
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
	testFile := app.Arg("testFile", "File with tests").Required().String()

	if _, err := app.Parse(args[1:]); err != nil { // Skip args[0], which is the name of the program, not an argument.
		app.FatalUsage("could not parse arguments: %v\n", err)
	}

	return &parsedArgs{
		testFile: *testFile,
		Verbose:  *verbose,
		Trace:    *trace,
	}
}

func parseTests(testFile string) []*Benchmark {
	file, err := os.Open(testFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	benchmarks := make([]*Benchmark, 0)
	for scanner.Scan() {
		split := strings.Split(scanner.Text(), " ")
		msgSize, err := strconv.ParseInt(split[0], 10, 64)
		if err != nil {
			panic(err)
		}
		duration, err := strconv.ParseInt(split[1], 10, 64)
		if err != nil {
			panic(err)
		}
		algorithm := split[2]

		benchmarks = append(benchmarks, &Benchmark{
			messageSize: msgSize,
			duration:    time.Duration(duration * time.Second.Nanoseconds()),
			algorithm:   algorithm,
		})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return benchmarks
}
