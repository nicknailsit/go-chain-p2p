package wallet

import "C"
import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	. "github.com/mr-tron/base58"
	"golang.org/x/crypto/argon2"
	"log"
	"strings"
	"time"
)

type Address struct {
	prefix byte
	nVersion []byte
	pubKey []byte
	privKey []byte
	timeStamp int64
}

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}
var Name, _ = hex.DecodeString("17BDFD4AE")
var Mainnet, _  = hex.DecodeString("94636F142A")
var Testnet, _ = hex.DecodeString("94636F1439")
var DevPublic, _ = hex.DecodeString("94636F1416")
var AddrVersion = byte(0x1B)
var PrivKeyVersion = byte(0x9F)



func createSeed(n uint32) (buf []byte, err error) {

	buf = make([]byte, n)
	_, err = rand.Read(buf)
	return

}

var p params

func init() {

	p = params{
		memory: 64 * 1024,
		iterations: 3,
		parallelism: 2,
		saltLength: 16,
		keyLength: 32,
	}

}

func generateECDSAKeyPair() (priv *ecdsa.PrivateKey, pub ecdsa.PublicKey){

	priv, _ = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)

	pub = priv.PublicKey
	return

}



func (a *Address) New() *Address {

	priv, _ := generateECDSAKeyPair()
	addr := &Address{
		prefix: AddrVersion,
		nVersion: DevPublic,
		pubKey: a.Encode(priv.Y.Bytes(), p),
		privKey: a.Encode(priv.X.Bytes(), p),
		timeStamp: time.Now().UTC().UnixNano(),
	}



	return addr

}

func (a *Address) Serialize() []byte {

	buf := make([]byte, 45)
	copy(buf, []byte(string(a.prefix)))
	copy(buf[1:6], a.nVersion)
	copy(buf[7:38],a.pubKey)

	hasher := sha256.New()
	hasher.Write(buf)
	chksum := hasher.Sum(hasher.Sum(nil))[:4]
	copy(buf[39:42], chksum)



	return buf


}


func (a *Address) Encode(data []byte, p params) (encoded []byte) {

	salt, err := createSeed(p.saltLength)
	if err != nil {
		panic(err)
	}

	encoded = argon2.IDKey([]byte(data), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	return


}

func CreateNewAddress() {

	addr := &Address{}
	a := addr.New()

	encoded := base64.StdEncoding.EncodeToString(a.Serialize())

	log.Printf("address: %s", strings.ToUpper(Encode([]byte(encoded))))


}