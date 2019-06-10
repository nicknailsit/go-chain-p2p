package swaggchain

import (
	"github.com/google/uuid"
	chainTypes "swaggp2p/core"
	logging "github.com/ipfs/go-log"
	"swaggp2p/core/pb"
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


/**
type SwaggChain struct {
	ID string
	IsMain bool
	nodes []*SwaggNode
	currentTime *time.Time
	epoch uint64
	versionB byte
	magicNumber byte
	genesisTime int64
	height uint32
	coinbase *Coinbase
	addressBook *AddressBook
	Blocks []*pb.Block
	lastHash []byte
	lastBlockIndex []byte
	lastReward uint64
	difficulty uint64
	logoAddress []byte
	dna []byte

}
 */

func (bc *Blockchain) Create() *Blockchain {

	//get new chain id
	chainID := uuid.New().String()

	//create new genesis block
	genesis := new(GenesisBlock)
	g := genesis.Create(chainID)

	//init blockchain type
	BC := &Blockchain{}
	chain := BC.swaggchain



	chainBlocks := make([]*pb.Block, 1)

	//append genesis block to the chain
	chainBlocks = append(chainBlocks, g)

	BC.swaggchain = chain

	return BC



}


type Coinbase chainTypes.Coinbase

func (co *Coinbase) Create() {

}