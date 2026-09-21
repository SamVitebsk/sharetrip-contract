package service

import (
	"context"
	"fmt"
	"time"

	"github.com/SamVitebsk/sharetrip-contract/internal/domain"
	"github.com/google/uuid"
)

type SignContractCommand struct {
	ContractID uuid.UUID
}

type SignContractResult struct {
	ID         uuid.UUID
	ClientID   uuid.UUID
	Status     string
	ValidFrom  time.Time
	ValidUntil *time.Time
	Services   []ContractService
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (s *Service) SignContract(ctx context.Context, command SignContractCommand) (SignContractResult, error) {
	contractID, err := domain.NewContractID(command.ContractID)
	if err != nil {
		return SignContractResult{}, mapSignContractError(err)
	}

	var signedContract domain.Contract
	err = s.txRunner(ctx, func(transactionCtx context.Context, repoTx RepositoryTx) error {
		contract, err := repoTx.GetForUpdate(transactionCtx, contractID)
		if err != nil {
			return fmt.Errorf("get contract for signing: %w", err)
		}

		response, err := contract.SignContract(domain.SignContractRequest{
			Now: time.Now().UTC().Truncate(time.Microsecond),
		})
		if err != nil {
			return mapSignContractError(err)
		}
		savedContract, err := repoTx.SaveSigned(transactionCtx, response.Contract)
		if err != nil {
			return fmt.Errorf("save signed contract: %w", err)
		}
		signedContract = savedContract
		return nil
	})
	if err != nil {
		return SignContractResult{}, fmt.Errorf("sign contract: %w", err)
	}

	return toSignContractResult(signedContract), nil
}
