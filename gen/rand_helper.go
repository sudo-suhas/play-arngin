package gen

import (
	"math/rand"
	"time"
)

type randHelper struct {
	*rand.Rand
}

func newRandHelper() randHelper {
	return randHelper{
		Rand: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

// SampleInts generates a slice with random number of values taken from
// the input slice. The number of values does not exceed the given
// limit.
func (r randHelper) SampleInts(a []int, limit int) []int {
	// See https://github.com/golang/go/issues/23717#issuecomment-363464050
	switch len(a) {
	case 0:
		return nil
	case 1:
		return []int{a[0]}
	}

	// Create a copy of the list and shuffle it.
	tmp := append([]int(nil), a...)
	r.Shuffle(len(a), func(i, j int) { tmp[i], tmp[j] = tmp[j], tmp[i] })

	// Pick an index at random within the given limit, ensure it is not out of
	// bounds.
	k := r.Intn(limit-1) + 1
	if limit > len(a) {
		k = r.Intn(len(a)-1) + 1
	}

	// Copy the first k elements from the randomised list copy.
	sample := make([]int, k)
	copy(sample, tmp)
	return sample
}

// SampleStrings generates a slice with random number of values taken
// from the input slice. The number of values does not exceed the given
// limit.
func (r randHelper) SampleStrings(a []string, limit int) []string {
	switch len(a) {
	case 0:
		return nil
	case 1:
		return []string{a[0]}
	}

	tmp := append([]string(nil), a...)
	r.Shuffle(len(a), func(i, j int) { tmp[i], tmp[j] = tmp[j], tmp[i] })

	k := r.Intn(limit-1) + 1
	if limit > len(a) {
		k = r.Intn(len(a)-1) + 1
	}

	sample := make([]string, k)
	copy(sample, tmp)
	return sample
}

// Generate a random boolean, integer and float value.

func (r randHelper) Bool() bool               { return r.Float32() < 0.5 }
func (r randHelper) Floatn(n float64) float64 { return r.Float64() * n }

// Pick an element at random from the given slice.

func (r randHelper) IntEl(a []int) int          { return a[r.Intn(len(a))] }
func (r randHelper) StringEl(a []string) string { return a[r.Intn(len(a))] }
