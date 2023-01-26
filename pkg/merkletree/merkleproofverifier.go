package merkletree

import (
	"fmt"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/commonpb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	t "github.com/filecoin-project/mir/pkg/types"
	"github.com/multivactech/MultiVAC/model/merkle"
	"github.com/pkg/errors"
)

type MerkleProofVerifier struct{}

func NewVerifier() *MerkleProofVerifier {
	return &MerkleProofVerifier{}
}

func (m *MerkleProofVerifier) ApplyEvents(eventsIn *events.EventList) (*events.EventList, error) {
	return modules.ApplyEventsConcurrently(eventsIn, m.ApplyEvent)
}

func (m *MerkleProofVerifier) ApplyEvent(event *eventpb.Event) (*events.EventList, error) {
	switch e := event.Type.(type) {
	case *eventpb.Event_Init:
		// no actions on init
		return events.EmptyList(), nil
	case *eventpb.Event_MerkelBuildRequest:
		hashes := make([]*merkle.MerkleHash, len(e.MerkelBuildRequest.Messages))
		for i, msg := range e.MerkelBuildRequest.Messages {
			hashes[i] = merkle.ComputeMerkleHash(msg)
		}
		tree := merkle.NewFullMerkleTree(hashes...)
		paths := make([]*commonpb.MerklePath, len(e.MerkelBuildRequest.Messages))
		for i, hash := range hashes {
			path, err := tree.GetMerklePath(hash)
			if err != nil {
				return nil, errors.Wrap(err, "Unexpected hash")
			}
			tmp := make([][]byte, len(path.Hashes))
			for j, h := range path.Hashes {
				tmp[j] = (*h)[:]
			}
			paths[i] = &commonpb.MerklePath{
				Hashes: tmp,
				Proof:  path.ProofPath,
			}
		}
		//tree.GetMerklePath()
		return events.ListOf(
			events.MerkleBuildResult(t.ModuleID(e.MerkelBuildRequest.Origin.Module), (*tree.MerkleRoot)[:], paths, e.MerkelBuildRequest.Origin),
		), nil
	case *eventpb.Event_MerkelVerifyRequest:
		// TODO
		hashes := make([]*merkle.MerkleHash, len(e.MerkelVerifyRequest.Proof.Hashes))
		for i, hash := range e.MerkelVerifyRequest.Proof.Hashes {
			hashes[i] = (*merkle.MerkleHash)(hash)
		}
		merklePath := merkle.MerklePath{
			Hashes:    hashes,
			ProofPath: e.MerkelVerifyRequest.Proof.Proof,
		}
		err := merklePath.Verify((*merkle.MerkleHash)(e.MerkelVerifyRequest.RootHash))

		var result bool
		if err != nil {
			result = false
		} else {
			result = true
		}
		return events.ListOf(
			events.MerkelProofVerifyResult(t.ModuleID(e.MerkelVerifyRequest.Origin.Module), result, e.MerkelVerifyRequest.Origin),
		), nil
	default:
		// Complain about all other incoming event types.
		return nil, fmt.Errorf("unexpected event: %T", event.Type)
	}
}

// The ImplementsModule method only serves the purpose of indicating that this is a Module and must not be called.
func (m *MerkleProofVerifier) ImplementsModule() {}
