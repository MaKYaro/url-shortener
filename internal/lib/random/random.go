package random

import (
	"math/rand"
	"strings"
	"time"
)

const chars = "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"

type AliasGenerator struct {
	aliasLen int
}

func NewGenerator(aliasLen int) *AliasGenerator {
	return &AliasGenerator{aliasLen: aliasLen}
}

func (g *AliasGenerator) Generate() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	sb := strings.Builder{}
	sb.Grow(g.aliasLen)
	for i := 0; i < g.aliasLen; i++ {
		sb.WriteByte(chars[rnd.Intn(len(chars))])
	}
	return sb.String()
}
