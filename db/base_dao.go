package db

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/wbrush/go-common/config"
	"github.com/wbrush/go-common/datamodels"

	"sync"

	"github.com/go-pg/migrations/v7"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/sharding/v7"
	"github.com/sirupsen/logrus"
)

const (
	MigrationsBaseDir  = "base_migrations"
	MigrationsShardDir = "shard_migrations"
)

const (
	ErrNoShardsYet = "no shards was created yet"
)

const (
	IdleTimeout        = 30
	IdleCheckFrequency = 15
	ReadTimeout        = 50
	WriteTimeout       = 50
	MaxRetries         = 5
	PoolTimeout        = ReadTimeout + 1
)

type (
	BaseDataAccessObject interface {
		BeginTx(shardId ...int64) (Transaction, error)
		Migrate(path string, shardId ...int64) error
		Init(cfg *config.DbParams) error
		InitCluster() error
		UpdateCluster() error
		ValidateCluster(shard int64) bool
		AddNewShard(shard *datamodels.Shard) error
		GetShardList() ([]int64, error)
		GetShardName(shardId int64) (string, error)
		UpdateShardPropertyData(property *datamodels.Property) error
	}

	BasePgDAO struct {
		BaseDB             *pg.DB
		LockMux            sync.Mutex
		Cluster            *sharding.Cluster
		shardsCount        int
		baseMigrationsPath string
	}
)

func NewBasePgDAO(baseMigrationsPath string) BasePgDAO {
	return BasePgDAO{
		baseMigrationsPath: baseMigrationsPath,
	}
}

func (d *BasePgDAO) Init(cfg *config.DbParams) error {
	err := d.initBaseDB(cfg)
	if err != nil {
		return err
	}

	err = d.InitCluster()
	if err != nil {
		return err
	}

	return nil
}

func (d *BasePgDAO) createShardIfNotExists(shardId int64) error {
	logrus.Tracef("createShardIfNotExists %d", shardId)
	_, err := d.Cluster.Shard(shardId).Exec(`CREATE SCHEMA IF NOT EXISTS ?shard`)
	if err != nil {
		return err
	}

	return nil
}

func (d *BasePgDAO) Migrate(path string, shardId ...int64) error {
	c := migrations.NewCollection()
	c.DisableSQLAutodiscover(true)
	err := c.DiscoverSQLMigrations(path)

	var DB *pg.DB
	if len(shardId) > 0 {
		logrus.Tracef("Using migration with a shard id %d", shardId[0])

		DB = d.Cluster.Shard(shardId[0])

		if DB == nil {
			return errors.New("no such shard was found")
		}

		c.SetTableName("?shard.gopg_migrations")

	} else {
		logrus.Tracef("Using base (no-shard) migration")
		DB = d.BaseDB
	}

	var oldVersion, newVersion int64
	oldVersion, newVersion, err = c.Run(DB, flag.Args()...)
	if err != nil {
		//init migrations mechanism for the 1st run
		if newVersion == 0 {
			_, _, err = c.Run(DB, "init")
			if err != nil {
				return err
			}

			oldVersion, newVersion, err = c.Run(DB)
		}

		if err != nil { //this check is because it can be changed by new c.Run(db) two lines upper
			return err
		}
	}

	if newVersion != oldVersion {
		logrus.Infof("DB migrated from version %d to %d for shard %d\n", oldVersion, newVersion, shardId)
	} else {
		logrus.Debugf("DB version is %d for shard %d\n", oldVersion, shardId)
	}

	return nil
}

func (d *BasePgDAO) initBaseDB(cfg *config.DbParams) error {
	conns, _ := strconv.Atoi(cfg.NumConns)
	if conns < 3 {
		conns = 3
	}

	d.BaseDB = pg.Connect(&pg.Options{
		Addr:               fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		User:               cfg.User,
		Password:           cfg.Password,
		Database:           cfg.Database,
		PoolSize:           conns,
		PoolTimeout:        PoolTimeout * time.Second,
		IdleTimeout:        IdleTimeout * time.Second,
		IdleCheckFrequency: IdleCheckFrequency * time.Second,
		ReadTimeout:        ReadTimeout * time.Second,
		WriteTimeout:       WriteTimeout * time.Second,
		MaxRetries:         MaxRetries,
	})
	d.BaseDB.AddQueryHook(LoggerHook{})
	err := d.Migrate(d.baseMigrationsPath + "/" + MigrationsBaseDir)
	if err != nil {
		return fmt.Errorf("cannot migrate: %s", err.Error())
	}

	return nil
}

func (d *BasePgDAO) InitCluster() error {
	//d.LockMux.Lock() //use mutex for cluster init to lock all shards calls
	//defer d.LockMux.Unlock()

	//load all known shards
	var (
		shards   []int64
		original int
		err      error
	)

	original = d.shardsCount
	shards, err = d.setClusterSize()

	if err != nil && err != pg.ErrNoRows {
		return err
	}
	d.shardsCount = len(shards)
	if d.shardsCount == 0 {
		return errors.New(ErrNoShardsYet)
	}

	// why not just create the new ones?
	//  since we are ordering by descending, we should be able to do this. if it changes, this will be broken
	realCount := d.shardsCount
	if original > 0 {
		realCount = d.shardsCount - original + 1
	}
	if realCount <= 0 {
		realCount = 1
	}
	if realCount > len(shards) {
		realCount = len(shards)
	}

	// Create database schema for our logical shards
	logrus.Infof("InitCluster(): Creating %d shards", realCount)
	for i := 0; i < realCount; i++ {
		shardId := shards[i]
		if err := d.createShardIfNotExists(shardId); err != nil {
			return err
		}

		//run migrations for every shard
		err := d.Migrate(d.baseMigrationsPath+"/"+MigrationsShardDir, shardId)
		if err != nil {
			return fmt.Errorf("cannot migrate a shard %d: %s", shardId, err.Error())
		}
	}

	return nil
}

//  this allows us to update the clusters without doing the migration (which shouldn't be needed for multi-instance operation)
func (d *BasePgDAO) UpdateCluster() error {

	_, err := d.setClusterSize()

	return err
}

func (d *BasePgDAO) setClusterSize() ([]int64, error) {
	dbs := []*pg.DB{d.BaseDB}

	//load all known shards
	var (
		shards []int64
		err    error
	)

	shards, err = d.GetShardList()
	if err != nil && err != pg.ErrNoRows {
		return shards, err
	}
	d.shardsCount = len(shards)
	if d.shardsCount == 0 {
		return shards, errors.New(ErrNoShardsYet)
	}

	//workaround to use shard id as given (not in array from 0 to count)
	for _, shardId := range shards { //find max ID
		if d.shardsCount < int(shardId) {
			d.shardsCount = int(shardId)
		}
	}
	d.shardsCount++ //increment to have a way to use max id (shards counter starts from 0)

	d.Cluster = sharding.NewCluster(dbs, d.shardsCount)

	return shards, nil
}

func (d *BasePgDAO) ValidateCluster(shard int64) bool {
	if d.Cluster == nil {
		return false
	}

	// cleaning up due to multi-instance operation
	if len(d.Cluster.Shards(d.BaseDB)) <= int(shard) {
		logrus.Debugf("ValidateCluster(): updating cluster due to size mismatch; len(size): %d, shardId: %d", len(d.Cluster.Shards(d.BaseDB))-1, shard)
		d.UpdateCluster()
	}

	return true
}

func (d *BasePgDAO) AddNewShard(shard *datamodels.Shard) error {
	tx, err := d.BeginTx()
	if err != nil {
		return fmt.Errorf("begin tx error: %s", err.Error())
	}
	defer tx.Rollback() // Rollback tx on error

	err = tx.GetTx().Insert(shard)
	if err != nil {
		if CheckIfDuplicateError(err) {
			return nil //duplicate is not an error here, at all
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	err = d.InitCluster()
	if err != nil {
		return err
	}

	return nil
}

func (d *BasePgDAO) GetShardList() ([]int64, error) {
	//load all known shards
	var (
		shardsList []datamodels.Shard
		err        error
	)

	d.shardsCount, err = d.BaseDB.Model(&shardsList).Order("shard_id DESC").SelectAndCount()
	if err != nil && err != pg.ErrNoRows {
		return nil, err
	}

	if d.shardsCount == 0 {
		return nil, errors.New(ErrNoShardsYet)
	}

	list := make([]int64, 0)
	for _, shard := range shardsList {
		list = append(list, shard.ShardId)
	}
	return list, nil
}

func (d *BasePgDAO) GetShardName(shardId int64) (string, error) {
	//load all known shards
	var (
		shardsList []datamodels.Shard
		err        error
		name       string
	)

	d.shardsCount, err = d.BaseDB.Model(&shardsList).SelectAndCount()
	if err != nil && err != pg.ErrNoRows {
		return name, err
	}

	if d.shardsCount == 0 {
		return name, errors.New(ErrNoShardsYet)
	}

	for _, shard := range shardsList {
		if shard.ShardId == shardId {
			name = shard.PropertyName
			break
		}
	}

	return name, err
}

func (d *BasePgDAO) UpdateShardPropertyData(property *datamodels.Property) error {
	var (
		shardModel datamodels.Shard
		err        error
	)

	_, err = d.BaseDB.Model(&shardModel).
		Set("property_name=?", property.DisplayName).
		Where("shard_id=?", property.PropertyId).
		Update()

	return err
}
