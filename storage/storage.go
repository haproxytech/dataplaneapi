package storage

import (
	"github.com/haproxytech/dataplaneapi/storagetype"
)

// type constraint for storable objects, add more when need arises with |
type Storable interface {
	storagetype.ClusterModeData | storagetype.ConsulData | storagetype.AWSRegionData
}

type Storage[T Storable] interface {
	Get() (T, error)
	Store(data T) error
}
