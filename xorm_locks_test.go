package xorm_locks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalXormLockFactory_InitiallyNil(t *testing.T) {
	assert.Nil(t, GlobalXormLockFactory)
}
