package db

import (
	"github.com/go-pg/pg/v9"
	"github.com/sirupsen/logrus"
)

type (
	Transaction interface {
		Rollback()
		Commit() error
		GetTx() *pg.Tx
	}

	PgTx struct {
		dbTx        *pg.Tx
		isCommitted bool
	}
)

func (d *BasePgDAO) BeginTx(shardId ...int64) (Transaction, error) {
	var (
		tx  *pg.Tx
		err error
	)

	if len(shardId) > 0 {
		tx, err = d.Cluster.Shard(shardId[0]).Begin()
	} else {
		tx, err = d.BaseDB.Begin()
	}
	if err != nil {
		return nil, err
	}

	return &PgTx{
		dbTx:        tx,
		isCommitted: false,
	}, nil
}

func (tx *PgTx) GetTx() *pg.Tx {
	return tx.dbTx
}

func (tx *PgTx) Rollback() {
	if !tx.isCommitted {
		err := tx.dbTx.Rollback()
		if err != nil {
			logrus.Fatalf("cannot rollback the DB transaction: %s", err.Error())
		}
	}
}

func (tx *PgTx) Commit() error {
	err := tx.dbTx.Commit()
	if err != nil {
		return err
	}

	tx.isCommitted = true
	return nil
}
