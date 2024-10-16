package coreConfig

import (
	"sync"
	"yfapi/typedef"
)

type HotConf struct {
	Conf           typedef.Config
	L              sync.RWMutex
	LastModifyTime int64
}
