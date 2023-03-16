package merkletree

import (
	"crypto"
	_ "crypto/sha1"
	"fmt"
	"github.com/filecoin-project/mir/pkg/events"
	"github.com/filecoin-project/mir/pkg/modules"
	"github.com/filecoin-project/mir/pkg/pb/commonpb"
	"github.com/filecoin-project/mir/pkg/pb/eventpb"
	t "github.com/filecoin-project/mir/pkg/types"
	merkle "github.com/wealdtech/go-merkletree"
)

type MerkleProofVerifier struct {
	hashImpl *HashImpl
}

func NewVerifier() *MerkleProofVerifier {
	return &MerkleProofVerifier{
		hashImpl: &HashImpl{
			hasher: crypto.SHA1,
		},
	}
}

type HashImpl struct {
	hasher crypto.Hash
}

func (h *HashImpl) Hash(data []byte) []byte {
	hs := h.hasher.New()
	hs.Write(data)
	hash := hs.Sum(nil)
	return hash
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
		tree, err := merkle.NewUsing(e.MerkleBuildRequest.Messages, m.hashImpl, nil)
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
		return events.ListOf(
			events.MerkleBuildResult(t.ModuleID(e.MerkleBuildRequest.Origin.Module), tree.Root(), proofs, e.MerkleBuildRequest.Origin),
		), nil
	case *eventpb.Event_MerkleVerifyRequest:
		proof, err := merkle.VerifyProofUsing(e.MerkleVerifyRequest.Chunk, &merkle.Proof{
			Hashes: e.MerkleVerifyRequest.Proof.Hashes,
			Index:  e.MerkleVerifyRequest.Proof.Index,
		}, e.MerkleVerifyRequest.RootHash, m.hashImpl, nil)
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
