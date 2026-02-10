1. Бинарный поиск при вставке группы или тега, бинарный поиск тегов по группам в writeTag
BenchmarkLoggerGroups-12   11483868	107.0 ns/op      96 B/op       2 allocs/op
BenchmarkCtxGroups-12      7030424	184.7 ns/op      72 B/op       2 allocs/op
BenchmarkWriteTags-12      8216127	154.4 ns/op     130 B/op       0 allocs/op

2 maps в тегах и группах
BenchmarkLoggerGroups-12   68840890	 17.80 ns/op       0 B/op       0 allocs/op
BenchmarkCtxGroups-12      11911579	 98.34 ns/op      48 B/op       1 allocs/op
BenchmarkWriteTags-12       4614613	 246.6 ns/op     116 B/op       0 allocs/op

3 как было изначально - добавление через slices.Index и writeTag через slices.Index
BenchmarkLoggerGroups-12   11681947      96.80 ns/op      96 B/op       2 allocs/op
BenchmarkCtxGroups-12       7727312      137.9 ns/op      72 B/op       2 allocs/op
BenchmarkWriteTags-12      10563631      115.1 ns/op     101 B/op       0 allocs/op