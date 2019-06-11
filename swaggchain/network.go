package swaggchain

import (
	"encoding/hex"
)

type Network struct {
	name string
	symbol string
	pubkey byte
	privkey byte
	hdkey [4]byte
	hdkeypub [4]byte
}

// 63 73 23 38 38 = swagg decimal hex = 0x17BDFD4AE
// 63 73 23 38 38 50-51 = SWAGGM (mainnet) hex = 0x94636F142A
// + 65-66 SWAGGT (testnet) = 0x94636F1439
// + 30-31 SWAGGD (devpublic) = 0x94636F1416
// Symbol SWG


//Network Versions Hex

var Name, _ = hex.DecodeString("17BDFD4AE")
var Mainnet, _  = hex.DecodeString("94636F142A")
var Testnet, _ = hex.DecodeString("94636F1439")
var DevPublic, _ = hex.DecodeString("94636F1416")


//
// swagg it address prefix
//

var AddrVersion = byte(0x1B)
var PrivKeyVersion = byte(0x9F)

