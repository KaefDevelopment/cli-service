package repository

import (
	"gorm.io/gorm"

	"github.com/jaroslav1991/cli-service/internal/service"
)

type TxProvider struct {
	db *gorm.DB
}

func NewTxProvider(db *gorm.DB) *TxProvider {
	return &TxProvider{db: db}
}

func (p *TxProvider) Transaction(f func(txProvider service.TxProvider) error) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		txProvider := &TxProvider{db: tx}
		return f(txProvider)
	})
}

func (p *TxProvider) getTx() *gorm.DB {
	return p.db
}
