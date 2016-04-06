package api

import (
	"fmt"
	"math"
	"math/big"
)

const (
	MaxRank = math.MaxInt64
	MinRank = math.MinInt64
)

func NewRanks(n int) ([]int64, error) {
	rs := make([]int64, n)
	err := newRanks(rs, MinRank, MaxRank)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func newRanks(ranks []int64, prev, next int64) error {
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

func NewRank(prev, next int64) (int64, error) {
	sum := new(big.Int).Add(big.NewInt(prev), big.NewInt(next))
	quo := new(big.Int).Div(sum, big.NewInt(2))
	r := quo.Int64()
	if r == prev || r == next {
		return 0, fmt.Errorf("cannot create new rank between %d and %d", prev, next)
	}
	return r, nil
}
