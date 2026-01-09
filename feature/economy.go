package feature

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// EconomyCalculator performs economic calculations
type EconomyCalculator struct{}

func NewEconomyCalculator() *EconomyCalculator {
	return &EconomyCalculator{}
}

func (ec *EconomyCalculator) ID() string {
	return "economy"
}

func (ec *EconomyCalculator) Name() string {
	return "Economy Calculator"
}

func (ec *EconomyCalculator) Description() string {
	return "Calculate inflation, compound interest, and financial metrics"
}

func (ec *EconomyCalculator) Help() string {
	return `Economy Calculator helps you perform common financial calculations:

Supported calculations:
• Inflation adjustment: "inflation 1000 3.5 10" (amount, rate%, years)
• Compound interest: "compound 5000 4.2 15" (principal, rate%, years)
• Future value: "future 1000 5 20" (amount, rate%, years)

All rates should be provided as percentages (e.g., 3.5 for 3.5%)`
}

func (ec *EconomyCalculator) Examples() []Example {
	return []Example{
		{Input: "inflation 1000 3.5 10", Description: "Calculate $1000 after 10 years of 3.5% inflation"},
		{Input: "compound 5000 4.2 15", Description: "Calculate compound interest on $5000 at 4.2% for 15 years"},
		{Input: "future 1000 5 20", Description: "Calculate future value of $1000 at 5% annual return"},
	}
}

func (ec *EconomyCalculator) Execute(input string) (string, error) {
	parts := strings.Fields(input)
	if len(parts) < 4 {
		return "", fmt.Errorf("format: <command> <amount> <rate%%> <years>")
	}

	command := parts[0]
	amount, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return "", fmt.Errorf("invalid amount: %v", err)
	}

	rate, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return "", fmt.Errorf("invalid rate: %v", err)
	}

	years, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return "", fmt.Errorf("invalid years: %v", err)
	}

	var result float64
	var description string

	switch command {
	case "inflation", "compound", "future":
		result = amount * math.Pow(1+rate/100, years)
		description = fmt.Sprintf("After %.1f years at %.2f%% annual rate", years, rate)
	default:
		return "", fmt.Errorf("unknown command: %s", command)
	}

	output := fmt.Sprintf(`Initial Amount: $%.2f
Annual Rate:    %.2f%%
Time Period:    %.1f years
%s

Final Value:    $%.2f
Total Gain:     $%.2f (%.1f%%)`,
		amount, rate, years, description, result, result-amount, ((result-amount)/amount)*100)

	return output, nil
}
