automerge-benchmark
----------------------------------------

`automerge-benchmark` provides a benchmark for automerge.  

Separated from `automerge` library to avoid burdening `automerge` with the size
of the data set files.

See [https://github.com/savaki/automerge](https://github.com/savaki/automerge).

```text
=== RUN   TestPerformance
applying edits ...
 25000:  62388 bytes, 3.0 µs/op
 50000: 126464 bytes, 3.1 µs/op
 75000: 189055 bytes, 4.3 µs/op
100000: 252513 bytes, 3.6 µs/op
125000: 314784 bytes, 2.8 µs/op
150000: 376359 bytes, 7.1 µs/op
175000: 439171 bytes, 3.5 µs/op
200000: 502771 bytes, 6.7 µs/op
225000: 569985 bytes, 7.5 µs/op
250000: 634307 bytes, 4.1 µs/op

edits -> 259778
bytes -> 659432
--- PASS: TestPerformance (1.52s)
PASS
```
