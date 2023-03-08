package merkletree

import (
	"fmt"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/commonpb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	t "github.com/filecoin-project/mir/pkg/types"
	merkle "github.com/wealdtech/go-merkletree"
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
	case *eventpb.Event_MerkleBuildRequest:
		tree, err := merkle.New(e.MerkleBuildRequest.Messages)
		proofs := make([]*commonpb.MerklePath, 0)
		if err != nil {
			return nil, err
		}
		for _, message := range e.MerkleBuildRequest.Messages {
			proof, err := tree.GenerateProof(message)
			if err != nil {
				return nil, err
			}
			proofs = append(proofs, &commonpb.MerklePath{
				Hashes: proof.Hashes,
				Index:  proof.Index,
			})
		}
		//hashes := make([]*merkle.MerkleHash, len(e.MerkleBuildRequest.Messages))
		//for i, msg := range e.MerkleBuildRequest.Messages {
		//	hashes[i] = merkle.ComputeMerkleHash(msg)
		//}
		//tree := merkle.NewFullMerkleTree(hashes...)
		//paths := make([]*commonpb.MerklePath, len(e.MerkleBuildRequest.Messages))
		//for i, hash := range hashes {
		//	path, err := tree.GetMerklePath(hash)
		//	if err != nil {
		//		return nil, errors.Wrap(err, "Unexpected hash")
		//	}
		//	tmp := make([][]byte, len(path.Hashes))
		//	for j, h := range path.Hashes {
		//		tmp[j] = (*h)[:]
		//	}
		//	paths[i] = &commonpb.MerklePath{
		//		Hashes: tmp,
		//		Proof:  path.ProofPath,
		//	}
		//}
		//tree.GetMerklePath()
		return events.ListOf(
			events.MerkleBuildResult(t.ModuleID(e.MerkleBuildRequest.Origin.Module), tree.Root(), proofs, e.MerkleBuildRequest.Origin),
		), nil
	case *eventpb.Event_MerkleVerifyRequest:
		proof, err := merkle.VerifyProof(e.MerkleVerifyRequest.Chunk, &merkle.Proof{
			Hashes: e.MerkleVerifyRequest.Proof.Hashes,
			Index:  e.MerkleVerifyRequest.Proof.Index,
		}, e.MerkleVerifyRequest.RootHash)
		if err != nil {
			return nil, err
		}

		return events.ListOf(
			events.MerkleProofVerifyResult(t.ModuleID(e.MerkleVerifyRequest.Origin.Module), proof, e.MerkleVerifyRequest.Origin),
		), nil
	default:
		// Complain about all other incoming event types.
		return nil, fmt.Errorf("unexpected event: %T", event.Type)
	}
}

// The ImplementsModule method only serves the purpose of indicating that this is a Module and must not be called.
func (m *MerkleProofVerifier) ImplementsModule() {}
