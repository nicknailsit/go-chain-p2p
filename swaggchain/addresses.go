package swaggchain


type Address []byte

type AddressBook []*Address


func NewAddress() *Address {

	return nil

}


func (a *Address) Generate(pubkey []byte) {



}

func (a *Address) Validate(address []byte) {

}

func (a *Address) GetPublicKeyFromBytes(address []byte) {

}

func (a *Address) Lock() error {

	return nil

}

func (a *Address) Unlock() error {

	return nil
}

func (a *Address) Encode() error {
	return nil
}

func (a *Address) Decode() error {
	return nil
}


func (A *AddressBook) GetAddressBook() AddressBook {

	return nil
}

func (A *AddressBook) add(addr *Address) error {
	return nil
}

func (A *AddressBook) remove(addr Address) error {
	return nil
}

func (a Address) String() string {
	return ""
}

