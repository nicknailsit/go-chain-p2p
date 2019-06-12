package swaggchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/cbergoon/merkletree"
	"github.com/google/uuid"
	"swaggp2p/pb"
	"sync"
	"time"
)

const GC_TTL = time.Minute;


type Blockchain struct {
	sync.Mutex
	swaggchain *SwaggChain
	Swagg
	MerkleRoot []byte
	MerkleString string
}



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

	BC.swaggchain.lastHash = genesis.Header.Hash;

	BC.swaggchain = chain

	return BC


}

type MerkleContent struct {
	x string
}

func (m MerkleContent) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(m.x)); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func (m MerkleContent) Equals(other merkletree.Content) (bool, error) {
	return m.x == other.(MerkleContent).x, nil
}



func (bc Blockchain) ComputeMerkleTree() {


	var list []merkletree.Content
	for _, block := range bc.swaggchain.Blocks {

		ser, _ := json.Marshal(block)

		list = append(list, MerkleContent{x: hex.EncodeToString(ser)})

	}

	t, err := merkletree.NewTree(list)
	if err != nil {
		log.Fatal(err)
	}

	merkleroot := t.MerkleRoot()
	bc.MerkleRoot = merkleroot
	bc.MerkleString = t.String()
}

func (bc Blockchain) VerifyMerkleTree(t *merkletree.MerkleTree) {

	_, err := t.VerifyTree()
	if err != nil {
		log.Fatal(err)
	}

}

func (bc Blockchain) VerifyMerkleTreeContent(t *merkletree.MerkleTree, toVerify MerkleContent) {
	_, err := t.VerifyContent(toVerify)
	if err != nil {
		log.Fatal(err)
	}
}

func (bc Blockchain) getLastHash() []byte {
	return bc.swaggchain.lastHash;
}

