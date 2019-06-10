package wire

import (
	chainTypes "swaggp2p/chain"
	logging "github.com/ipfs/go-log"
	"sync"
	"time"
)

const GC_TTL = time.Minute;
var log = logging.Logger("swaggchain")

type Blockchain struct {
	sync.Mutex
	swaggchain *chainTypes.SwaggChain
	chainTypes.Swagg
}


func (bc *Blockchain) Create() {

	

}