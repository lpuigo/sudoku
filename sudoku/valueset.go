package sudoku

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type ValueSet map[int]struct{}

func NewValueSet(values ...int) ValueSet {
	pv := make(ValueSet)
	for _, v := range values {
		pv[v] = struct{}{}
	}
	return pv
}

func (pv ValueSet) String() string {
	res := []string{}
	for v, _ := range pv {
		res = append(res, strconv.Itoa(v))
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

func (pv ValueSet) Contains(ovs ValueSet) bool {
	for ov := range ovs {
		if _, found := pv[ov]; !found {
			return false
		}
	}
	return true
}

// RemoveSet removes all elements of the given ValueSet `ovs` from the receiver ValueSet `pv`.
func (pv ValueSet) RemoveSet(ovs ValueSet) {
	for ov := range ovs {
		delete(pv, ov)
	}
}

// RemoveButSet removes all elements not included in the given ValueSet `ovs` from the receiver ValueSet `pv`.
func (pv ValueSet) RemoveButSet(ovs ValueSet) {
	for v, _ := range pv {
		if _, found := ovs[v]; !found {
			delete(pv, v)
		}
	}
}

func (pv ValueSet) RemoveValue(value int) {
	delete(pv, value)
}
