// +build database

package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := setupIntegrationTests()

	if err != nil {
		fmt.Println(fmt.Sprintf("error setting up integration tests: %s", err))
		os.Exit(-1)
	}
	os.Exit(m.Run())
}
