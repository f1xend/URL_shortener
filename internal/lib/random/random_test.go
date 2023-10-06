package random

import "testing"

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "lenght = 1",
			length: 1,
		},
		{
			name:   "lenght = 3",
			length: 3,
		},
		{
			name:   "lenght 5",
			length: 5,
		},
		{
			name:   "lenght 10",
			length: 10,
		},
		{
			name:   "lenght 15",
			length: 15,
		},
		{
			name:   "lenght 20",
			length: 20,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			str := NewRandomString(v.length)
			if len(str) != v.length {
				t.Errorf("got %d, want %d", len(str), v.length)
			}
		})
	}
}
