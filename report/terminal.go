package report

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/logrusorgru/aurora"
	"github.com/dwillist/spec"
)

// Terminal reports specs via stdout.
type Terminal struct{}

func (Terminal) Start(_ *testing.T, plan spec.Plan) {
	fmt.Println("Suite:", plan.Text)
	fmt.Printf("Total: %d | Focused: %d | Pending: %d\n", plan.Total, plan.Focused, plan.Pending)
	if plan.HasRandom {
		fmt.Println("Random seed:", plan.Seed)
	}
	if plan.HasFocus {
		fmt.Println("Focus is active.")
	}
}

func (Terminal) Specs(_ *testing.T, specs <-chan spec.Spec) {
	var passed, failed, skipped int
	for s := range specs {
		switch {
		case s.Failed:
			failed++
			if !testing.Verbose() {
				fmt.Print("x")
			} else {
				if out, err := ioutil.ReadAll(s.Out); err == nil {
					fmt.Println(">" + string(out)  + "<")
					//fmt.Println(">" + err.Error()  + "<")
					fmt.Print(aurora.Red(fmt.Sprintf("%s", out)))
					//fmt.Print(aurora.Red(fmt.Sprintf("Failed: %d | ", failed)))
				}
			}
		case s.Skipped:
			skipped++
			if !testing.Verbose() {
				fmt.Print("s")
			}
		default:
			passed++
			if !testing.Verbose() {
				fmt.Print(".")
			}
		}
	}
	fmt.Println()
	fmt.Print(aurora.Green(fmt.Sprintf("Passed: %d | ", passed)))
	fmt.Print(aurora.Red(fmt.Sprintf("Failed: %d | ", failed)))
	fmt.Print(aurora.Blue(fmt.Sprintf("Skipped: %d ", skipped)))
	fmt.Println()
	//fmt.Printf("\nPassed: %d | Failed: %d | Skipped: %d\n\n", passed, failed, skipped)
}
