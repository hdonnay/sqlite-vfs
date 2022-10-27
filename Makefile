testdata/test.db: Makefile
	sqlite3 -echo $@ "CREATE TABLE test(key, value); INSERT INTO test (key, value) VALUES ('a', 'b');"
