package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/filecoin-project/mir/pkg/pb/brbdxrpb"
	"os"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
)

type controlModule struct {
	eventsOut chan *events.EventList
	isLeader  bool
}

func newControlModule(isLeader bool) modules.ActiveModule {
	return &controlModule{
		eventsOut: make(chan *events.EventList),
		isLeader:  isLeader,
	}
}

func (m *controlModule) ImplementsModule() {}

func (m *controlModule) ApplyEvents(ctx context.Context, events *events.EventList) error {
	iter := events.Iterator()
	for event := iter.Next(); event != nil; event = iter.Next() {
		switch event.Type.(type) {

		case *eventpb.Event_Init:
			if m.isLeader {
				go func() {
					err := m.readMessageFromConsole()
					if err != nil {
						panic(err)
					}
				}()
			} else {
				fmt.Println("Waiting for the message...")
			}

		case *eventpb.Event_Brbdxr:
			brbEvent := event.Type.(*eventpb.Event_Brbdxr).Brbdxr
			switch brbEvent.Type.(type) {

			case *brbdxrpb.Event_Deliver:
				deliverEvent := brbEvent.Type.(*brbdxrpb.Event_Deliver).Deliver
				fmt.Println("Leader says: ", string(deliverEvent.Data))

			default:
				return fmt.Errorf("unknown brb event type: %T", brbEvent.Type)
			}

		default:
			return fmt.Errorf("unknown event type: %T", event.Type)
		}
	}

	return nil
}

func (m *controlModule) EventsOut() <-chan *events.EventList {
	return m.eventsOut
}

func (m *controlModule) readMessageFromConsole() error {
	// Read the user input
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Type in a message and press Enter: ")
	scanner.Scan()
	if scanner.Err() != nil {
		return fmt.Errorf("error reading from console: %w", scanner.Err())
	}

	m.eventsOut <- events.ListOf(&eventpb.Event{
		DestModule: "brbdxr",
		Type: &eventpb.Event_Brbdxr{
			Brbdxr: &brbdxrpb.Event{
				Type: &brbdxrpb.Event_Request{
					Request: &brbdxrpb.BroadcastRequest{
						Data: []byte(scanner.Text()),
					},
				},
			},
		},
	})

	return nil
}
