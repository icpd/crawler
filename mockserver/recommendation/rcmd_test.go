package recommendation

import (
	"math"
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNextGuess(t *testing.T) {
	tests := []struct {
		name string
		id   int64
		want []int64
	}{
		{
			name: "normal",
			id:   12345,
			want: []int64{12335, 12340, 12350, 12355},
		},
		{
			name: "very small",
			id:   8,
			want: []int64{3, 13, 18},
		},
		{
			name: "very large",
			id:   math.MaxInt64 - 7,
			want: []int64{9223372036854775790, 9223372036854775795, 9223372036854775805},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rand.Seed(23465)
			cl := Client{}
			got := cl.NextGuess(tt.id)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("NextGuess(): diff: -want + got\n%s", diff)
			}
		})
	}
}
