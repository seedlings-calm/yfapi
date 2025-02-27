package coreSnowflake

import (
	"sync"
	snowflake2 "yfapi/core/coreTools/snowflake"
)

var (
	once      sync.Once
	snowflake *snowflake2.Node
	err       error
)

func New(node int64) (*snowflake2.Node, error) {
	once.Do(func() {
		snowflake, err = snowflake2.NewNode(node)
	})
	return snowflake, err
}

func GetSnowId() string {
	return snowflake.Generate().String()
}
