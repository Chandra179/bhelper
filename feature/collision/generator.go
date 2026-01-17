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
