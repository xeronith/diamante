package search_test

import (
	"testing"

	"github.com/xeronith/diamante/utility/search"
)

func TestMatchAny(test *testing.T) {
	type arguments struct {
		input    string
		criteria string
	}

	testCases := []struct {
		name        string
		expectation bool
		arguments   arguments
	}{
		{
			"Case1",
			false,
			arguments{
				input:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
				criteria: "not found",
			},
		},
		{
			"Case2",
			true,
			arguments{
				input:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
				criteria: "iPs am   ",
			},
		},
		{
			"Case3",
			true,
			arguments{
				input:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
				criteria: " sEC not found",
			},
		},
		{
			"Case4",
			false,
			arguments{
				input:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
				criteria: "ip_sum   nope",
			},
		},
		{
			"Case4",
			true,
			arguments{
				input:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
				criteria: "",
			},
		},
	}

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			if result := search.MatchAny(testCase.arguments.input, testCase.arguments.criteria); result != testCase.expectation {
				test.Errorf("MatchAny() = %v, expected %v", result, testCase.expectation)
			}
		})
	}
}

func BenchmarkMatchAny(benchmark *testing.B) {
	for i := 0; i < benchmark.N; i++ {
		if !search.MatchAny("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Lorem ipsum dolor sit amet, consectetur adipiscing elit. target", "looking for the target") {
			benchmark.FailNow()
		}
	}
}
