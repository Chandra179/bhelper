package collision

import (
	"bhelper/feature"
	"fmt"
	"math"
	"math/big"
	"time"
)

type MathResult struct {
	TotalSpace         uint64
	TotalIDs           int64
	Probability        *big.Float
	ExpectedCollisions int64
	TimeToCollision    *TimeResult
}

type TimeResult struct {
	P50  time.Duration
	P01  time.Duration
	P001 time.Duration
}

type SimResult struct {
	Collisions  int
	Iterations  int
	Probability float64
}

type CollisionAnalyzer struct {
	registry *GeneratorRegistry
}

func NewCollisionAnalyzer() *CollisionAnalyzer {
	reg := NewGeneratorRegistry()
	gen64, _ := NewBase64Generator(8)
	gen62, _ := NewBase62Generator(10)
	genSnow, _ := NewSnowflakeGenerator()

	reg.Register(gen64)
	reg.Register(gen62)
	reg.Register(genSnow)

	return &CollisionAnalyzer{
		registry: reg,
	}
}

func (c *CollisionAnalyzer) ID() string {
	return "collision"
}

func (c *CollisionAnalyzer) Name() string {
	return "Collision Analyzer"
}

func (c *CollisionAnalyzer) Description() string {
	return "Analyze collision probability for ID generation systems"
}

func (c *CollisionAnalyzer) Help() string {
	return `Analyzes the probability of ID collisions for various generation schemes.

Input format: format:length:rate/unit

Rate units: sec (second), min (minute), ms (millisecond), ns (nanosecond)

Examples:
  base64:10:1000/sec     - Base64 IDs, length 10, 1000/sec
  base62:8:500/min       - Base62 IDs, length 8, 500/min
  snowflake:0:10000/ms    - Snowflake IDs, 10000/ms

The analysis includes:
- Mathematical calculation using birthday paradox
- Actual simulation with generated IDs
- Time to collision at 50%, 1%, and 0.1% probabilities
- Comparison between theoretical and empirical results`
}

func (c *CollisionAnalyzer) Examples() []feature.Example {
	return []feature.Example{
		{
			Input:       "base64:10:1000/sec",
			Description: "Analyze 10-character Base64 IDs at 1000/sec",
		},
		{
			Input:       "base62:8:500/min",
			Description: "Analyze 8-character Base62 IDs at 500/min",
		},
	}
}

func (c *CollisionAnalyzer) Execute(input string) (string, error) {
	config, err := ParseInput(input)
	if err != nil {
		return "", err
	}

	var gen IDGenerator
	switch config.Format {
	case "base64":
		gen, err = NewBase64Generator(config.Length)
	case "base62":
		gen, err = NewBase62Generator(config.Length)
	case "snowflake":
		gen, err = NewSnowflakeGenerator()
	default:
		return "", fmt.Errorf("unknown generator: %s", config.Format)
	}

	if err != nil {
		return "", fmt.Errorf("failed to create generator: %v", err)
	}

	var ratePerSec int64
	switch config.RateUnit {
	case "sec":
		ratePerSec = config.Rate
	case "min":
		ratePerSec = config.Rate / 60
	case "ms":
		ratePerSec = config.Rate * 1000
	case "ns":
		ratePerSec = config.Rate * 1000000000
	default:
		return "", fmt.Errorf("unsupported rate unit: %s", config.RateUnit)
	}

	totalIDs := ratePerSec

	mathResult, err := CalculateProbability(totalIDs, gen.TotalSpace())
	if err != nil {
		return "", fmt.Errorf("calculation error: %v", err)
	}

	simIterations := 1000000
	simResult, err := SimulateCollisions(gen, simIterations)
	if err != nil {
		return "", fmt.Errorf("simulation error: %v", err)
	}

	return FormatResult(config.Format, config.Length, ratePerSec, mathResult, simResult), nil
}

func CalculateProbability(n int64, N uint64) (*MathResult, error) {
	if n <= 0 || N == 0 {
		return nil, fmt.Errorf("invalid inputs: n=%d, N=%d", n, N)
	}

	nFloat := big.NewFloat(float64(n))
	nSquared := new(big.Float).Mul(nFloat, nFloat)

	twoN := big.NewFloat(2.0)
	NFloat := big.NewFloat(float64(N))
	twoN_N := new(big.Float).Mul(twoN, NFloat)

	exponent := new(big.Float).Quo(nSquared, twoN_N)
	exponent.Neg(exponent)

	e_pow_x := new(big.Float).SetPrec(100)
	expFloat, _ := exponent.Float64()
	e_pow_x.SetFloat64(math.Exp(expFloat))

	probability := new(big.Float).Sub(big.NewFloat(1.0), e_pow_x)

	expectedCollisions := new(big.Float).Mul(probability, nFloat)
	expectedCollisionsInt, _ := expectedCollisions.Int(nil)

	return &MathResult{
		TotalSpace:         N,
		TotalIDs:           n,
		Probability:        probability,
		ExpectedCollisions: expectedCollisionsInt.Int64(),
		TimeToCollision:    CalculateTimeToCollision(N, 1),
	}, nil
}

func CalculateTimeToCollision(N uint64, rate int64) *TimeResult {
	NFloat := float64(N)
	rateFloat := float64(rate)

	p50 := math.Sqrt(-2*NFloat*math.Log(0.5)) / rateFloat
	p01 := math.Sqrt(-2*NFloat*math.Log(0.99)) / rateFloat
	p001 := math.Sqrt(-2*NFloat*math.Log(0.999)) / rateFloat

	return &TimeResult{
		P50:  time.Duration(p50 * float64(time.Second)),
		P01:  time.Duration(p01 * float64(time.Second)),
		P001: time.Duration(p001 * float64(time.Second)),
	}
}

func SimulateCollisions(gen IDGenerator, maxIterations int) (*SimResult, error) {
	if maxIterations <= 0 {
		return nil, fmt.Errorf("maxIterations must be positive")
	}

	seen := make(map[string]bool)
	collisions := 0

	for i := 0; i < maxIterations; i++ {
		id := gen.Generate()
		if seen[id] {
			collisions++
		} else {
			seen[id] = true
		}
	}

	probability := float64(collisions) / float64(maxIterations)

	return &SimResult{
		Collisions:  collisions,
		Iterations:  maxIterations,
		Probability: probability,
	}, nil
}
