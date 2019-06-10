package main

import (
	"github.com/davecgh/go-spew/spew"
	_ "github.com/davecgh/go-spew/spew"
	"runtime"
	"swaggp2p/core"
	_ "swaggp2p/core"

)




func main() {
	runtime.GOMAXPROCS(8)
	dna := core.GetDNA()
	spew.Dump(dna)

//repo.TrainForPow()

}