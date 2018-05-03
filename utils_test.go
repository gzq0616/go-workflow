package go_workflow

import "testing"

func TestParse(t *testing.T) {
	condition := "( { value1 } > 88 and { value2 } != true ) && { value3 } == ok "
	parse(condition)
}