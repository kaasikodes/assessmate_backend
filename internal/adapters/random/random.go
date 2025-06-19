package randomadapter

import (
	"crypto/rand"
	"math/big"
	"strings"

	randomidgenerator "github.com/kaasikodes/assessmate_backend/internal/ports/outbound/random-id-generator"
)

type RandomIdAdapter struct {
	charset string
}

// NewRandomIdAdapter creates a new instance with default charset [A-Za-z0-9]
func NewRandomIdAdapter() randomidgenerator.RandomIdGenerator {
	return &RandomIdAdapter{
		charset: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
	}
}

// Create generates a random ID with the given prefix and length
func (g *RandomIdAdapter) Create(prefix string, length int) string {
	var sb strings.Builder
	sb.WriteString(prefix)

	charsetLen := big.NewInt(int64(len(g.charset)))

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			// fallback if crypto fails
			sb.WriteByte(g.charset[i%len(g.charset)])
		} else {
			sb.WriteByte(g.charset[n.Int64()])
		}
	}

	return sb.String()
}
