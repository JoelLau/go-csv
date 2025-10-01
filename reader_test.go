package gocsv_test

import (
	"testing"

	gocsv "github.com/JoelLau/go-csv"
	"github.com/stretchr/testify/require"
)

func TestCSVParser(t *testing.T) {
	t.Parallel()

	given := []byte("id,fruit\n1,apple\n2,banana\n3,cherry")
	got, err := gocsv.ReadAll(given)
	want := [][]string{
		{"id", "fruit"},
		{"1", "apple"},
		{"2", "banana"},
		{"3", "cherry"},
	}

	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestReadRow(t *testing.T) {
	t.Parallel()

	given := "asdf,qwer,123,k!@#"
	got, err := gocsv.ReadRow(given)
	want := []string{
		"asdf",
		"qwer",
		"123",
		"k!@#",
	}

	require.NoError(t, err)
	require.Equal(t, want, got)
}
