package primality

import (
	"lukechampine.com/uint128"
	"testing"
)

func TestPrimeValidation(t *testing.T) {
	p1 := PrimeCheck(uint128.From64(22616281), false)
	if !p1 {
		t.Errorf("22616281 is prime but wasn't detected!!")
	}

	p2 := PrimeCheck(uint128.From64(75234869), false)
	if !p2 {
		t.Errorf("75234869 is prime but wasn't detected!!")
	}

	p3 := PrimeCheck(uint128.From64(177857809), false)
	if !p3 {
		t.Errorf("177857809 is prime but wasn't detected!!")
	}

	np1 := PrimeCheck(uint128.From64(22617281), false)
	if np1 {
		t.Errorf("22617281 isn't prime but was flagged!!")
	}

	np2 := PrimeCheck(uint128.From64(75235869), false)
	if np2 {
		t.Errorf("75235869 isn't prime but was flagged!!")
	}

	np3 := PrimeCheck(uint128.From64(177867809), false)
	if np3 {
		t.Errorf("177867809 isn't prime but was flagged!!")
	}
}

func TestTrivialPrimeValidation(t *testing.T) {
	if !PrimeCheck(uint128.From64(2), false) {
		t.Errorf("2 is prime but wasn't detected!!")
	}

	if !PrimeCheck(uint128.From64(3), false) {
		t.Errorf("3 is prime but wasn't detected!!")
	}

	if PrimeCheck(uint128.From64(4), false) {
		t.Errorf("4 isn't prime but was flagged!!")
	}

	if !PrimeCheck(uint128.From64(5), false) {
		t.Errorf("5 is prime but wasn't detected!!")
	}

	if PrimeCheck(uint128.From64(6), false) {
		t.Errorf("6 isn't prime but was flagged!!")
	}

	if !PrimeCheck(uint128.From64(7), false) {
		t.Errorf("7 is prime but wasn't detected!!")
	}

	if PrimeCheck(uint128.From64(8), false) {
		t.Errorf("8 isn't prime but was flagged!!")
	}

	if PrimeCheck(uint128.From64(9), false) {
		t.Errorf("9 isn't prime but was flagged!!")
	}
}
