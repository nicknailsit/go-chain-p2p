package swaggchain

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"os"
)

type ChainKeyPair struct {

	PrivKey *rsa.PrivateKey
	PubKey rsa.PublicKey

}


const BITSIZE = 4096



func (c *ChainKeyPair) Generate()  {

	reader := rand.Reader
	key, err := rsa.GenerateKey(reader, BITSIZE)
	checkError(err)

	publicKey := key.PublicKey

	saveGobKey("private.key", key)
	savePEMKey("private.pem", key)
	saveGobKey("public.key", publicKey)
	savePublicPEMKey("public.pem", publicKey)

}

func LoadKeys() (interface{}, error) {

	priv, err := os.Open("private.pem")
	checkError(err)

	pemfileinfo, _ := priv.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)
	buffer := bufio.NewReader(priv)
	_, err = buffer.Read(pembytes)
	data, _ := pem.Decode([]byte(pembytes))

	_ = priv.Close()

	privImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	checkError(err)

	pubKey := privImported.PublicKey

	return &ChainKeyPair{
		privImported,pubKey,
	}, nil



}




func saveGobKey(fileName string, key interface{}) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	encoder := gob.NewEncoder(outFile)
	err = encoder.Encode(key)
	checkError(err)
}

func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	checkError(err)
}

func savePublicPEMKey(fileName string, pubkey rsa.PublicKey) {
	asn1Bytes, err := asn1.Marshal(pubkey)
	checkError(err)

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemfile, err := os.Create(fileName)
	checkError(err)
	defer pemfile.Close()

	err = pem.Encode(pemfile, pemkey)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}