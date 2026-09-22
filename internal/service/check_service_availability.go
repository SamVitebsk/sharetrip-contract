package service

import (
	"context"
	"errors"
	"time"

	"github.com/SamVitebsk/sharetrip-contract/internal/domain"
	"github.com/google/uuid"
)

type CheckServiceAvailabilityCommand struct {
	ClientID    uuid.UUID
	ServiceCode string
}

type CheckServiceAvailabilityResult struct {
	Allowed    bool
	Reason     string
	ContractID *uuid.UUID
}

func (s *Service) CheckServiceAvailability(ctx context.Context, cmd CheckServiceAvailabilityCommand) (CheckServiceAvailabilityResult, error) {
	clientID, err := domain.NewClientID(cmd.ClientID)
	if err != nil {
		return CheckServiceAvailabilityResult{}, mapCheckServiceAvailabilityError(err)
	}

	serviceCode := domain.ContractServiceCode(cmd.ServiceCode)
	if err := serviceCode.Validate(); err != nil {
		return CheckServiceAvailabilityResult{}, mapCheckServiceAvailabilityError(err)
	}

	var result CheckServiceAvailabilityResult
	err = s.txRunner(ctx, func(txCtx context.Context, repo RepositoryTx) error {
		contract, repoErr := repo.GetClientActiveContract(txCtx, clientID)
		if repoErr != nil {
			if errors.Is(repoErr, ErrNotFound) {
				result = CheckServiceAvailabilityResult{
					Allowed: false,
					Reason:  string(domain.ReasonContractNotFound),
				}
				return nil
			}
			return repoErr
		}

		allowed, reason := contract.CheckService(serviceCode, time.Now())

		contractID := contract.ID.Value()
		result = CheckServiceAvailabilityResult{
			Allowed:    allowed,
			Reason:     string(reason),
			ContractID: &contractID,
		}

		return nil
	})

	if err != nil {
		return CheckServiceAvailabilityResult{}, err
	}

	return result, nil
}
