package queue

import (
	"aplication-design-test-task/internal/core/util"
	"testing"
)

func TestCheckTopicsForDuplicates(t *testing.T) {
	err := util.CheckSliceForDuplicates(AllTopics[:])
	if err != nil {
		t.Errorf("AllTopics array contain duplicates, check `queue.AllTopics`")
	}
}
