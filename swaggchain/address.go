package swaggchain

import "C"
import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/dongri/go-mnemonic"
	"golang.org/x/crypto/argon2"

	"time"
)

type Address struct {
	prefix byte
	nVersion []byte
	pubKey []byte
	privKey []byte
	timeStamp int64
}

type paramsAddr struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}




func createSeed(n uint32) (buf []byte, err error) {

	buf = make([]byte, n)
	_, err = rand.Read(buf)
	return

}

var p paramsAddr

func init() {

	p = paramsAddr{
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



func (a *Address) New() (*Address, string, string) {

	priv, _ := generateECDSAKeyPair()
	m, _ := mnemonic.GenerateMnemonic(256, mnemonic.LanguageEnglish)
	salt:= mnemonic.ToSeedHex(m, hex.EncodeToString(priv.X.Bytes()))
	addr := &Address{
		prefix: AddrVersion,
		nVersion: DevPublic,
		pubKey: []byte(base58.Encode(a.Encode(priv.Y.Bytes(), []byte(salt), p))),
		timeStamp: time.Now().UTC().UnixNano(),
	}




	return addr, salt, m

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


func (a *Address) Encode(data []byte, salt []byte, p paramsAddr) (encodedHash []byte) {



	encoded := argon2.IDKey([]byte(data), salt, p.iterations, p.memory, p.parallelism, p.keyLength)
	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(encoded)

	// Return a string using the standard encoded hash representation.
	encodedHash1 := fmt.Sprintf("%s",  bytes.Join([][]byte{
		[]byte(b64Hash),
		[]byte(b64Salt),
	},nil))
	encodedHash = []byte(encodedHash1)


	return


}


//CreateNewAddress @returns string address, mnemonic

func CreateNewAddress() (string, string) {

	addr := &Address{}
	a, _, m := addr.New()

	return base58.Encode(a.Serialize()), m

}