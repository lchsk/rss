package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSortedMigrationFiles(t *testing.T) {
	file_1 := "11_test.sql"
	file_2 := "test.sql"
	file_3 := "2_test.sql"

	filenames := sortMigrationFiles([]string{file_1, file_2, file_3})

	assert.Equal(t, []string{file_2, file_3, file_1}, filenames)
}

func TestValidateFilenames__no_number_in_filename(t *testing.T) {
	err := validateFilenames([]string{"1_test.sql", "test.sql"})

	assert.NotNil(t, err)
}

func TestValidateFilenames__success(t *testing.T) {
	err := validateFilenames([]string{"1_test.sql", "2_test.sql"})

	assert.Nil(t, err)
}
