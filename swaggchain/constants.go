package swaggchain

import (
	"encoding/hex"
	"time"
)

var Name, _ = hex.DecodeString("17BDFD4AE")
var Mainnet, _  = hex.DecodeString("94636F142A")
var Testnet, _ = hex.DecodeString("94636F1439")
var DevPublic, _ = hex.DecodeString("94636F1416")
var AddrVersion = byte(0x1B)
var PrivKeyVersion = byte(0x9F)
var BlockVersion = byte(0x1D)

var GenerateNewBlockEach = time.Minute
var CoinbaseMaturity = 1000
var GenerateOrphanBlocks = true
