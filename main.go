package main

import (
	"github.com/davecgh/go-spew/spew"
	_ "github.com/davecgh/go-spew/spew"
	"runtime"
	"swaggp2p/chain"
	_ "swaggp2p/chain"

)




func main() {
	runtime.GOMAXPROCS(8)
	dna := chain.GetDNA()
	spew.Dump(dna)

//repo.TrainForPow()

}