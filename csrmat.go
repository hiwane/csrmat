package csrmat

import (
	"fmt"
)

/*
 * Compressed Sparse Row
 *
 * - 最初 Append() だけやる
 * - Compress() する
 * - その後は　Get(), Set() のみが許される
 *
 * 「行 >> 列」である疎な行列であり，
 * 同じ成分のみを上書き参照する場合を想定
 */
type CSRMatrix struct {
	v         []float64
	row       []int
	col       []int
	compressed bool
}

func NewCSRMatrix() *CSRMatrix {
	m := new(CSRMatrix)
	return m
}

/*
 * i 行 j 列に要素 val を追加する.
 * 受検者 i
 * 項目 j
 */
func (m *CSRMatrix) Append(i, j int, val float64) error {
	if m.compressed {
		return fmt.Errorf("CSRMatrix is compressed")
	}

	if len(m.row) == 0 || m.row[len(m.row)-1] < i ||
		m.row[len(m.row)-1] == i && m.col[len(m.col)-1] < j {
		m.v = append(m.v, val)
		m.row = append(m.row, i)
		m.col = append(m.col, j)
		return nil
	}
	// 追加する場所を探す.
	// row は昇順になっていて，最後の行のものが追加されるはず
	var k int
	for k = len(m.row) - 1; k >= 0; k-- {
		if m.row[k] < i {
			break
		}
	}
	for ; ; k++ {
		if m.row[k+1] > i || m.col[k+1] > j {
			k++
			m.row = append(m.row, 1)
			copy(m.row[k+1:], m.row[k:])
			m.row[k] = i
			m.col = append(m.col, 1)
			copy(m.col[k+1:], m.col[k:])
			m.col[k] = j
			m.v = append(m.v, 1)
			copy(m.v[k+1:], m.v[k:])
			m.v[k] = val
			return nil
		}
	}
}

func (m *CSRMatrix) Compress() error {
	if m.compressed {
		return nil
	}
	for i := 1; i < len(m.row); i++ {
		if m.row[i-1] > m.row[i] || m.row[i-1] == m.row[i] && m.col[i-1] >= m.col[i] {
			return fmt.Errorf("invalid i=%d, (%d,%d) (%d,%d)", i-1, m.row[i-1], m.col[i-1], m.row[i], m.col[i])
		}
	}

	// row 行を圧縮
	nrow := m.row[len(m.row)-1] + 1
	row := make([]int, nrow+1)
	j := int(0)
	for i := int(0); i < nrow; i++ {
		row[i] = -1
		for m.row[j] < i {
			j++
		}
		if m.row[j] == i {
			row[i] = j
			continue
		}
	}
	row[nrow] = int(len(m.row))
	p := row[nrow]
	for i := nrow - 1; i >= 0; i-- {
		if row[i] == -1 {
			row[i] = p
		} else {
			p = row[i]
		}
	}
	m.row = row
	m.compressed = true
	return nil
}

func (m *CSRMatrix) index(i, j int) int {
	l := m.row[i]
	h := m.row[i+1] - 1
	for l <= h {
		mid := (h + l) / 2
		k := m.col[mid]
		if k == j {
			return mid
		} else if k > j {
			h = mid - 1
		} else {
			l = mid + 1
		}
	}
	panic(fmt.Sprintf("not found index(%d,%d) range=%d..%d", i, j, m.row[i], m.row[i+1]))
}

func (m *CSRMatrix) Get(i, j int) float64 {
	return m.v[m.index(i, j)]
}

func (m *CSRMatrix) Set(i, j int, v float64) {
	m.v[m.index(i, j)] = v
}
