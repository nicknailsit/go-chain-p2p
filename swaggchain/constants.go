package core

import (
	"encoding/hex"
	"github.com/mr-tron/base58"
)

var PUBLIC, _ = hex.DecodeString(base58.Encode([]byte("737767D7")))

var PRIVATE, _ = hex.DecodeString(base58.Encode([]byte("0x73776774")))
var DEV, _ = hex.DecodeString(base58.Encode([]byte("0x73776764")))
var TEST, _ = hex.DecodeString(base58.Encode([]byte("0x73776774")))

