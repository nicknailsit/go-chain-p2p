package swaggchain

import (
	"encoding/hex"
	"encoding/json"
)

//COIN DEFINITION


type SWAGG Coin


func (s *SWAGG) Create() *SWAGG {

	coin := new(SWAGG)
	coin.Name = "SWAGG"
	coin.IsDefault = true
	coin.IsInMainChain = true
	coin.Mineable = true
	coin.MaxAllowed = -1
	coin.MinAllowed = 1e-8
	coin.Exchangeable = true
	coin.Spendable = true
	coin.MaxAvailable = 1e8
	coin.TotalAvailable = 1e8
	coin.Symbol = "SWG"
	coin.Unlocked = 0
	coin.CoinbaseAddr = ""

	return coin

}

func (c *SWAGG) buy() error {

	return nil

}

func (c *SWAGG) mint() error {
	return nil
}

func (c *SWAGG) burn() error {
	return nil
}

func (c *SWAGG) exchange() error {
	return nil
}

func (c *SWAGG) swap() error {
	return nil
}


func (cb *Coinbase) Create() *Coinbase {

	addr, _ := CreateNewAddress()
	addrBytes, _ := hex.DecodeString(addr)
	swagg := &SWAGG{}

	cbase := &Coinbase{
		Address: addrBytes,
		Coin:swagg.Create(),

	}

	return cbase

}

func (cb *Coinbase) Serialize() []byte {

	cs, _ := json.Marshal(cb)
	return cs


}

func (cb *Coinbase) AppendToChain(bc SwaggChain) {

	bc.coinbase = cb

}