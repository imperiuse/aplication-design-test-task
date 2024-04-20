package model

import (
	"testing"

	"aplication-design-test-task/internal/core/util"
)

func TestCheckStatusForDuplicates(t *testing.T) {
	err := util.CheckSliceForDuplicates(allStatuses[:])
	if err != nil {
		t.Errorf("not all vaue in status are uniq, check Status enums")
	}
}
