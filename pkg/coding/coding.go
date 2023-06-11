package coding

import (
	"fmt"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	rs "github.com/filecoin-project/mir/pkg/rs_ezpwd"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/pkg/errors"
)

type Coder struct {
}

func (c *Coder) ApplyEvents(eventsIn *events.EventList) (*events.EventList, error) {
	return modules.ApplyEventsConcurrently(eventsIn, c.ApplyEvent)
}

func (c *Coder) ApplyEvent(event *eventpb.Event) (*events.EventList, error) {
	switch e := event.Type.(type) {
	case *eventpb.Event_Init:
		// no actions on init
		return events.EmptyList(), nil
	case *eventpb.Event_EncodeRequest:
		encoded := rs.EncodeWrapper(e.EncodeRequest.TotalShards, (e.EncodeRequest.TotalShards-e.EncodeRequest.DataShards)/2, e.EncodeRequest.PaddedData)

		if len(encoded) == 0 {
			return nil, errors.New("Failed to encode")
		}

		return events.ListOf(events.EncodeResult(t.ModuleID(e.EncodeRequest.Origin.Module), encoded, e.EncodeRequest.Origin)), nil
	case *eventpb.Event_DecodeRequest:
		shares := make([][]byte, e.DecodeRequest.TotalShards)
		missing := make([]int, 0)
		for _, share := range e.DecodeRequest.Shares {
			shares[share.Number] = share.Chunk
		}
		for i := range shares {
			if shares[i] == nil {
				missing = append(missing, i)
			}
		}
		decoded := rs.DecodeWrapper(e.DecodeRequest.TotalShards, (e.DecodeRequest.TotalShards-e.DecodeRequest.DataShards)/2, shares, missing)

		if len(decoded) == 0 {
			return events.ListOf(events.DecodeResult(t.ModuleID(e.DecodeRequest.Origin.Module), false, nil, e.DecodeRequest.Origin)), nil
		} else {
			return events.ListOf(events.DecodeResult(t.ModuleID(e.DecodeRequest.Origin.Module), true, decoded, e.DecodeRequest.Origin)), nil
		}
	//case *eventpb.Event_RebuildRequest:
	//	encoder, err := getFEC(int(e.RebuildRequest.TotalShards), int(e.RebuildRequest.DataShards), c.coders)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	readys := make([]rs.Share, 0)
	//	for _, rd := range e.RebuildRequest.Shares {
	//		readys = append(readys, rs.Share{
	//			Number: int(rd.Number),
	//			Data:   rd.Chunk,
	//		})
	//	}
	//
	//	output := make([]byte, len(readys[0].Data)*int(e.RebuildRequest.DataShards))
	//
	//	err = encoder.Rebuild(readys, func(s rs.Share) {
	//		copy(output[s.Number*len(s.Data):], s.Data)
	//	})
	//	if err != nil {
	//		return events.ListOf(events.RebuildResult(t.ModuleID(e.RebuildRequest.Origin.Module), false, nil, e.RebuildRequest.Origin)), nil
	//	} else {
	//		return events.ListOf(events.RebuildResult(t.ModuleID(e.RebuildRequest.Origin.Module), true, output, e.RebuildRequest.Origin)), nil
	//	}
	default:
		// Complain about all other incoming event types.
		return nil, fmt.Errorf("unexpected event: %T", event.Type)
	}
}

func (c *Coder) ImplementsModule() {}

func NewModule() *Coder {
	return &Coder{}
}
