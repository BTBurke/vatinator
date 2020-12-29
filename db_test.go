package vatinator

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabase(t *testing.T) {
	st := "CREATE TABLE IF NOT EXISTS foo (id INTEGER PRIMARY KEY AUTOINCREMENT, name STRING, value INTEGER, blob BLOB);"

	tmpdir, err := ioutil.TempDir("", "db-test")
	require.NoError(t, err)
	db, err := NewDB(filepath.Join(tmpdir, "test.db"))
	assert.NoError(t, err)
	defer db.Close()
	defer os.RemoveAll(tmpdir)

	if _, err := db.Exec(st); err != nil {
		require.NoError(t, err)
	}
	res, err := db.Exec("INSERT INTO foo (name, value, blob) VALUES ($1, $2, $3)", "Test", 2, []byte("hello"))
	assert.NoError(t, err)
	id, err := res.LastInsertId()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)

	// tests conversion of BLOB to []byte
	var b []byte
	assert.NoError(t, db.Get(&b, "SELECT blob FROM foo WHERE id=$1;", id))
	assert.Equal(t, []byte("hello"), b)

	// tests conversion of account ids into database
	if _, err := db.Exec("INSERT INTO foo (value) VALUES ($1);", AccountID(2)); err != nil {
		assert.NoError(t, err)
	}

	type row struct {
		ID    AccountID `sql:"id"`
		Name  string
		Value int
		Blob  []byte
	}
	var r row
	assert.NoError(t, db.Get(&r, "SELECT * FROM foo where id=$1", id))
	assert.Equal(t, row{AccountID(1), "Test", 2, []byte("hello")}, r)
}

func TestDatabaseLocks(t *testing.T) {
	st := "CREATE TABLE IF NOT EXISTS foo (id INTEGER PRIMARY KEY AUTOINCREMENT, value INTEGER);"
	tmpdir, err := ioutil.TempDir("", "db-test")
	require.NoError(t, err)
	db, err := NewDB(filepath.Join(tmpdir, "test.db"))
	assert.NoError(t, err)
	defer db.Close()
	defer os.RemoveAll(tmpdir)
	if _, err := db.Exec(st); err != nil {
		require.NoError(t, err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(me int, t *testing.T) {
			defer wg.Done()
			for j := 0; j <= 100; j++ {
				//t.Logf("Routine %d: Inserting %d", me, j)
				if _, err := db.Exec("INSERT INTO foo (value) VALUES ($1)", j); err != nil {
					t.Fatalf("Got error in routine %d: %+v", me, err)
				}
			}
		}(i, t)
	}
	wg.Wait()
}

func TestMigration(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "test-migration")
	require.NoError(t, err)
	path := filepath.Join(tmpdir, "test.db")
	defer os.RemoveAll(tmpdir)

	db, err := NewDB(path)
	require.NoError(t, err)
	assert.NoError(t, Migrate(db, "1.sql"))
}
