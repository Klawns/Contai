package ids

import (
	"contai/internal/commitments/app/ports"
	"contai/internal/commitments/domain"

	"github.com/google/uuid"
)

var _ ports.CommitmentIDGenerator = UUIDCommitmentIDGenerator{}

type UUIDCommitmentIDGenerator struct{}

func NewUUIDCommitmentIDGenerator() UUIDCommitmentIDGenerator {
	return UUIDCommitmentIDGenerator{}
}

func (generator UUIDCommitmentIDGenerator) NewCommitmentID() domain.CommitmentID {
	return domain.CommitmentID(uuid.NewString())
}
