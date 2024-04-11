package cn

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	clientnative "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/client-native/v6/configuration"
	configuration_options "github.com/haproxytech/client-native/v6/configuration/options"
	runtime_api "github.com/haproxytech/client-native/v6/runtime"
	runtime_options "github.com/haproxytech/client-native/v6/runtime/options"
	parser "github.com/haproxytech/config-parser/v5"
	"github.com/haproxytech/config-parser/v5/types"

	dataplaneapi_config "github.com/haproxytech/dataplaneapi/configuration"
	"github.com/haproxytech/dataplaneapi/log"
	"github.com/haproxytech/dataplaneapi/misc"
)

var (
	socketsList   = map[int]string{} //nolint:unused
	muSocketsList sync.Mutex
)

func ConfigureConfigurationClient(haproxyOptions dataplaneapi_config.HAProxyConfiguration, mWorker bool) (configuration.Configuration, error) {
	confClient, err := configuration.New(context.Background(),
		configuration_options.ConfigurationFile(haproxyOptions.ConfigFile),
		configuration_options.HAProxyBin(haproxyOptions.HAProxy),
		configuration_options.Backups(haproxyOptions.BackupsNumber),
		configuration_options.BackupsDir(haproxyOptions.BackupsDir),
		configuration_options.UsePersistentTransactions,
		configuration_options.TransactionsDir(haproxyOptions.TransactionDir),
		configuration_options.ValidateCmd(haproxyOptions.ValidateCmd),
		configuration_options.MasterWorker,
		configuration_options.UseMd5Hash,
	)
	if err != nil {
		return nil, fmt.Errorf("error setting up configuration client: %s", err.Error())
	}

	p := confClient.Parser()
	comments, err := p.Get(parser.Comments, parser.CommentsSectionName, "#")
	insertDisclaimer := false
	if err != nil {
		insertDisclaimer = true
	}
	data, ok := comments.([]types.Comments)
	if !ok {
		insertDisclaimer = true
	} else if len(data) == 0 || data[0].Value != "Dataplaneapi managed File" {
		insertDisclaimer = true
	}
	if insertDisclaimer {
		commentsNew := types.Comments{Value: "Dataplaneapi managed File"}
		err = p.Insert(parser.Comments, parser.CommentsSectionName, "#", commentsNew, 0)
		if err != nil {
			return nil, fmt.Errorf("error setting up configuration client: %s", err.Error())
		}
		commentsNew = types.Comments{Value: "changing file directly can cause a conflict if dataplaneapi is running"}
		err = p.Insert(parser.Comments, parser.CommentsSectionName, "#", commentsNew, 1)
		if err != nil {
			return nil, fmt.Errorf("error setting up configuration client: %s", err.Error())
		}
	}

	return confClient, nil
}

func ConfigureRuntimeClient(ctx context.Context, confClient configuration.Configuration, haproxyOptions dataplaneapi_config.HAProxyConfiguration) runtime_api.Runtime {
	mapsDir := runtime_options.MapsDir(haproxyOptions.MapsDir)
	var runtimeClient runtime_api.Runtime

	_, globalConf, err := confClient.GetGlobalConfiguration("")
	waitForRuntimeOption := runtime_options.AllowDelayedStart(haproxyOptions.DelayedStartMax, haproxyOptions.DelayedStartTick)

	// First try to setup master runtime socket
	if err == nil {
		var err error
		// If master socket is set and a valid unix socket, use only this
		if haproxyOptions.MasterRuntime != "" && misc.IsUnixSocketAddr(haproxyOptions.MasterRuntime) {
			masterSocket := haproxyOptions.MasterRuntime
			// if nbproc is set, set nbproc sockets
			if globalConf.Nbproc > 0 {
				nbproc := int(globalConf.Nbproc)
				ms := runtime_options.MasterSocket(masterSocket, nbproc)
				runtimeClient, err = runtime_api.New(ctx, mapsDir, ms, waitForRuntimeOption)
				if err == nil {
					return runtimeClient
				}
				log.Warningf("Error setting up runtime client with master socket: %s : %s", masterSocket, err.Error())
			} else {
				// if nbproc is not set, use master socket with 1 process
				ms := runtime_options.MasterSocket(masterSocket, 1)
				runtimeClient, err = runtime_api.New(ctx, mapsDir, ms, waitForRuntimeOption)
				if err == nil {
					return runtimeClient
				}
				log.Warningf("Error setting up runtime client with master socket: %s : %s", masterSocket, err.Error())
			}
		}
		runtimeAPIs := globalConf.RuntimeAPIs
		// if no master socket set, read from first valid socket if nbproc <= 1
		if globalConf.Nbproc <= 1 {
			sockets := make(map[int]string)
			for _, r := range runtimeAPIs {
				if misc.IsUnixSocketAddr(*r.Address) {
					sockets[1] = *r.Address
					socketsL := runtime_options.Sockets(sockets)
					runtimeClient, err = runtime_api.New(ctx, mapsDir, socketsL, waitForRuntimeOption)
					if err == nil {
						muSocketsList.Lock()
						socketsList = sockets
						muSocketsList.Unlock()
						return runtimeClient
					}
					log.Warningf("Error setting up runtime client with socket: %s : %s", *r.Address, err.Error())
				}
			}
		} else {
			// else try to find process specific sockets and set them up
			sockets := make(map[int]string)
			for _, r := range runtimeAPIs {
				//nolint:govet
				if misc.IsUnixSocketAddr(*r.Address) && r.Process != "" {
					process, err := strconv.ParseInt(r.Process, 10, 64)
					if err == nil {
						sockets[int(process)] = *r.Address
					}
				}
			}
			// no process specific settings found, Issue a warning and return empty runtime client
			if len(sockets) == 0 {
				log.Warning("Runtime API not configured, found multiple processes and no stats sockets bound to them.")
				return runtimeClient
				// use only found process specific sockets issue a warning if not all processes have a socket configured
			}
			if len(sockets) < int(globalConf.Nbproc) {
				log.Warning("Runtime API not configured properly, there are more processes then configured sockets")
			}

			socketLst := runtime_options.Sockets(sockets)
			runtimeClient, err = runtime_api.New(ctx, mapsDir, socketLst, waitForRuntimeOption)
			if err == nil {
				return runtimeClient
			}
			log.Warningf("Error setting up runtime client with sockets: %v : %s", sockets, err.Error())
		}
		if err != nil {
			log.Warning("Runtime API not configured, not using it: " + err.Error())
		} else {
			log.Warning("Runtime API not configured, not using it")
		}
		return runtimeClient
	}
	log.Warning("Cannot read runtime API configuration, not using it")
	return runtimeClient
}

// ReconfigureRuntime check if runtime client need be reconfigured by comparing the current configuration with the old
// one (i.e. runtimeAPIsOld) and returns a callback that ReloadAgent can use to reconfigure the runtime client.
func ReconfigureRuntime(client clientnative.HAProxyClient, runtimeAPIsOld []*models.RuntimeAPI) (callbackNeeded bool, callback func(), err error) {
	cfg, err := client.Configuration()
	if err != nil {
		return false, nil, err
	}
	_, globalConf, err := cfg.GetGlobalConfiguration("")
	if err != nil {
		return false, nil, err
	}
	runtimeAPIsNew := globalConf.RuntimeAPIs
	reconfigureRuntime := false
	if len(runtimeAPIsOld) != len(runtimeAPIsNew) {
		reconfigureRuntime = true
	} else {
		for _, runtimeOld := range runtimeAPIsOld {
			if runtimeOld.Address == nil {
				continue
			}
			found := false
			for _, runtimeNew := range runtimeAPIsNew {
				if runtimeNew.Address == nil {
					continue
				}
				if *runtimeNew.Address == *runtimeOld.Address {
					found = true
					break
				}
			}
			if !found {
				reconfigureRuntime = true
				break
			}
		}
	}

	if reconfigureRuntime {
		dpapiCfg := dataplaneapi_config.Get()
		haproxyOptions := dpapiCfg.HAProxy
		return true, func() {
			var rnt runtime_api.Runtime
			i := 1
			for i < 10 {
				rnt = ConfigureRuntimeClient(context.Background(), cfg, haproxyOptions)
				if rnt != nil {
					break
				}
				time.Sleep(time.Duration(i) * time.Second)
				i += i // exponential backoof
			}
			client.ReplaceRuntime(rnt)
			if rnt == nil {
				log.Debugf("reload callback completed, no runtime API")
			} else {
				log.Debugf("reload callback completed, runtime API reconfigured")
			}
		}, nil
	}

	return false, nil, nil
}
