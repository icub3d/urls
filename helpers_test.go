package gourl

import (
	"testing"
)

func TestValidID(t *testing.T) {
	tests := []struct {
		id       string
		expected bool
	}{
		// Test an empty string.
		{
			id:       "",
			expected: false,
		},

		// Test one with an invalid character.
		{
			id:       "98Nb-a",
			expected: false,
		},

		// Test a valid id.
		{
			id:       "9A8b1hdkalKAHgkalhg",
			expected: false,
		},
	}

	for k, test := range tests {
		result := ValidID(test.id)
		if result != test.expected {
			t.Errorf("Test %v: expected %v from ValidID(%v), but got %v",
				k, test.expected, test.id, result)
		}
	}
}

func TestDigit(t *testing.T) {
	tests := []struct {
		id       int
		expected string
	}{
		// Test values in int range.
		{
			id:       0,
			expected: "0",
		},
		{
			id:       5,
			expected: "5",
		},
		{
			id:       9,
			expected: "9",
		},

		// Test values in A-Z range.
		{
			id:       10,
			expected: "A",
		},
		{
			id:       25,
			expected: "P",
		},
		{
			id:       35,
			expected: "Z",
		},

		// Test values in a-z range.
		{
			id:       36,
			expected: "a",
		},
		{
			id:       51,
			expected: "p",
		},
		{
			id:       61,
			expected: "z",
		},

		// Test an invalid id.
		{
			id:       62,
			expected: "",
		},
	}

	for k, test := range tests {
		result := digit(test.id)
		if result != test.expected {
			t.Errorf("Test %v: expected %v from digit(%v), but got %v",
				k, test.expected, test.id, result)
		}
	}

}

func TestIntToShort(t *testing.T) {

	tests := []struct {
		id       int
		expected string
	}{
		// Test some known cases.
		{
			id:       -1,
			expected: "0",
		},
		{
			id:       0,
			expected: "0",
		},

		// Test some random cases.
		{
			id:       25883599,
			expected: "1kbVP",
		},
		{
			id:       13992492,
			expected: "wi5M",
		},
		{
			id:       24025617,
			expected: "1co9x",
		},
		{
			id:       10702,
			expected: "2mc",
		},
		{
			id:       56351835,
			expected: "3oRgZ",
		},
		{
			id:       60799965,
			expected: "476qb",
		},
		{
			id:       30313280,
			expected: "23Brs",
		},
		{
			id:       11612025,
			expected: "mioj",
		},
		{
			id:       19568944,
			expected: "1K6m8",
		},
		{
			id:       23213121,
			expected: "1ZOnB",
		},
		{
			id:       36808014,
			expected: "2URRe",
		},
		{
			id:       73173873,
			expected: "4x1rl",
		},
		{
			id:       33140829,
			expected: "2F3RV",
		},
		{
			id:       49448467,
			expected: "3LTnv",
		},
		{
			id:       23215655,
			expected: "1ZPS3",
		},
		{
			id:       31857603,
			expected: "29fcJ",
		},
		{
			id:       83196636,
			expected: "5d5FE",
		},
		{
			id:       64932637,
			expected: "4ORwb",
		},
		{
			id:       62056666,
			expected: "4CNly",
		},
		{
			id:       53483481,
			expected: "3cPUn",
		},
		{
			id:       6655522,
			expected: "RvP8",
		},
		{
			id:       20680798,
			expected: "1Om1G",
		},
		{
			id:       41750143,
			expected: "2pB7P",
		},
		{
			id:       3023948,
			expected: "CgfM",
		},
		{
			id:       10194536,
			expected: "gm40",
		},
		{
			id:       60294934,
			expected: "44zSw",
		},
		{
			id:       1985607,
			expected: "8KXv",
		},
		{
			id:       89591901,
			expected: "63uwf",
		},
		{
			id:       21746119,
			expected: "1TF9r",
		},
		{
			id:       66617413,
			expected: "4VWEP",
		},
		{
			id:       14147405,
			expected: "xMNx",
		},
		{
			id:       50948628,
			expected: "3Rm44",
		},
		{
			id:       67471159,
			expected: "4Z6KV",
		},
		{
			id:       57769279,
			expected: "3uOQZ",
		},
		{
			id:       80925232,
			expected: "5TYLg",
		},
		{
			id:       51234358,
			expected: "3SyOc",
		},
		{
			id:       21762868,
			expected: "1TJW0",
		},
		{
			id:       85775710,
			expected: "5nuBC",
		},
		{
			id:       2068155,
			expected: "8g1L",
		},
		{
			id:       38720804,
			expected: "2cT36",
		},
		{
			id:       40975171,
			expected: "2lvVr",
		},
		{
			id:       75084343,
			expected: "552rn",
		},
		{
			id:       69963207,
			expected: "4jYcp",
		},
		{
			id:       68132074,
			expected: "4bsGQ",
		},
		{
			id:       77810961,
			expected: "5GUBV",
		},
		{
			id:       49361381,
			expected: "3L79J",
		},
		{
			id:       68236415,
			expected: "4cJPL",
		},
		{
			id:       67175877,
			expected: "4XrVt",
		},
		{
			id:       65592921,
			expected: "4RDiL",
		},
		{
			id:       41238522,
			expected: "2n21S",
		},
		{
			id:       94283385,
			expected: "6NbPl",
		},
		{
			id:       41956665,
			expected: "2q2qP",
		},
		{
			id:       37707752,
			expected: "2YDVY",
		},
		{
			id:       68610149,
			expected: "4dsdJ",
		},
		{
			id:       28209392,
			expected: "1uMYC",
		},
		{
			id:       38691855,
			expected: "2cLWB",
		},
		{
			id:       99522769,
			expected: "6jaPx",
		},
		{
			id:       69682314,
			expected: "4iNYI",
		},
		{
			id:       16162743,
			expected: "15ofP",
		},
		{
			id:       38386419,
			expected: "2b43n",
		},
		{
			id:       42666958,
			expected: "2t1ck",
		},
		{
			id:       98849355,
			expected: "6glER",
		},
		{
			id:       17320289,
			expected: "1AfnV",
		},
		{
			id:       25808476,
			expected: "1kHxk",
		},
		{
			id:       91304382,
			expected: "6B6RK",
		},
		{
			id:       33144228,
			expected: "2F4KK",
		},
		{
			id:       61813865,
			expected: "4BMbp",
		},
		{
			id:       90135250,
			expected: "66CIM",
		},
		{
			id:       84605551,
			expected: "5izlf",
		},
		{
			id:       75446667,
			expected: "56Z7j",
		},
		{
			id:       35347002,
			expected: "2OJMw",
		},
		{
			id:       84988268,
			expected: "5kbKW",
		},
		{
			id:       344523,
			expected: "1Rcp",
		},
		{
			id:       34542600,
			expected: "2Kw6i",
		},
		{
			id:       11579636,
			expected: "maOK",
		},
		{
			id:       50145891,
			expected: "3OPEh",
		},
		{
			id:       97880514,
			expected: "6chBy",
		},
		{
			id:       44866383,
			expected: "32FnL",
		},
		{
			id:       29717072,
			expected: "20glc",
		},
		{
			id:       39625641,
			expected: "2gGRF",
		},
		{
			id:       9977484,
			expected: "frbA",
		},
		{
			id:       38415131,
			expected: "2bBWt",
		},
		{
			id:       25011862,
			expected: "1gwj8",
		},
		{
			id:       4742753,
			expected: "Jto1",
		},
		{
			id:       30922277,
			expected: "25kIP",
		},
		{
			id:       18561924,
			expected: "1Fsns",
		},
		{
			id:       86782082,
			expected: "5s7z0",
		},
		{
			id:       42781204,
			expected: "2tVLQ",
		},
		{
			id:       21860274,
			expected: "1Tir4",
		},
		{
			id:       82510071,
			expected: "5aCdb",
		},
		{
			id:       37078850,
			expected: "2VZty",
		},
		{
			id:       55022721,
			expected: "3irvF",
		},
		{
			id:       51293216,
			expected: "3TDhw",
		},
		{
			id:       90557044,
			expected: "67y1U",
		},
		{
			id:       94867994,
			expected: "6Q3Uw",
		},
		{
			id:       26714036,
			expected: "1o5XY",
		},
		{
			id:       11013256,
			expected: "kD3A",
		},
		{
			id:       37169131,
			expected: "2VxO7",
		},
		{
			id:       7988373,
			expected: "XW8j",
		},
		{
			id:       52386136,
			expected: "3Xo1g",
		},
	}

	for k, test := range tests {
		result := IntToShort(test.id)
		if result != test.expected {
			t.Errorf("Test %v: expected %v from IntToShort(%v), but got %v",
				k, test.expected, test.id, result)
		}
	}
}
