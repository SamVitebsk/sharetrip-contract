package domain

import "errors"

var (
	ErrInvalidContractTime   = errors.New("время операции не должно быть пустым")
	ErrInvalidContractStart  = errors.New("дата начала договора обязательна")
	ErrInvalidContractPeriod = errors.New("дата окончания договора должна быть позже начала")

	ErrEmptyContractServices    = errors.New("договор должен содержать хотя бы одну услугу")
	ErrDuplicateContractService = errors.New("услуги договора не должны повторяться")

	ErrContractNotDraft = errors.New("подписать можно только черновик договора")
	ErrContractExpired  = errors.New("срок действия договора истёк")

	ErrInvalidClientID   = errors.New("идентификатор клиента не должен быть пустым")
	ErrInvalidContractID = errors.New("идентификатор договора не должен быть пустым")

	ErrInvalidContractServiceCode = errors.New("неизвестный код услуги")

	ErrInvalidContractStatus = errors.New("неизвестный статус договора")
)
