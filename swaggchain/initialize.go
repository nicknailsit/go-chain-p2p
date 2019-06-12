package swaggchain

import "sync"

var mu sync.Mutex
var wg sync.WaitGroup
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

