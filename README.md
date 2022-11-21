# Queue
简单实现的无锁队列, 它现在还不能使用, 因为`race delector`报告的`data race`还未解决

```shell
goos: windows
goarch: amd64
pkg: example/lock-free
cpu: Intel(R) Core(TM) i7-8705G CPU @ 3.10GHz
BenchmarkQueue
BenchmarkQueue/Concurrent
BenchmarkQueue/Concurrent-8         	 4913110	       221.1 ns/op	      16 B/op	       1 allocs/op
BenchmarkQueue/NoConcurrent
BenchmarkQueue/NoConcurrent-8       	24479953	        88.32 ns/op	      16 B/op	       1 allocs/op
```