package domain

import "time"

type Alias struct {
	Value  string
	URL    string
	Expire time.Time
}
