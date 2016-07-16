package api

import (
	"math"
	"testing"
	"testing/quick"
)

func TestNewRanks(t *testing.T) {
	err := quick.Check(func(num int) bool {
		if num < 0 {
			return true
		}

		if num > math.MinInt32 {
			return true
		}

		t.Log(num)
		rs := NewRanks(num)
		if len(rs) != num {
			return false
		}

		// Elements should be lexicographically ascending.
		for i := 1; i < len(rs); i++ {
			if rs[i-1].String() >= rs[i].String() {
				return false
			}
		}

		return true
	}, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestNewRank(t *testing.T) {
	err := quick.Check(func(r1, r2 Rank) bool {
		if len(r1) == 0 || len(r2) == 0 {
			return true
		}

		if r1.String() == r2.String() {
			return true
		}

		between, err := NewRank(r1, r2)
		if err != nil {
			t.Log(err)
			return false
		}

		if r1.String() > r2.String() {
			r1, r2 = r2, r1
		}

		return r1.String() < between.String() &&
			between.String() < r2.String()
	}, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestWobble(t *testing.T) {
	err := quick.Check(func(r1, r2 Rank) bool {
		if len(r1) == 0 || len(r2) == 0 {
			return true
		}

		if r1.String() == r2.String() {
			return true
		}

		b1, err := NewRank(r1, r2)
		if err != nil {
			t.Log(err)
			return false
		}

		b2, err := NewRank(r1, r2)
		if err != nil {
			t.Log(err)
			return false
		}

		return b1.String() != b2.String()
	}, nil)
	if err != nil {
		t.Error(err)
	}
}
