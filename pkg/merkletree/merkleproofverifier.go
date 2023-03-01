package merkletree

import (
	"fmt"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	t "github.com/filecoin-project/mir/pkg/types"
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
		//events.MerkleBuildResult(t.ModuleID(e.MerkleBuildRequest.Origin.Module), (*tree.MerkleRoot)[:], paths, e.MerkleBuildRequest.Origin),
		), nil
	case *eventpb.Event_MerkleVerifyRequest:
		// TODO
		//hashes := make([]*merkle.MerkleHash, len(e.MerkleVerifyRequest.Proof.Hashes))
		//for i, hash := range e.MerkleVerifyRequest.Proof.Hashes {
		//	hashes[i] = (*merkle.MerkleHash)(hash)
		//}
		//merklePath := merkle.MerklePath{
		//	Hashes:    hashes,
		//	ProofPath: e.MerkleVerifyRequest.Proof.Proof,
		//}
		//err := merklePath.Verify((*merkle.MerkleHash)(e.MerkleVerifyRequest.RootHash))

		//var result bool
		//if err != nil {
		//	result = false
		//} else {
		//	result = true
		//}
		return events.ListOf(
			events.MerkleProofVerifyResult(t.ModuleID(e.MerkleVerifyRequest.Origin.Module), true, e.MerkleVerifyRequest.Origin),
		), nil
	default:
		// Complain about all other incoming event types.
		return nil, fmt.Errorf("unexpected event: %T", event.Type)
	}
}

// The ImplementsModule method only serves the purpose of indicating that this is a Module and must not be called.
func (m *MerkleProofVerifier) ImplementsModule() {}
