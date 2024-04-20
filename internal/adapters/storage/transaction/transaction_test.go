package transaction

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransaction_Execute(t *testing.T) {
	ctx := context.Background()
	tr := New(ctx)

	// Test adding a successful operation
	tr.Execute(func() error {
		return nil // simulate successful operation
	}, func() error {
		return errors.New("rollback should not be called") // simulate rollback (should not be called)
	})

	assert.Len(t, tr.operations, 1, "Should add the operation to the list")
	assert.Len(t, tr.rollbacks, 1, "Should add the rollback function to the list")

	// Test adding a failing operation
	tr.Execute(func() error {
		return errors.New("operation error") // simulate operation failure
	}, func() error {
		return nil // rollback for the failed operation (should be irrelevant here)
	})

	assert.Len(t, tr.operations, 2, "Should still add the failing operation for later execution")
	assert.Len(t, tr.rollbacks, 2, "Should still add the rollback function")
}

func TestTransaction_Commit_Success(t *testing.T) {
	ctx := context.Background()
	tr := New(ctx)

	// Setup operations
	tr.Execute(func() error { return nil }, func() error { return errors.New("rollback failed") })
	tr.Execute(func() error { return nil }, func() error { return nil })

	// Attempt to commit the transaction
	err := tr.Commit()
	assert.NoError(t, err, "Commit should succeed without any errors")
	assert.Len(t, tr.operations, 0, "Commit should clear operations after successful execution")
	assert.Len(t, tr.rollbacks, 0, "Commit should clear rollbacks after successful execution")
}

func TestTransaction_Commit_Failure(t *testing.T) {
	ctx := context.Background()
	tr := New(ctx)

	// Setup operations, one of which will fail
	tr.Execute(func() error { return nil }, func() error { return nil })
	tr.Execute(func() error { return errors.New("operation failed") }, func() error { return nil })

	// Attempt to commit the transaction
	err := tr.Commit()
	assert.Error(t, err, "Commit should fail due to an operation error")
	assert.Contains(t, err.Error(), "operation failed", "Error message should include the operation failure")
	assert.Len(t, tr.operations, 0, "Commit should clear operations regardless of outcome")
	assert.Len(t, tr.rollbacks, 0, "Commit should clear rollbacks regardless of outcome")
}

func TestTransaction_Rollback(t *testing.T) {
	ctx := context.Background()
	tr := New(ctx)

	// Setup rollback functions
	tr.Execute(func() error { return nil }, func() error { return errors.New("rollback failed 1") })
	tr.Execute(func() error { return nil }, func() error { return nil }) // this rollback should succeed

	// Attempt to rollback the transaction
	err := tr.Rollback()
	assert.Error(t, err, "Rollback should fail due to one failing rollback function")
	assert.Contains(t, err.Error(), "rollback failed 1", "Error message should include the rollback failure")
	assert.Len(t, tr.rollbacks, 0, "Rollback should clear all rollback functions after execution")
}

func TestTransaction_RollbackFrom(t *testing.T) {
	ctx := context.Background()
	tr := New(ctx)

	// Set up a series of rollbacks where some fail and some succeed
	var rollbacksCalled []int
	tr.rollbacks = append(tr.rollbacks,
		func() error {
			rollbacksCalled = append(rollbacksCalled, 1)
			return nil // This rollback will succeed
		},
		func() error {
			rollbacksCalled = append(rollbacksCalled, 2)
			return errors.New("rollback failed 2") // This will fail
		},
		func() error {
			rollbacksCalled = append(rollbacksCalled, 3)
			return errors.New("rollback failed 3") // This will also fail
		},
		func() error {
			rollbacksCalled = append(rollbacksCalled, 4)
			return errors.New("rollback failed 4") // This will also fail
		},
	)

	// Calling rollbackFrom with 3 means we expect to rollback the last 3 operations only
	err := tr.rollbackFrom(3)

	// We are using require to ensure we do have an error before making more assertions
	require.Error(t, err, "rollbackFrom should return an error due to rollback failures")

	// The error message should be a concatenation of the individual rollback errors
	expectedErrMsg := "rollback failed 3; rollback failed 2"
	assert.Equal(t, expectedErrMsg, err.Error(), "Error message should contain all rollback errors in order")

	expectedRollbacksCalled := []int{3, 2, 1}
	assert.Equal(t, expectedRollbacksCalled, rollbacksCalled, "Only the last three of the rollback functions should have been called")

	// After rollback, there should be no rollback functions left in the slice
	assert.Len(t, tr.rollbacks, 0, "All rollback functions should be cleared after rollbackFrom is called")
}
