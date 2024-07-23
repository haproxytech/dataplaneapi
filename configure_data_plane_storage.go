package dataplaneapi

import (
	"os"
	"path/filepath"

	client_native "github.com/haproxytech/client-native/v6"
	native_misc "github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/dataplaneapi/log"
)

func initDataplaneStorage(path string, client client_native.HAProxyClient) {
	if path == "" {
		return
	}
	if _, err := native_misc.CheckOrCreateWritableDirectory(path); err != nil {
		log.Fatalf("error initializing dataplane internal storage: %v", err)
	}
	if _, err := os.Stat(filepath.Join(path, "logs.json")); err != nil && os.IsNotExist(err) {
		generalStorage, storageErr := client.GeneralStorage()
		if storageErr != nil {
			log.Warningf("failed to get general storage: %v", storageErr)
			return
		}
		src, _, getErr := generalStorage.Get("logs.json")
		if getErr == nil {
			dst := filepath.Join(path, "logs.json")
			if err := os.Rename(src, dst); err != nil {
				log.Warningf("Failed to move logs.json to dataplane internal storage: %v", err)
			}
		}
	}
}
