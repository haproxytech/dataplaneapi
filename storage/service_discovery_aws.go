package storage

import (
	"slices"
	"sync"

	"github.com/haproxytech/client-native/v6/models"
	"github.com/haproxytech/dataplaneapi/storagetype"
)

const (
	AWSRegionFileName = "service_discovery/aws.json"
)

type SDAWStorage interface {
	AWS

	Load() error
	Store() error
}

type AWS interface {
	// GetAWS does not load the AWSRegions from the storage. Use Load() to load the them if needed.
	GetAWSRegions() storagetype.AWSRegions
	// AddAWSRegionsAndStore adds a list of AWSRegions to the storage and stores in the storage file.
	AddAWSRegionsAndStore(AWS storagetype.AWSRegions) error
	// AddAWSRegionAndStore adds a new AWSRegion to the storage and stores in the storage file.
	AddAWSRegionAndStore(AWS *models.AwsRegion) error
	// RemoveAWSAndStore removes a AWSRegion from the storage and stores in the storage file.
	RemoveAWSRegionAndStore(AWS *models.AwsRegion) error
	// ReplaceAllAWSRegionsAndStore replaces the list of AWSRegions in the storage and stores in the storage file.
	ReplaceAllAWSRegionsAndStore(AWS storagetype.AWSRegions) error
}

type AWSRegionStorageImpl struct {
	AWSData storagetype.AWSRegionData
	storage Storage[storagetype.AWSRegionData]
	mu      sync.RWMutex
}

// NewSDAWSRegionStorage creates a new AWSRegionStorageImpl with initial configuration from a file path.
func NewSDAWSRegionStorage(path string) (SDAWStorage, error) {
	fs := &fileStorage[storagetype.AWSRegionData]{path}
	if err := fs.initFile(); err != nil {
		return nil, err
	}
	AWSData, err := fs.Get()
	if err != nil {
		return nil, err
	}
	return &AWSRegionStorageImpl{
		storage: fs,
		AWSData: AWSData,
	}, nil
}

func (cs *AWSRegionStorageImpl) GetAWSRegions() storagetype.AWSRegions {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.AWSData.AWSRegions
}

func (cs *AWSRegionStorageImpl) Load() error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	AWS, err := cs.load()
	if err != nil {
		return err
	}
	cs.AWSData = AWS

	return nil
}

func (cs *AWSRegionStorageImpl) Store() error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if err := cs.store(); err != nil {
		return err
	}
	return nil
}

func (cs *AWSRegionStorageImpl) AddAWSRegionAndStore(region *models.AwsRegion) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.addAWSRegion(region)

	// In case something went wrong whil storing, delete the AWS from the store
	if err := cs.store(); err != nil {
		cs.removeAWSRegion(region)
		return err
	}
	return nil
}

func (cs *AWSRegionStorageImpl) AddAWSRegionsAndStore(regions storagetype.AWSRegions) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.addAWSRegions(regions)

	// In case something went wrong whil storing, delete the AWS from the store
	if err := cs.store(); err != nil {
		cs.removeAWSRegions(regions)
		return err
	}
	return nil
}

func (cs *AWSRegionStorageImpl) ReplaceAllAWSRegionsAndStore(regions storagetype.AWSRegions) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	oldAWS := cs.AWSData.AWSRegions

	cs.AWSData.AWSRegions = regions
	if err := cs.store(); err != nil {
		cs.AWSData.AWSRegions = oldAWS
		return err
	}
	return nil
}

func (cs *AWSRegionStorageImpl) RemoveAWSRegionAndStore(region *models.AwsRegion) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.removeAWSRegion(region)

	// In case something went wrong while storing, add the AWS to the store
	if err := cs.store(); err != nil {
		cs.addAWSRegion(region)
		return err
	}
	return nil
}

func (cs *AWSRegionStorageImpl) addAWSRegion(region *models.AwsRegion) {
	cs.AWSData.AWSRegions = append(cs.AWSData.AWSRegions, region)
}

func (cs *AWSRegionStorageImpl) addAWSRegions(regions storagetype.AWSRegions) {
	cs.AWSData.AWSRegions = append(cs.AWSData.AWSRegions, regions...)
}

func (cs *AWSRegionStorageImpl) removeAWSRegion(region *models.AwsRegion) {
	for i, u := range cs.AWSData.AWSRegions {
		if u.Name == region.Name {
			cs.AWSData.AWSRegions = slices.Delete(cs.AWSData.AWSRegions, i, i+1)
			break
		}
	}
}

func (cs *AWSRegionStorageImpl) removeAWSRegions(regions storagetype.AWSRegions) {
	for _, AWS := range cs.AWSData.AWSRegions {
		cs.removeAWSRegion(AWS)
	}
}

func (cs *AWSRegionStorageImpl) setAWSRegions(regions storagetype.AWSRegions) {
	cs.AWSData.AWSRegions = regions
}

func (cs *AWSRegionStorageImpl) store() error {
	return cs.storage.Store(cs.AWSData)
}

func (cs *AWSRegionStorageImpl) load() (storagetype.AWSRegionData, error) {
	return cs.storage.Get()
}
