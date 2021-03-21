package probcmp

import (
	"math"
	"strings"

	"github.com/go-dedup/megophone"
)

// Matcher returns the result of the comparison as a float64.
type Matcher interface {
	Match(a Comparable, b Comparable) (float64, error)
}

// Comparable provides a way for Matcher to read the comparison data.
type Comparable interface {
	GetField(name string) (string, bool)
}

// Deterministic defines the keys to be used for deterministic comparisons.
type Deterministic struct {
	Keys []string
}

// Match is a Matcher implementation of Deterministic.
// The result will be 0 or 1.
func (d *Deterministic) Match(a Comparable, b Comparable) (float64, error) {
	var result float64
	var err error
	for _, name := range d.Keys {
		var voa, vob string
		var ok bool
		if voa, ok = a.GetField(name); !ok {
			continue
		}
		if vob, ok = b.GetField(name); !ok {
			continue
		}
		if voa == vob {
			result = 1
			return result, err
		}
	}
	return result, err
}

// ComparablePatient is the patient data to be compared.
type ComparablePatient struct {
	Data map[string]string
}

// GetField makes ComparablePatient comparable to Matcher.
func (c *ComparablePatient) GetField(name string) (string, bool) {
	var result string
	var ok bool
	if result, ok = c.Data[name]; !ok {
		return result, false
	}
	return result, true
}

// Probability defines the Match-Probability and Unmatch-Probability of a Key
type Probability struct {
	Key         string
	MatchProb   float64
	UnmatchProb float64
	MatchFunc   func(a string, b string) bool
}

// Probabilistic defines the probability of a key to perform a probabilistic comparison
type Probabilistic struct {
	Probabilities []Probability
}

// Match is a Matcher implementation of Probabilistic.
func (p *Probabilistic) Match(a Comparable, b Comparable) (float64, error) {
	var result float64
	var err error
	for _, v := range p.Probabilities {
		var voa, vob string
		var ok bool
		if voa, ok = a.GetField(v.Key); !ok {
			continue
		}
		if vob, ok = b.GetField(v.Key); !ok {
			continue
		}

		var match bool
		if v.MatchFunc != nil {
			match = v.MatchFunc(voa, vob)
		} else {
			match = strings.ToLower(voa) == strings.ToLower(vob)
		}

		if match {
			result += math.Log(v.MatchProb / v.UnmatchProb)
		} else {
			result += math.Log((1 - v.MatchProb) / (1 - v.UnmatchProb))
		}
	}
	return result, err
}

// NameMatch compares with phonetic values
func NameMatch(a string, b string) bool {
	ap1, ap2 := megophone.DoubleMetaphone(a)
	bp1, bp2 := megophone.DoubleMetaphone(b)
	return ap1 == bp1 && ap2 == bp2
}
