package lock_free

import "testing"

func TestQueue(t *testing.T) {
	queue := NewQueue[int]()
	for i := 0; i < 1000; i++ {
		queue.Push(i)
	}
	for i := 0; i < 1000; i++ {
		t.Log(queue.Pop())
	}
	queue.Pop()
	queue.Push(100)
}

func BenchmarkQueue(b *testing.B) {
	queue := NewQueue[int]()
	b.Run("Concurrent", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				queue.Push(100)
			}
		})
	})
	b.Run("NoConcurrent", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			queue.Push(100)
		}
	})
}
