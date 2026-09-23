package service

import (
	"context"
	"errors"
	"time"

	"github.com/SamVitebsk/sharetrip-contract/internal/domain"
	"github.com/google/uuid"
)

type CheckServiceAvailabilityRequest struct {
	ClientID    uuid.UUID
	ServiceCode string
}

type CheckServiceAvailabilityResponse struct {
	Allowed    bool
	Reason     string
	ContractID *uuid.UUID
}

func (s *Service) CheckServiceAvailability(ctx context.Context, req CheckServiceAvailabilityRequest) (CheckServiceAvailabilityResponse, error) {
	clientID, err := domain.NewClientID(req.ClientID)
	if err != nil {
		return CheckServiceAvailabilityResponse{}, mapCheckServiceAvailabilityError(err)
	}

	serviceCode := domain.ContractServiceCode(req.ServiceCode)
	if err := serviceCode.Validate(); err != nil {
		return CheckServiceAvailabilityResponse{}, mapCheckServiceAvailabilityError(err)
	}

	var response CheckServiceAvailabilityResponse
	err = s.txRunner(ctx, func(txCtx context.Context, repo RepositoryTx) error {
		contract, repoErr := repo.GetClientActiveContract(txCtx, clientID)
		if repoErr != nil {
			if errors.Is(repoErr, ErrNotFound) {
				response = CheckServiceAvailabilityResponse{
					Allowed: false,
					Reason:  string(domain.ReasonContractNotFound),
				}
				return nil
			}
			return repoErr
		}

		allowed, reason := contract.CheckService(serviceCode, time.Now())

		contractID := contract.ID.Value()
		response = CheckServiceAvailabilityResponse{
			Allowed:    allowed,
			Reason:     string(reason),
			ContractID: &contractID,
		}

		return nil
	})

	if err != nil {
		return CheckServiceAvailabilityResponse{}, err
	}

	return response, nil
}
