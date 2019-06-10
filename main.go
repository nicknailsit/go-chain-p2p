package main

import (
	"github.com/davecgh/go-spew/spew"
	"swaggp2p/swaggchain"
)




func main() {


	addr := swaggchain.GetNewAddress()
	spew.Dump(addr)

//repo.TrainForPow()

}