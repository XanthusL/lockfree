package list

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

type (
	node struct {
		value *atomic.Value
		next  unsafe.Pointer
		prev  unsafe.Pointer
	}
	// LinkedList ...
	LinkedList struct {
		root unsafe.Pointer
	}
)

// NewList ..
func NewList() *LinkedList {
	l := new(LinkedList)
	l.Init()
	return l
}

// Init ..
func (l *LinkedList) Init() {
	if l.root != nil {
		return
	}
	l.root = unsafe.Pointer(createElem(struct{}{}))
}

// Walk ...
func (l *LinkedList) Walk(f func(v interface{})) {
	hd := l.getRoot().next
	for hd != nil && hd != l.root {
		ele := (*node)(hd)
		f(ele.value.Load())
		hd = ele.next
	}
}

func (l *LinkedList) getRoot() *node {
	return (*node)(l.root)
}

// Add first element to the list
// elem.prev = root
// elem.next = root
// root.prev = elem
// root.next = elem
func (l *LinkedList) addFirst(n *node) bool {
	if !atomic.CompareAndSwapPointer(&n.next, nil, l.root) {
		return false
	}
	if !atomic.CompareAndSwapPointer(&n.prev, nil, l.root) {
		return false
	}
	elem := unsafe.Pointer(n)
	root := l.getRoot()
	ok := atomic.CompareAndSwapPointer(&root.next, nil, elem)
	if !ok {
		return false
	}
	return atomic.CompareAndSwapPointer(&root.prev, nil, elem)
}

func createElem(v interface{}) *node {
	n := &node{
		value: new(atomic.Value),
	}
	n.value.Store(v)
	return n
}

// PushFront TODO
func (l *LinkedList) PushFront(v interface{}) {

}

// PushBack appends an element to the tail
// elem.prev = tail
// elem.next = root
// tail.next = elem
// root.prev = elem
func (l *LinkedList) PushBack(v interface{}) {
	for {
		if l.insert((*node)(l.root).prev,
			l.root, v) {
			break
		}
		runtime.Gosched()
	}
}

func (l *LinkedList) insert(prev, next unsafe.Pointer, v interface{}) bool {
	n := createElem(v)
	if prev == nil && (*node)(next).prev == nil ||
		next == nil && (*node)(prev).next == nil {
		return l.addFirst(n)
	}
	if prev == nil || next == nil {
		return false
	}
	elem := unsafe.Pointer(n)
	if !atomic.CompareAndSwapPointer(
		&(((*node)(prev)).next),
		next, elem) {
		return false
	}
	if !atomic.CompareAndSwapPointer(&n.prev, nil, prev) {
		return false
	}
	if !atomic.CompareAndSwapPointer(&n.next, nil, next) {
		return false
	}
	if !atomic.CompareAndSwapPointer(&((*node)(next).prev), prev, elem) {
		return false
	}
	return true
}
