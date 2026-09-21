package service

import "errors"

type Service struct {
	txRunner TxRunner
}

func New(txRunner TxRunner) (*Service, error) {
	if txRunner == nil {
		return nil, errors.New("функция выполнения транзакций договора обязательна")
	}
	return &Service{txRunner: txRunner}, nil
}
