package sudoku

import (
	"fmt"
	"sort"
)

type Option struct {
	row, col int
	option   ValueSet
}

func (o Option) Length() int {
	return len(o.option)
}

func (o Option) GetValues() []int {
	return o.option.GetValues()
}

func (o Option) String() string {
	return fmt.Sprintf("%s_%s", o.option.String(), o.posString())
}

func (o Option) ValueString(value int) string {
	return fmt.Sprintf("%d_%s", value, o.posString())
}

func (o Option) posString() string {
	col := 'A' + o.col
	return fmt.Sprintf("%c%d", col, o.row+1)
}

type Options []Option

func (o Options) String() string {
	res := ""
	for i, option := range o {
		if i > 0 {
			res += fmt.Sprintf(", ")
		}
		res += option.String()
	}
	return res
}

func (o Options) SortByLength() {
	sort.Slice(o, func(i, j int) bool {
		return len(o[i].option) < len(o[j].option)
	})
}

func (o Options) Filter(keep func(opt Option) bool) Options {
	res := Options{}
	for _, option := range o {
		if keep(option) {
			res = append(res, option)
		}
	}
	res.SortByLength()
	return res
}

func FilterRowFunc(row int) func(opt Option) bool {
	return func(opt Option) bool { return opt.row == row }
}

func FilterColFunc(col int) func(opt Option) bool {
	return func(opt Option) bool { return opt.col == col }
}

func FilterSubSquareFunc(row, col int) func(opt Option) bool {
	rMin := row / 3 * 3
	rMax := rMin + 2
	cMin := col / 3 * 3
	cMax := cMin + 2
	return func(opt Option) bool { return opt.row >= rMin && opt.row <= rMax && opt.col >= cMin && opt.col <= cMax }
}

func (o Options) GetRowOptions(row int) Options {
	return o.Filter(FilterRowFunc(row))
}

func (o Options) GetColumnOptions(col int) Options {
	return o.Filter(FilterColFunc(col))
}

func (o Options) GetSubScareOptions(row, col int) Options {
	return o.Filter(FilterSubSquareFunc(row, col))
}
