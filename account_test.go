package vatinator

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	password := "t3stp4ssw0rd"
	email := "Test@test.com"

	tmpdir, err := ioutil.TempDir("", "test-account")
	require.NoError(t, err)
	//t.Logf("temp dir: %s", tmpdir)
	defer os.RemoveAll(tmpdir)

	db, err := NewDB(filepath.Join(tmpdir, "test.db"))
	require.NoError(t, err)
	assert.NoError(t, Migrate(db, "1.sql"))

	// create a new account
	as := NewAccountService(db)
	id, err := as.Create(email, password)
	assert.NoError(t, err)
	assert.Equal(t, AccountID(1), id)

	// check the password
	id2, err := as.CheckPassword("test@test.com", password)
	assert.NoError(t, err)
	assert.Equal(t, AccountID(1), id2)

	// change and check password again
	newPassword := "pass"
	if err := as.SetPassword(email, newPassword); err != nil {
		assert.NoError(t, err)
	}
	id3, err := as.CheckPassword("test@test.com", newPassword)
	assert.NoError(t, err)
	assert.Equal(t, AccountID(1), id3)

	// store some form data
	fd := FormData{
		FirstName:    "Test",
		LastName:     "Guy",
		FullName:     "Test Guy",
		DiplomaticID: "B999900000",
		Embassy:      "US Embassy",
		Address:      "Kentmanni 20",
		Bank:         "Test bank",
		BankName:     "Test Bank",
		Account:      "EE20000000000000",
	}
	require.True(t, fd.IsValid())
	b, err := json.Marshal(fd)
	require.NoError(t, err)
	assert.NoError(t, as.UpdateFormData(AccountID(1), b))

	// read and unmarshal some form data
	b2, err := as.GetFormData(AccountID(1))
	assert.NoError(t, err)
	fd2, err := UnmarshalFormData(b2)
	assert.NoError(t, err)
	assert.Equal(t, fd, fd2)

}
