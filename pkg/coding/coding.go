package coding

import (
	"fmt"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/codingpb"
	codingpbdsl "github.com/filecoin-project/mir/pkg/pb/codingpb/dsl"
	t "github.com/filecoin-project/mir/pkg/types"
	rs "github.com/vivint/infectious"
)

type ModuleConfig struct {
	Self t.ModuleID // id of this module
}

func NewModule(mc *ModuleConfig) modules.PassiveModule {
	m := dsl.NewModule(mc.Self)

	encoders := make(map[string]*rs.FEC)

	codingpbdsl.UponEncodeRequest(m, func(totalShards int64, dataShards int64, paddedData []uint8, origin *codingpb.EncodeOrigin) error {
		encoder, err := getFEC(int(totalShards), int(dataShards), encoders)
		if err != nil {
			return err
		}

		encoded := make([][]byte, totalShards)

		output := func(s rs.Share) {
			encoded[s.Number] = make([]byte, len(s.Data))
			copy(encoded[s.Number], s.Data)
		}

		err = encoder.Encode(paddedData, output)

		if err != nil {
			return err
		}

		codingpbdsl.EncodeResult(m, t.ModuleID(origin.Module), encoded, origin)

		return nil
	})

	codingpbdsl.UponDecodeRequest(m, func(totalShards int64, dataShards int64, shares []*codingpb.Share, origin *codingpb.DecodeOrigin) error {
		encoder, err := getFEC(int(totalShards), int(dataShards), encoders)
		if err != nil {
			return err
		}

		readys := make([]rs.Share, 0)
		for _, rd := range shares {
			readys = append(readys, rs.Share{
				Number: int(rd.Number),
				Data:   rd.Chunk,
			})
		}

		res, err := encoder.Decode(nil, readys)

		if err != nil {
			codingpbdsl.DecodeResult(m, t.ModuleID(origin.Module), false, nil, origin)
		} else {
			codingpbdsl.DecodeResult(m, t.ModuleID(origin.Module), true, res, origin)
		}

		return nil
	})

	codingpbdsl.UponRebuildRequest(m, func(totalShards int64, dataShards int64, shares []*codingpb.Share, origin *codingpb.RebuildOrigin) error {
		encoder, err := getFEC(int(totalShards), int(dataShards), encoders)
		if err != nil {
			return err
		}

		readys := make([]rs.Share, 0)
		for _, rd := range shares {
			readys = append(readys, rs.Share{
				Number: int(rd.Number),
				Data:   rd.Chunk,
			})
		}

		output := make([]byte, len(readys[0].Data)*int(dataShards))

		err = encoder.Rebuild(readys, func(s rs.Share) {
			copy(output[s.Number*len(s.Data):], s.Data)
		})
		if err != nil {
			codingpbdsl.RebuildResult(m, t.ModuleID(origin.Module), false, nil, origin)
		} else {
			codingpbdsl.RebuildResult(m, t.ModuleID(origin.Module), true, output, origin)
		}

		return nil
	})

	return m
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
