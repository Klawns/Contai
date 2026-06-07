package ports

import "contai/internal/commitments/domain"

type CommitmentIDGenerator interface {
	NewCommitmentID() domain.CommitmentID
}
