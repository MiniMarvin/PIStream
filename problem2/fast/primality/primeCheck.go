package primality

import (
	"lukechampine.com/uint128"
	"math/big"
	"strings"
)

var bounds []uint128.Uint128
var primes []uint128.Uint128
var didInit bool

func setPrimes() {
	primes = []uint128.Uint128{}
	primes = append(primes,
		uint128.From64(2),
		uint128.From64(3),
		uint128.From64(5),
		uint128.From64(7),
		uint128.From64(11),
		uint128.From64(13),
		uint128.From64(17),
		uint128.From64(19),
		uint128.From64(23),
		uint128.From64(29),
		uint128.From64(31),
		uint128.From64(37),
		uint128.From64(41))
}

func setBounds() {
	bounds = []uint128.Uint128{}
	u0, _ := uint128.FromString(strings.Replace("2_047", "_", "", -1))
	u1, _ := uint128.FromString(strings.Replace("1_373_653", "_", "", -1))
	u2, _ := uint128.FromString(strings.Replace("25_326_001", "_", "", -1))
	u3, _ := uint128.FromString(strings.Replace("3_215_031_751", "_", "", -1))
	u4, _ := uint128.FromString(strings.Replace("2_152_302_898_747", "_", "", -1))
	u5, _ := uint128.FromString(strings.Replace("3_474_749_660_383", "_", "", -1))
	u6, _ := uint128.FromString(strings.Replace("341_550_071_728_321", "_", "", -1))
	u7, _ := uint128.FromString(strings.Replace("1", "_", "", -1))
	u8, _ := uint128.FromString(strings.Replace("3_825_123_056_546_413_051", "_", "", -1))
	u9, _ := uint128.FromString(strings.Replace("1", "_", "", -1))
	u10, _ := uint128.FromString(strings.Replace("1", "_", "", -1))
	u11, _ := uint128.FromString(strings.Replace("318_665_857_834_031_151_167_461", "_", "", -1))
	u12, _ := uint128.FromString(strings.Replace("3_317_044_064_679_887_385_961_981", "_", "", -1))

	bounds = append(bounds, u0, u1, u2, u3, u4, u5, u6, u7, u8, u9, u10, u11, u12)
}

func lastMayBePrime(n uint128.Uint128) bool {
	return n.Mod(uint128.From64(10)).Equals(uint128.From64(1)) ||
		n.Mod(uint128.From64(10)).Equals(uint128.From64(3)) ||
		n.Mod(uint128.From64(10)).Equals(uint128.From64(7)) ||
		n.Mod(uint128.From64(10)).Equals(uint128.From64(9))
}

func trivialPrimeCheck(n uint128.Uint128, allowProbable bool) bool {
	if (n.Cmp(uint128.From64(2)) < 0) || n.Mod(uint128.From64(2)).Equals(uint128.From64(0)) {
		return false
	}
	// check if is divisible by 2 or 5
	if n.Cmp(uint128.From64(5)) > 0 && !lastMayBePrime(n) {
		return false
	}
	if n.Cmp(bounds[12]) > 0 && !allowProbable {
		// returning false to avoid raise an error
		return false
	}

	return true
}

func millerRabin(n uint128.Uint128) bool {
	var usedPrimes []uint128.Uint128

	for i := 0; i < len(bounds); i++ {
		if n.Cmp(bounds[i]) < 0 {
			usedPrimes = primes[:i+1]
			break
		}
	}

	d := n.Sub(uint128.From64(1))
	var s uint64
	s = 0

	for d.Mod(uint128.From64(2)).Equals(uint128.From64(0)) {
		d = d.Div(uint128.From64(2))
		s += 1
	}

	for _, prime := range usedPrimes {
		pr := false
		var r uint64
		for r = 0; r < s; r++ {
			bigPrime := prime.Big()
			bigD := d.Big()
			bigR := big.NewInt(int64(r))
			bigN := n.Big()
			big2 := big.NewInt(2)
			bigM := big.NewInt(0)
			bigTmp := big.NewInt(0)
			bigMulRes := big.NewInt(0)

			bigTmp.Exp(big2, bigR, nil)
			bigMulRes.Mul(bigD, bigTmp)
			bigM.Exp(bigPrime, bigMulRes, bigN)

			m := uint128.FromBig(bigM)

			if (r == 0 && m.Equals(uint128.From64(1))) ||
				(m.Add(uint128.From64(1)).Mod(n).Equals(uint128.From64(0))) {
				pr = true
				break
			}

		}

		if pr {
			continue
		}
		return false
	}

	return true
}

/**
 * This is a function that exists in order
 * to validate if a number is prime or not
 * and it's based on the miller_rabin approach
 * implemented in python by Nathan Damon,
 * @bizzfitch on github.
 * This is a deterministic Miller-Rabin algorithm
 * for primes ~< 3.32e24.
 **/
func PrimeCheck(n uint128.Uint128, allowProbable bool) bool {
	if !didInit {
		didInit = true
		setBounds()
		setPrimes()
	}

	if n.Cmp(uint128.From64(2)) == 0 {
		return true
	}

	mayBePrime := trivialPrimeCheck(n, allowProbable)
	if !mayBePrime {
		return false
	}

	return millerRabin(n)
}
