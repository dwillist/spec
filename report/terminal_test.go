package report

import (
	"github.com/dwillist/spec"
	"testing"
)

func TestUnitTerminal(t *testing.T) {
	spec.Run(t, "build", testTerminal, spec.Report(Terminal{}))
}

func testTerminal(t *testing.T, when spec.G, it spec.S) {
	when("NewNode", func() {

		it("returns true if a build plan exists", func() {

		})

		it("fails", func() {
			t.Fatal("Failing test")
		})

	})
}
