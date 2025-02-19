package main

import (
	"context"
	"crypto"
	_ "crypto/sha1"
	"fmt"
	"github.com/filecoin-project/mir/pkg/brbdxr"
	"github.com/filecoin-project/mir/pkg/coding"
	"github.com/filecoin-project/mir/pkg/merkletree"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/filecoin-project/mir"
	mirCrypto "github.com/filecoin-project/mir/pkg/crypto"
	"github.com/filecoin-project/mir/pkg/logging"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/net/grpc"
	t "github.com/filecoin-project/mir/pkg/types"
	grpctools "github.com/filecoin-project/mir/pkg/util/grpc"
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

	coder := coding.NewModule()

	brbModule, err := brbdxr.NewModule(
		&brbdxr.ModuleConfig{
			Self:     "brbdxr",
			Consumer: "control",
			Net:      "net",
			Crypto:   "crypto",
			Hasher:   "hasher",
			Coder:    "coder",
			//MerkleProofVerifier: "merkle",
		},
		&brbdxr.ModuleParams{
			InstanceUID: []byte("testing instance"),
			AllNodes:    nodeIDs,
			Leader:      nodeIDs[leaderNode],
		},
		args.OwnID,
		"optimized",
	)

	if err != nil {
		return nil
	}

	// control module reads the user input from the console and processes it.
	control := newControlModule( /*isLeader=*/ args.OwnID == nodeIDs[leaderNode])
	hasher := mirCrypto.NewHasher(crypto.SHA1)

	m := map[t.ModuleID]modules.Module{
		"net":     transportModule,
		"crypto":  mirCrypto.New(&mirCrypto.DummyCrypto{DummySig: []byte{0}}),
		"brbdxr":  brbModule,
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

// Parses the command-line arguments and returns them in a params struct.
func parseArgs(args []string) *parsedArgs {
	app := kingpin.New("chat-demo", "Small chat application to demonstrate the usage of the Mir library.")
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
