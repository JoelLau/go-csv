package gocsv_test

import (
	"testing"

	gocsv "github.com/JoelLau/go-csv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshal_HeaderBased(t *testing.T) {
	t.Parallel()

	type M struct {
		ID    string `csv:"id"`
		Fruit string `csv:"fruit"`
	}

	given := []byte(`csv,fruit
aaaa,apple
bbbb,banana
cccc,cherry`)

	got := []M{}
	want := []M{
		{ID: "aaaa", Fruit: "apple"},
		{ID: "bbbb", Fruit: "banana"},
		{ID: "cccc", Fruit: "cherry"},
	}

	err := gocsv.Unmarshal(given, &got)
	require.NoError(t, err)

	assert.NotNil(t, got)
	assert.Len(t, got, 3)
	require.ElementsMatch(t, got, want)
}
