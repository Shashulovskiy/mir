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
	"github.com/filecoin-project/mir/pkg/coding"
	mirCrypto "github.com/filecoin-project/mir/pkg/crypto"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/merkletree"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/net/grpc"
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

	ByzantineBehavior string
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
	var byzantine bool

	// print the IP addresses of all nodes
	for _, node := range nodes.Items {
		ips = append(ips, node.Status.PodIP)
		fmt.Printf("Node IP: %s\n", node.Status.PodIP)
	}

	slices.Sort(ips)

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

	for i, ip := range ips {
		if ip == thisIp {
			ownID = t.NewNodeIDFromInt(i)
			if i > len(ips)-len(ips)/3 && args.ByzantineBehavior != "none" {
				byzantine = true
				fmt.Printf("Enabling byzantine behavior\n")
			} else {
				byzantine = false
				fmt.Printf("Disabling byzantine behavior\n")
			}
		}
	}

	fmt.Printf("Own ID: %s\n", ownID)

	println("Node starting..")

	if args.ByzantineBehavior != "none" {
		fmt.Printf("Byzantine behavior in system: %s\n", args.ByzantineBehavior)
	} else {
		fmt.Println("No byzantine behavior in system")
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

	hasher := mirCrypto.NewHasher(crypto.SHA1)
	merkle := merkletree.NewVerifier()

	brbBrachaModule := configureBracha(byzantine, args.ByzantineBehavior, nodeIDs, leaderNode, ownID)
	brbCtModule, err := configureCT(byzantine, args.ByzantineBehavior, nodeIDs, leaderNode, ownID)
	brbDxrModule, err := configureDXR(byzantine, args.ByzantineBehavior, nodeIDs, leaderNode, ownID)

	coder := coding.NewModule()

	if err != nil {
		return nil
	}

	counter := 0

	control := newControlModule(
		/*isLeader=*/ ownID == nodeIDs[leaderNode],
		func(id int64, message *[]byte, algorithm string) *events.EventList {
			return events.ListOf(&eventpb.Event{
				DestModule: algorithm,
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
		},
		func(bytes []byte) {
			counter++
		},
		tests,
	)

	m := map[t.ModuleID]modules.Module{
		"net":     transportModule,
		"crypto":  mirCrypto.New(&mirCrypto.DummyCrypto{DummySig: []byte{0}}),
		"brb":     brbBrachaModule,
		"brbct":   brbCtModule,
		"brbdxr":  brbDxrModule,
		"control": control,
		"hasher":  hasher,
		"merkle":  merkle,
		"coder":   coder,
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

func configureBracha(byzantine bool, byzantineStrategy string, nodeIDs []t.NodeID, leaderNode int, ownID t.NodeID) modules.PassiveModule {
	moduleConfig := &brb.ModuleConfig{
		Self:     "brb",
		Consumer: "control",
		Net:      "net",
		Crypto:   "crypto",
	}
	moduleParams := &brb.ModuleParams{
		InstanceUID: []byte("testing instance"),
		AllNodes:    nodeIDs,
		Leader:      nodeIDs[leaderNode],
	}
	if byzantine {
		return brb.NewByzantineModule(
			moduleConfig,
			moduleParams,
			ownID,
			byzantineStrategy)
	} else {
		return brb.NewModule(
			moduleConfig,
			moduleParams,
			ownID,
		)
	}
}

func configureCT(byzantine bool, byzantineStrategy string, nodeIDs []t.NodeID, leaderNode int, ownID t.NodeID) (modules.PassiveModule, error) {
	moduleConfig := &brbct.ModuleConfig{
		Self:                "brbct",
		Consumer:            "control",
		Net:                 "net",
		Crypto:              "crypto",
		Hasher:              "hasher",
		MerkleProofVerifier: "merkle",
		Coder:               "coder",
	}
	moduleParams := &brbct.ModuleParams{
		InstanceUID: []byte("testing instance"),
		AllNodes:    nodeIDs,
		Leader:      nodeIDs[leaderNode],
	}
	if byzantine {
		return brbct.NewByzantineModule(moduleConfig, moduleParams, ownID, byzantineStrategy)
	} else {
		return brbct.NewModule(moduleConfig, moduleParams, ownID)
	}
}

func configureDXR(byzantine bool, byzantineStrategy string, nodeIDs []t.NodeID, leaderNode int, ownID t.NodeID) (modules.PassiveModule, error) {
	moduleConfig := &brbdxr.ModuleConfig{
		Self:     "brbdxr",
		Consumer: "control",
		Net:      "net",
		Crypto:   "crypto",
		Hasher:   "hasher",
		Coder:    "coder",
	}
	moduleParams := &brbdxr.ModuleParams{
		InstanceUID: []byte("testing instance"),
		AllNodes:    nodeIDs,
		Leader:      nodeIDs[leaderNode],
	}
	if byzantine {
		return brbdxr.NewByzantineModule(moduleConfig, moduleParams, ownID, byzantineStrategy)
	} else {
		return brbdxr.NewModule(moduleConfig, moduleParams, ownID)
	}
}

func parseArgs(args []string) *parsedArgs {
	app := kingpin.New("brb-channel", "BRB Chanel for continuously sending messages")
	verbose := app.Flag("verbose", "Verbose mode.").Short('v').Bool()
	trace := app.Flag("trace", "Very verbose mode.").Bool()
	testFile := app.Arg("testFile", "File with tests").Required().String()
	byzantineBehavior := app.Arg("byzantineBehavior", "Enable byzantine behavior").String()

	if _, err := app.Parse(args[1:]); err != nil { // Skip args[0], which is the name of the program, not an argument.
		app.FatalUsage("could not parse arguments: %v\n", err)
	}

	return &parsedArgs{
		testFile:          *testFile,
		Verbose:           *verbose,
		Trace:             *trace,
		ByzantineBehavior: *byzantineBehavior,
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
