package crypto

import "math/big"


var tt256m1 = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

type Number struct {
	num *big.Int
	limit func(n *Number) *Number
}

func Uint256(n int64) *Number {
	return &Number{big.NewInt(n), limitUnsigned256}
}


func limitUnsigned256(x *Number) *Number {
	x.num.And(x.num, tt256m1)
	return x
}

func (i *Number) Uint256() *Number {
	return Uint(0).Set(i)
}

var (
	Zero       = Uint(0)
	One        = Uint(1)
	Two        = Uint(2)

	// "typedefs"
	Uint = Uint256

)