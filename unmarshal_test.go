package gocsv_test

import (
	"testing"

	gocsv "github.com/JoelLau/go-csv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	type Model struct {
		Int    int     `csv:"int"`
		Float  float32 `csv:"float"`
		String string  `csv:"string"`
	}

	given := []byte(`int,float,string
0,0.9,asdf
1,6.2,lorem ipsum
2,1.3,sdfdf1
`)

	got := []Model{}
	want := []Model{
		{Int: 0, String: "asdf", Float: 0.9},
		{Int: 1, String: "lorem ipsum", Float: 6.2},
		{Int: 2, String: "sdfdf1", Float: 1.3},
	}

	err := gocsv.Unmarshal(given, &got)
	require.NoError(t, err)

	assert.NotNil(t, got)
	assert.Len(t, got, 3)
	require.ElementsMatch(t, got, want)
}
