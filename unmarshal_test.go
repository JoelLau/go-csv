package gocsv_test

import (
	"fmt"
	"testing"
	"time"

	gocsv "github.com/JoelLau/go-csv"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshal_BasicTypes(t *testing.T) {
	t.Parallel()

	type Model struct {
		Int     int     `csv:"int"`
		Uint    uint    `csv:"uint"`
		Float32 float32 `csv:"float32"`
		Float64 float64 `csv:"float64"`
		Bool    bool    `csv:"bool"`
		String  string  `csv:"string"`
	}

	given := []byte(`int,uint,float32,float64,bool,string,unmapped
0,5,0.9,1.111,true,asdf,
1,6,6.2,22.22,false,lorem ipsum,
2,7,1.3,333.3,true,sdfdf1,unmapped
`)

	got := []Model{}
	want := []Model{
		{Int: 0, Uint: 5, Float32: 0.9, Float64: 1.111, Bool: true, String: "asdf"},
		{Int: 1, Uint: 6, Float32: 6.2, Float64: 22.22, Bool: false, String: "lorem ipsum"},
		{Int: 2, Uint: 7, Float32: 1.3, Float64: 333.3, Bool: true, String: "sdfdf1"},
	}

	err := gocsv.Unmarshal(given, &got)
	require.NoError(t, err)

	assert.NotNil(t, got)
	assert.Len(t, got, 3)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unmarshall() mismatch (-want +got):\n%s", diff)
	}
}

type date_yyyyMMdd time.Time

func (d *date_yyyyMMdd) Unmarshal(data []byte) error {
	s := string(data)
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return fmt.Errorf("failed to parse date '%q': %w", s, err)
	}

	*d = date_yyyyMMdd(t)
	return nil
}

var _ gocsv.CSVUnmarshaller = &date_yyyyMMdd{}

func TestUnmarshal_CustomTypes(t *testing.T) {
	t.Parallel()

	type StructWithCustomTypes struct {
		Date date_yyyyMMdd `csv:"date"`
	}

	given := []byte(`date
2025-10-11
1993-03-30
`)
	got := make([]StructWithCustomTypes, 0)
	want := []StructWithCustomTypes{
		{Date: date_yyyyMMdd(newDate(t, "2025-10-11"))},
		{Date: date_yyyyMMdd(newDate(t, "1993-03-30"))},
	}

	err := gocsv.Unmarshal(given, &got)
	require.NoError(t, err)
	require.Len(t, got, len(want))
	require.WithinDuration(t, time.Time(got[0].Date), time.Time(want[0].Date), 0)
	require.WithinDuration(t, time.Time(got[1].Date), time.Time(want[1].Date), 0)
}

func newDate(t *testing.T, str string) time.Time {
	t.Helper()

	r, err := time.Parse(time.DateOnly, str)
	assert.NoError(t, err)

	return r
}
