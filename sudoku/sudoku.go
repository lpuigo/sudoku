package sudoku

import (
	"fmt"
	"sort"
	"strings"
)

type Sudoku struct {
	size   int
	values []int
}

const (
	valueUndef int = 0
	valueError int = -1
)

func New(size int) Sudoku {
	s := Sudoku{
		size:   size,
		values: make([]int, size*size),
	}

	return s
}

// Clone returns a deep copy of receiver
func (s Sudoku) Clone() Sudoku {
	nsv := make([]int, len(s.values))
	for i, value := range s.values {
		nsv[i] = value
	}

	return Sudoku{
		size:   s.size,
		values: nsv,
	}
}

func (s *Sudoku) SetValue(value, row, col int) {
	if !(row >= 0 && row <= s.size) {
		return
	}
	if !(col >= 0 && col <= s.size) {
		return
	}
	s.values[col+row*s.size] = value
}

func (s Sudoku) GetValue(row, col int) int {
	if !(row >= 0 && row <= s.size) {
		return valueError
	}
	if !(col >= 0 && col <= s.size) {
		return valueError
	}
	return s.values[col+row*s.size]
}

func (s Sudoku) getValue(row, col int) int {
	return s.values[col+row*s.size]
}

func (s Sudoku) getSubScareBounds(row, col int) (rowMin, rowMax, colMin, colMax int) {
	rowMin = row / 3 * 3
	colMin = col / 3 * 3
	return rowMin, rowMin + 2, colMin, colMin + 2
}

// IsValid returns true if value at position (row, col) is legit
func (s Sudoku) IsValid(value, row, col int) bool {
	// check for row
	for i := 0; i < s.size; i++ {
		if i == row {
			continue
		}
		if value == s.getValue(i, col) {
			return false
		}
	}
	// check for col
	for j := 0; j < s.size; j++ {
		if j == col {
			continue
		}
		if value == s.getValue(row, j) {
			return false
		}
	}
	// check for subScare
	rMin, rMax, cMin, cMax := s.getSubScareBounds(row, col)
	for r := rMin; r <= rMax; r++ {
		for c := cMin; c <= cMax; c++ {
			if r == row && c == col {
				continue
			}
			if value == s.getValue(r, c) {
				return false
			}
		}
	}
	return true
}

func (s Sudoku) String() string {
	res := strings.Builder{}
	res.WriteString("\n")
	for r := 0; r < s.size; r++ {
		if r > 0 && r%3 == 0 {
			res.WriteString("----------+-----------+----------\n")
		}
		for c := 0; c < s.size; c++ {
			if c > 0 && c%3 == 0 {
				res.WriteString(" | ")
			}
			v := s.getValue(r, c)
			if v == valueUndef {
				res.WriteString("   ")
			} else {
				res.WriteString(fmt.Sprintf(" %d ", v))
			}
		}
		res.WriteString("\n")
	}
	return res.String()
}

// GetValid returns a ValueSet of all possibles values at given position
func (s Sudoku) GetValid(row, col int) ValueSet {
	res := make(ValueSet)
	for _, v := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9} {
		if s.IsValid(v, row, col) {
			res[v] = struct{}{}
		}
	}
	return res
}

// GetAllOptions returns a slice of Option, giving, for each undef position, all possible values
//
// Result slice is sorted from option with the fewest possible values first
func (s Sudoku) GetAllOptions() []Option {
	res := []Option{}
	for r := 0; r < s.size; r++ {
		for c := 0; c < s.size; c++ {
			if s.getValue(r, c) != valueUndef {
				continue
			}
			opt := Option{
				row:    r,
				col:    c,
				option: s.GetValid(r, c),
			}
			if len(opt.option) == 0 {
				continue
			}
			res = append(res, opt)
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return len(res[i].option) < len(res[j].option)
	})
	return res
}

// Completed returns true if receiver has no undefined values (all values are set)
func (s Sudoku) Completed() bool {
	for _, value := range s.values {
		if value == valueUndef {
			return false
		}
	}
	return true
}

func (s *Sudoku) Solve(depth int) bool {
	fmt.Printf(s.String())
	// loop on obvious solutions
	options := s.GetAllOptions()
	//fmt.Printf("(depth = %d) Found %d options\n", depth, len(options))
	for {
		// if no options found, Sudoku is solved
		if len(options) == 0 {
			completed := s.Completed()
			fmt.Printf("(depth = %d) No other options, sudoku completed=%v %s", depth, completed, s.String())
			return completed
		}

		nbObvious := 0
		for _, option := range options {
			// set all obvious option (that is option with only 1 possible value)
			if len(option.option) != 1 {
				break
			}
			nbObvious++
			//fmt.Printf("(depth = %d) Set obvious %s\n", depth, option.String())
			s.SetValue(option.GetValues()[0], option.row, option.col)
		}
		// no obvious solution found, exit current loop to switch to another strategy
		if nbObvious == 0 {
			break
		} else {
			options = s.GetAllOptions()
			//fmt.Printf("(depth = %d) Found %d options\n", depth, len(options))
		}
	}

	// try first non-trivial options with recursive strategy
	option := options[0]
	for _, value := range option.GetValues() {
		fmt.Printf("(depth = %d) Set possible %d of %s\n", depth, value, option.String())
		s2 := s.Clone()
		s2.SetValue(value, option.row, option.col)
		if s2.Solve(depth + 1) {
			// this option/value was OK, accept result and exit successfully
			s.values = s2.values
			return true
		}
	}
	return s.Completed()
}
