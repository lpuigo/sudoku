package sudoku

import "fmt"

type Option struct {
	row, col int
	option   ValueSet
}

func (o Option) GetValues() []int {
	return o.option.GetValues()
}

func (o Option) String() string {
	return fmt.Sprintf("%s at (%d, %d)", o.option.String(), o.row, o.col)
}

func (o Option) ValueString(value int) string {
	return fmt.Sprintf("%d at (%d, %d)", value, o.row, o.col)
}
