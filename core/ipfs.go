package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	config "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	coreiface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
)

// returns a go-ipfs node backend CoreAPI instance
func CreateIpfsCoreApi(ctx context.Context, cidComputeOnly bool) (coreiface.CoreAPI, error) {
	cfg := core.BuildCfg{
		Online:    false, //networking off
		Permanent: false, // want node to be lightweight
	}

	var err error
	if cidComputeOnly {
		cfg.NilRepo = true
	} else {
		var repoPath string
		if repoPath, err = initIpfsRepo(); err != nil {
			return nil, err
		}
		if err := setupPlugins(repoPath); err != nil {
			return nil, err
		}
		var ipfsConfig *config.Config
		if ipfsConfig, err = config.Init(os.Stdout, 2048); err != nil {
			return nil, err
		}

		ipfsConfig.Datastore = config.DefaultDatastoreConfig()
		if err = fsrepo.Init(repoPath, ipfsConfig); err != nil {
			return nil, err
		}
		if cfg.Repo, err = fsrepo.Open(repoPath); err != nil {
			return nil, err
		}
	}

	var nnode *core.IpfsNode
	if nnode, err = core.NewNode(ctx, &cfg); err != nil {
		return nil, err
	}
	return coreapi.NewCoreAPI(nnode, options.Api.Offline(true), options.Api.FetchBlocks(false))
}

func initIpfsRepo() (string, error) {
	pathRoot, err := config.PathRoot() // IFPS path root, can be changed via env variable too
	if err != nil {
		return "", fmt.Errorf("error getting path root: %s", err)
	}
	if err = os.Mkdir(pathRoot, 0777); err != nil {
		return "", fmt.Errorf("can't create ipfs repo directory: %s", err)
	}

	return pathRoot, nil
}

func setupPlugins(externalPluginsPath string) error {
	// Load any external plugins if available on externalPluginsPath
	plugins, err := loader.NewPluginLoader(filepath.Join(externalPluginsPath, "plugins"))
	if err != nil {
		return fmt.Errorf("error loading plugins: %s", err)
	}

	// Load preloaded and external plugins
	if err := plugins.Initialize(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	return nil
}
