package view

func idSliceToIntSlice(ids []Id) []int64 {
	ints := make([]int64, len(ids))
	for idx, i := range ids {
		ints[idx] = int64(i)
	}
	return ints
}
