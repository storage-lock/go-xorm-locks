package xorm_locks

import "xorm.io/xorm"

var GlobalXormLockFactory *XormLockFactory

func InitGlobalXormLockFactory(engine *xorm.Engine) error {
	factory, err := NewXormLockFactory(engine)
	if err != nil {
		return err
	}
	GlobalXormLockFactory = factory
	return nil
}
