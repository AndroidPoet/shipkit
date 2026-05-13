package guide

import (
	"bytes"
	"strings"
	"testing"
)

func Test_Run_printsSetupPlanFromAnswers(t *testing.T) {
	input := strings.NewReader("LaunchKit\nboth\nyes\ngithub\n")
	var output bytes.Buffer

	answers, err := Run(input, &output)

	if err != nil {
		t.Fatal(err)
	}
	if answers.AppName != "LaunchKit" {
		t.Fatalf("AppName = %q", answers.AppName)
	}
	for _, want := range []string{
		`shipkit init "LaunchKit"`,
		"gpc setup --auto",
		"rc login",
		"asc auth login",
		"shipkit release all",
	} {
		if !strings.Contains(output.String(), want) {
			t.Fatalf("output missing %q:\n%s", want, output.String())
		}
	}
}

func Test_Run_usesDefaultsForBlankAnswers(t *testing.T) {
	input := strings.NewReader("\n\n\n\n")
	var output bytes.Buffer

	answers, err := Run(input, &output)

	if err != nil {
		t.Fatal(err)
	}
	if answers.AppName != "My App" || answers.Platform != "both" || !answers.RevenueCat || answers.CIProvider != "github" {
		t.Fatalf("answers = %#v", answers)
	}
}
