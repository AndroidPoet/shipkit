package guide

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Answers struct {
	AppName    string
	Platform   string
	RevenueCat bool
	CIProvider string
}

func Run(stdin io.Reader, stdout io.Writer) (Answers, error) {
	reader := bufio.NewReader(stdin)
	fmt.Fprintln(stdout, "Shipkit Guide")
	fmt.Fprintln(stdout, "Answer a few questions and Shipkit will give you the shortest setup path.")
	fmt.Fprintln(stdout)

	appName, err := ask(reader, stdout, "App name", "My App")
	if err != nil {
		return Answers{}, err
	}
	platform, err := askChoice(reader, stdout, "Platforms", []string{"both", "android", "ios"}, "both")
	if err != nil {
		return Answers{}, err
	}
	revenueCatAnswer, err := askChoice(reader, stdout, "Use RevenueCat", []string{"yes", "no"}, "yes")
	if err != nil {
		return Answers{}, err
	}
	ciProvider, err := askChoice(reader, stdout, "CI provider", []string{"github", "local"}, "github")
	if err != nil {
		return Answers{}, err
	}

	answers := Answers{
		AppName:    appName,
		Platform:   platform,
		RevenueCat: revenueCatAnswer == "yes",
		CIProvider: ciProvider,
	}
	PrintPlan(stdout, answers)
	return answers, nil
}

func PrintPlan(stdout io.Writer, answers Answers) {
	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, "Recommended setup")
	fmt.Fprintf(stdout, "  shipkit init %q\n", answers.AppName)
	fmt.Fprintln(stdout, "  shipkit install")
	fmt.Fprintln(stdout, "  shipkit doctor")
	if answers.CIProvider == "github" {
		fmt.Fprintln(stdout, "  shipkit ci github")
	}

	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, "Provider auth")
	if answers.Platform == "both" || answers.Platform == "android" {
		fmt.Fprintln(stdout, "  gpc setup --auto")
	}
	if answers.RevenueCat {
		fmt.Fprintln(stdout, "  rc login")
	}
	if answers.Platform == "both" || answers.Platform == "ios" {
		fmt.Fprintln(stdout, "  asc auth login")
	}

	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, "Release commands")
	if answers.Platform == "both" || answers.Platform == "android" {
		fmt.Fprintln(stdout, "  shipkit release android")
	}
	if answers.Platform == "both" || answers.Platform == "ios" {
		fmt.Fprintln(stdout, "  shipkit release ios")
	}
	if answers.Platform == "both" {
		fmt.Fprintln(stdout, "  shipkit release all")
	}
}

func ask(reader *bufio.Reader, stdout io.Writer, label, fallback string) (string, error) {
	fmt.Fprintf(stdout, "%s [%s]: ", label, fallback)
	value, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback, nil
	}
	return value, nil
}

func askChoice(reader *bufio.Reader, stdout io.Writer, label string, choices []string, fallback string) (string, error) {
	for {
		value, err := ask(reader, stdout, fmt.Sprintf("%s (%s)", label, strings.Join(choices, "/")), fallback)
		if err != nil {
			return "", err
		}
		for _, choice := range choices {
			if strings.EqualFold(value, choice) {
				return choice, nil
			}
		}
		fmt.Fprintf(stdout, "Choose one of: %s\n", strings.Join(choices, ", "))
	}
}
