package xorm_locks

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"xorm.io/xorm"
	"xorm.io/xorm/core"
)

func testXormEngine(t *testing.T) *xorm.Engine {
	t.Helper()
	engine, err := xorm.NewEngine("sqlite3", "file:xorm_test.db?mode=memory")
	assert.Nil(t, err)
	t.Cleanup(func() { engine.Close() })
	return engine
}

func TestNewXormConnectionManager(t *testing.T) {
	engine := testXormEngine(t)
	m := NewXormConnectionManager(engine)
	assert.NotNil(t, m)
	assert.Equal(t, "xorm-connection-manager", m.Name())

	ctx := context.Background()
	got, err := m.Take(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, got)

	assert.Nil(t, m.Return(ctx, got))
	assert.Nil(t, m.Shutdown(ctx))
}

func TestNewXormLockFactory(t *testing.T) {
	engine := testXormEngine(t)
	factory, err := NewXormLockFactory(engine)
	assert.Nil(t, err)
	assert.NotNil(t, factory)

	// 端到端加解锁
	ctx := context.Background()
	lock, err := factory.CreateLock("test-xorm-lock")
	assert.Nil(t, err)
	assert.NotNil(t, lock)
	assert.Nil(t, lock.Lock(ctx, "owner-1"))
	assert.Nil(t, lock.UnLock(ctx, "owner-1"))
}

func TestCreateStorageForXorm_Errors(t *testing.T) {
	engine := testXormEngine(t)
	m := NewXormConnectionManager(engine)
	s, err := CreateStorageForXorm(engine, m)
	assert.Nil(t, err)
	assert.NotNil(t, s)

	// mysql 驱动名走第一条路（mysql storage 在 sqlmock 上建表失败），
	// 回退 NewStorageBySqlDb 识别 sqlmock 驱动失败——覆盖两个回退分支
	mockDB, _, err := sqlmock.New()
	assert.Nil(t, err)
	defer mockDB.Close()
	bogus, err := xorm.NewEngineWithDB("mysql", "root:rootpass@tcp(127.0.0.1:13307)/storage_lock_test", core.FromDB(mockDB))
	assert.Nil(t, err)
	defer bogus.Close()
	bogusM := NewXormConnectionManager(bogus)
	_, err = CreateStorageForXorm(bogus, bogusM)
	assert.NotNil(t, err)

	_, err = NewXormLockFactory(bogus)
	assert.NotNil(t, err)

	assert.NotNil(t, InitGlobalXormLockFactory(bogus))
	assert.Nil(t, GlobalXormLockFactory)
}

func TestInitGlobalXormLockFactory(t *testing.T) {
	engine := testXormEngine(t)
	assert.Nil(t, InitGlobalXormLockFactory(engine))
	assert.NotNil(t, GlobalXormLockFactory)
	GlobalXormLockFactory = nil
}
