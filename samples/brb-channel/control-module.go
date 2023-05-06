package main

import (
	"context"
	"fmt"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/brbpb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
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
	broadcastRequestGenerator func(int64, int64, *[]byte, string) *events.EventList
	broadcastDeliverValidator func([]byte)
	tests                     []*Benchmark
	lastTest                  int
	testStats                 map[string][]testResult
}

type testResult struct {
	operations  int64
	duration    float64
	messageSize int64
}

func newControlModule(isLeader bool, broadcastRequestGenerator func(int64, int64, *[]byte, string) *events.EventList, broadcastDeliverValidator func([]byte), tests []*Benchmark) modules.ActiveModule {
	return &controlModule{
		eventsOut:                 make(chan *events.EventList),
		lastId:                    0,
		isLeader:                  isLeader,
		broadcastRequestGenerator: broadcastRequestGenerator,
		broadcastDeliverValidator: broadcastDeliverValidator,
		tests:                     tests,
		lastTest:                  0,
		testStats:                 make(map[string][]testResult),
	}
}

func (m *controlModule) ImplementsModule() {}

func (m *controlModule) ApplyEvents(ctx context.Context, events *events.EventList) error {
	iter := events.Iterator()
	for event := iter.Next(); event != nil; event = iter.Next() {
		switch event.Type.(type) {

		case *eventpb.Event_Init:
			m.newIteration()

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
					m.printTestResults()
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
				//p := profile.Start(profile.CPUProfile, profile.ProfilePath(fmt.Sprintf("./%s_%d/", m.currentBenchmark.algorithm, m.currentBenchmark.messageSize)))
				//pmem := profile.Start(profile.MemProfile, profile.ProfilePath(fmt.Sprintf("./%s_%d_mem/", m.currentBenchmark.algorithm, m.currentBenchmark.messageSize)))
				m.lastId++
				m.sentMessages++
				m.eventsOut <- m.broadcastRequestGenerator(m.lastId, m.currentBenchmark.n, &m.currentBenchmark.message, m.currentBenchmark.algorithm)
				go func() {
					time.Sleep(m.currentBenchmark.duration)
					fmt.Printf("Total Iterations for %s msgSize=%d: %d\n", m.currentBenchmark.algorithm, m.currentBenchmark.messageSize, m.sentMessages)
					fmt.Printf("  Iterations/sec for %s msgSize=%d: %f\n\n", m.currentBenchmark.algorithm, m.currentBenchmark.messageSize, float64(m.sentMessages)/m.currentBenchmark.duration.Seconds())
					if _, ok := m.testStats[m.currentBenchmark.algorithm]; !ok {
						m.testStats[m.currentBenchmark.algorithm] = make([]testResult, 0)
					}
					m.testStats[m.currentBenchmark.algorithm] = append(m.testStats[m.currentBenchmark.algorithm], testResult{
						operations:  m.sentMessages,
						duration:    m.currentBenchmark.duration.Seconds(),
						messageSize: m.currentBenchmark.messageSize,
					})
					//p.Stop()
					//pmem.Stop()
					m.currentBenchmark = nil
					m.sentMessages = 0
				}()
			} else {
				m.lastId++
				m.sentMessages++
				m.eventsOut <- m.broadcastRequestGenerator(m.lastId, m.currentBenchmark.n, &m.currentBenchmark.message, m.currentBenchmark.algorithm)
			}
		}()
	}
}

func (m *controlModule) printTestResults() {
	for algorithm, testResults := range m.testStats {
		messageSizes := make([]string, 0)
		operationsPerSecond := make([]string, 0)
		totalDataTransmitted := make([]string, 0)
		for _, result := range testResults {
			messageSizes = append(messageSizes, strconv.FormatInt(result.messageSize, 10))
			operationsPerSecond = append(operationsPerSecond, fmt.Sprint(float64(result.operations)/result.duration))
			totalDataTransmitted = append(totalDataTransmitted, fmt.Sprint(float64(result.operations)/result.duration*float64(result.messageSize)))
		}
		fmt.Println()
		fmt.Printf("Algorithm %s: messageSizes: %s\n", algorithm, strings.Join(messageSizes, ","))
		fmt.Printf("        Algorithm %s: Op/S: %s\n", algorithm, strings.Join(operationsPerSecond, ","))
		fmt.Printf(" Algorithm %s: transmitted: %s\n", algorithm, strings.Join(totalDataTransmitted, ","))
	}
}
