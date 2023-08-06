package xorm_locks

import (
	"database/sql"
	sqldb_storage "github.com/storage-lock/go-sqldb-storage"
	"github.com/storage-lock/go-storage"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
	"xorm.io/xorm"
)

type XormLockFactory struct {
	engine *xorm.Engine
	*storage_lock_factory.StorageLockFactory[*sql.DB]
}

func NewXormLockFactory(engine *xorm.Engine) (*XormLockFactory, error) {
	connectionManager := NewXormConnectionManager(engine)

	storage, err := CreateStorageForXorm(engine, connectionManager)
	if err != nil {
		return nil, err
	}

	factory := storage_lock_factory.NewStorageLockFactory[*sql.DB](storage, connectionManager)

	return &XormLockFactory{
		engine:             engine,
		StorageLockFactory: factory,
	}, nil
}

// CreateStorageForXorm 尝试从xorm创建Storage
func CreateStorageForXorm(engine *xorm.Engine, connectionManager storage.ConnectionManager[*sql.DB]) (storage.Storage, error) {

	// 先尝试根据驱动名称创建
	storage, err := sqldb_storage.NewStorageByDriverName(engine.DriverName(), connectionManager)
	if storage != nil && err == nil {
		return storage, err
	}

	// 再然后根据识别出来的名称创建
	return sqldb_storage.NewStorageBySqlDb(engine.DB().DB, connectionManager)
}
