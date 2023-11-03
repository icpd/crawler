// Package recommendation provides a simple implementation mimicing a recommendation system.
package recommendation

import "math/rand"

// If given id is 20, we may generate 10, 15, 25, ...
// Intentionally desgined to increase duplication within recommendation cycles,
// to mimic real-world behvior.
const (
	step     = 5
	offset   = -10
	maxGuess = 7
)

// Client defines the recommendation client.
type Client struct{}

// NextGuess returns the next batch of guesses for the given id.
func (Client) NextGuess(id int64) []int64 {
	var res []int64

	// We use a global random function because:
	// 1. We want recommendation result to be different each time.
	// 2. Only global random function is thread-safe.
	guessCnt := rand.Intn(maxGuess)

	for i, r := 0, id+offset; i < guessCnt; {
		if r > 0 && r != id {
			res = append(res, r)
		}
		i++
		r += step
	}
	return res
}
