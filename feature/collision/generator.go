package collision

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"

	"github.com/bwmarrin/snowflake"
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

func NewBase64Generator(length int) (*Base64Generator, error) {
	if length <= 0 {
		return nil, fmt.Errorf("length must be positive, got %d", length)
	}
	return &Base64Generator{
		length: length,
		chars:  []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"),
	}, nil
}

func (g *Base64Generator) Generate() string {
	result := make([]byte, g.length)
	max := big.NewInt(int64(len(g.chars)))

	for i := 0; i < g.length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(fmt.Sprintf("crypto/rand.Int failed: %v", err))
		}
		result[i] = g.chars[n.Int64()]
	}

	return string(result)
}

func (g *Base64Generator) TotalSpace() uint64 {
	space := uint64(1)
	for i := 0; i < g.length; i++ {
		if space > math.MaxUint64/64 {
			return math.MaxUint64
		}
		space *= 64
	}
	return space
}

func (g *Base64Generator) Name() string {
	return "base64"
}

type Base62Generator struct {
	length int
	chars  []byte
}

func NewBase62Generator(length int) (*Base62Generator, error) {
	if length <= 0 {
		return nil, fmt.Errorf("length must be positive, got %d", length)
	}
	return &Base62Generator{
		length: length,
		chars:  []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"),
	}, nil
}

func (g *Base62Generator) Generate() string {
	result := make([]byte, g.length)
	max := big.NewInt(int64(len(g.chars)))

	for i := 0; i < g.length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(fmt.Sprintf("crypto/rand.Int failed: %v", err))
		}
		result[i] = g.chars[n.Int64()]
	}

	return string(result)
}

func (g *Base62Generator) TotalSpace() uint64 {
	space := uint64(1)
	for i := 0; i < g.length; i++ {
		if space > math.MaxUint64/62 {
			return math.MaxUint64
		}
		space *= 62
	}
	return space
}

func (g *Base62Generator) Name() string {
	return "base62"
}

type SnowflakeGenerator struct {
	node *snowflake.Node
}

func NewSnowflakeGenerator() (*SnowflakeGenerator, error) {
	node, err := snowflake.NewNode(0)
	if err != nil {
		return nil, fmt.Errorf("failed to create snowflake node: %w", err)
	}
	return &SnowflakeGenerator{
		node: node,
	}, nil
}

func (g *SnowflakeGenerator) Generate() string {
	return g.node.Generate().String()
}

func (g *SnowflakeGenerator) TotalSpace() uint64 {
	return uint64(1) << 63
}

func (g *SnowflakeGenerator) Name() string {
	return "snowflake"
}
