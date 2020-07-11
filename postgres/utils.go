package postgres

import (
	"github.com/georgysavva/generic-wallet/config"
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"net"
	"strconv"
	"time"
)

func connect(settings *config.Postgres) (*pg.DB, error) {
	timeout := time.Millisecond * time.Duration(settings.Timeout)
	db := pg.Connect(&pg.Options{
		User:         settings.User,
		Password:     settings.Password,
		Database:     settings.Database,
		Addr:         net.JoinHostPort(settings.Host, strconv.Itoa(settings.Port)),
		MaxRetries:   settings.RetriesNum,
		DialTimeout:  timeout,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	})
	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		return nil, errors.Wrap(err, "connection failed")
	}
	return db, nil
}
