package core

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	config "github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/core/node/libp2p"
	"github.com/ipfs/kubo/plugin/loader"
	"github.com/ipfs/kubo/repo/fsrepo"
)

// returns a go-ipfs node backend CoreAPI instance
func CreateIpfsNode(ctx context.Context, cidComputeOnly bool) (*core.IpfsNode, error) {
	cfg := core.BuildCfg{
		Online:    !cidComputeOnly, // networking
		Permanent: !cidComputeOnly, // data persists across restarts?
	}

	var err error
	if cidComputeOnly {
		cfg.Repo = nil
	} else {
		cfg.Routing = libp2p.DHTOption
		cfg.Host = libp2p.DefaultHostOption

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

	return nnode, err
}

func initIpfsRepo() (string, error) {
	pathRoot, err := config.PathRoot() // IFPS path root, can be changed via env variable too
	if err != nil {
		return "", fmt.Errorf("error getting path root: %s", err)
	}
	if err = os.MkdirAll(pathRoot, fs.ModeDir); err != nil {
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
