package state

import "time"

func mustTime(in string) time.Time {
	t, err := time.Parse(time.RFC3339, in)
	if err != nil {
		panic(err)
	}
	return t
}
