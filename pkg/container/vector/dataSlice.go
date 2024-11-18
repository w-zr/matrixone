package vector

import "github.com/matrixorigin/matrixone/pkg/container/types"

type DataSlice[T any] interface {
	MustToCol(*Vector, uint32)
	At(i, j int) T
	Length(i int) int
	Size() int
}

type FixedDataSlice[T any] struct {
	mustColFunc func(*Vector) []T

	cols [][]T
}

func NewFixedDataSlice[T any](mustColFunc func(*Vector) []T, size int) *FixedDataSlice[T] {
	return &FixedDataSlice[T]{
		mustColFunc: mustColFunc,
		cols:        make([][]T, size),
	}
}

func NewFixedDataSliceWithSlices[T any](slices [][]T) *FixedDataSlice[T] {
	return &FixedDataSlice[T]{
		cols: slices,
	}
}

func (s *FixedDataSlice[T]) MustToCol(v *Vector, i uint32) {
	s.cols[i] = s.mustColFunc(v)
}

func (s *FixedDataSlice[T]) At(i, j int) T {
	return s.cols[i][j]
}
func (s *FixedDataSlice[T]) Length(i int) int {
	return len(s.cols[i])
}

func (s *FixedDataSlice[T]) Size() int {
	return len(s.cols)
}

type VarlenaDataSlice struct {
	cols []struct {
		data []types.Varlena
		area []byte
	}
}

func NewVarlenaDataSlice(size int) DataSlice[string] {
	return &VarlenaDataSlice{cols: make([]struct {
		data []types.Varlena
		area []byte
	}, size)}
}

func NewVarlenaDataSliceWithSlices(data [][]types.Varlena, area []byte) DataSlice[string] {
	cols := make([]struct {
		data []types.Varlena
		area []byte
	}, len(data))

	for i := range data {
		cols[i] = struct {
			data []types.Varlena
			area []byte
		}{data: data[i], area: area}
	}
	return &VarlenaDataSlice{cols: cols}
}

func (s *VarlenaDataSlice) MustToCol(v *Vector, i uint32) {
	data, area := MustVarlenaRawData(v)
	s.cols[i] = struct {
		data []types.Varlena
		area []byte
	}{data: data, area: area}
}

func (s *VarlenaDataSlice) At(i, j int) string {
	return s.cols[i].data[j].UnsafeGetString(s.cols[i].area)
}

func (s *VarlenaDataSlice) Length(i int) int {
	return len(s.cols[i].data)
}

func (s *VarlenaDataSlice) Size() int {
	return len(s.cols)
}
