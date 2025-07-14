package mysql

import (
	"context"
	"gorm.io/gorm"
)

type Migrator interface {
	Migrate(ctx context.Context) error
}

type migrator struct {
	db *gorm.DB
}

func (m migrator) Migrate(ctx context.Context) error {
	if err := m.db.AutoMigrate(
		new(Account),
		new(Event),
		new(Ticket),
	); err != nil {
		return err
	}
	return nil
}

func NewMigrator(db *gorm.DB) Migrator {
	return &migrator{db: db}
}
