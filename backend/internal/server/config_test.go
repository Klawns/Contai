package server

import "testing"

func TestGetDBAutoMigrateDefaultsToEnabledOutsideProduction(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("GIN_MODE", "")
	t.Setenv("DB_AUTO_MIGRATE", "")

	if !getDBAutoMigrate() {
		t.Fatal("expected auto migrate to be enabled outside production")
	}
}

func TestGetDBAutoMigrateDefaultsToDisabledInProduction(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("GIN_MODE", "")
	t.Setenv("DB_AUTO_MIGRATE", "")

	if getDBAutoMigrate() {
		t.Fatal("expected auto migrate to be disabled in production")
	}
}

func TestGetDBAutoMigrateCanBeEnabledInProduction(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("GIN_MODE", "")
	t.Setenv("DB_AUTO_MIGRATE", "true")

	if !getDBAutoMigrate() {
		t.Fatal("expected auto migrate to be enabled by DB_AUTO_MIGRATE=true")
	}
}

func TestGetDBAutoMigrateCanBeDisabledOutsideProduction(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("GIN_MODE", "")
	t.Setenv("DB_AUTO_MIGRATE", "false")

	if getDBAutoMigrate() {
		t.Fatal("expected auto migrate to be disabled by DB_AUTO_MIGRATE=false")
	}
}
