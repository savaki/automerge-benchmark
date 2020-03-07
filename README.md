automerge-benchmark
----------------------------------------

`automerge-benchmark` provides a benchmark for automerge.  

Separated from `automerge` library to avoid burdening `automerge` with the size
of the data set files.

See [https://github.com/savaki/automerge](https://github.com/savaki/automerge).

```text
applying edits ...
 25000:  55640 bytes, 7.5 µs/op
 50000: 112304 bytes, 7.4 µs/op
 75000: 168763 bytes, 8.1 µs/op
100000: 225596 bytes, 8.6 µs/op
125000: 281101 bytes, 8.5 µs/op
150000: 337149 bytes, 8.4 µs/op
175000: 393577 bytes, 9.3 µs/op
200000: 450640 bytes, 9.2 µs/op
225000: 508737 bytes, 10.1 µs/op
250000: 565644 bytes, 10.4 µs/op

edits:    259778
bytes:    588102
elapsed:  2.287s
```

