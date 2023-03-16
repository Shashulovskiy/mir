package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/filecoin-project/mir/pkg/pb/brbpb"
	"os"

	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
)

type controlModule struct {
	eventsOut chan *events.EventList
	isLeader  bool
	lastId    int
}

func newControlModule(isLeader bool) modules.ActiveModule {
	return &controlModule{
		eventsOut: make(chan *events.EventList),
		isLeader:  isLeader,
		lastId:    1,
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

		case *eventpb.Event_Brb:
			brbEvent := event.Type.(*eventpb.Event_Brb).Brb
			switch brbEvent.Type.(type) {

			case *brbpb.Event_Deliver:
				deliverEvent := brbEvent.Type.(*brbpb.Event_Deliver).Deliver
				fmt.Println("Leader says: ", string(deliverEvent.Data))
				m.lastId += 1
				if m.isLeader {
					m.readMessageFromConsole()
				}

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
		DestModule: "brbct",
		Type: &eventpb.Event_Brb{
			Brb: &brbpb.Event{
				Type: &brbpb.Event_Request{
					Request: &brbpb.BroadcastRequest{
						Id:   int64(m.lastId),
						Data: []byte(scanner.Text()),
					},
				},
			},
		},
	})

	return nil
}
