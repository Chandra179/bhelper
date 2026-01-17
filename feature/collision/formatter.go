package collision

import (
	"fmt"
	"math/big"
	"strings"
	"time"
)

func FormatResult(format string, length int, rate int64, mathResult *MathResult, simResult *SimResult) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Collision Analysis: %s length %d\n\n", format, length))

	sb.WriteString("Mathematical Results:\n")
	sb.WriteString(fmt.Sprintf("  - Total ID Space: %s\n", formatNumber(mathResult.TotalSpace)))
	sb.WriteString(fmt.Sprintf("  - Generation Rate: %d/sec\n", rate))
	sb.WriteString(fmt.Sprintf("  - Collision Probability (1 sec): %s\n", formatProbability(mathResult.Probability)))
	sb.WriteString(fmt.Sprintf("  - Expected Collisions (1 sec): %d\n\n", mathResult.ExpectedCollisions))

	sb.WriteString("  Time to Collision:\n")
	sb.WriteString(fmt.Sprintf("  - 50%% probability: %s\n", formatDuration(mathResult.TimeToCollision.P50)))
	sb.WriteString(fmt.Sprintf("  - 1%% probability: %s\n", formatDuration(mathResult.TimeToCollision.P01)))
	sb.WriteString(fmt.Sprintf("  - 0.1%% probability: %s\n\n", formatDuration(mathResult.TimeToCollision.P001)))

	sb.WriteString("Simulation Results:\n")
	sb.WriteString(fmt.Sprintf("  - Collisions Found: %d\n", simResult.Collisions))
	sb.WriteString(fmt.Sprintf("  - Measured Probability: %s (%d in %d)\n",
		formatProbabilityFloat(simResult.Probability),
		int(simResult.Probability*10000), 10000))

	mathProb := mathResult.Probability
	simProb := big.NewFloat(simResult.Probability)
	diff := new(big.Float).Sub(simProb, mathProb)
	sb.WriteString(fmt.Sprintf("  - Difference: %s\n", formatProbability(diff)))

	return sb.String()
}

func formatNumber(n uint64) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	} else if n < 1000000 {
		return fmt.Sprintf("%d,%03d", n/1000, n%1000)
	} else if n < 1000000000 {
		return fmt.Sprintf("%d,%03d,%03d", n/1000000, (n/1000)%1000, n%1000)
	} else {
		return fmt.Sprintf("%d,%03d,%03d,%03d", n/1000000000, (n/1000000)%1000, (n/1000)%1000, n%1000)
	}
}

func formatProbability(p *big.Float) string {
	if p == nil {
		return "N/A"
	}

	percent := new(big.Float).Mul(p, big.NewFloat(100))
	percentStr := fmt.Sprintf("%.4f", percent)

	if percent.Cmp(big.NewFloat(0.0001)) >= 0 {
		oneIn := new(big.Float).Quo(big.NewFloat(100), percent)
		oneInInt, _ := oneIn.Int64()
		return fmt.Sprintf("%s%% (1 in %d)", percentStr, oneInInt)
	}
	return fmt.Sprintf("%s%% (<1 in 10,000)", percentStr)
}

func formatProbabilityFloat(p float64) string {
	percent := p * 100
	if percent >= 0.0001 {
		oneIn := 100.0 / percent
		return fmt.Sprintf("%.4f%% (1 in %d)", percent, int(oneIn))
	}
	return fmt.Sprintf("%.4f%% (<1 in 10,000)", percent)
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.1f seconds", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.1f minutes", d.Minutes())
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%.1f hours", d.Hours())
	} else if d < 365*24*time.Hour {
		return fmt.Sprintf("%.1f days", d.Hours()/24)
	} else {
		return fmt.Sprintf("%.1f years", d.Hours()/24/365)
	}
}
