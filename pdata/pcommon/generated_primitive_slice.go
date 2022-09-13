// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by "model/internal/cmd/pdatagen/main.go". DO NOT EDIT.
// To regenerate this file run "go run model/internal/cmd/pdatagen/main.go".

package pcommon

import "go.opentelemetry.io/collector/pdata/internal"

// ByteSlice represents a []byte slice.
// The instance of ByteSlice can be assigned to multiple objects since it's immutable.
//
// Must use NewByteSlice function to create new instances.
// Important: zero-initialized instance is not valid for use.
type ByteSlice internal.ByteSlice

func (ms ByteSlice) getOrig() *[]byte {
	return internal.GetOrigByteSlice(internal.ByteSlice(ms))
}

// NewByteSlice creates a new empty ByteSlice.
func NewByteSlice() ByteSlice {
	orig := []byte(nil)
	return ByteSlice(internal.NewByteSlice(&orig))
}

// AsRaw returns a copy of the []byte slice.
func (ms ByteSlice) AsRaw() []byte {
	return copyByteSlice(nil, *ms.getOrig())
}

// FromRaw copies raw []byte into the slice ByteSlice.
func (ms ByteSlice) FromRaw(val []byte) {
	*ms.getOrig() = copyByteSlice(*ms.getOrig(), val)
}

// Len returns length of the []byte slice value.
func (ms ByteSlice) Len() int {
	return len(*ms.getOrig())
}

// At returns an item from particular index.
func (ms ByteSlice) At(i int) byte {
	return (*ms.getOrig())[i]
}

// SetAt sets byte item at particular index.
func (ms ByteSlice) SetAt(i int, val byte) {
	(*ms.getOrig())[i] = val
}

// MoveTo moves ByteSlice to another instance.
func (ms ByteSlice) MoveTo(dest ByteSlice) {
	*dest.getOrig() = *ms.getOrig()
	*ms.getOrig() = nil
}

// CopyTo copies ByteSlice to another instance.
func (ms ByteSlice) CopyTo(dest ByteSlice) {
	*dest.getOrig() = copyByteSlice(*dest.getOrig(), *ms.getOrig())
}

func copyByteSlice(dst, src []byte) []byte {
	dst = dst[:0]
	return append(dst, src...)
}

// Float64Slice represents a []float64 slice.
// The instance of Float64Slice can be assigned to multiple objects since it's immutable.
//
// Must use NewFloat64Slice function to create new instances.
// Important: zero-initialized instance is not valid for use.
type Float64Slice internal.Float64Slice

func (ms Float64Slice) getOrig() *[]float64 {
	return internal.GetOrigFloat64Slice(internal.Float64Slice(ms))
}

// NewFloat64Slice creates a new empty Float64Slice.
func NewFloat64Slice() Float64Slice {
	orig := []float64(nil)
	return Float64Slice(internal.NewFloat64Slice(&orig))
}

// AsRaw returns a copy of the []float64 slice.
func (ms Float64Slice) AsRaw() []float64 {
	return copyFloat64Slice(nil, *ms.getOrig())
}

// FromRaw copies raw []float64 into the slice Float64Slice.
func (ms Float64Slice) FromRaw(val []float64) {
	*ms.getOrig() = copyFloat64Slice(*ms.getOrig(), val)
}

// Len returns length of the []float64 slice value.
func (ms Float64Slice) Len() int {
	return len(*ms.getOrig())
}

// At returns an item from particular index.
func (ms Float64Slice) At(i int) float64 {
	return (*ms.getOrig())[i]
}

// SetAt sets float64 item at particular index.
func (ms Float64Slice) SetAt(i int, val float64) {
	(*ms.getOrig())[i] = val
}

// MoveTo moves Float64Slice to another instance.
func (ms Float64Slice) MoveTo(dest Float64Slice) {
	*dest.getOrig() = *ms.getOrig()
	*ms.getOrig() = nil
}

// CopyTo copies Float64Slice to another instance.
func (ms Float64Slice) CopyTo(dest Float64Slice) {
	*dest.getOrig() = copyFloat64Slice(*dest.getOrig(), *ms.getOrig())
}

func copyFloat64Slice(dst, src []float64) []float64 {
	dst = dst[:0]
	return append(dst, src...)
}

// UInt64Slice represents a []uint64 slice.
// The instance of UInt64Slice can be assigned to multiple objects since it's immutable.
//
// Must use NewUInt64Slice function to create new instances.
// Important: zero-initialized instance is not valid for use.
type UInt64Slice internal.UInt64Slice

func (ms UInt64Slice) getOrig() *[]uint64 {
	return internal.GetOrigUInt64Slice(internal.UInt64Slice(ms))
}

// NewUInt64Slice creates a new empty UInt64Slice.
func NewUInt64Slice() UInt64Slice {
	orig := []uint64(nil)
	return UInt64Slice(internal.NewUInt64Slice(&orig))
}

// AsRaw returns a copy of the []uint64 slice.
func (ms UInt64Slice) AsRaw() []uint64 {
	return copyUInt64Slice(nil, *ms.getOrig())
}

// FromRaw copies raw []uint64 into the slice UInt64Slice.
func (ms UInt64Slice) FromRaw(val []uint64) {
	*ms.getOrig() = copyUInt64Slice(*ms.getOrig(), val)
}

// Len returns length of the []uint64 slice value.
func (ms UInt64Slice) Len() int {
	return len(*ms.getOrig())
}

// At returns an item from particular index.
func (ms UInt64Slice) At(i int) uint64 {
	return (*ms.getOrig())[i]
}

// SetAt sets uint64 item at particular index.
func (ms UInt64Slice) SetAt(i int, val uint64) {
	(*ms.getOrig())[i] = val
}

// MoveTo moves UInt64Slice to another instance.
func (ms UInt64Slice) MoveTo(dest UInt64Slice) {
	*dest.getOrig() = *ms.getOrig()
	*ms.getOrig() = nil
}

// CopyTo copies UInt64Slice to another instance.
func (ms UInt64Slice) CopyTo(dest UInt64Slice) {
	*dest.getOrig() = copyUInt64Slice(*dest.getOrig(), *ms.getOrig())
}

func copyUInt64Slice(dst, src []uint64) []uint64 {
	dst = dst[:0]
	return append(dst, src...)
}
