package utils

import "github.com/underscorenico/tobcast/internal/data"

func ArrayDifference(a, b []data.Message) (diff []data.Message) {
	m := make(map[data.Message]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}
