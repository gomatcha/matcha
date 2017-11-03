package internal

func Diff(keysA, keysB []interface{}) (addedA []int, matchingAB map[int]int, removedB []int) {
	matching := map[int]int{}
	added := make([]bool, len(keysA))
	removed := make([]bool, len(keysB))

	for idxA, a := range keysA {
		found := false
		for idxB, b := range keysB {
			if a == b && !removed[idxB] {
				removed[idxB] = true
				matching[idxA] = idxB
				found = true
				break
			}
		}

		if !found {
			added[idxA] = true
		}
	}

	added2 := []int{}
	removed2 := []int{}
	for idx, i := range removed {
		if !i {
			removed2 = append(removed2, idx)
		}
	}
	for idx, i := range added {
		if i {
			added2 = append(added2, idx)
		}
	}

	return added2, matching, removed2
}
