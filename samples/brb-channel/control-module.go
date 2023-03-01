package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/brbctpb"
	"github.com/filecoin-project/mir/pkg/pb/brbdxrpb"
	"github.com/filecoin-project/mir/pkg/pb/brbpb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	"os"
	"strconv"
	"strings"
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
}

func newControlModule(isLeader bool,
	broadcastRequestGenerator func(int64, *[]byte, string) *events.EventList,
	broadcastDeliverValidator func([]byte),
) modules.ActiveModule {
	return &controlModule{
		eventsOut:                 make(chan *events.EventList),
		lastId:                    0,
		isLeader:                  isLeader,
		broadcastRequestGenerator: broadcastRequestGenerator,
		broadcastDeliverValidator: broadcastDeliverValidator,
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
				reader := bufio.NewReader(os.Stdin)
				println("Input benchmark params: [msgSize] [duration] [algorithm]")
				input, err := reader.ReadString('\n')
				if err != nil {
					return
				}
				split := strings.Split(input[:len(input)-1], " ")
				msgSize, err := strconv.ParseInt(split[0], 10, 64)
				if err != nil {
					return
				}
				iterations, err := strconv.ParseInt(split[1], 10, 64)
				if err != nil {
					return
				}
				data := make([]byte, msgSize)
				for i := int64(0); i < msgSize; i++ {
					data[i] = byte(42)
				}
				m.currentBenchmark = &Benchmark{
					message:   &data,
					duration:  time.Duration(iterations * time.Second.Nanoseconds()),
					algorithm: split[2],
				}
				println("Starting benchmark...")
				//p := profile.Start(profile.MemProfile, profile.ProfilePath(fmt.Sprintf("./%s_%d_%d/", m.currentBenchmark.algorithm, iterations, msgSize)))
				m.lastId++
				m.sentMessages++
				m.eventsOut <- m.broadcastRequestGenerator(m.lastId, m.currentBenchmark.message, m.currentBenchmark.algorithm)
				go func() {
					time.Sleep(m.currentBenchmark.duration)
					fmt.Printf("Total Iterations for %s msgSize=%d: %d\n", split[2], msgSize, m.sentMessages)
					fmt.Printf("  Iterations/sec for %s msgSize=%d: %f\n\n", split[2], msgSize, float64(m.sentMessages)/float64(iterations))
					//p.Stop()
					m.currentBenchmark = nil
					m.sentMessages = 0
				}()
			} else {
				m.lastId++
				m.sentMessages++
				m.eventsOut <- m.broadcastRequestGenerator(m.lastId, m.currentBenchmark.message, m.currentBenchmark.algorithm)
			}
		}()
	}
}

//128 60 brb
//256 60 brb
//512 60 brb
//1024 60 brb
//2048 60 brb
//4096 60 brb
//8192 60 brb
//16384 60 brb
//32768 60 brb
//32768 60 brb
//65536 60 brb
//131072 60 brb
