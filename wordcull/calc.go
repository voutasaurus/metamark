package wordcull

import (
	"math"
)

// BetterThan takes in three known values in order to calculate
// just how good a passphrase is in terms of some other password
// values.
//
// Here mw represents the number of words the user chooses to
// include in their passphrase; nw represents the number of words
// in the database; and nc represents the number of characters in
// the character set of the other password values.
// These are, for some examples:
// 	Alphabetic      52
//	Alphanumeric    62
//	Printable       95
func BetterThan(mw, nw, nc int) float64 {
	return float64(mw) * math.Log2(float64(nw)) / math.Log2(float64(nc))
}

// TODO(Luke): Write hooks for form passing
