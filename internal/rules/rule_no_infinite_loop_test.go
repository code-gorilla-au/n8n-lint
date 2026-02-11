package rules

import (
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestRule_no_infinite_loop(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("should do something", func(t *testing.T) {
			// TODO
		}).
		Run()
	odize.AssertNoError(t, err)
}
