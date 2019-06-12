package swaggchain

import (
	"swaggp2p/pb"
	"sync"
	"time"
)

type Swagg interface {
	Version() string
	Height() uint32
	Genesis() GenesisBlock
	LastBlockHash() string
	LastBlockFound() uint32
	Validate() bool
	CurrentTime() *time.Time
	Difficulty() uint64
	Coinbase() Coinbase
	Logo() []byte
	Create() *SwaggChain

}

type GenesisBlock struct {
	sync.Mutex
	*pb.Block
	codes GenesisSpecialCodes // special codes not queryable othen than from the core who did create it

}

type GenesisSpecialCodes struct {
	canMine bool
	initialReward float64
	initialState []byte
	magicNumber byte
	routingAddresses []byte
	ancestorsDNA []byte
}

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

type Coinbase struct {
	Address []byte
	Coin interface{}
}

type AddressBook struct {

}


type Ancestor struct {
	mother []byte
	father []byte
	genePool []Organism
	generation int

}

type DNASolution struct {
	a []byte
	o []byte
	D []byte
	t []byte
}

type Organism struct {
	DNA []byte
	Fitness float64
	ParentA []byte
	ParentB []byte
}

//Population fitness of the core (blocks) based on the natural selection matheuristics by
type ChainDNA interface {
	createOrganism(target []byte) (organism Organism)
	createPopulation(target []byte) (population []Organism)
	createGenePool(population []Organism, target []byte, maxFitness float64) (pool []Organism)
	calculateFitness(target []byte)
	naturalSelection(pool []Organism, population []Organism, target []byte) []Organism
	crossover(d1 Organism, d2 Organism) Organism
	mutate()
	generateTarget() string
	validate() bool
	lock()


}

type DNAParams struct {
	mutationRate float64
	populationSize uint32
	target []byte
}

type TX interface {
	Spend(coinType []byte, amount float64, address string, to string) (*pb.Transaction, error)
	CreateOutput(input *pb.Transaction_Input, index int, spendable bool, solvable bool, useMaxSigInput bool) (*pb.Transaction_Output, error)
	CreateInput(tx *pb.Transaction, isCoinbase bool) (*pb.Transaction_Input, error)
	Validate(tx *pb.Transaction) error
	GetCoinType(coinType []byte) (*Coin, error)
}

type Coin struct {
	Name string
	Symbol string
	MinAllowed float64
	MaxAllowed float64
	MaxAvailable uint32
	Mineable bool
	Spendable bool
	Exchangeable bool
	TotalAvailable float64
	Unlocked float64
	CoinbaseAddr string
	IsInMainChain bool
	IsDefault bool
}

