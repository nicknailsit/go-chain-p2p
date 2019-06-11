package wallet

import (
	"swaggp2p/core"
	"swaggp2p/core/pb"
)

type WalletTX struct {
	Transaction *pb.Transaction
	core.TX
}

func (w *WalletTX) Spend(coinType []byte, amount float64, address string, to string) (*pb.Transaction, error) {



}

func (w *WalletTX) CreateOutput(input *pb.Transaction_Input, index int, spendable bool, solvable bool, useMaxSigInput bool) (*pb.Transaction_Output, error) {

}

func (w *WalletTX) CreateInput(tx *pb.Transaction, isCoinbase bool) (*pb.Transaction_Input, error) {

}

func (w *WalletTX) Validate(tx *pb.Transaction) error {

}

func (w *WalletTX) GetCoinType(coinType []byte) (*core.Coin, error) {

}