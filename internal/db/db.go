package db

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	"github.com/wfen/go-grpc-services-course/internal/rocket"
)

type Store struct {
	db *sqlx.DB
}

// New - returns a new store ore error
func New() (Store, error) {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbTable := os.Getenv("DB_TABLE")
	sslMode := os.Getenv("SSL_MODE")
	sslRootCert := os.Getenv("SSL_ROOTCERT")
	options := os.Getenv("OPTIONS")

	urlValues := url.Values{"sslmode": []string{sslMode}}
	if sslRootCert != "" {
		urlValues["sslrootcert"] = []string{sslRootCert}
	}
	if options != "" {
		urlValues["options"] = []string{options}
	}
	connectionString := url.URL{
		User:     url.UserPassword(dbUsername, dbPassword),
		Scheme:   "postgresql",
		Host:     fmt.Sprintf("%s:%s", dbHost, dbPort),
		Path:     dbTable,
		RawQuery: urlValues.Encode(),
	}

	db, err := sqlx.Connect("postgres", connectionString.String())
	if err != nil {
		return Store{}, err
	}
	return Store{
		db: db,
	}, nil
}

// GetRocketByID - retrieves a rocket from the database by id
func (s Store) GetRocketByID(id string) (rocket.Rocket, error) {
	var rkt rocket.Rocket
	row := s.db.QueryRow(`
SELECT id, type, name FROM rockets where id=$1;`,
		id,
	)
	err := row.Scan(&rkt.ID, &rkt.Type, &rkt.Name)
	if err != nil {
		log.Print(err.Error())
		return rocket.Rocket{}, err
	}
	return rkt, nil
}

// InsertRocket - inserts a rocket into the rockets table.
func (s Store) InsertRocket(rkt rocket.Rocket) (rocket.Rocket, error) {
	_, err := s.db.NamedQuery(`
INSERT INTO rockets
(id, name, type)
VALUES (:id, :name, :type)`,
		rkt,
	)
	if err != nil {
		return rocket.Rocket{}, errors.New("failed to insert into database")
	}
	return rocket.Rocket{
		ID:   rkt.ID,
		Type: rkt.Type,
		Name: rkt.Name,
	}, nil
}

// DeleteRocket - attempts to delete a rocket from the database, returning err if error
func (s Store) DeleteRocket(id string) error {
	uid, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
DELETE FROM rockets WHERE id = $1;`,
		uid,
	)
	if err != nil {
		return err
	}
	return nil
}
