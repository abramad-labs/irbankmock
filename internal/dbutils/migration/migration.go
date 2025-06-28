package migration

import (
	"log"

	"gorm.io/gorm"
)

type Migration func(gorm.Migrator) error

var migrations []MigrationNamePair

type MigrationNamePair struct {
	Name      string
	Migration Migration
}

func init() {
	migrations = make([]MigrationNamePair, 0)
}

func RegisterMigration(name string, migration Migration) {
	migrations = append(migrations, MigrationNamePair{
		Name:      name,
		Migration: migration,
	})
}

func ApplyMigrations(migrator gorm.Migrator) error {
	for _, mig := range migrations {
		log.Printf("auto migration: %s", mig.Name)
		err := mig.Migration(migrator)
		if err != nil {
			return err
		}
	}
	migrations = nil
	return nil
}
