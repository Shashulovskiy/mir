package coding

import (
	"fmt"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
)

type Coder struct{}

func NewCoder() *Coder {
	return &Coder{}
}

func (c *Coder) ApplyEvents(eventsIn *events.EventList) (*events.EventList, error) {
	return modules.ApplyEventsConcurrently(eventsIn, c.ApplyEvent)
}

func (c *Coder) ApplyEvent(event *eventpb.Event) (*events.EventList, error) {
	switch e := event.Type.(type) {
	case *eventpb.Event_Init:
		// no actions on init
		return events.EmptyList(), nil

	case *eventpb.Event_DataEncodeRequest:
		// TODO
		println(e)
		return events.ListOf(
		//events.MerkleProofVerifyResult(t.ModuleID(e.DataEncodeRequest.Origin.Module), result, e.MerkleVerifyRequest.Origin),
		), nil
	default:
		// Complain about all other incoming event types.
		return nil, fmt.Errorf("unexpected event: %T", event.Type)
	}
}

// The ImplementsModule method only serves the purpose of indicating that this is a Module and must not be called.
func (c *Coder) ImplementsModule() {}
