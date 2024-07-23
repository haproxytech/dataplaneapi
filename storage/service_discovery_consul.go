package storage

import (
	"slices"
	"sync"

	"github.com/haproxytech/client-native/v6/models"
	"github.com/haproxytech/dataplaneapi/storagetype"
)

const (
	ConsulFileName = "service_discovery/consul.json"
)

type SDConsulStorage interface {
	Consul

	Load() error
	Store() error
}

type Consul interface {
	// GetConsuls does not load the consuls from the storage. Use Load() to load the them if needed.
	GetConsuls() storagetype.Consuls
	// AddConsulsAndStore adds a list of consuls to the storage and stores in the storage file.
	AddConsulsAndStore(consuls storagetype.Consuls) error
	// AddConsulAndStore adds a new consul to the storage and stores in the storage file.
	AddConsulAndStore(consul *models.Consul) error
	// RemoveConsulAndStore removes a consul from the storage and stores in the storage file.
	RemoveConsulAndStore(consul *models.Consul) error
	// ReplaceAllConsulsAndStore replaces the list of consuls in the storage and stores in the storage file.
	ReplaceAllConsulsAndStore(consuls storagetype.Consuls) error
}

type consulStorageImpl struct {
	ConsulData storagetype.ConsulData
	storage    Storage[storagetype.ConsulData]
	mu         sync.RWMutex
}

// NewSDConsulStorage creates a new consulStorageImpl with initial configuration from a file path.
func NewSDConsulStorage(path string) (SDConsulStorage, error) {
	fs := &fileStorage[storagetype.ConsulData]{path}
	if err := fs.initFile(); err != nil {
		return nil, err
	}
	consulData, err := fs.Get()
	if err != nil {
		return nil, err
	}
	return &consulStorageImpl{
		storage:    fs,
		ConsulData: consulData,
	}, nil
}

func (cs *consulStorageImpl) GetConsuls() storagetype.Consuls {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.ConsulData.Consuls
}

func (cs *consulStorageImpl) Load() error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	consul, err := cs.load()
	if err != nil {
		return err
	}
	cs.ConsulData = consul

	return nil
}

func (cs *consulStorageImpl) Store() error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if err := cs.store(); err != nil {
		return err
	}
	return nil
}

func (cs *consulStorageImpl) AddConsulAndStore(consul *models.Consul) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.addConsul(consul)

	// In case something went wrong whil storing, delete the consul from the store
	if err := cs.store(); err != nil {
		cs.removeConsul(consul)
		return err
	}
	return nil
}

func (cs *consulStorageImpl) AddConsulsAndStore(consuls storagetype.Consuls) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.addConsuls(consuls)

	// In case something went wrong whil storing, delete the consul from the store
	if err := cs.store(); err != nil {
		cs.removeConsuls(consuls)
		return err
	}
	return nil
}

func (cs *consulStorageImpl) ReplaceAllConsulsAndStore(consuls storagetype.Consuls) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	oldConsuls := cs.ConsulData.Consuls

	cs.ConsulData.Consuls = consuls
	if err := cs.store(); err != nil {
		cs.ConsulData.Consuls = oldConsuls
		return err
	}
	return nil
}

func (cs *consulStorageImpl) RemoveConsulAndStore(consul *models.Consul) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.removeConsul(consul)

	// In case something went wrong while storing, add the consul to the store
	if err := cs.store(); err != nil {
		cs.addConsul(consul)
		return err
	}
	return nil
}

func (cs *consulStorageImpl) addConsul(consul *models.Consul) {
	cs.ConsulData.Consuls = append(cs.ConsulData.Consuls, consul)
}

func (cs *consulStorageImpl) addConsuls(consuls storagetype.Consuls) {
	cs.ConsulData.Consuls = append(cs.ConsulData.Consuls, consuls...)
}

func (cs *consulStorageImpl) removeConsul(consul *models.Consul) {
	for i, u := range cs.ConsulData.Consuls {
		if u.Name == consul.Name {
			cs.ConsulData.Consuls = slices.Delete(cs.ConsulData.Consuls, i, i+1)
			break
		}
	}
}

func (cs *consulStorageImpl) removeConsuls(consuls storagetype.Consuls) {
	for _, consul := range consuls {
		cs.removeConsul(consul)
	}
}

func (cs *consulStorageImpl) setConsuls(consuls storagetype.Consuls) {
	cs.ConsulData.Consuls = consuls
}

func (cs *consulStorageImpl) store() error {
	return cs.storage.Store(cs.ConsulData)
}

func (cs *consulStorageImpl) load() (storagetype.ConsulData, error) {
	return cs.storage.Get()
}
