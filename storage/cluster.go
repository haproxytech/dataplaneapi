package storage

import (
	"slices"
	"sync"

	"github.com/haproxytech/dataplaneapi/storagetype"
)

const (
	ClusterModeDataFileName = "cluster.json"
)

// ClusterModeStorage persists dataplane-managed users in the dataplane storage
// file (historically cluster.json).
type ClusterModeStorage interface {
	Users

	Load() error
	Store() error
}

type Users interface {
	// GetUsers does not load the users from the storage. Use Load() to load the users if needed.
	GetUsers() storagetype.Users
	// AddUsersAndStore adds a list of users to the storage and stores in the storage file.
	AddUsersAndStore(users storagetype.Users) error
	// AddUserAndStore adds a new user to the storage and stores in the storage file.
	AddUserAndStore(user storagetype.User) error
	// RemoveUserAndStore removes a user from the storage and stores in the storage file.
	RemoveUserAndStore(user storagetype.User) error
	// ReplaceAllUsersAndStore replaces the list of users in the storage and stores in the storage file.
	ReplaceAllUsersAndStore(users storagetype.Users) error
}

type clusterModeStorageImpl struct {
	ClusterModeData storagetype.ClusterModeData
	storage         Storage[storagetype.ClusterModeData]
	mu              sync.RWMutex
}

// NewClusterModeStorage creates a new clusterModeStorageImpl with initial configuration from a file path.
func NewClusterModeStorage(path string) (ClusterModeStorage, error) {
	fs := &fileStorage[storagetype.ClusterModeData]{path}
	if err := fs.initFile(); err != nil {
		return nil, err
	}
	clusterStorage, err := fs.Get()
	if err != nil {
		return nil, err
	}
	return &clusterModeStorageImpl{
		storage:         fs,
		ClusterModeData: clusterStorage,
	}, nil
}

func (cs *clusterModeStorageImpl) Load() error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	data, err := cs.load()
	if err != nil {
		return err
	}
	cs.ClusterModeData = data

	return nil
}

func (cs *clusterModeStorageImpl) Store() error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if err := cs.store(); err != nil {
		return err
	}
	return nil
}

func (cs *clusterModeStorageImpl) GetUsers() storagetype.Users {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.ClusterModeData.Users
}

func (cs *clusterModeStorageImpl) AddUserAndStore(user storagetype.User) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.addUser(user)

	// In case something went wrong whil storing, delete the user from the store
	if err := cs.store(); err != nil {
		cs.removeUser(user)
		return err
	}
	return nil
}

func (cs *clusterModeStorageImpl) AddUsersAndStore(users storagetype.Users) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.addUsers(users)

	// In case something went wrong whil storing, delete the user from the store
	if err := cs.store(); err != nil {
		cs.removeUsers(users)
		return err
	}
	return nil
}

func (cs *clusterModeStorageImpl) ReplaceAllUsersAndStore(users storagetype.Users) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	oldUsers := cs.ClusterModeData.Users

	cs.ClusterModeData.Users = users
	if err := cs.store(); err != nil {
		cs.ClusterModeData.Users = oldUsers
		return err
	}
	return nil
}

func (cs *clusterModeStorageImpl) RemoveUserAndStore(user storagetype.User) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.removeUser(user)

	// In case something went wrong while storing, add the user to the store
	if err := cs.store(); err != nil {
		cs.addUser(user)
		return err
	}
	return nil
}

func (cs *clusterModeStorageImpl) addUser(user storagetype.User) {
	cs.ClusterModeData.Users = append(cs.ClusterModeData.Users, user)
}

func (cs *clusterModeStorageImpl) addUsers(users storagetype.Users) {
	cs.ClusterModeData.Users = append(cs.ClusterModeData.Users, users...)
}

func (cs *clusterModeStorageImpl) removeUser(user storagetype.User) {
	for i, u := range cs.ClusterModeData.Users {
		if u.Name == user.Name {
			cs.ClusterModeData.Users = slices.Delete(cs.ClusterModeData.Users, i, i+1)
			break
		}
	}
}

func (cs *clusterModeStorageImpl) removeUsers(users storagetype.Users) {
	for _, user := range users {
		cs.removeUser(user)
	}
}

func (cs *clusterModeStorageImpl) store() error {
	return cs.storage.Store(cs.ClusterModeData)
}

func (cs *clusterModeStorageImpl) load() (storagetype.ClusterModeData, error) {
	return cs.storage.Get()
}
