package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExpired(t *testing.T) {
	now := time.Now()

	cases := []struct {
		name  string
		alias *Alias
		want  bool
	}{
		{
			"expired two hours ago",
			&Alias{Expire: now.Add(-2 * time.Hour)},
			true,
		},
		{
			"expired a minute ago",
			&Alias{Expire: now.Add(-1 * time.Minute)},
			true,
		},
		{
			"expired 3 seconds ago",
			&Alias{Expire: now.Add(-3 * time.Second)},
			true,
		},
		{
			"valid for two hours",
			&Alias{Expire: now.Add(2 * time.Hour)},
			false,
		},
		{
			"valid for a minute",
			&Alias{Expire: now.Add(1 * time.Minute)},
			false,
		},
		{
			"valid for 3 seconds",
			&Alias{Expire: now.Add(3 * time.Second)},
			false,
		},
		{
			"valid for 0 seconds",
			&Alias{Expire: now},
			true,
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			got := tcase.alias.Expired()
			require.Equal(t, tcase.want, got)
		})
	}
}
