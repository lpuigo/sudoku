package sudoku

import (
	"fmt"
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
	res.WriteString("       A  B  C  .  D  E  F  .  G  H  I\n")
	for r := 0; r < s.size; r++ {
		if r > 0 && r%3 == 0 {
			res.WriteString("   -  ----------+-----------+----------\n")
		}
		res.WriteString(fmt.Sprintf("   %d  ", r+1))
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
func (s Sudoku) GetAllOptions() Options {
	res := Options{}
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

	res.SortByLength()
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

func (s *Sudoku) Solve(depth int) (int, bool) {
	fmt.Printf(s.String())
	mdepth := depth
	var options Options
	// loop on obvious solutions
	//fmt.Printf("(depth = %d) Found %d options\n", depth, len(options))

	// get all available options
	options = s.GetAllOptions()
	for len(options) > 0 {

		nbObvious, result := s.ResolveObviousOptions(options)
		if nbObvious > 0 { // some obvious solution found, update options
			fmt.Println(result)
			options = s.GetAllOptions()
			continue
		}

		nbHiddenSingletons, result := s.ResolveHiddenSingletonsOptions(options)
		if nbHiddenSingletons > 0 {
			fmt.Println(result)
			options = s.GetAllOptions()
			continue
		}

		nbHiddenPairs, result := s.ResolveHiddenPairsOptions(options)
		if nbHiddenPairs > 0 {
			fmt.Println(result)
			continue
		}

		nbHiddenTriplets, result := s.ResolveHiddenTripletsOptions(options)
		if nbHiddenTriplets > 0 {
			fmt.Println(result)
			continue
		}

		nakedTriplets, result := s.ResolveNakedTripletOptions(options)
		if nakedTriplets > 0 {
			fmt.Println(result)
			continue
		}

		nakedPairs, result := s.ResolveNakedPairOptions(options)
		if nakedPairs > 0 {
			fmt.Println(result)
			continue
		}

		// no obvious solution found, exit current loop to switch to another strategy
		if nbObvious == 0 {
			break
		}
	}

	// if no options found, Sudoku is solved
	if len(options) == 0 {
		completed := s.Completed()
		fmt.Printf("No other options, sudoku completed=%v\n%s", completed, s.String())
		return depth, completed
	}

	// finally try remaining options with recursive strategy
	option := options[0]
	for _, value := range option.GetValues() {
		fmt.Printf("(depth = %d/%d) Set possible %d of %s\n", depth, mdepth, value, option.String())
		s2 := s.Clone()
		s2.SetValue(value, option.row, option.col)

		ld, completed := s2.Solve(depth + 1)
		if ld > mdepth {
			mdepth = ld
		}
		if completed {
			// this option/value was OK, accept result and exit successfully
			s.values = s2.values
			return mdepth, true
		}
	}
	return mdepth, s.Completed()
}

// ResolveObviousOptions sets all obvious option (that is option with only 1 possible value)
func (s *Sudoku) ResolveObviousOptions(options Options) (int, string) {
	obvOpts := Options{}
	res := "Obvious Options:"
	for _, option := range options {
		if option.Length() != 1 {
			continue
		}
		obvOpts = append(obvOpts, option)
		s.SetValue(option.GetValues()[0], option.row, option.col)
	}
	nbObvious := len(obvOpts)
	if nbObvious == 0 {
		res += " None"
	} else {
		list := make([]string, nbObvious)
		for i, opt := range obvOpts {
			list[i] = fmt.Sprintf("%s", opt.String())
		}
		res += fmt.Sprintf(" x%d (%s)", nbObvious, strings.Join(list, ", "))
	}
	return nbObvious, res
}

// ResolveNakedPairOptions based on https://sudoku.com/fr/regles-du-sudoku/paires-nues
func (s Sudoku) ResolveNakedPairOptions(options Options) (int, string) {
	// for each subsquare
	actions := []string{}
	for c := 0; c < s.size; c += 3 {
		for r := 0; r < s.size; r += 3 {
			// get options for current subsquare
			subSquareFilter := FilterSubSquareFunc(r, c)
			keep := func(opt Option) bool { return subSquareFilter(opt) && opt.Length() >= 2 }
			localOptions := options.Filter(keep)
			if len(localOptions) < 4 { // not enough options for naked pair technic
				continue
			}

			// first and second localOptions must be a pair, otherwise no solution => skip to next subsquare
			if localOptions[0].Length() != 2 || localOptions[1].Length() != 2 {
				continue
			}

			// check if same pair, otherwise no solution => skip to next subsquare
			pair := localOptions[0].option
			if !localOptions[1].option.Contains(pair) {
				continue
			}

			// we found our two pairs, remove them from remaining options
			for _, option := range localOptions[2:] {
				if option.option.Contains(pair) {
					actions = append(actions, fmt.Sprintf("%s from %s", pair.String(), option.String()))
					option.option.RemoveSet(pair)
				}
			}
		}
	}
	nbNakedPairs := len(actions)
	res := "Naked Pairs: "
	if nbNakedPairs == 0 {
		res += "    None"
	} else {
		res += fmt.Sprintf("   x%d (%s)", nbNakedPairs, strings.Join(actions, ", "))
	}

	return nbNakedPairs, res
}

// ResolveNakedTripletOptions based on https://sudoku.com/fr/regles-du-sudoku/triplets-nus
func (s Sudoku) ResolveNakedTripletOptions(options Options) (int, string) {
	// for each subsquare
	actions := []string{}

	controlTriplets := func(localOpts Options) {
		//fmt.Printf("DEBUG TRIPLET controlTriplets: %d options: %s\n", len(localOpts), localOpts.String())
		if len(localOpts) < 4 { // not enough options for naked triplets technic
			return
		}

		// localOptions are sorted by ascending length. Three first localOptions must be a pair, otherwise no solution => skip to next subsquare
		if localOpts[2].Length() != 2 {
			return
		}

		// get possibles numbers given by three first localOptions
		possibleNumbers := make(map[int]int)
		for _, pair := range localOpts[:3] {
			for _, v := range pair.GetValues() {
				possibleNumbers[v]++
			}
		}
		// check if there are 3 possibles numbers, each two times
		if len(possibleNumbers) > 3 {
			return
		}
		for _, nb := range possibleNumbers {
			if nb != 2 {
				return
			}
		}

		// we found our three pairs, remove possibles numbers from remaining options
		for _, option := range localOpts[3:] {
			for _, pair := range localOpts[:3] {
				if option.option.Contains(pair.option) {
					actions = append(actions, fmt.Sprintf("%s from %s", pair.option.String(), option.String()))
					option.option.RemoveSet(pair.option)
				}
			}
		}
	}

	for c := 0; c < s.size; c += 3 {
		for r := 0; r < s.size; r += 3 {
			// get options for current subsquare
			subSquareFilter := FilterSubSquareFunc(r, c)
			keep := func(opt Option) bool { return subSquareFilter(opt) && opt.Length() >= 2 }
			localOptions := options.Filter(keep)
			controlTriplets(localOptions)
		}
	}

	nbNakedTriplets := len(actions)
	res := "Naked Triplets: "
	if nbNakedTriplets == 0 {
		res += "    None"
	} else {
		res += fmt.Sprintf(" x%d (%s)", nbNakedTriplets, strings.Join(actions, ", "))
	}

	return nbNakedTriplets, res
}

// ResolveHiddenSingletonsOptions based on https://sudoku.com/fr/regles-du-sudoku/singletons-caches
func (s Sudoku) ResolveHiddenSingletonsOptions(options Options) (int, string) {
	// for each subsquare
	actions := []string{}

	controlHiddenSingleton := func(localOpts Options) {
		//fmt.Printf("DEBUG HiddenSingleton control: %d options: %s\n", len(localOpts), localOpts.String())
		if len(localOpts) < 1 { // not enough options for naked triplets technic
			return
		}

		// get possibles numbers given by localOptions
		possibleNumbers := make(map[int]int)
		for _, pair := range localOpts {
			for _, v := range pair.GetValues() {
				possibleNumbers[v]++
			}
		}

		// search and keep only singleton
		for n, nb := range possibleNumbers {
			if nb > 1 {
				delete(possibleNumbers, n)
			}
		}

		// process singleton
		for _, option := range localOpts {
			for n, _ := range possibleNumbers {
				// if singleton is within this option, apply it
				if _, found := option.option[n]; found {
					option.option = ValueSet{n: struct{}{}}
					actions = append(actions, fmt.Sprintf("%s", option.String()))
					s.SetValue(n, option.row, option.col)
					break
				}
			}
		}
	}

	for c := 0; c < s.size; c += 3 {
		for r := 0; r < s.size; r += 3 {
			// get options for current subsquare
			subSquareFilter := FilterSubSquareFunc(r, c)
			keep := func(opt Option) bool { return subSquareFilter(opt) && opt.Length() >= 2 }
			localOptions := options.Filter(keep)
			controlHiddenSingleton(localOptions)
		}
	}

	nbHiddenSingletons := len(actions)
	res := "Hidden Singletons: "
	if nbHiddenSingletons == 0 {
		res += "    None"
	} else {
		res += fmt.Sprintf(" x%d (%s)", nbHiddenSingletons, strings.Join(actions, ", "))
	}

	return nbHiddenSingletons, res
}

// ResolveHiddenPairsOptions based on https://sudoku.com/fr/regles-du-sudoku/paires-cachees/
func (s Sudoku) ResolveHiddenPairsOptions(options Options) (int, string) {
	// for each subsquare
	actions := []string{}

	controlHiddenPairs := func(localOpts Options) {
		//fmt.Printf("DEBUG HiddenPairs control: %d options: %s\n", len(localOpts), localOpts.String())
		if len(localOpts) < 2 { // not enough options for hidden pairs technic
			return
		}

		// get possible numbers given by localOptions
		possibleNumbers := make(map[int]int)
		for _, option := range localOpts {
			for _, v := range option.GetValues() {
				possibleNumbers[v]++
			}
		}

		// search and keep only numbers with 2 possibles
		for n, nb := range possibleNumbers {
			if nb != 2 {
				delete(possibleNumbers, n)
			}
		}
		//fmt.Printf("DEBUG HiddenPairs control: possibles %v\n", possibleNumbers)
		// check if there are at least 2 numbers with 2 possibles
		if len(possibleNumbers) < 2 {
			return
		}

		// search for available pairs of numbers
		var targetOptions Options
		var pair ValueSet
	Search:
		for i, _ := range possibleNumbers {
			for j, _ := range possibleNumbers {
				if i == j {
					continue
				}
				pair = ValueSet{i: struct{}{}, j: struct{}{}}
				// find options with this pair
				targetOptions = Options{}
				for _, option := range localOpts {
					if option.option.Contains(pair) {
						targetOptions = append(targetOptions, option)
					}
				}
				if len(targetOptions) == 2 {
					break Search
				}
			}
		}
		if len(targetOptions) != 2 {
			return
		}
		//fmt.Printf("DEBUG HiddenPairs control: found pair %s in %s\n", pair.String(), targetOptions.String())

		// remove all others possibles than pair from target options
		actualOptions := []string{}
		for _, option := range targetOptions {
			if option.Length() < 3 { // do not process options with less than 3 possibles
				continue
			}
			actualOptions = append(actualOptions, option.String())
			option.option.RemoveButSet(pair)
		}
		if len(actualOptions) > 0 {
			actions = append(actions, fmt.Sprintf("%s from %s", pair.String(), strings.Join(actualOptions, " / ")))
		}
	}

	for c := 0; c < s.size; c += 3 {
		for r := 0; r < s.size; r += 3 {
			// get options for current subsquare
			subSquareFilter := FilterSubSquareFunc(r, c)
			keep := func(opt Option) bool { return subSquareFilter(opt) && opt.Length() >= 2 }
			localOptions := options.Filter(keep)
			controlHiddenPairs(localOptions)
		}
	}

	nbHiddenPairs := len(actions)
	res := "Hidden Pairs: "
	if nbHiddenPairs == 0 {
		res += "    None"
	} else {
		res += fmt.Sprintf(" x%d (%s)", nbHiddenPairs, strings.Join(actions, ", "))
	}

	return nbHiddenPairs, res
}

// ResolveHiddenTripletsOptions based on https://sudoku.com/fr/regles-du-sudoku/triplets-caches/
func (s Sudoku) ResolveHiddenTripletsOptions(options Options) (int, string) {
	// for each subsquare
	actions := []string{}

	controlHiddenTriplets := func(localOpts Options) {
		//fmt.Printf("DEBUG HiddenTriplets control: %d options: %s\n", len(localOpts), localOpts.String())
		if len(localOpts) < 3 { // not enough options for hidden triplet technic
			return
		}

		// get possible numbers given by localOptions
		possibleNumbers := make(map[int]int)
		for _, option := range localOpts {
			for _, v := range option.GetValues() {
				possibleNumbers[v]++
			}
		}

		// search and keep only numbers with 3 possibles
		for n, nb := range possibleNumbers {
			if nb != 3 {
				delete(possibleNumbers, n)
			}
		}
		//fmt.Printf("DEBUG HiddenTriplets control: possibles %v\n", possibleNumbers)
		// check if there are at least 3 numbers with 3 possibles
		if len(possibleNumbers) < 3 {
			return
		}

		// search for available triplet of numbers
		var targetOptions Options
		var triplet ValueSet
	Search:
		for i, _ := range possibleNumbers {
			for j, _ := range possibleNumbers {
				if i == j {
					continue
				}
				for k, _ := range possibleNumbers {
					if i == k || j == k {
						continue
					}
					triplet = ValueSet{i: struct{}{}, j: struct{}{}, k: struct{}{}}
					// find options with this pair
					targetOptions = Options{}
					for _, option := range localOpts {
						if option.option.Contains(triplet) {
							targetOptions = append(targetOptions, option)
						}
					}
					if len(targetOptions) == 3 {
						break Search
					}
				}
			}
		}
		if len(targetOptions) != 3 {
			return
		}
		//fmt.Printf("DEBUG HiddenTriplets control: found triplet %s in %s\n", triplet.String(), targetOptions.String())

		// remove all others possibles than triplet from target options
		actualOptions := []string{}
		for _, option := range targetOptions {
			if option.Length() < 4 { // do not process options with less than 4 possibles
				continue
			}
			actualOptions = append(actualOptions, option.String())
			option.option.RemoveButSet(triplet)
		}
		if len(actualOptions) > 0 {
			actions = append(actions, fmt.Sprintf("%s from %s", triplet.String(), strings.Join(actualOptions, " / ")))
		}
	}

	for c := 0; c < s.size; c += 3 {
		for r := 0; r < s.size; r += 3 {
			// get options for current subsquare
			subSquareFilter := FilterSubSquareFunc(r, c)
			keep := func(opt Option) bool { return subSquareFilter(opt) && opt.Length() >= 2 }
			localOptions := options.Filter(keep)
			controlHiddenTriplets(localOptions)
		}
	}

	nbHiddenTriplets := len(actions)
	res := "Hidden Triplets: "
	if nbHiddenTriplets == 0 {
		res += "    None"
	} else {
		res += fmt.Sprintf(" x%d (%s)", nbHiddenTriplets, strings.Join(actions, ", "))
	}

	return nbHiddenTriplets, res
}
