package domain

import "time"

type Alias struct {
	Value  string
	URL    string
	Expire time.Time
}

func (a *Alias) IsExpired() bool {
	return a.Expire.Before(time.Now())
}
