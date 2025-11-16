package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/avito/internship/pr-service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

const pingTimeout = 500 * time.Millisecond

func InitDBPool(cnf *config.Config) *pgxpool.Pool {
	pgxConf := createPgxConf(cnf)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbPool, err := pgxpool.NewWithConfig(ctx, pgxConf)
	if err != nil {
		log.Fatalf("unable to create connection pool: %v\n", err)
	}

	pingAttempts := 3
	pingConnection(pingAttempts, dbPool)

	log.Println("Database connection pool established")
	return dbPool
}

func createPgxConf(cnf *config.Config) *pgxpool.Config {
	connString := fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		cnf.Database.Driver,
		cnf.Database.User,
		cnf.Database.Password,
		cnf.Database.Host,
		cnf.Database.Port,
		cnf.Database.DB,
	)
	pgxConf, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatalf("unable to parse connection string: %v\n", err)
	}
	return pgxConf
}

func pingConnection(pingAttempts int, dbPool *pgxpool.Pool) {
	var pingErr error
	for i := 1; i <= pingAttempts; i++ {
		pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
		pingErr = dbPool.Ping(pingCtx)
		pingCancel()
		if pingErr == nil {
			pingErr = nil
			break
		}
		log.Printf("ping attempt %d/%d failed: %v. Retrying...\n", i, pingAttempts, pingErr)
		if i < pingAttempts {
			time.Sleep(pingTimeout)
		}
	}
	if pingErr != nil {
		log.Fatalf("unable to ping database: %v\n", pingErr)
	}
}
