package db

import (
	"fmt"
	"net/url"
	"os"

	"github.com/wfen/go-grpc-services-course/internal/rocket"

	"github.com/jmoiron/sqlx"
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

func (s Store) GetRocketByID(id string) (rocket.Rocket, error) {
	return rocket.Rocket{}, nil
}

func (s Store) InsertRocket(rkt rocket.Rocket) (rocket.Rocket, error) {
	return rocket.Rocket{}, nil
}

func (s Store) DeleteRocket(id string) error {
	return nil
}
