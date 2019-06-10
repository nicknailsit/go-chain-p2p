package swaggchain

import (

	"crypto"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/minio/blake2b-simd"
	"swaggp2p/core"

	"strconv"
	"swaggp2p/core/pb"

	"time"
)

type GenesisBlock pb.Block
type GBlockHeader pb.BlockHeader


const G_BLOCKINDEX = 0
var G_VERSION, _ = hex.DecodeString(string(core.PUBLIC))
const G_REWARD = 1000000
var G_TXCount = 0
var G_TXVALUE = float32(0)
var G_MERKLE_ROOT = []byte(nil)
var TIMESTAMP = time.Now().UTC().UnixNano()


/*
GenesisBlock instance of pb.Block

message Block {
    uint32 blockindex = 1;
    bytes version = 2;
    BlockHeader header =3;
    uint32 size = 4;
    float reward = 5;
    float txvalues = 6;
    uint32 txcount = 7;
    bytes merkleroot = 8;
    repeated Transaction transactions = 9;
    int64 timestamp = 10;
}

message BlockHeader {
    uint32 chainid = 1;
    uint32 height = 2;
    uint32 difficulty = 3;
    bytes nonce = 4;
    bytes hash = 5;
    bytes signature = 6;
    int64 timestamp = 7;
    bytes blocktype=8;
    bytes blockversion=9;
    bytes blocknetwork=10;
}
 */

func (g *GenesisBlock) Create(chainID string) *pb.Block {


	keypair := new(ChainKeyPair)
	keypair.Generate()

	BH := new(GBlockHeader)
	header := BH.Create(chainID)

	priv := keypair.PrivKey

	B := new(pb.Block)
	B.Header = header
	B.Blockindex = G_BLOCKINDEX
	B.Version = []byte{0x73}
	B.Reward = G_REWARD
	B.Txcount = uint32(G_TXCount)
	B.Txvalues = G_TXVALUE
	B.Merkleroot = G_MERKLE_ROOT
	B.Timestamp = TIMESTAMP

	pbBytes, _ := json.Marshal(B)

	hashByte := blake2b.New512()
	hashByte.Reset()
	hashByte.Write(pbBytes)
	reader := rand.Reader


	signature, _ := priv.Sign(reader, hashByte.Sum(nil), crypto.BLAKE2b_512)
	B.Header.Signature = signature

	return B

}

func (gh *GBlockHeader) Create(chainID string) *pb.BlockHeader {

	chainIDint, _ := strconv.Atoi(chainID)

	header := &pb.BlockHeader{

		Chainid: uint32(chainIDint),
		Height: 1,
		Difficulty: 127,
		Nonce: make([]byte, 0),
		Signature: make([]byte, 0),
		Timestamp: TIMESTAMP,
		Blocktype: make([]byte, 0),
		Blockversion: make([]byte, 0),
		Blocknetwork: make([]byte, 0),
	}

	return header

}

