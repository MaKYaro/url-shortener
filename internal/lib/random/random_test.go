package random

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	cases := []struct {
		name string
		size int
	}{
		{
			name: "size 1",
			size: 1,
		},
		{
			name: "size 5",
			size: 5,
		},
		{
			name: "size 10",
			size: 10,
		},
		{
			name: "size 20",
			size: 20,
		},
		{
			name: "size 30",
			size: 30,
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			g := NewGenerator(tcase.size)
			str1 := g.Generate()
			str2 := g.Generate()

			require.Len(t, str1, tcase.size)
			require.Len(t, str2, tcase.size)

			require.NotEqual(t, str1, str2)
		})
	}
}
