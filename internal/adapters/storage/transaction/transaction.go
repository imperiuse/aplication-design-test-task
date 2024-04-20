package transaction

import (
	"context"
	"fmt"
)

// NOT go-routine safe!
// super simple transaction logic just for concept illustration!
// operation "hidden" until not "Commit" them.
type transaction struct {
	ctx context.Context

	operations []operation
	rollbacks  []rollbackFunc
}

type (
	operation    = func() error
	rollbackFunc = func() error
)

// New creates a new transaction with its own context.
func New(ctx context.Context) *transaction {
	return &transaction{
		ctx:        ctx,
		rollbacks:  make([]rollbackFunc, 0),
		operations: make([]operation, 0),
	}
}

func (t *transaction) Execute(op operation, rb rollbackFunc) {
	t.operations = append(t.operations, op)
	t.rollbacks = append(t.rollbacks, rb)
}

// Rollback executes all rollback functions in reverse order and cancels the context.
func (t *transaction) Rollback() error {
	return t.rollbackFrom(len(t.rollbacks))
}

func (t *transaction) rollbackFrom(lastOperationNumber int) error {
	var sumErr error
	for i := lastOperationNumber - 1; i >= 0; i-- {
		if err := t.rollbacks[i](); err != nil {
			if sumErr == nil {
				sumErr = err
			} else {
				sumErr = fmt.Errorf("%v; %w", sumErr, err)
			}
		}
	}

	t.rollbacks = nil

	return sumErr
}

// Commit clears the rollback stack and effectively commits the transaction.
func (t *transaction) Commit() error {
	defer func() {
		t.operations = nil
		t.rollbacks = nil
	}()

	for i, op := range t.operations {
		if err := op(); err != nil {
			rErr := t.rollbackFrom(i)
			return fmt.Errorf("Operation error: %v; Rollback err status: %w", err, rErr)
		}
	}

	return nil
}
