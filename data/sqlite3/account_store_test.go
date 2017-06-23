package sqlite3_test

import (
	"testing"

	"github.com/keratin/authn/data"
	"github.com/keratin/authn/data/sqlite3"
	"github.com/keratin/authn/tests"
)

func TestCreate(t *testing.T) {
	db, err := sqlite3.TempDB()
	if err != nil {
		panic(err)
	}
	store := sqlite3.AccountStore{db}

	account, err := store.Create("authn@keratin.tech", []byte("password"))
	if err != nil {
		t.Error(err)
	}
	tests.RefuteEqual(t, 0, account.Id)
	tests.AssertEqual(t, "authn@keratin.tech", account.Username)

	account, err = store.Create("authn@keratin.tech", []byte("password"))
	if account != nil {
		tests.RefuteEqual(t, nil, account)
	}
	if !data.IsUniquenessError(err) {
		t.Errorf("expected uniqueness error, got %T %v", err, err)
	}
}