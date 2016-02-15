package gopdf

import (
	"bytes"
	"testing"
)

func TestWriteUInt32(t *testing.T) {
	var buff bytes.Buffer
	err := WriteUInt32(&buff, 65536)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	b := buff.Bytes()
	if b[0] != 0 || b[1] != 1 || b[2] != 0 || b[3] != 0 {
		t.Errorf("WriteUInt64 fail")
	}
}

/*
func TestFloat(t *testing.T){
	a := float64(1170.08)
    b := float64(1013.08)
    c := float64(54673.00)
    d := float64(131588.00)
    e := float64(54236.52)
	sum := a + b + c + d + e
	if sum != 242680.68 {
		t.Errorf("!242680.68")
	}
}*/
