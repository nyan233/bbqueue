package lock_free

import (
	"sync/atomic"
	"unsafe"
)

type node[Value any] struct {
	next *node[Value]
	val  Value
}

type LFQueue[Value any] struct {
	prev   *node[Value]
	tail   *node[Value]
	length int64
}

func NewQueue[Value any]() *LFQueue[Value] {
	return &LFQueue[Value]{}
}

func (q *LFQueue[Value]) Push(v Value) {
	newNode := &node[Value]{
		val: v,
	}
	if q.Len() == 0 && q.cas(&q.prev, nil, newNode) {
		q.inc()
		q.swap(&q.tail, newNode)
		return
	}
	for {
		if q.load(&q.tail) == nil {
			continue
		}
		swapped := q.cas(&q.tail.next, nil, newNode)
		if swapped {
			q.inc()
			q.store(&q.tail, q.tail.next)
			return
		}
	}
}

func (q *LFQueue[Value]) Pop() (Value, bool) {
	if q.load(&q.prev) == nil {
		return *new(Value), false
	}
	for {
		old := q.load(&q.prev)
		if old == nil {
			return *new(Value), false
		}
		if q.cas(&q.prev, old, old.next) {
			q.dec()
			return old.val, true
		}
	}
}

func (q *LFQueue[Value]) load(rootNode **node[Value]) *node[Value] {
	return (*node[Value])(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(rootNode))))
}

func (q *LFQueue[Value]) store(rootNode **node[Value], node *node[Value]) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(rootNode)), unsafe.Pointer(node))
}

func (q *LFQueue[Value]) cas(rootNode **node[Value], old, new *node[Value]) bool {
	return atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(rootNode)), unsafe.Pointer(old), unsafe.Pointer(new))
}

func (q *LFQueue[Value]) swap(rootNode **node[Value], new *node[Value]) (old *node[Value]) {
	return (*node[Value])(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(rootNode)), unsafe.Pointer(new)))
}

func (q *LFQueue[Value]) Len() int {
	return int(atomic.LoadInt64(&q.length))
}

func (q *LFQueue[Value]) inc() {
	atomic.AddInt64(&q.length, 1)
}

func (q *LFQueue[Value]) dec() {
	atomic.AddInt64(&q.length, -1)
}
