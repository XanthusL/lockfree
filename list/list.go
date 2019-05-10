package list

import (
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
	n := createElem(v)
	root := l.getRoot()
	for {
		if root.next == nil {
			ok := l.addFirst(n)
			if ok {
				break
			}
			continue
		}
		tail := root.prev
		elem := unsafe.Pointer(n)
		if !atomic.CompareAndSwapPointer(
			&(((*node)(tail)).next),
			l.root, elem) {
			continue
		}
		if !atomic.CompareAndSwapPointer(&n.prev, nil, tail) {
			continue
		}
		if !atomic.CompareAndSwapPointer(&n.next, nil, l.root) {
			continue
		}
		if !atomic.CompareAndSwapPointer(&root.prev, tail, elem) {
			continue
		}
		break
	}
}
