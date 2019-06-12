package swaggchain

import (
	"crypto/rsa"
	"sync"
)

var mu sync.Mutex
var wg sync.WaitGroup
var PrivateKey *rsa.PrivateKey
var PublicKey *rsa.PublicKey

/*
func init() {


	wg.Add(1)

	defer wg.Done()
	go func() {

		//generate blockchain keypairs
		NewChainKeyPair()

	} ()

	wg.Wait()

	go func() {



	}()


}
*/

