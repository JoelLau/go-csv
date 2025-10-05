package gocsv_test

import (
	"testing"

	gocsv "github.com/JoelLau/go-csv"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshal_BasicTypes(t *testing.T) {
	t.Parallel()

	type Model struct {
		Int    int     `csv:"int"`
		Uint   uint    `csv:"uint"`
		Float  float32 `csv:"float"`
		Bool   bool    `csv:"bool"`
		String string  `csv:"string"`
	}

	given := []byte(`int,uint,float,bool,string
0,5,0.9,true,asdf
1,6,6.2,false,lorem ipsum
2,7,1.3,true,sdfdf1
`)

	got := []Model{}
	want := []Model{
		{Int: 0, Uint: 5, Float: 0.9, Bool: true, String: "asdf"},
		{Int: 1, Uint: 6, Float: 6.2, Bool: false, String: "lorem ipsum"},
		{Int: 2, Uint: 7, Float: 1.3, Bool: true, String: "sdfdf1"},
	}

	err := gocsv.Unmarshal(given, &got)
	require.NoError(t, err)

	assert.NotNil(t, got)
	assert.Len(t, got, 3)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
	}
}
