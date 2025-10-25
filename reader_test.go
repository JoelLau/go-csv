package gocsv_test

import (
	"testing"

	gocsv "github.com/JoelLau/go-csv"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestCSVParser(t *testing.T) {
	t.Parallel()

	given := []byte(`ID,Chosen Name,Birth Date
1,"Barrack Obama",1961-08-04
2,"""Stone Cold"", Steve Austin", 1964-12-08
`)

	p := gocsv.Parser{Delimeter: ','}
	got, err := p.ReadAll(given)
	want := [][]string{
		{"ID", "Chosen Name", "Birth Date"},
		{"1", "Barrack Obama", "1961-08-04"},
		{"2", "\"Stone Cold\", Steve Austin", "1964-12-08"},
	}

	require.NoError(t, err)

	diff := cmp.Diff(got, want)
	require.Emptyf(t, diff, "%+v", diff)
}
