package probcmp

import (
	"fmt"
	"testing"
)

var empty = ComparablePatient{Data: map[string]string{}}
var katy1 = ComparablePatient{Data: map[string]string{"address": `{"street": "9993 Lola St", "zip": "73924"}`, "dob": "1980-10-04", "fname": "Katy", "lname": "Framingham", "medicalRecordNumber": "55427", "ssn": "710359155", "visitNumber": "442878"}}
var katy2 = ComparablePatient{Data: map[string]string{"dob": "1980-10-04", "fname": "Katie", "lname": "Framingham", "heartRate": "116 bpm"}}
var brian = ComparablePatient{Data: map[string]string{"address": `{"street": "4761 Blossom Glens", "zip": "49867"}`, "dob": "2017-11-24", "fname": "Brian", "lname": "Lang", "medicalRecordNumber": "78831", "ssn": "467813477", "visitNumber": "825889"}}
var julian = ComparablePatient{Data: map[string]string{"address": `{"street": "4761 Blossom Glens", "zip": "49867"}`, "dob": "2017-11-24", "fname": "Julian", "lname": "Lang", "medicalRecordNumber": "78832", "ssn": "681812714"}}
var jerry = ComparablePatient{Data: map[string]string{"address": `{"street": "4729 Esteban Hills", "zip": "59030"}`, "dob": "1958-02-06", "fname": "Jerry", "lname": "Johnson", "medicalRecordNumber": "34785", "ssn": "335345129", "visitNumber": "761325"}}
var lukas1 = ComparablePatient{Data: map[string]string{"address": `{"street": "1066 Maple Rd", "zip": "82403"}`, "dob": "1998-09-25", "fname": "Lukas", "lname": "Johnson", "medicalRecordNumber": "34788", "ssn": "335345129", "visitNumber": "761329"}}
var lukas2 = ComparablePatient{Data: map[string]string{"address": `{"street": "1066 Maple Rd", "zip": "82403"}`, "dob": "1998-09-25", "fname": "Lukas", "lname": "Johnson", "medicalRecordNumber": "34788", "ssn": "283737684", "visitNumber": "761329"}}

func TestDeterministic_Match(t *testing.T) {
	assertMatch := func(want float64, m Matcher, a Comparable, b Comparable) {
		got, err := m.Match(a, b)
		if err != nil {
			t.Error(err)
			return
		}
		if got != want {
			t.Errorf("got %f want %f", got, want)
		}
	}

	m := Deterministic{Keys: []string{"ssn", "visitNumber", "medicalRecordNumber"}}

	for _, tt := range []struct {
		test    string
		want    float64
		matcher Deterministic
		a       ComparablePatient
		b       ComparablePatient
	}{
		{test: "katy", want: 0, matcher: m, a: katy1, b: katy2},
		{test: "same_ssn", want: 1, matcher: m, a: jerry, b: lukas1},
		{test: "different_ssn", want: 1, matcher: m, a: lukas1, b: lukas2},
	} {
		t.Run(tt.test, func(t *testing.T) {
			assertMatch(tt.want, &tt.matcher, &tt.a, &tt.b)
		})
	}
}

func TestProbabilistic_Match(t *testing.T) {
	assertMatch := func(want float64, m Matcher, a Comparable, b Comparable) {
		got, err := m.Match(a, b)
		if err != nil {
			t.Error(err)
			return
		}
		if fmt.Sprintf("%g", got) != fmt.Sprintf("%g", want) {
			fmt.Printf("%g\n", got)
			fmt.Printf("%g\n", want)
			t.Errorf("got %f want %f", got, want)
		}
	}

	m := Probabilistic{
		Probabilities: []Probability{
			{
				Key:         "address",
				MatchProb:   0.90,
				UnmatchProb: 0.10,
			},
			{
				Key:         "dob",
				MatchProb:   0.95,
				UnmatchProb: 0.01,
			},
			{
				Key:         "fname",
				MatchProb:   0.60,
				UnmatchProb: 0.20,
				MatchFunc:   NameMatch,
			},
			{
				Key:         "lname",
				MatchProb:   0.90,
				UnmatchProb: 0.10,
				MatchFunc:   NameMatch,
			},
			{
				Key:         "medicalRecordNumber",
				MatchProb:   0.95,
				UnmatchProb: 0.05,
			},
			{
				Key:         "ssn",
				MatchProb:   0.98,
				UnmatchProb: 0.01,
			},
			{
				Key:         "visitNumber",
				MatchProb:   0.95,
				UnmatchProb: 0.05,
			},
		},
	}

	for _, tt := range []struct {
		test    string
		want    float64
		matcher Deterministic
		a       ComparablePatient
		b       ComparablePatient
	}{
		{test: "no_match", want: 0, a: katy1, b: empty},
		{test: "completely_same", want: 20.520783771944544, a: katy1, b: katy1},
		{test: "completely_different", want: -17.864128900840395, a: katy1, b: lukas1},
		{test: "katy", want: 7.849713757604871, a: katy1, b: katy2},
		{test: "twins", want: 1.4087672169719516, a: brian, b: julian},
		{test: "same_ssn", want: -4.98273959792274, a: jerry, b: lukas1},
		{test: "different_ssn", want: 12.033843623699326, a: lukas2, b: lukas1},
	} {
		t.Run(tt.test, func(t *testing.T) {
			assertMatch(tt.want, &m, &tt.a, &tt.b)
		})
	}
}
