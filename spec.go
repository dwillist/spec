package spec

import (
	"strings"
	"testing"
)

type G func(string, func(), ...Option)

func (g G) Pend(text string, f func(), _ ...Option) {
	g(text, f, func(c *config) { c.pend = true })
}

func (g G) Focus(text string, f func(), opts ...Option) {
	g(text, f, append(opts, func(c *config) { c.focus = true })...)
}

type S func(string, func(), ...Option)

func (s S) Before(f func()) {
	s("", f, func(c *config) { c.before = true })
}

func (s S) After(f func()) {
	s("", f, func(c *config) { c.after = true })
}

func (s S) Pend(text string, f func(), _ ...Option) {
	s(text, f, func(c *config) { c.pend = true })
}

func (s S) Focus(text string, f func(), opts ...Option) {
	s(text, f, append(opts, func(c *config) { c.focus = true })...)
}

func Run(t *testing.T, f func(*testing.T, G, S), opts ...Option) bool {
	specs, focused, seed := oldParse(f, opts...)

	t.Logf("Running %d specs.", len(specs))
	if seed != 0 {
		t.Logf("Random seed: %d", seed)
	}

	nodes := parse(f, opts...)

	return nodes.run(t, nil, func(t *testing.T, groups []int, n node) {
		switch {
		case n.pend, focused && !n.focus:
			t.SkipNow()
		case n.order == orderParallel:
			t.Parallel()
		}
		var (
			spec          func()
			before, after []func()
		)
		f(t, func(_ string, f func(), _ ...Option) {
			switch {
			case len(groups) == 0:
			case groups[0] > 0:
				groups[0]--
			default:
				groups = groups[1:]
				f()
			}
		}, func(_ string, f func(), opts ...Option) {
			cfg := options(opts).apply()
			switch {
			case cfg.before:
				before = append(before, f)
			case cfg.after:
				after = append([]func(){f}, after...)
			case spec != nil || len(groups) > 0:
			case n.index > 0:
				n.index--
			default:
				spec = f
			}
		})

		if spec == nil {
			t.Fatal("Failed to parse.")
		}

		run(before...)
		defer run(after...)
		run(spec)
	})
}

func run(fs ...func()) {
	for _, f := range fs {
		f()
	}
}