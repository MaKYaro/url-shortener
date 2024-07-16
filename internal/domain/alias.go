package domain

import "time"

type Alias struct {
	Value  string
	URL    string
	Expire time.Time
}

func (a *Alias) Expired() bool {
	return a.Expire.Before(time.Now())
}
