package core

import (
	"swaggp2p/core/pb"
	"swaggp2p/swaggchain/wallet"
)

//coin definition

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
	coin.MaxAvailable = 1e14
	coin.TotalAvailable = 1e10
	coin.Symbol = "SWG"
	coin.Unlocked = 0
	coin.CoinbaseAddr = ""

	return coin

}

func createCoinbase() *Coinbase {

	address, _ := wallet.CreateNewAddress()
	c := new(Coinbase)
	c.Address = []byte(address)
	swagg := createSWAGG()
	swagg.CoinbaseAddr = address
	c.Coin = swagg

	return c

}

func createSWAGG() *SWAGG {

	S := new(SWAGG)
	return S.Create()

}

func CoinbaseTX() *pb.Transaction {
	return nil
}

func (c *Coinbase) Sign() error {
	return nil
}

func (c *Coinbase) Validate() error {
	return nil
}

func (c *Coinbase) Release() error {
	return nil
}

func (c *Coinbase) Reject() {

}

func (c *Coinbase) Refresh() {

}

func (c *Coinbase) GetCoinInfo() {

}

func (c *Coinbase) GetAddress() string {
	return ""
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








