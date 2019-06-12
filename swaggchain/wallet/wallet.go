package wallet

import (
	"go/types"
	"swaggp2p/pb"
)

type WalletInterface interface {
	AddWalletOptions(opts types.Object) bool
	UnloadWallet(wallet *Wallet) bool
	RemoveWallet(wallet *Wallet) bool
	GetWallets() []*Wallet
	GetWallet(name string) *Wallet
	LoadWallet(chainID string, location string) error
	LockWallet(w *Wallet, secret string) error
	UnlockWallet(w *Wallet, secret string) error
	Create(chainID string, version []byte, prefix byte, checksum []byte, coinType []byte) *Wallet
	CreateMultisig(numSig int, hdpubkey []byte, hdprivkey []byte) error
	Serialize() []byte
	Unserialize() *struct{}
	GetUTXOS() []pb.Transaction_Output
	SignTX(tx *pb.Transaction)
	SignUTXOS(tx []pb.Transaction)
}

type Wallet struct {
	index int64
	pubKey *pubKey
	internal bool
}

type pubKey []byte


type KeypoolInterface interface {
	//reserve key from keypool
	GetReservedKey(pubkey []byte, internal bool) error
	//return key to the keypool
	ReturnKey() error
}

type Keypool []pubKey

type AddressBookInterface interface {


}

type AddressBook struct {
	Name string
	Purpose string
	destData map[string]string
}

type Recipient struct {
	scriptPubKey []byte
	Amount float64
	AmountMinusFees float64
}
