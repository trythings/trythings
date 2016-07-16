package api

import (
	"bytes"
	"fmt"
	"math/big"
)

type Rank []byte

func (r Rank) String() string {
	return fmt.Sprintf("Rank{%#x}", []byte(r))
}

var (
	MinRank = Rank{0x00}
	MaxRank = Rank{0xff}
)

func NewRanks(num int) []Rank {
	if num == 0 {
		return nil
	}

	rs := make([]Rank, num)
	err := newRanks(rs, MinRank, MaxRank)
	if err != nil {
		// This could only happen if MinRank == MaxRank.
		panic(err)
	}
	return rs
}

func newRanks(ranks []Rank, prev, next Rank) error {
	if len(ranks) == 0 {
		return nil
	}

	r, err := NewRank(prev, next)
	if err != nil {
		return err
	}

	// Middle
	mid := len(ranks) / 2
	ranks[mid] = r

	// Left
	err = newRanks(ranks[:mid], prev, r)
	if err != nil {
		return err
	}

	// Right
	err = newRanks(ranks[mid+1:], r, next)
	if err != nil {
		return err
	}

	return nil
}

// wobble reduces the likelihood of two ranks colliding
// if they are generated between the same prev and next.
var wobble byte

func NewRank(prev, next Rank) (Rank, error) {
	// We want to align the most significant (leftmost) digits
	// of prev and next before we add them.
	// This is equivalent to aligning the radix points at beginning of the numbers.
	numBytes := len(prev)
	if len(next) > numBytes {
		numBytes = len(next)
	}
	// We add padding because we don't want to lose a trailing one during integer division.
	// We will trim trailing zeros after division.
	numBytes++

	paddedPrev := make([]byte, numBytes)
	copy(paddedPrev, prev)
	p := new(big.Int).SetBytes(paddedPrev)

	paddedNext := make([]byte, numBytes)
	copy(paddedNext, next)
	n := new(big.Int).SetBytes(paddedNext)

	if p.Cmp(n) == 0 {
		return nil, fmt.Errorf("cannot generate rank between %s and %s", prev, next)
	}

	// Find the average of the two ranks.
	sum := new(big.Int).Add(p, n)
	quo := new(big.Int).Div(sum, big.NewInt(2))

	r := quo.Bytes()
	r = bytes.TrimRight(r, "\000")
	if wobble > 0 {
		r = append(r, wobble)
	}
	// It is intended that wobble overflow and wrap from 255 to 0.
	wobble++

	return r, nil
}
