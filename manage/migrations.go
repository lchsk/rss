package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/lchsk/rss/db"
)

const MigrationsPath = "../sql/migrations"

func sortMigrationFiles(filenames []string) []string {
	sort.Slice(filenames, func(i, j int) bool {
		firstSplit := strings.Split(filenames[i], "_")
		secondSplit := strings.Split(filenames[j], "_")

		firstNumber, _ := strconv.Atoi(firstSplit[0])
		secondNumber, _ := strconv.Atoi(secondSplit[0])

		return firstNumber < secondNumber

	})
	return filenames
}

func validateFilenames(filenames []string) error {
	for _, f := range filenames {
		split := strings.Split(f, "_")

		if len(split) < 2 {
			return errors.New("Filenames should be in format <number>_<name>.sql")
		}

		_, err := strconv.Atoi(split[0])

		if err != nil {
			return err
		}
	}

	return nil
}

func getSortedMigrationFiles(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Printf("Error loading migration files: %s", err)
		return []string{}, err
	}

	var filenames []string

	for _, f := range files {
		name := f.Name()
		ext := f.Name()[len(name)-3:]

		if ext == "sql" {
			filenames = append(filenames, name)
		}
	}

	if err := validateFilenames(filenames); err != nil {
		return []string{}, err
	}

	return sortMigrationFiles(filenames), nil
}

func Migrate(dba *db.DbAccess) {
	rows, err := dba.DB.Query("select filename from migrations")

	if err != nil {
		log.Printf("Error loading migrations from db: %s", err)
		return
	}

	defer rows.Close()

	migrationsRan := make(map[string]struct{})

	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			log.Printf("Error getting migration filename: %s", err)
			return
		}

		log.Printf("Migration already ran: %s", filename)
		migrationsRan[filename] = struct{}{}
	}

	filenames, err := getSortedMigrationFiles(MigrationsPath)

	if err != nil {
		log.Printf("Error preparing sorted migration filenames: %s", err)
		return
	}

	for _, filename := range filenames {
		data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", MigrationsPath, filename))
		if err != nil {
			log.Printf("Error reading migration file: %s %s", filename, err)
			return
		}

		if _, ok := migrationsRan[filename]; !ok {
			log.Printf("Running migration %s", filename)

			txn, err := dba.DB.Begin()

			if err != nil {
				log.Printf("Error opening transaction: %s %s", filename, err)
				return
			}

			_, err = txn.Exec(string(data))

			rollback := true

			if err != nil {
				log.Printf("Error reading migration sql: %s %s", filename, err)
			} else {
				log.Printf("Migration %s ran successfully", filename)
				_, err := txn.Exec("insert into migrations (id, filename) values ($1, $2)", uuid.New().String(), filename)

				if err != nil {
					log.Printf("Error recording successful migration: %s %s", filename, err)
				} else {
					if err := txn.Commit(); err != nil {
						log.Printf("Error commiting transaction: %s", filename)
					}

					rollback = false
				}
			}

			if rollback {
				if err := txn.Rollback(); err != nil {
					log.Printf("Error rolling back transaction: %s", filename)
				}
			}
		}
	}
}
