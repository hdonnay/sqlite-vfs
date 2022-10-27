package main

import (
	"database/sql"
	"embed"
	"io/fs"
	"os"
	"testing"

	_ "modernc.org/sqlite" // import driver
	"modernc.org/sqlite/vfs"
)

func TestVFS(t *testing.T) {
	t.Run("Trivial", func(t *testing.T) {
		sys := os.DirFS(".")
		name, vfs, err := vfs.New(sys)
		if err != nil {
			t.Fatal(err)
		}
		defer vfs.Close()
		db, err := sql.Open(`sqlite`, `testdata/test.db?immutable=1&vfs=`+name)
		if err != nil {
			t.Fatal(err)
		}
		if err := db.Ping(); err != nil {
			t.Fatal(err)
		}
		var want, got string
		want = `b`
		if err := db.QueryRow(`SELECT value FROM test WHERE key = 'a'`).Scan(&got); err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("%q != %q", got, want)
		}
	})
	t.Run("Prefix", func(t *testing.T) {
		sys := os.DirFS("testdata")
		name, vfs, err := vfs.New(sys)
		if err != nil {
			t.Fatal(err)
		}
		defer vfs.Close()
		db, err := sql.Open(`sqlite`, `test.db?immutable=1&vfs=`+name)
		if err != nil {
			t.Fatal(err)
		}
		if err := db.Ping(); err != nil {
			t.Fatal(err)
		}
		var want, got string
		want = `b`
		if err := db.QueryRow(`SELECT value FROM test WHERE key = 'a'`).Scan(&got); err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("%q != %q", got, want)
		}
	})
	t.Run("Embed", func(t *testing.T) {
		name, vfs, err := vfs.New(testfs)
		if err != nil {
			t.Fatal(err)
		}
		defer vfs.Close()
		db, err := sql.Open(`sqlite`, `testdata/test.db?immutable=1&vfs=`+name)
		if err != nil {
			t.Fatal(err)
		}
		if err := db.Ping(); err != nil {
			t.Fatal(err)
		}
		var want, got string
		want = `b`
		if err := db.QueryRow(`SELECT value FROM test WHERE key = 'a'`).Scan(&got); err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("%q != %q", got, want)
		}
	})
	t.Run("EmbedSub", func(t *testing.T) {
		sys, err := fs.Sub(testfs, `testdata`)
		if err != nil {
			t.Fatal(err)
		}
		name, vfs, err := vfs.New(sys)
		if err != nil {
			t.Fatal(err)
		}
		defer vfs.Close()
		db, err := sql.Open(`sqlite`, `test.db?immutable=1&vfs=`+name)
		if err != nil {
			t.Fatal(err)
		}
		if err := db.Ping(); err != nil {
			t.Fatal(err)
		}
		var want, got string
		want = `b`
		if err := db.QueryRow(`SELECT value FROM test WHERE key = 'a'`).Scan(&got); err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("%q != %q", got, want)
		}
	})
}

//go:embed testdata/test.db
var testfs embed.FS
