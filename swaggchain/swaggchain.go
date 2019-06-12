package swaggchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/cbergoon/merkletree"
	"github.com/google/uuid"
	"github.com/minio/blake2b-simd"
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



func (bc Blockchain) ComputeMerkleTree() *merkletree.MerkleTree {


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
	return t
}

func (bc Blockchain) VerifyMerkleTree(t *merkletree.MerkleTree) error {

	_, err := t.VerifyTree()
	if err != nil {
		return err
	}

	return nil

}

func (bc Blockchain) VerifyMerkleTreeContent(t *merkletree.MerkleTree, toVerify MerkleContent) error {
	_, err := t.VerifyContent(toVerify)
	if err != nil {
		return err
	}
	return nil
}

func (bc Blockchain) getLastHash() []byte {
	return bc.swaggchain.lastHash;
}

func CompareChains(chain...  *Blockchain) *Blockchain {

		var bestCandidate *Blockchain

		bestCandidate = nil

		for i, _ := range chain {


			//compare lengths
			if len(chain[i].swaggchain.Blocks) == len(chain[i+1].swaggchain.Blocks) {
					bestCandidate = chain[i+1]
			} else {
				if len(chain[i].swaggchain.Blocks) > len(chain[i+1].swaggchain.Blocks) {
					bestCandidate = chain[i]
				} else if len(chain[i].swaggchain.Blocks) < len(chain[i+1].swaggchain.Blocks)  {
					bestCandidate = chain[i+1]
				}
			}

		}

		return bestCandidate

}

func ValidateMerkle(chain ...*Blockchain) *Blockchain {

	var bestCandidate *Blockchain

	bestCandidate = nil

	for i, _ := range chain {

		ms := chain[i].MerkleString
		ms2 := chain[i+1].MerkleString
		if ms == ms2 {
			bestCandidate = chain[i]
		} else if len(ms) > len(ms2) {
			err := chain[i].VerifyMerkleTree(chain[i].ComputeMerkleTree())
			if err == nil {
				bestCandidate = chain[i]
			}
			err = chain[i+1].VerifyMerkleTree(chain[i+1].ComputeMerkleTree())
			if err == nil {
				bestCandidate = chain[i+1]
			}

			if err != nil {
				i = i+2
				bestCandidate = chain[i]
			}

		}



	}

	return bestCandidate

}

func ValidateChain(chain ...*Blockchain) *Blockchain {

	var BestChains []*Blockchain

	for i, _ := range(chain) {
		BestChains = append(BestChains, CompareChains(chain[i], chain[i+1]))
	}

	var BestChain *Blockchain
	for i, _ := range(BestChains) {
		BestChain = ValidateMerkle(chain[i], chain[i+1])
	}

	return BestChain


}


func ValidateBlock(chain *Blockchain, b *pb.Block) error {


		ser, _ := json.Marshal(b)
		hasher := blake2b.New512()
		hasher.Write(ser)
		hash := hasher.Sum(nil)

		if bytes.Compare(b.Header.Hash, hash) == 1 {

			return nil

	} else {
		return errors.New("invalid block hash")
	}

}

