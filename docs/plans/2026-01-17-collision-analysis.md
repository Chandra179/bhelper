# Collision Analysis Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a pluggable collision analysis feature that calculates probability of ID collisions using both mathematical formulas and actual simulation.

**Architecture:** Pluggable IDGenerator interface with registry-based generator management, dual-engine analysis (math + simulation), integrated into existing Feature system.

**Tech Stack:** Go, existing bubbletea CLI framework, math/big for precision calculations, Go's built-in crypto/rand for random generation.

---

### Task 1: Create IDGenerator Interface

**Files:**
- Create: `feature/collision/generator.go`

**Step 1: Write the failing test**

Create `feature/collision/generator_test.go`:
```go
package collision

import "testing"

func TestIDGeneratorInterface(t *testing.T) {
    var g IDGenerator
    if g == nil {
        t.Fatal("IDGenerator interface exists")
    }
}

func TestGenerate(t *testing.T) {
    gen := NewBase64Generator(8)
    id := gen.Generate()
    if len(id) != 8 {
        t.Errorf("Expected length 8, got %d", len(id))
    }
}

func TestTotalSpace(t *testing.T) {
    gen := NewBase64Generator(8)
    space := gen.TotalSpace()
    expected := uint64(1)
    for i := 0; i < 8; i++ {
        expected *= 64
    }
    if space != expected {
        t.Errorf("Expected %d, got %d", expected, space)
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./feature/collision/ -v`
Expected: FAIL with "undefined: IDGenerator", "undefined: NewBase64Generator"

**Step 3: Write minimal implementation**

Create `feature/collision/generator.go`:
```go
package collision

import (
    "crypto/rand"
    "math/big"
)

type IDGenerator interface {
    Generate() string
    TotalSpace() uint64
    Name() string
}

type Base64Generator struct {
    length int
    chars  []byte
}

func NewBase64Generator(length int) *Base64Generator {
    return &Base64Generator{
        length: length,
        chars:  []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"),
    }
}

func (g *Base64Generator) Generate() string {
    result := make([]byte, g.length)
    max := big.NewInt(int64(len(g.chars)))
    
    for i := 0; i < g.length; i++ {
        n, _ := rand.Int(rand.Reader, max)
        result[i] = g.chars[n.Int64()]
    }
    
    return string(result)
}

func (g *Base64Generator) TotalSpace() uint64 {
    space := uint64(1)
    for i := 0; i < g.length; i++ {
        space *= 64
    }
    return space
}

func (g *Base64Generator) Name() string {
    return "base64"
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./feature/collision/ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add feature/collision/generator.go feature/collision/generator_test.go
git commit -m "feat(collision): add IDGenerator interface and Base64Generator"
```

---

### Task 2: Add Base62 Generator

**Files:**
- Modify: `feature/collision/generator.go`
- Modify: `feature/collision/generator_test.go`

**Step 1: Write the failing test**

Add to `feature/collision/generator_test.go`:
```go
func TestBase62Generator(t *testing.T) {
    gen := NewBase62Generator(10)
    id := gen.Generate()
    if len(id) != 10 {
        t.Errorf("Expected length 10, got %d", len(id))
    }
    
    space := gen.TotalSpace()
    expected := uint64(1)
    for i := 0; i < 10; i++ {
        expected *= 62
    }
    if space != expected {
        t.Errorf("Expected %d, got %d", expected, space)
    }
    
    if gen.Name() != "base62" {
        t.Errorf("Expected 'base62', got '%s'", gen.Name())
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./feature/collision/ -v`
Expected: FAIL with "undefined: NewBase62Generator"

**Step 3: Write minimal implementation**

Add to `feature/collision/generator.go`:
```go
type Base62Generator struct {
    length int
    chars  []byte
}

func NewBase62Generator(length int) *Base62Generator {
    return &Base62Generator{
        length: length,
        chars:  []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"),
    }
}

func (g *Base62Generator) Generate() string {
    result := make([]byte, g.length)
    max := big.NewInt(int64(len(g.chars)))
    
    for i := 0; i < g.length; i++ {
        n, _ := rand.Int(rand.Reader, max)
        result[i] = g.chars[n.Int64()]
    }
    
    return string(result)
}

func (g *Base62Generator) TotalSpace() uint64 {
    space := uint64(1)
    for i := 0; i < g.length; i++ {
        space *= 62
    }
    return space
}

func (g *Base62Generator) Name() string {
    return "base62"
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./feature/collision/ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add feature/collision/generator.go feature/collision/generator_test.go
git commit -m "feat(collision): add Base62Generator"
```

---

### Task 3: Add Snowflake Generator

**Files:**
- Modify: `feature/collision/generator.go`
- Modify: `feature/collision/generator_test.go`

**Step 1: Write the failing test**

Add to `feature/collision/generator_test.go`:
```go
func TestSnowflakeGenerator(t *testing.T) {
    gen := NewSnowflakeGenerator()
    id1 := gen.Generate()
    id2 := gen.Generate()
    
    if id1 == id2 {
        t.Error("Expected different IDs")
    }
    
    space := gen.TotalSpace()
    if space == 0 {
        t.Error("Expected non-zero space")
    }
    
    if gen.Name() != "snowflake" {
        t.Errorf("Expected 'snowflake', got '%s'", gen.Name())
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./feature/collision/ -v`
Expected: FAIL with "undefined: NewSnowflakeGenerator"

**Step 3: Write minimal implementation**

Add to `feature/collision/generator.go`:
```go
import (
    "crypto/rand"
    "math/big"
    "sync"
    "time"
)

type SnowflakeGenerator struct {
    mu      sync.Mutex
    lastMS  int64
    counter int64
}

func NewSnowflakeGenerator() *SnowflakeGenerator {
    return &SnowflakeGenerator{
        lastMS:  time.Now().UnixMilli(),
        counter: 0,
    }
}

func (g *SnowflakeGenerator) Generate() string {
    g.mu.Lock()
    defer g.mu.Unlock()
    
    ms := time.Now().UnixMilli()
    if ms == g.lastMS {
        g.counter++
    } else {
        g.lastMS = ms
        g.counter = 0
    }
    
    // Simplified snowflake: 41-bit timestamp + 22-bit counter
    id := (ms << 22) | (g.counter & 0x3FFFFF)
    return fmt.Sprintf("%d", id)
}

func (g *SnowflakeGenerator) TotalSpace() uint64 {
    // 41-bit timestamp + 22-bit counter = 63 bits total
    return uint64(1) << 63
}

func (g *SnowflakeGenerator) Name() string {
    return "snowflake"
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./feature/collision/ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add feature/collision/generator.go feature/collision/generator_test.go
git commit -m "feat(collision): add SnowflakeGenerator"
```

---

### Task 4: Create Generator Registry

**Files:**
- Create: `feature/collision/registry.go`
- Create: `feature/collision/registry_test.go`

**Step 1: Write the failing test**

Create `feature/collision/registry_test.go`:
```go
package collision

import "testing"

func TestNewRegistry(t *testing.T) {
    reg := NewGeneratorRegistry()
    if reg == nil {
        t.Fatal("Expected non-nil registry")
    }
}

func TestRegistryRegister(t *testing.T) {
    reg := NewGeneratorRegistry()
    gen := NewBase64Generator(8)
    
    reg.Register(gen)
    
    got, ok := reg.Get("base64")
    if !ok {
        t.Error("Expected generator to be registered")
    }
    if got.Name() != "base64" {
        t.Errorf("Expected 'base64', got '%s'", got.Name())
    }
}

func TestRegistryGetNotFound(t *testing.T) {
    reg := NewGeneratorRegistry()
    
    _, ok := reg.Get("unknown")
    if ok {
        t.Error("Expected false for unknown generator")
    }
}

func TestRegistryList(t *testing.T) {
    reg := NewGeneratorRegistry()
    reg.Register(NewBase64Generator(8))
    reg.Register(NewBase62Generator(10))
    
    list := reg.List()
    if len(list) != 2 {
        t.Errorf("Expected 2 generators, got %d", len(list))
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./feature/collision/ -v`
Expected: FAIL with "undefined: NewGeneratorRegistry"

**Step 3: Write minimal implementation**

Create `feature/collision/registry.go`:
```go
package collision

type GeneratorRegistry struct {
    generators map[string]IDGenerator
    order      []string
}

func NewGeneratorRegistry() *GeneratorRegistry {
    return &GeneratorRegistry{
        generators: make(map[string]IDGenerator),
        order:      make([]string, 0),
    }
}

func (r *GeneratorRegistry) Register(gen IDGenerator) {
    name := gen.Name()
    r.generators[name] = gen
    r.order = append(r.order, name)
}

func (r *GeneratorRegistry) Get(name string) (IDGenerator, bool) {
    gen, ok := r.generators[name]
    return gen, ok
}

func (r *GeneratorRegistry) List() []IDGenerator {
    result := make([]IDGenerator, 0, len(r.order))
    for _, name := range r.order {
        result = append(result, r.generators[name])
    }
    return result
}

func (r *GeneratorRegistry) Names() []string {
    return r.order
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./feature/collision/ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add feature/collision/registry.go feature/collision/registry_test.go
git commit -m "feat(collision): add GeneratorRegistry"
```

---

### Task 5: Create Input Parser

**Files:**
- Create: `feature/collision/parser.go`
- Create: `feature/collision/parser_test.go`

**Step 1: Write the failing test**

Create `feature/collision/parser_test.go`:
```go
package collision

import "testing"

func TestParseInput(t *testing.T) {
    input := "base64:10:1000/sec:1day"
    
    config, err := ParseInput(input)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    if config.Format != "base64" {
        t.Errorf("Expected 'base64', got '%s'", config.Format)
    }
    if config.Length != 10 {
        t.Errorf("Expected 10, got %d", config.Length)
    }
    if config.Rate != 1000 {
        t.Errorf("Expected 1000, got %d", config.Rate)
    }
    if config.Duration != time.Hour*24 {
        t.Errorf("Expected 24 hours, got %v", config.Duration)
    }
}

func TestParseInputWithMinute(t *testing.T) {
    input := "base62:8:500/min:2hour"
    
    config, err := ParseInput(input)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    if config.Rate != 500 {
        t.Errorf("Expected 500, got %d", config.Rate)
    }
    if config.Duration != time.Hour*2 {
        t.Errorf("Expected 2 hours, got %v", config.Duration)
    }
}

func TestParseInputError(t *testing.T) {
    input := "invalid"
    
    _, err := ParseInput(input)
    if err == nil {
        t.Error("Expected error for invalid input")
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./feature/collision/ -v`
Expected: FAIL with "undefined: ParseInput", "undefined: Config"

**Step 3: Write minimal implementation**

Create `feature/collision/parser.go`:
```go
package collision

import (
    "errors"
    "fmt"
    "strconv"
    "strings"
    "time"
)

type Config struct {
    Format   string
    Length   int
    Rate     int64
    RateUnit string
    Duration time.Duration
}

func ParseInput(input string) (*Config, error) {
    parts := strings.Split(input, ":")
    if len(parts) != 4 {
        return nil, errors.New("invalid format: expected 'format:length:rate/unit:duration'")
    }
    
    format := parts[0]
    length, err := strconv.Atoi(parts[1])
    if err != nil {
        return nil, fmt.Errorf("invalid length: %v", err)
    }
    
    rateParts := strings.Split(parts[2], "/")
    if len(rateParts) != 2 {
        return nil, errors.New("invalid rate format: expected 'rate/unit'")
    }
    
    rate, err := strconv.ParseInt(rateParts[0], 10, 64)
    if err != nil {
        return nil, fmt.Errorf("invalid rate: %v", err)
    }
    rateUnit := rateParts[1]
    
    duration, err := parseDuration(parts[3])
    if err != nil {
        return nil, fmt.Errorf("invalid duration: %v", err)
    }
    
    return &Config{
        Format:   format,
        Length:   length,
        Rate:     rate,
        RateUnit: rateUnit,
        Duration: duration,
    }, nil
}

func parseDuration(s string) (time.Duration, error) {
    if len(s) < 2 {
        return 0, errors.New("duration too short")
    }
    
    numStr := s[:len(s)-1]
    unit := s[len(s)-1:]
    
    num, err := strconv.Atoi(numStr)
    if err != nil {
        return 0, fmt.Errorf("invalid duration number: %v", err)
    }
    
    switch unit {
    case "s":
        return time.Duration(num) * time.Second, nil
    case "m":
        return time.Duration(num) * time.Minute, nil
    case "h":
        return time.Duration(num) * time.Hour, nil
    case "d":
        return time.Duration(num) * 24 * time.Hour, nil
    case "y":
        return time.Duration(num) * 365 * 24 * time.Hour, nil
    default:
        return 0, fmt.Errorf("unknown duration unit: %s", unit)
    }
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./feature/collision/ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add feature/collision/parser.go feature/collision/parser_test.go
git commit -m "feat(collision): add input parser"
```

---

### Task 6: Create Mathematical Calculator

**Files:**
- Create: `feature/collision/analyzer.go`
- Create: `feature/collision/analyzer_test.go`

**Step 1: Write the failing test**

Create `feature/collision/analyzer_test.go`:
```go
package collision

import (
    "math/big"
    "testing"
    "time"
)

func TestCalculateProbability(t *testing.T) {
    // Simple case: 10 IDs from space of 100
    result, err := CalculateProbability(10, 100)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    if result.Probability == nil {
        t.Fatal("Expected non-nil probability")
    }
    
    // Probability should be small but not zero
    zero := big.NewFloat(0)
    if result.Probability.Cmp(zero) <= 0 {
        t.Error("Expected positive probability")
    }
}

func TestCalculateTimeToCollision(t *testing.T) {
    result := CalculateTimeToCollision(100, 1)
    
    if result.P50 == 0 {
        t.Error("Expected non-zero P50")
    }
    if result.P01 == 0 {
        t.Error("Expected non-zero P01")
    }
    if result.P001 == 0 {
        t.Error("Expected non-zero P001")
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./feature/collision/ -v`
Expected: FAIL with "undefined: CalculateProbability", "undefined: CalculateTimeToCollision"

**Step 3: Write minimal implementation**

Create `feature/collision/analyzer.go`:
```go
package collision

import (
    "fmt"
    "math"
    "math/big"
)

type MathResult struct {
    TotalSpace      uint64
    TotalIDs        int64
    Probability     *big.Float
    ExpectedCollisions int64
    TimeToCollision *TimeResult
}

type TimeResult struct {
    P50   time.Duration
    P01   time.Duration
    P001  time.Duration
}

func CalculateProbability(n int64, N uint64) (*MathResult, error) {
    if n <= 0 || N == 0 {
        return nil, fmt.Errorf("invalid inputs: n=%d, N=%d", n, N)
    }
    
    // P(collision) = 1 - e^(-n²/(2N))
    nFloat := big.NewFloat(float64(n))
    nSquared := new(big.Float).Mul(nFloat, nFloat)
    
    twoN := big.NewFloat(2.0)
    NFloat := big.NewFloat(float64(N))
    twoN_N := new(big.Float).Mul(twoN, NFloat)
    
    exponent := new(big.Float).Quo(nSquared, twoN_N)
    exponent.Neg(exponent)
    
    e_pow_x := new(big.Float).SetPrec(100)
    e_pow_x.SetFloat64(math.Exp(exponent.Float64()))
    
    probability := new(big.Float).Sub(big.NewFloat(1.0), e_pow_x)
    
    expectedCollisions := new(big.Float).Mul(probability, nFloat)
    expectedCollisionsInt, _ := expectedCollisions.Int(nil)
    
    return &MathResult{
        TotalSpace:        N,
        TotalIDs:          n,
        Probability:       probability,
        ExpectedCollisions: expectedCollisionsInt.Int64(),
        TimeToCollision:   CalculateTimeToCollision(N, 1),
    }, nil
}

func CalculateTimeToCollision(N uint64, rate int64) *TimeResult {
    NFloat := float64(N)
    rateFloat := float64(rate)
    
    // t = sqrt(-2 * N * ln(1-P)) / rate
    
    p50 := math.Sqrt(-2 * NFloat * math.Log(0.5)) / rateFloat
    p01 := math.Sqrt(-2 * NFloat * math.Log(0.99)) / rateFloat
    p001 := math.Sqrt(-2 * NFloat * math.Log(0.999)) / rateFloat
    
    return &TimeResult{
        P50:  time.Duration(p50 * float64(time.Second)),
        P01:  time.Duration(p01 * float64(time.Second)),
        P001: time.Duration(p001 * float64(time.Second)),
    }
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./feature/collision/ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add feature/collision/analyzer.go feature/collision/analyzer_test.go
git commit -m "feat(collision): add mathematical calculator"
```

---

### Task 7: Create Simulator

**Files:**
- Modify: `feature/collision/analyzer.go`
- Modify: `feature/collision/analyzer_test.go`

**Step 1: Write the failing test**

Add to `feature/collision/analyzer_test.go`:
```go
func TestSimulateCollisions(t *testing.T) {
    gen := NewBase64Generator(6)
    iterations := 10000
    
    result, err := SimulateCollisions(gen, iterations)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    if result.Collisions < 0 {
        t.Errorf("Expected non-negative collisions, got %d", result.Collisions)
    }
    
    if result.Iterations != iterations {
        t.Errorf("Expected %d iterations, got %d", iterations, result.Iterations)
    }
}

func TestSimulateCollisionsWithSmallSpace(t *testing.T) {
    gen := NewBase64Generator(1) // Very small space
    iterations := 1000
    
    result, err := SimulateCollisions(gen, iterations)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    // Should have many collisions with small space
    if result.Collisions == 0 {
        t.Error("Expected some collisions with small space")
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./feature/collision/ -v`
Expected: FAIL with "undefined: SimulateCollisions"

**Step 3: Write minimal implementation**

Add to `feature/collision/analyzer.go`:
```go
type SimResult struct {
    Collisions int
    Iterations int
    Probability float64
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
        Collisions: collisions,
        Iterations: maxIterations,
        Probability: probability,
    }, nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./feature/collision/ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add feature/collision/analyzer.go feature/collision/analyzer_test.go
git commit -m "feat(collision): add collision simulator"
```

---

### Task 8: Create Output Formatter

**Files:**
- Create: `feature/collision/formatter.go`
- Create: `feature/collision/formatter_test.go`

**Step 1: Write the failing test**

Create `feature/collision/formatter_test.go`:
```go
package collision

import (
    "strings"
    "testing"
    "time"
)

func TestFormatResult(t *testing.T) {
    mathResult := &MathResult{
        TotalSpace: 1000,
        TotalIDs: 100,
        Probability: big.NewFloat(0.005),
        ExpectedCollisions: 1,
        TimeToCollision: &TimeResult{
            P50: time.Hour,
            P01: time.Minute,
            P001: time.Second,
        },
    }
    
    simResult := &SimResult{
        Collisions: 2,
        Iterations: 1000,
        Probability: 0.002,
    }
    
    output := FormatResult("base64", 8, 1000, time.Hour, mathResult, simResult)
    
    if output == "" {
        t.Error("Expected non-empty output")
    }
    
    if !strings.Contains(output, "Collision Analysis") {
        t.Error("Expected 'Collision Analysis' in output")
    }
    
    if !strings.Contains(output, "Mathematical") {
        t.Error("Expected 'Mathematical' in output")
    }
    
    if !strings.Contains(output, "Simulation") {
        t.Error("Expected 'Simulation' in output")
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./feature/collision/ -v`
Expected: FAIL with "undefined: FormatResult"

**Step 3: Write minimal implementation**

Create `feature/collision/formatter.go`:
```go
package collision

import (
    "fmt"
    "math/big"
    "strings"
)

func FormatResult(format string, length int, rate int64, duration time.Duration, mathResult *MathResult, simResult *SimResult) string {
    var sb strings.Builder
    
    // Header
    sb.WriteString(fmt.Sprintf("📊 Collision Analysis: %s length %d\n\n", format, length))
    
    // Mathematical Results
    sb.WriteString("📈 Mathematical Results:\n")
    sb.WriteString(fmt.Sprintf("  - Total ID Space: %s\n", formatNumber(uint64(mathResult.TotalSpace))))
    sb.WriteString(fmt.Sprintf("  - Generation Rate: %d/sec\n", rate))
    sb.WriteString(fmt.Sprintf("  - Duration: %s\n", duration))
    sb.WriteString(fmt.Sprintf("  - Total IDs: %s\n", formatNumber(uint64(mathResult.TotalIDs))))
    sb.WriteString(fmt.Sprintf("  - Collision Probability: %s\n", formatProbability(mathResult.Probability)))
    sb.WriteString(fmt.Sprintf("  - Expected Collisions: %d\n\n", mathResult.ExpectedCollisions))
    
    // Time to Collision
    sb.WriteString("  Time to Collision:\n")
    sb.WriteString(fmt.Sprintf("  - 50%% probability: %s\n", formatDuration(mathResult.TimeToCollision.P50)))
    sb.WriteString(fmt.Sprintf("  - 1%% probability: %s\n", formatDuration(mathResult.TimeToCollision.P01)))
    sb.WriteString(fmt.Sprintf("  - 0.1%% probability: %s\n\n", formatDuration(mathResult.TimeToCollision.P001)))
    
    // Simulation Results
    sb.WriteString("🎯 Simulation Results:\n")
    sb.WriteString(fmt.Sprintf("  - Collisions Found: %d\n", simResult.Collisions))
    sb.WriteString(fmt.Sprintf("  - Measured Probability: %s (%d in %d)\n", 
        formatProbabilityFloat(simResult.Probability), 
        int(simResult.Probability*10000), 10000))
    
    // Difference
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
    
    // Calculate "1 in X"
    if percent.Cmp(big.NewFloat(0.0001)) >= 0 {
        oneIn := new(big.Float).Quo(big.NewFloat(100), percent)
        return fmt.Sprintf("%s%% (1 in %d)", percentStr, int(oneIn.Int64()))
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
```

**Step 4: Run test to verify it passes**

Run: `go test ./feature/collision/ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add feature/collision/formatter.go feature/collision/formatter_test.go
git commit -m "feat(collision): add output formatter"
```

---

### Task 9: Create Main CollisionAnalyzer Feature

**Files:**
- Create: `feature/collision.go`
- Create: `feature/collision_test.go`

**Step 1: Write the failing test**

Create `feature/collision_test.go`:
```go
package feature

import "testing"

func TestNewCollisionAnalyzer(t *testing.T) {
    analyzer := NewCollisionAnalyzer()
    
    if analyzer == nil {
        t.Fatal("Expected non-nil analyzer")
    }
    
    if analyzer.ID() != "collision" {
        t.Errorf("Expected ID 'collision', got '%s'", analyzer.ID())
    }
}

func TestCollisionAnalyzerExecute(t *testing.T) {
    analyzer := NewCollisionAnalyzer()
    
    result, err := analyzer.Execute("base64:10:1000/sec:1day")
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    if result == "" {
        t.Error("Expected non-empty result")
    }
}

func TestCollisionAnalyzerExecuteError(t *testing.T) {
    analyzer := NewCollisionAnalyzer()
    
    _, err := analyzer.Execute("invalid")
    if err == nil {
        t.Error("Expected error for invalid input")
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./feature/ -v`
Expected: FAIL with "undefined: NewCollisionAnalyzer"

**Step 3: Write minimal implementation**

Create `feature/collision.go`:
```go
package feature

import (
    "bhelper/feature/collision"
)

type CollisionAnalyzer struct {
    registry *collision.GeneratorRegistry
}

func NewCollisionAnalyzer() *CollisionAnalyzer {
    reg := collision.NewGeneratorRegistry()
    reg.Register(collision.NewBase64Generator(0))
    reg.Register(collision.NewBase62Generator(0))
    reg.Register(collision.NewSnowflakeGenerator())
    
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

Input format: format:length:rate/unit:duration

Examples:
  base64:10:1000/sec:1day    - Base64 IDs, length 10, 1000/sec for 1 day
  base62:8:500/min:2hour      - Base62 IDs, length 8, 500/min for 2 hours
  snowflake:0:10000/ms:1year  - Snowflake IDs, 10000/ms for 1 year

The analysis includes:
- Mathematical calculation using birthday paradox
- Actual simulation with generated IDs
- Time to collision at 50%, 1%, and 0.1% probabilities
- Comparison between theoretical and empirical results`
}

func (c *CollisionAnalyzer) Examples() []Example {
    return []Example{
        {
            Input:       "base64:10:1000/sec:1day",
            Description: "Analyze 10-character Base64 IDs at 1000/sec for 1 day",
        },
        {
            Input:       "base62:8:500/min:2hour",
            Description: "Analyze 8-character Base62 IDs at 500/min for 2 hours",
        },
    }
}

func (c *CollisionAnalyzer) Execute(input string) (string, error) {
    config, err := collision.ParseInput(input)
    if err != nil {
        return "", err
    }
    
    gen, ok := c.registry.Get(config.Format)
    if !ok {
        return "", fmt.Errorf("unknown generator: %s", config.Format)
    }
    
    // Create generator with correct length
    switch config.Format {
    case "base64":
        gen = collision.NewBase64Generator(config.Length)
    case "base62":
        gen = collision.NewBase62Generator(config.Length)
    case "snowflake":
        gen = collision.NewSnowflakeGenerator()
    }
    
    // Calculate total IDs
    var ratePerSec int64
    switch config.RateUnit {
    case "sec":
        ratePerSec = config.Rate
    case "min":
        ratePerSec = config.Rate / 60
    case "ms":
        ratePerSec = config.Rate * 1000
    default:
        return "", fmt.Errorf("unsupported rate unit: %s", config.RateUnit)
    }
    
    totalIDs := int64(config.Duration.Seconds()) * ratePerSec
    
    // Mathematical analysis
    mathResult, err := collision.CalculateProbability(totalIDs, gen.TotalSpace())
    if err != nil {
        return "", err
    }
    
    // Simulation (limited iterations for performance)
    simIterations := 1000000
    simResult, err := collision.SimulateCollisions(gen, simIterations)
    if err != nil {
        return "", err
    }
    
    return collision.FormatResult(config.Format, config.Length, ratePerSec, config.Duration, mathResult, simResult), nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./feature/ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add feature/collision.go feature/collision_test.go
git commit -m "feat(collision): add main CollisionAnalyzer feature"
```

---

### Task 10: Integrate into Main CLI

**Files:**
- Modify: `main.go`

**Step 1: Verify integration needed**

Check current main.go has feature registration pattern.

**Step 2: Register CollisionAnalyzer feature**

Modify `main.go`:
```go
func main() {
    // Register all features
    registry := feature.NewFeatureRegistry()
    registry.Register(feature.NewCharacterAnalyzer())
    registry.Register(feature.NewTimezoneAnalyzer())
    registry.Register(feature.NewCollisionAnalyzer())  // Add this line
    // registry.Register(NewWeatherForecast())
    // ... register 100 features here

    // Start CLI with all registered features
    p := tea.NewProgram(NewCLI(registry))
    if _, err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

**Step 3: Test integration**

Run: `go build && ./bhelper`

Verify "Collision Analyzer" appears in the feature list when you run the CLI.

**Step 4: Commit**

```bash
git add main.go
git commit -m "feat: integrate CollisionAnalyzer into main CLI"
```

---

### Task 11: Run All Tests and Verify

**Step 1: Run all tests**

Run: `go test ./... -v`

Expected: All tests pass

**Step 2: Build and test manually**

Run: `go build && ./bhelper`

Navigate to "Collision Analyzer" and test with:
- `base64:10:1000/sec:1day`
- `base62:8:500/min:2hour`
- `snowflake:0:10000/sec:1year`

**Step 3: Commit if everything works**

```bash
git add .
git commit -m "test: verify collision analyzer integration"
```

---

## Testing Strategy

**Unit Tests:**
- Each generator tested individually
- Parser tested with valid and invalid inputs
- Math calculator verified with known values
- Simulator tested with small spaces for predictable collisions
- Formatter tested for output format

**Integration Tests:**
- End-to-end feature execution
- Error handling for edge cases
- Performance verification (simulation doesn't hang)

**Manual Testing:**
- Test all three generators with various configurations
- Verify output readability and accuracy
- Test error messages for invalid inputs

---

## Notes for Implementation

1. **Precision:** Use math/big for very small probabilities to avoid underflow
2. **Performance:** Limit simulation iterations to prevent long-running operations
3. **Extensibility:** Adding new generators only requires implementing IDGenerator interface and registering it
4. **Error Messages:** Provide clear, actionable error messages with examples
5. **Documentation:** Help text should show example inputs for all generators

---

## Checklist for Completion

- [ ] All unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing successful
- [ ] All generators work correctly
- [ ] Error handling comprehensive
- [ ] Output is readable and accurate
- [ ] Performance acceptable (simulation completes in reasonable time)
- [ ] Code follows project conventions
- [ ] Feature integrates smoothly into existing CLI
