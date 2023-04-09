package coding

import (
	"fmt"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	t "github.com/filecoin-project/mir/pkg/types"
	rs "github.com/vivint/infectious"
)

type Coder struct {
	coders map[string]*rs.FEC
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
		encoder, err := getFEC(int(e.EncodeRequest.TotalShards), int(e.EncodeRequest.DataShards), c.coders)
		if err != nil {
			return nil, err
		}

		encoded := make([][]byte, e.EncodeRequest.TotalShards)

		output := func(s rs.Share) {
			encoded[s.Number] = make([]byte, len(s.Data))
			copy(encoded[s.Number], s.Data)
		}

		err = encoder.Encode(e.EncodeRequest.PaddedData, output)

		if err != nil {
			return nil, err
		}

		return events.ListOf(events.EncodeResult(t.ModuleID(e.EncodeRequest.Origin.Module), encoded, e.EncodeRequest.Origin)), nil
	case *eventpb.Event_DecodeRequest:
		encoder, err := getFEC(int(e.DecodeRequest.TotalShards), int(e.DecodeRequest.DataShards), c.coders)
		if err != nil {
			return nil, err
		}

		readys := make([]rs.Share, 0)
		for _, rd := range e.DecodeRequest.Shares {
			readys = append(readys, rs.Share{
				Number: int(rd.Number),
				Data:   rd.Chunk,
			})
		}

		res, err := encoder.Decode(nil, readys)

		if err != nil {
			return events.ListOf(events.DecodeResult(t.ModuleID(e.DecodeRequest.Origin.Module), false, nil, e.DecodeRequest.Origin)), nil
		} else {
			return events.ListOf(events.DecodeResult(t.ModuleID(e.DecodeRequest.Origin.Module), true, res, e.DecodeRequest.Origin)), nil
		}
	case *eventpb.Event_RebuildRequest:
		encoder, err := getFEC(int(e.RebuildRequest.TotalShards), int(e.RebuildRequest.DataShards), c.coders)
		if err != nil {
			return nil, err
		}

		readys := make([]rs.Share, 0)
		for _, rd := range e.RebuildRequest.Shares {
			readys = append(readys, rs.Share{
				Number: int(rd.Number),
				Data:   rd.Chunk,
			})
		}

		output := make([]byte, len(readys[0].Data)*int(e.RebuildRequest.DataShards))

		err = encoder.Rebuild(readys, func(s rs.Share) {
			copy(output[s.Number*len(s.Data):], s.Data)
		})
		if err != nil {
			return events.ListOf(events.RebuildResult(t.ModuleID(e.RebuildRequest.Origin.Module), false, nil, e.RebuildRequest.Origin)), nil
		} else {
			return events.ListOf(events.RebuildResult(t.ModuleID(e.RebuildRequest.Origin.Module), true, output, e.RebuildRequest.Origin)), nil
		}
	default:
		// Complain about all other incoming event types.
		return nil, fmt.Errorf("unexpected event: %T", event.Type)
	}
}

func (c *Coder) ImplementsModule() {}

func NewModule() *Coder {
	return &Coder{coders: make(map[string]*rs.FEC)}
}

func getFEC(n, k int, fecs map[string]*rs.FEC) (*rs.FEC, error) {
	stringRepr := fmt.Sprintf("%d_%d", n, k)
	if fec, ok := fecs[stringRepr]; ok {
		return fec, nil
	} else {
		fec, err := rs.NewFEC(k, n)
		fecs[stringRepr] = fec
		return fec, err
	}
}
