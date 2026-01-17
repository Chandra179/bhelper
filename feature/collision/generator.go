package collision

import (
	"crypto/rand"
	"fmt"
	"math"
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
