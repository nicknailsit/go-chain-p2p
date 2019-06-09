package repo

import (
	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p-core/crypto"
	lvldb "github.com/ipfs/go-ds-leveldb"
	"github.com/libp2p/go-libp2p-core/peer"
	homedir "github.com/mitchellh/go-homedir"
	"os"
	pth "path/filepath"
	"runtime"
)

type Repo struct {
	path string
	store datastore.Batching
	privKey crypto.PrivKey
	bootstrapPeers []peer.AddrInfo
}

func NewRepo(path string) (*Repo, error) {

	var err error

	if path == "" {
		path, err = defaultRepoPath()
		if err != nil {
			return nil, err
		}
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(pth.Join(path, "datastore"), os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	privKey, err := loadPrivKey(path)
	if err != nil {
		return nil, err
	}

	dstore, err := lvldb.NewDatastore(pth.Join(path, "datastore"), nil)
	if err != nil {
		return nil, err
	}

	bootstrapPeers, err := ParseBootstrapPeers(defaultBootstrapPeers)
	if err != nil {
		return nil, err
	}

	return &Repo{
		path: path,
		privKey: privKey,
		store: dstore,
		bootstrapPeers:bootstrapPeers,
	},nil
}

func (r *Repo) Path() string {
	return r.path
}

func (r *Repo) PrivKey() crypto.PrivKey {
	return r.privKey
}

func (r *Repo) Datastore() datastore.Batching {
	return r.store
}

func (r *Repo) BootstrapPeers() []peer.AddrInfo {
	return r.bootstrapPeers
}

func defaultRepoPath() (string, error) {
	path := "~"
	directoryName := "swaggchain"
	switch runtime.GOOS {
	case "linux":
		directoryName = ".swaggchain"
	case "darwin":
		path = "~/Library/Application Support"
	}

	fullPath, err := homedir.Expand(pth.Join(path, directoryName))
	if err != nil {
		return "", err
	}
	return pth.Clean(fullPath), nil
}
