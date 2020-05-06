package csrmat

import (
	"testing"
)

func TestAppend(t *testing.T) {
	vv := []struct {
		i, j int
		v    float64
	}{
		{0, 0, 0.0},
		{0, 3, 0.3},
		{0, 4, 0.4},
		{2, 6, 2.6},
		{2, 5, 2.5},
		{2, 7, 2.7},
		{3, 6, 3.6},
		{3, 7, 3.7},
		{3, 1, 3.1}}

	m := NewCSRMatrix()
	for _, v := range vv {
		m.Append(v.i, v.j, v.v)
	}
	m.Compress()
	var u float64
	for _, v := range vv {
		u = m.Get(v.i, v.j)
		if u != v.v {
			t.Errorf("m.Get(%d,%d): expect %f, actual %f", v.i, v.j, v.v, u)
		}
	}

	m.Set(0, 3, -0.3)
	m.Set(2, 6, -2.6)

	vv[1].v = -vv[1].v
	vv[3].v = -vv[3].v

	for _, v := range vv {
		u = m.Get(v.i, v.j)
		if u != v.v {
			t.Errorf("m.Set(%d,%d): expect %f, actual %f", v.i, v.j, v.v, u)
		}
	}
}
