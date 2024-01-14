package sudoku

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type ValueSet map[int]struct{}

func (pv ValueSet) String() string {
	res := []string{}
	for _, v := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9} {
		if _, found := pv[v]; found {
			res = append(res, strconv.Itoa(v))
		}
	}
	sort.Strings(res)
	return fmt.Sprintf("[%s]", strings.Join(res, ", "))
}

func (pv ValueSet) GetValues() []int {
	res := []int{}
	for v, _ := range pv {
		res = append(res, v)
	}
	sort.Ints(res)
	return res
}
