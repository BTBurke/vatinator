package vatinator

import (
	"errors"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gorilla/securecookie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionService(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "test-session")
	require.NoError(t, err)
	ss, err := NewSessionService(filepath.Join(tmpdir, "session.db"), securecookie.GenerateRandomKey(32))
	require.NoError(t, err)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	_, errNoSession := ss.Get(w, r)
	assert.True(t, errors.Is(errNoSession, SessionNotValid))

	if err := ss.New(w, r, AccountID(1)); err != nil {
		require.NoError(t, err)
	}
	assert.Greater(t, len(w.Header().Get("Set-Cookie")), 0)

	r2 := httptest.NewRequest("GET", "/", nil)
	r2.Header.Set("Cookie", w.Header().Get("Set-Cookie"))
	id, err := ss.Get(w, r2)
	assert.Equal(t, AccountID(1), id)
	assert.NoError(t, err)
}

func TestSessionKeys(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "test-session")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	db, err := NewDB(filepath.Join(tmpdir, "test.db"))
	require.NoError(t, err)
	assert.NoError(t, Migrate(db, "1.sql"))

	// should create a session key on startup
	keys, err := GetSessionKeys(db)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(keys))

	var b []byte
	if err := db.Get(&b, "SELECT key FROM keys WHERE id = $1", 1); err != nil {
		require.NoError(t, err)
	}
	assert.Equal(t, len(b), 32)
}
