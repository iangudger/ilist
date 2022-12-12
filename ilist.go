// Copyright 2018 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package ilist provides the implementation of intrusive linked lists.
package ilist

// Linker is the interface that objects must implement if they want to be added
// to and/or removed from List objects.
type Linker[T any] interface {
	Next() T
	Prev() T
	SetNext(T)
	SetPrev(T)
}

// Element the item that is used at the API level.
type Element[T any] interface {
	*T
	Linker[*T]
}

// List is an intrusive list. Entries can be added to or removed from the list
// in O(1) time and with no additional memory allocations.
//
// The zero value for List is an empty list ready to use.
//
// To iterate over a list (where l is a List):
//
//	for e := l.Front(); e != nil; e = e.Next() {
//		// do something with e.
//	}
type List[T any, U Element[T]] struct {
	head *T
	tail *T
}

// Reset resets list l to the empty state.
func (l *List[T, U]) Reset() {
	l.head = nil
	l.tail = nil
}

// Empty returns true iff the list is empty.
//
//go:nosplit
func (l *List[T, U]) Empty() bool {
	return l.head == nil
}

// Front returns the first element of list l or nil.
//
//go:nosplit
func (l *List[T, U]) Front() *T {
	return l.head
}

// Back returns the last element of list l or nil.
//
//go:nosplit
func (l *List[T, U]) Back() *T {
	return l.tail
}

// Len returns the number of elements in the list.
//
// NOTE: This is an O(n) operation.
//
//go:nosplit
func (l *List[T, U]) Len() (count int) {
	for e := l.Front(); e != nil; e = U(e).Next() {
		count++
	}
	return count
}

// PushFront inserts the element e at the front of list l.
//
//go:nosplit
func (l *List[T, U]) PushFront(e *T) {
	U(e).SetNext(l.head)
	U(e).SetPrev(nil)
	if l.head != nil {
		U(l.head).SetPrev(e)
	} else {
		l.tail = e
	}

	l.head = e
}

// PushFrontList inserts list m at the start of list l, emptying m.
//
//go:nosplit
func (l *List[T, U]) PushFrontList(m *List[T, U]) {
	if l.head == nil {
		l.head = m.head
		l.tail = m.tail
	} else if m.head != nil {
		U(l.head).SetPrev(m.tail)
		U(m.tail).SetNext(l.head)

		l.head = m.head
	}
	m.head = nil
	m.tail = nil
}

// PushBack inserts the element e at the back of list l.
//
//go:nosplit
func (l *List[T, U]) PushBack(e *T) {
	U(e).SetNext(nil)
	U(e).SetPrev(l.tail)
	if l.tail != nil {
		U(l.tail).SetNext(e)
	} else {
		l.head = e
	}

	l.tail = e
}

// PushBackList inserts list m at the end of list l, emptying m.
//
//go:nosplit
func (l *List[T, U]) PushBackList(m *List[T, U]) {
	if l.head == nil {
		l.head = m.head
		l.tail = m.tail
	} else if m.head != nil {
		U(l.tail).SetNext(m.head)
		U(m.head).SetPrev(l.tail)

		l.tail = m.tail
	}
	m.head = nil
	m.tail = nil
}

// InsertAfter inserts e after b.
//
//go:nosplit
func (l *List[T, U]) InsertAfter(b, e *T) {
	a := U(b).Next()

	U(e).SetNext(a)
	U(e).SetPrev(b)
	U(b).SetNext(e)

	if a != nil {
		U(a).SetPrev(e)
	} else {
		l.tail = e
	}
}

// InsertBefore inserts e before a.
//
//go:nosplit
func (l *List[T, U]) InsertBefore(a, e *T) {
	b := U(a).Prev()
	U(e).SetNext(a)
	U(e).SetPrev(b)
	U(a).SetPrev(e)

	if b != nil {
		U(b).SetNext(e)
	} else {
		l.head = e
	}
}

// Remove removes e from l.
//
//go:nosplit
func (l *List[T, U]) Remove(e *T) {
	prev := U(e).Prev()
	next := U(e).Next()

	if prev != nil {
		U(prev).SetNext(next)
	} else if l.head == e {
		l.head = next
	}

	if next != nil {
		U(next).SetPrev(prev)
	} else if l.tail == e {
		l.tail = prev
	}

	U(e).SetNext(nil)
	U(e).SetPrev(nil)
}

// Entry is a default implementation of Linker. Users can add anonymous fields
// of this type to their structs to make them automatically implement the
// methods needed by List.
type Entry[T any, U Element[T]] struct {
	next *T
	prev *T
}

// Next returns the entry that follows e in the list.
//
//go:nosplit
func (e *Entry[T, U]) Next() *T {
	return e.next
}

// Prev returns the entry that precedes e in the list.
//
//go:nosplit
func (e *Entry[T, U]) Prev() *T {
	return e.prev
}

// SetNext assigns 'entry' as the entry that follows e in the list.
//
//go:nosplit
func (e *Entry[T, U]) SetNext(elem *T) {
	e.next = elem
}

// SetPrev assigns 'entry' as the entry that precedes e in the list.
//
//go:nosplit
func (e *Entry[T, U]) SetPrev(elem *T) {
	e.prev = elem
}
