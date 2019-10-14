package restruct

import (
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/go-restruct/restruct/formats/png"
	"github.com/stretchr/testify/assert"
)

func readfile(fn string) []byte {
	if d, err := ioutil.ReadFile(fn); err == nil {
		return d
	} else {
		panic(err)
	}
}

func TestPNGGrad8RGB(t *testing.T) {
	EnableExprBeta()

	tests := []struct {
		format     interface{}
		expectdata []byte
		expectjson []byte
	}{
		{
			format:     png.File{},
			expectdata: readfile("testdata/pnggrad8rgb.png"),
			expectjson: readfile("testdata/pnggrad8rgb.json"),
		},
	}

	for _, test := range tests {
		f := reflect.New(reflect.TypeOf(test.format)).Interface()
		assert.Nil(t, json.Unmarshal(test.expectjson, f))
		data, err := Pack(binary.BigEndian, f)
		assert.Nil(t, err)
		assert.Equal(t, test.expectdata, data)

		f = reflect.New(reflect.TypeOf(test.format)).Interface()
		assert.Nil(t, Unpack(test.expectdata, binary.BigEndian, f))
		data, err = json.Marshal(f)
		assert.Nil(t, err)
		assert.JSONEq(t, string(test.expectjson), string(data))
	}
}
