package main

import (
	"context"
	"fmt"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/brbctpb"
	"github.com/filecoin-project/mir/pkg/pb/brbdxrpb"
	"github.com/filecoin-project/mir/pkg/pb/brbpb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"time"
)

type controlModule struct {
	eventsOut                 chan *events.EventList
	lastId                    int64
	sentMessages              int64
	isLeader                  bool
	currentBenchmark          *Benchmark
	broadcastRequestGenerator func(int64, *[]byte, string) *events.EventList
	broadcastDeliverValidator func([]byte)
	tests                     []*Benchmark
	lastTest                  int
}

func newControlModule(isLeader bool, broadcastRequestGenerator func(int64, *[]byte, string) *events.EventList, broadcastDeliverValidator func([]byte), tests []*Benchmark) modules.ActiveModule {
	return &controlModule{
		eventsOut:                 make(chan *events.EventList),
		lastId:                    0,
		isLeader:                  isLeader,
		broadcastRequestGenerator: broadcastRequestGenerator,
		broadcastDeliverValidator: broadcastDeliverValidator,
		tests:                     tests,
		lastTest:                  0,
	}
}

func (m *controlModule) ImplementsModule() {}

func (m *controlModule) ApplyEvents(ctx context.Context, events *events.EventList) error {
	iter := events.Iterator()
	for event := iter.Next(); event != nil; event = iter.Next() {
		switch event.Type.(type) {

		case *eventpb.Event_Init:
			m.newIteration()

		case *eventpb.Event_Brbct:
			brbEvent := event.Type.(*eventpb.Event_Brbct).Brbct
			switch brbEvent.Type.(type) {

			case *brbctpb.Event_Deliver:
				deliverEvent := brbEvent.Type.(*brbctpb.Event_Deliver).Deliver
				m.broadcastDeliverValidator(deliverEvent.Data)
				m.newIteration()
			default:
				return fmt.Errorf("unknown brb event type: %T", brbEvent.Type)
			}

		case *eventpb.Event_Brbdxr:
			brbEvent := event.Type.(*eventpb.Event_Brbdxr).Brbdxr
			switch brbEvent.Type.(type) {

			case *brbdxrpb.Event_Deliver:
				deliverEvent := brbEvent.Type.(*brbdxrpb.Event_Deliver).Deliver
				m.broadcastDeliverValidator(deliverEvent.Data)
				m.newIteration()
			default:
				return fmt.Errorf("unknown brb event type: %T", brbEvent.Type)
			}

		case *eventpb.Event_Brb:
			brbEvent := event.Type.(*eventpb.Event_Brb).Brb
			switch brbEvent.Type.(type) {

			case *brbpb.Event_Deliver:
				deliverEvent := brbEvent.Type.(*brbpb.Event_Deliver).Deliver
				m.broadcastDeliverValidator(deliverEvent.Data)
				m.newIteration()
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

func (m *controlModule) newIteration() {
	if m.isLeader {
		go func() {
			if m.currentBenchmark == nil {
				if m.lastTest >= len(m.tests) {
					return
				}
				test := m.tests[m.lastTest]
				m.lastTest++
				data := make([]byte, test.messageSize)
				for i := int64(0); i < test.messageSize; i++ {
					data[i] = byte(42)
				}
				m.currentBenchmark = test
				m.currentBenchmark.message = data

				println("Starting benchmark...")
				//p := profile.Start(profile.MemProfile, profile.ProfilePath(fmt.Sprintf("./%s_%d_%d/", m.currentBenchmark.algorithm, iterations, msgSize)))
				m.lastId++
				m.sentMessages++
				m.eventsOut <- m.broadcastRequestGenerator(m.lastId, &m.currentBenchmark.message, m.currentBenchmark.algorithm)
				go func() {
					time.Sleep(m.currentBenchmark.duration)
					fmt.Printf("Total Iterations for %s msgSize=%d: %d\n", m.currentBenchmark.algorithm, m.currentBenchmark.messageSize, m.sentMessages)
					fmt.Printf("  Iterations/sec for %s msgSize=%d: %f\n\n", m.currentBenchmark.algorithm, m.currentBenchmark.messageSize, float64(m.sentMessages)/m.currentBenchmark.duration.Seconds())
					//p.Stop()
					m.currentBenchmark = nil
					m.sentMessages = 0
				}()
			} else {
				m.lastId++
				m.sentMessages++
				m.eventsOut <- m.broadcastRequestGenerator(m.lastId, &m.currentBenchmark.message, m.currentBenchmark.algorithm)
			}
		}()
	}
}
