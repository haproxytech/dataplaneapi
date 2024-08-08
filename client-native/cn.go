package cn

import (
	"context"
	"fmt"
	"sync"
	"time"

	clientnative "github.com/haproxytech/client-native/v6"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/configuration"
	configuration_options "github.com/haproxytech/client-native/v6/configuration/options"
	runtime_api "github.com/haproxytech/client-native/v6/runtime"
	runtime_options "github.com/haproxytech/client-native/v6/runtime/options"

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
		configuration_options.PreferredTimeSuffix(haproxyOptions.PreferredTimeSuffix),
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
			// nbproc has been removed, use master socket with 1 process
			ms := runtime_options.MasterSocket(masterSocket)
			runtimeClient, err = runtime_api.New(ctx, mapsDir, ms, waitForRuntimeOption)
			if err == nil {
				return runtimeClient
			}
			log.Warningf("Error setting up runtime client with master socket: %s : %s", masterSocket, err.Error())
		}

		runtimeAPIs := globalConf.RuntimeAPIs
		// if no master socket set, read from first valid socket
		sockets := make(map[int]string)
		for _, r := range runtimeAPIs {
			if misc.IsUnixSocketAddr(*r.Address) {
				socketsL := runtime_options.Socket(*r.Address)
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
// one and returns a callback that ReloadAgent can use to reconfigure the runtime client.
func ReconfigureRuntime(client clientnative.HAProxyClient) (callbackNeeded bool, callback func(), err error) {
	cfg, err := client.Configuration()
	if err != nil {
		return false, nil, err
	}

	reconfigureRuntime := false

	// client Runtime is not reconfigured yet, so client.Runtime() return the "old" runtime
	oldRuntime, err := client.Runtime()
	// In case we have an error, we need to reconfigure
	if err != nil {
		reconfigureRuntime = true
	}

	if !reconfigureRuntime {
		// If we are not using stats socket (i.e. using master socket)
		// Return immediately and do not reconfigure
		// Do not try to compare the master socket path with stats socket paths
		if !oldRuntime.IsStatsSocket() {
			return false, nil, nil
		}

		_, globalConf, err := cfg.GetGlobalConfiguration("")
		if err != nil {
			return false, nil, err
		}
		runtimeAPIsNew := globalConf.RuntimeAPIs

		oldSocketPath := oldRuntime.SocketPath()

		// Now check if the new configuration contains the "old" runtime socket we are using
		// If yes, no need to reconfigure, socket path still exists.
		// If not found, only then reconfigure
		// This, only if stats socket, not for master socket
		found := false
		for _, runtimeNew := range runtimeAPIsNew {
			if runtimeNew.Address == nil {
				continue
			}
			if *runtimeNew.Address == oldSocketPath {
				found = true
				break
			}
		}
		if !found {
			reconfigureRuntime = true
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
