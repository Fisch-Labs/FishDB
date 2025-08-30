/*
 * FishDB
 *
// Copyright 2025 Fisch-labs
 *
*/

package storage

import (
	"errors"
	"fmt"

	"github.com/Fisch-Labs/common/pools"
)

/*
BufferPool is a pool of byte buffers.
*/
var BufferPool = pools.NewByteBufferPool()

/*
Common storage manager related errors.
*/
var (
	ErrSlotNotFound = errors.New("Slot not found")
	ErrNotInCache   = errors.New("No entry in cache")
)

/*
ManagerError is a storage manager related error.
*/
type ManagerError struct {
	Type        error
	Detail      string
	Managername string
}

/*
NewStorageManagerError returns a new StorageManager specific error.
*/
func NewStorageManagerError(smeType error, smeDetail string, smeManagername string) *ManagerError {
	return &ManagerError{smeType, smeDetail, smeManagername}
}

/*
Error returns a string representation of the error.
*/
func (e *ManagerError) Error() string {
	return fmt.Sprintf("%s (%s - %s)", e.Type.Error(), e.Managername, e.Detail)
}
