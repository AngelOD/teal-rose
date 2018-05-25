package main

import "testing"

func TestMysqlStoreDataRunner(t *testing.T) {
	initService()

	if !loadDotEnv(".test.env") {
		t.Error("Unable to load ENV-file.")
	}
}
