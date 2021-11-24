package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/nure-mates/api/src/config"
)

// Postgres bun connection.
type Postgres struct {
	DB  *bun.DB
	ctx context.Context
}

func New(ctx context.Context, wg *sync.WaitGroup, mainCfg *config.Postgres) (*Postgres, error) {
	pgconn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(fmt.Sprintf("%s:%s", mainCfg.Host, mainCfg.Port)),
		pgdriver.WithUser(mainCfg.User),
		pgdriver.WithPassword(mainCfg.Password),
		pgdriver.WithDatabase(mainCfg.Name),
		pgdriver.WithReadTimeout(mainCfg.ReadTimeout),
		pgdriver.WithWriteTimeout(mainCfg.WriteTimeout),
		pgdriver.WithInsecure(mainCfg.SSLMode),
	)
	sqldb := sql.OpenDB(pgconn)

	db := bun.NewDB(sqldb, pgdialect.New())

	p := &Postgres{DB: db, ctx: ctx}

	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()

		err := db.Close()
		if err != nil {
			log.Error("close db connection error:", err.Error())

			return
		}

		log.Info("close db connection")
	}()

	return p, nil
}

// Check checks db connection.
func (p *Postgres) Check() (err error) {
	return p.DB.Ping()
}
