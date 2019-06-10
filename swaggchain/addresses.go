package swaggchain

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/minio/sha256-simd"
	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
	"crypto/elliptic"
	"math/big"
	"swaggp2p/core"
)

var (
	Public []byte
	Private []byte
	Dev []byte
	Test []byte
)

type Address struct {

	V []byte // 4 bytes
	D uint16 // 1 byte
	F []byte // 4 bytes
	I []byte // 4 bytes
	C []byte //32 bytes
	K []byte //33 bytes

}

var curve = elliptic.P256()


const key = "VERY VERY SECRET SWAGG KEY" // will change for the network key and you won't know it ;) uh.. me either anyway
var master *Address
var masterpub *Address

type AddressBook []*Address


func init() {

	master = NewMaster([]byte(key))
	masterpub = master.Pub()


}

func GetNewAddress() string {

	_, pub := getKeysForAddress()
	address := pub.String()
	return address

}


func NewMaster(key []byte) *Address {

	mac := hmac.New(sha512.New, key)
	seed, _ := generateSeed(512)
	mac.Write(seed)
	I := mac.Sum(nil)
	secret := I[:len(I)/2]
	C := I[len(I)/2:]
	D := 0
	i := make([]byte, 4)
	F := make([]byte, 4)
	zero := make([]byte, 1)
	return &Address{core.PRIVATE, uint16(D), F, i, C, append(zero, secret...)}

}


func getKeysForAddress() (*Address, *Address) {

	childprv, err := master.Child(0)
	childpub, err := masterpub.Child(0)

	if err != nil {
		panic(err)
	}
	return childprv, childpub
}

func generateSeed(length int) ([]byte, error) {

	b := make([]byte, length)
	if length < 128 {
		return b, errors.New("length must be at least 128 bits")
	}
	_, err := rand.Read(b)
	return b, err
}

func (a *Address) Pub() *Address {
	if bytes.Compare(a.V, core.PUBLIC) == 0 {
		return &Address{a.V, a.D, a.F, a.I, a.C, a.K}
	} else {
		return &Address{core.PUBLIC, a.D, a.F, a.I, a.C, privToPub(a.K)}
	}
}

func (a *Address) Child(i uint32) (*Address, error) {

	var f, I, newkey []byte
	switch {
	case bytes.Compare(a.V, core.PRIVATE) == 0, bytes.Compare(a.V, core.TEST) == 0 :
		pub := privToPub(a.K)
		mac := hmac.New(sha512.New, a.C)
		if i >= uint32(0x80000000) {
			mac.Write(append(a.K, uint32ToByte(i)...))
		} else {
			mac.Write(append(pub, uint32ToByte(i)...))
		}
		I = mac.Sum(nil)
		iL := new(big.Int).SetBytes(I[:32])
		if iL.Cmp(curve.Params().N) >= 0 || iL.Sign() == 0 {
			return &Address{}, errors.New("invalid child")
		}
		newkey = addPrivKeys(I[:32], a.K)
		f = hash160(privToPub(a.K))[:4]

	case bytes.Compare(a.V, core.PUBLIC) == 0, bytes.Compare(a.V, core.TEST) == 0:
		mac := hmac.New(sha512.New, a.C)
		if i >= uint32(0x80000000) {
			return &Address{}, errors.New("Can't do Private derivation on Public key!")
		}
		mac.Write(append(a.K, uint32ToByte(i)...))
		I = mac.Sum(nil)
		iL := new(big.Int).SetBytes(I[:32])
		if iL.Cmp(curve.Params().N) >= 0 || iL.Sign() == 0 {
			return &Address{}, errors.New("invalid child")
		}
		newkey = addPubKeys(privToPub(I[:32]), a.K)
		f = hash160(a.K)[:4]

	}
	return &Address{a.V, a.D+1, f, uint32ToByte(i), I[32:], newkey}, nil

}


func (a *Address) Serialize() []byte {

	depth := uint16ToByte(uint16(a.D % 256))
	bindata := make([]byte, 78)
	copy(bindata, a.V)
	copy(bindata[4:], depth)
	copy(bindata[5:], a.F)
	copy(bindata[9:], a.I)
	copy(bindata[13:],a.C)
	copy(bindata[45:], a.K)
	chksum := dblSha256(bindata)[:4]
	return append(bindata, chksum...)

}

func (a *Address) String() string {
	return base58.FastBase58Encoding(a.Serialize())
}


func privToPub(key []byte) []byte {
	return compress(curve.ScalarBaseMult(key))
}

func compress(x, y *big.Int) []byte {
	two := big.NewInt(2)
	rem := two.Mod(y, two).Uint64()
	rem += 2
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(rem))
	rest := x.Bytes()
	pad := 32 - len(rest)
	if pad != 0 {
		zeroes := make([]byte, pad)
		rest = append(zeroes, rest...)
	}
	return append(b[1:], rest...)
}

func expand(key []byte) (*big.Int, *big.Int) {
	params := curve.Params()
	exp := big.NewInt(1)
	exp.Add(params.P, exp)
	exp.Div(exp, big.NewInt(4))
	x := big.NewInt(0).SetBytes(key[1:33])
	y := big.NewInt(0).SetBytes(key[:1])
	beta := big.NewInt(0)
	beta.Exp(x, big.NewInt(3), nil)
	beta.Add(beta, big.NewInt(7))
	beta.Exp(beta, exp, params.P)
	if y.Add(beta, y).Mod(y, big.NewInt(2)).Int64() == 0 {
		y = beta
	} else {
		y = beta.Sub(params.P, beta)
	}
	return x, y
}

func addPrivKeys(k1, k2 []byte) []byte {
	i1 := big.NewInt(0).SetBytes(k1)
	i2 := big.NewInt(0).SetBytes(k2)
	i1.Add(i1, i2)
	i1.Mod(i1, curve.Params().N)
	k := i1.Bytes()
	zero, _ := hex.DecodeString("00")
	return append(zero, k...)
}

func addPubKeys(k1, k2 []byte) []byte {
	x1, y1 := expand(k1)
	x2, y2 := expand(k2)
	return compress(curve.Add(x1, y1, x2, y2))
}

func uint32ToByte(i uint32) []byte {
	a := make([]byte, 4)
	binary.BigEndian.PutUint32(a, i)
	return a
}

func uint16ToByte(i uint16) []byte {
	a := make([]byte, 2)
	binary.BigEndian.PutUint16(a, i)
	return a[1:]
}

func byteToUint16(b []byte) uint16 {
	if len(b) == 1 {
		zero := make([]byte, 1)
		b = append(zero, b...)
	}
	return binary.BigEndian.Uint16(b)
}

func onCurve(x, y *big.Int) bool {
	return curve.IsOnCurve(x, y)
}

func hash160(data []byte) []byte {
	sha := sha256.New()
	ripe := ripemd160.New()
	sha.Write(data)
	ripe.Write(sha.Sum(nil))
	return ripe.Sum(nil)
}

func dblSha256(data []byte) []byte {
	sha1 := sha256.New()
	sha2 := sha256.New()
	sha1.Write(data)
	sha2.Write(sha1.Sum(nil))
	return sha2.Sum(nil)
}

