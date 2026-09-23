window.BENCHMARK_DATA = {
  "lastUpdate": 1790123481473,
  "repoUrl": "https://github.com/cplieger/ci",
  "entries": {
    "Benchmark": [
      {
        "commit": {
          "author": {
            "name": "Christopher Plieger",
            "username": "cplieger",
            "email": "917744+cplieger@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "9eeb5669fbddd771284ea173a3a244a552cbe22d",
          "message": "chore(sync): synced file(s) with cplieger/ci (#135)\n\nCo-authored-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2026-08-25T08:12:13Z",
          "url": "https://github.com/cplieger/runesafe/commit/9eeb5669fbddd771284ea173a3a244a552cbe22d"
        },
        "date": 1787698272090,
        "tool": "customSmallerIsBetter",
        "benches": [
          {
            "name": "BenchmarkAggregate/join_then_capped - B/op",
            "value": 532896,
            "range": "± 0.0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/join_then_capped - allocs/op",
            "value": 3,
            "range": "± 0.0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/join_then_capped",
            "value": 2598334.5,
            "range": "± 9492.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/streaming - B/op",
            "value": 624,
            "range": "± 0.0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/streaming - allocs/op",
            "value": 3,
            "range": "± 0.0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/streaming",
            "value": 1546,
            "range": "± 7.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkCapBytes - B/op",
            "value": 0,
            "range": "± 0.0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkCapBytes - allocs/op",
            "value": 0,
            "range": "± 0.0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkCapBytes",
            "value": 4.372,
            "range": "± 0.0085",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_256 - B/op",
            "value": 0,
            "range": "± 0.0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_256 - allocs/op",
            "value": 0,
            "range": "± 0.0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_256",
            "value": 809.65,
            "range": "± 2.7",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_4096 - B/op",
            "value": 0,
            "range": "± 0.0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_4096 - allocs/op",
            "value": 0,
            "range": "± 0.0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_4096",
            "value": 12789,
            "range": "± 33.0",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_65536 - B/op",
            "value": 0,
            "range": "± 0.0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_65536 - allocs/op",
            "value": 0,
            "range": "± 0.0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_65536",
            "value": 204667.5,
            "range": "± 456.0",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_utf8/bytes_4096 - B/op",
            "value": 0,
            "range": "± 0.0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_utf8/bytes_4096 - allocs/op",
            "value": 0,
            "range": "± 0.0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_utf8/bytes_4096",
            "value": 9752.5,
            "range": "± 21.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/invalid_utf8/bytes_4096 - B/op",
            "value": 11392,
            "range": "± 0.0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/invalid_utf8/bytes_4096 - allocs/op",
            "value": 2,
            "range": "± 0.0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/invalid_utf8/bytes_4096",
            "value": 28607,
            "range": "± 55.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_4096 - B/op",
            "value": 4864,
            "range": "± 0.0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_4096 - allocs/op",
            "value": 1,
            "range": "± 0.0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_4096",
            "value": 15221.5,
            "range": "± 44.0",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_65536 - B/op",
            "value": 73728,
            "range": "± 0.0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_65536 - allocs/op",
            "value": 1,
            "range": "± 0.0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_65536",
            "value": 242044.5,
            "range": "± 249.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_tail/bytes_4096 - B/op",
            "value": 4864,
            "range": "± 0.0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_tail/bytes_4096 - allocs/op",
            "value": 1,
            "range": "± 0.0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_tail/bytes_4096",
            "value": 14025.5,
            "range": "± 70.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitizeSingleLineBounded - B/op",
            "value": 208,
            "range": "± 0.0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitizeSingleLineBounded - allocs/op",
            "value": 1,
            "range": "± 0.0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitizeSingleLineBounded",
            "value": 12913,
            "range": "± 17.0",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkUntrustedLogValue - B/op",
            "value": 0,
            "range": "± 0.0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkUntrustedLogValue - allocs/op",
            "value": 0,
            "range": "± 0.0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkUntrustedLogValue",
            "value": 12800.5,
            "range": "± 22.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_1048576 - B/op",
            "value": 368,
            "range": "± 0.0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_1048576 - allocs/op",
            "value": 4,
            "range": "± 0.0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_1048576",
            "value": 748.2,
            "range": "± 2.7",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_4096 - B/op",
            "value": 368,
            "range": "± 0.0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_4096 - allocs/op",
            "value": 4,
            "range": "± 0.0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_4096",
            "value": 748.8,
            "range": "± 7.6",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_1048576 - B/op",
            "value": 4046981,
            "range": "± 2.5",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_1048576 - allocs/op",
            "value": 5,
            "range": "± 0.0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_1048576",
            "value": 7816406.5,
            "range": "± 236739.0",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_4096 - B/op",
            "value": 11520,
            "range": "± 0.0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_4096 - allocs/op",
            "value": 4,
            "range": "± 0.0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_4096",
            "value": 27739.5,
            "range": "± 406.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Christopher Plieger",
            "username": "cplieger",
            "email": "917744+cplieger@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "05960b82c7cbaaebc70b4f5e1946afc53266725e",
          "message": "chore(deps): update dependency go to v1.27.1 (#162)",
          "timestamp": "2026-09-01T21:11:18Z",
          "url": "https://github.com/cplieger/runesafe/commit/05960b82c7cbaaebc70b4f5e1946afc53266725e"
        },
        "date": 1788307994787,
        "tool": "customSmallerIsBetter",
        "benches": [
          {
            "name": "BenchmarkAggregate/join_then_capped - B/op",
            "value": 532896,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/join_then_capped - allocs/op",
            "value": 3,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/join_then_capped",
            "value": 1927630,
            "range": "± 6107",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/streaming - B/op",
            "value": 624,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/streaming - allocs/op",
            "value": 3,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/streaming",
            "value": 1346,
            "range": "± 3.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkCapBytes - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkCapBytes - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkCapBytes",
            "value": 4.1465,
            "range": "± 0.0795",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_256 - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_256 - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_256",
            "value": 741.3,
            "range": "± 2.45",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_4096 - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_4096 - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_4096",
            "value": 11762,
            "range": "± 58",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_65536 - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_65536 - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_65536",
            "value": 187908.5,
            "range": "± 1073.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_utf8/bytes_4096 - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_utf8/bytes_4096 - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_utf8/bytes_4096",
            "value": 8335.5,
            "range": "± 39",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/invalid_utf8/bytes_4096 - B/op",
            "value": 11392,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/invalid_utf8/bytes_4096 - allocs/op",
            "value": 2,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/invalid_utf8/bytes_4096",
            "value": 22365.5,
            "range": "± 174",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_4096 - B/op",
            "value": 4864,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_4096 - allocs/op",
            "value": 1,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_4096",
            "value": 12557,
            "range": "± 32",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_65536 - B/op",
            "value": 73728,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_65536 - allocs/op",
            "value": 1,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_65536",
            "value": 202786.5,
            "range": "± 1296.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_tail/bytes_4096 - B/op",
            "value": 4864,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_tail/bytes_4096 - allocs/op",
            "value": 1,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_tail/bytes_4096",
            "value": 12463.5,
            "range": "± 36",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitizeSingleLineBounded - B/op",
            "value": 208,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitizeSingleLineBounded - allocs/op",
            "value": 1,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitizeSingleLineBounded",
            "value": 12102.5,
            "range": "± 74.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkUntrustedLogValue - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkUntrustedLogValue - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkUntrustedLogValue",
            "value": 11788,
            "range": "± 91.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_1048576 - B/op",
            "value": 368,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_1048576 - allocs/op",
            "value": 4,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_1048576",
            "value": 623,
            "range": "± 6.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_4096 - B/op",
            "value": 368,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_4096 - allocs/op",
            "value": 4,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_4096",
            "value": 606.65,
            "range": "± 6.3",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_1048576 - B/op",
            "value": 4046979,
            "range": "± 1",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_1048576 - allocs/op",
            "value": 5,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_1048576",
            "value": 5969197,
            "range": "± 12470",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_4096 - B/op",
            "value": 11520,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_4096 - allocs/op",
            "value": 4,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_4096",
            "value": 22841.5,
            "range": "± 119",
            "unit": "ns/op",
            "extra": "10 samples, median"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Christopher Plieger",
            "username": "cplieger",
            "email": "917744+cplieger@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c7c2a56c10b85a4e0b5b155f46d75b4906546b04",
          "message": "chore(deps): update cplieger/ci digest to a2bb34b (#580)",
          "timestamp": "2026-09-09T00:02:08Z",
          "url": "https://github.com/cplieger/ci/commit/c7c2a56c10b85a4e0b5b155f46d75b4906546b04"
        },
        "date": 1788913091668,
        "tool": "customSmallerIsBetter",
        "benches": [
          {
            "name": "BenchmarkAggregate/join_then_capped - B/op",
            "value": 532896,
            "range": "± 0.5",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/join_then_capped - allocs/op",
            "value": 3,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/join_then_capped",
            "value": 2489632.5,
            "range": "± 16167",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/streaming - B/op",
            "value": 624,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/streaming - allocs/op",
            "value": 3,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/streaming",
            "value": 1630,
            "range": "± 23",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkCapBytes - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkCapBytes - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkCapBytes",
            "value": 4.68,
            "range": "± 0.0085",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_256 - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_256 - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_256",
            "value": 808.8,
            "range": "± 0.95",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_4096 - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_4096 - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_4096",
            "value": 12799.5,
            "range": "± 83.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_65536 - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_65536 - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_65536",
            "value": 204660,
            "range": "± 882.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_utf8/bytes_4096 - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_utf8/bytes_4096 - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_utf8/bytes_4096",
            "value": 9755,
            "range": "± 22.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/invalid_utf8/bytes_4096 - B/op",
            "value": 11392,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/invalid_utf8/bytes_4096 - allocs/op",
            "value": 2,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/invalid_utf8/bytes_4096",
            "value": 27155,
            "range": "± 63.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_4096 - B/op",
            "value": 4864,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_4096 - allocs/op",
            "value": 1,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_4096",
            "value": 15023.5,
            "range": "± 63",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_65536 - B/op",
            "value": 73728,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_65536 - allocs/op",
            "value": 1,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_65536",
            "value": 239552.5,
            "range": "± 249.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_tail/bytes_4096 - B/op",
            "value": 4864,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_tail/bytes_4096 - allocs/op",
            "value": 1,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_tail/bytes_4096",
            "value": 13375,
            "range": "± 41",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitizeSingleLineBounded - B/op",
            "value": 208,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitizeSingleLineBounded - allocs/op",
            "value": 1,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitizeSingleLineBounded",
            "value": 12867,
            "range": "± 18.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkUntrustedLogValue - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkUntrustedLogValue - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkUntrustedLogValue",
            "value": 12785,
            "range": "± 27.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_1048576 - B/op",
            "value": 368,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_1048576 - allocs/op",
            "value": 4,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_1048576",
            "value": 768.05,
            "range": "± 16.15",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_4096 - B/op",
            "value": 368,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_4096 - allocs/op",
            "value": 4,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_4096",
            "value": 747.55,
            "range": "± 13.85",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_1048576 - B/op",
            "value": 4046979,
            "range": "± 1.5",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_1048576 - allocs/op",
            "value": 5,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_1048576",
            "value": 7101424,
            "range": "± 101687",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_4096 - B/op",
            "value": 11520,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_4096 - allocs/op",
            "value": 4,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_4096",
            "value": 27825,
            "range": "± 388.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Christopher Plieger",
            "username": "cplieger",
            "email": "917744+cplieger@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "cbbd76bbe2dfe9a6df87963c66482ecc0de5ece4",
          "message": "chore(sync): synced file(s) with cplieger/ci (#188)\n\nCo-authored-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2026-09-15T11:11:49Z",
          "url": "https://github.com/cplieger/runesafe/commit/cbbd76bbe2dfe9a6df87963c66482ecc0de5ece4"
        },
        "date": 1789518163382,
        "tool": "customSmallerIsBetter",
        "benches": [
          {
            "name": "BenchmarkAggregate/join_then_capped - B/op",
            "value": 532896,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/join_then_capped - allocs/op",
            "value": 3,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/join_then_capped",
            "value": 1929153.5,
            "range": "± 15670.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/streaming - B/op",
            "value": 624,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/streaming - allocs/op",
            "value": 3,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/streaming",
            "value": 1219,
            "range": "± 4.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkCapBytes - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkCapBytes - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkCapBytes",
            "value": 3.825,
            "range": "± 0.0235",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_256 - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_256 - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_256",
            "value": 638.1,
            "range": "± 0.85",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_4096 - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_4096 - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_4096",
            "value": 10087,
            "range": "± 36",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_65536 - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_65536 - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_65536",
            "value": 161219.5,
            "range": "± 603",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_utf8/bytes_4096 - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_utf8/bytes_4096 - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_utf8/bytes_4096",
            "value": 7625,
            "range": "± 103",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/invalid_utf8/bytes_4096 - B/op",
            "value": 11392,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/invalid_utf8/bytes_4096 - allocs/op",
            "value": 2,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/invalid_utf8/bytes_4096",
            "value": 21054.5,
            "range": "± 50.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_4096 - B/op",
            "value": 4864,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_4096 - allocs/op",
            "value": 1,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_4096",
            "value": 11553,
            "range": "± 34",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_65536 - B/op",
            "value": 73728,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_65536 - allocs/op",
            "value": 1,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_65536",
            "value": 182120.5,
            "range": "± 662",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_tail/bytes_4096 - B/op",
            "value": 4864,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_tail/bytes_4096 - allocs/op",
            "value": 1,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_tail/bytes_4096",
            "value": 10827.5,
            "range": "± 152.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitizeSingleLineBounded - B/op",
            "value": 208,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitizeSingleLineBounded - allocs/op",
            "value": 1,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitizeSingleLineBounded",
            "value": 10191,
            "range": "± 170",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkUntrustedLogValue - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkUntrustedLogValue - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkUntrustedLogValue",
            "value": 10085,
            "range": "± 30",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_1048576 - B/op",
            "value": 368,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_1048576 - allocs/op",
            "value": 4,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_1048576",
            "value": 550.05,
            "range": "± 4.55",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_4096 - B/op",
            "value": 368,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_4096 - allocs/op",
            "value": 4,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_4096",
            "value": 551.95,
            "range": "± 1.7",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_1048576 - B/op",
            "value": 4046981,
            "range": "± 1",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_1048576 - allocs/op",
            "value": 5,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_1048576",
            "value": 5562041.5,
            "range": "± 60219",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_4096 - B/op",
            "value": 11520,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_4096 - allocs/op",
            "value": 4,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_4096",
            "value": 20613.5,
            "range": "± 234.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Christopher Plieger",
            "username": "cplieger",
            "email": "917744+cplieger@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "6a4c9a3a84522ecba1d2afc51b00c34c9f119e86",
          "message": "chore(sync): synced file(s) with cplieger/ci (#189)\n\nCo-authored-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2026-09-16T10:15:58Z",
          "url": "https://github.com/cplieger/runesafe/commit/6a4c9a3a84522ecba1d2afc51b00c34c9f119e86"
        },
        "date": 1789568023083,
        "tool": "customSmallerIsBetter",
        "benches": [
          {
            "name": "BenchmarkAggregate/join_then_capped - B/op",
            "value": 532896,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/join_then_capped - allocs/op",
            "value": 3,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/join_then_capped",
            "value": 2573836.5,
            "range": "± 4271",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/streaming - B/op",
            "value": 624,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/streaming - allocs/op",
            "value": 3,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/streaming",
            "value": 1628,
            "range": "± 3",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkCapBytes - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkCapBytes - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkCapBytes",
            "value": 4.3735,
            "range": "± 0.006",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_256 - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_256 - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_256",
            "value": 808.55,
            "range": "± 0.85",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_4096 - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_4096 - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_4096",
            "value": 12804.5,
            "range": "± 90.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_65536 - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_65536 - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_65536",
            "value": 204566.5,
            "range": "± 378.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_utf8/bytes_4096 - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_utf8/bytes_4096 - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_utf8/bytes_4096",
            "value": 9795.5,
            "range": "± 73",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/invalid_utf8/bytes_4096 - B/op",
            "value": 11392,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/invalid_utf8/bytes_4096 - allocs/op",
            "value": 2,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/invalid_utf8/bytes_4096",
            "value": 27497.5,
            "range": "± 76.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_4096 - B/op",
            "value": 4864,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_4096 - allocs/op",
            "value": 1,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_4096",
            "value": 15300.5,
            "range": "± 147",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_65536 - B/op",
            "value": 73728,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_65536 - allocs/op",
            "value": 1,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_65536",
            "value": 250705,
            "range": "± 1787",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_tail/bytes_4096 - B/op",
            "value": 4864,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_tail/bytes_4096 - allocs/op",
            "value": 1,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_tail/bytes_4096",
            "value": 13602,
            "range": "± 784.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitizeSingleLineBounded - B/op",
            "value": 208,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitizeSingleLineBounded - allocs/op",
            "value": 1,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitizeSingleLineBounded",
            "value": 12880.5,
            "range": "± 18.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkUntrustedLogValue - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkUntrustedLogValue - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkUntrustedLogValue",
            "value": 12802.5,
            "range": "± 40",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_1048576 - B/op",
            "value": 368,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_1048576 - allocs/op",
            "value": 4,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_1048576",
            "value": 769.25,
            "range": "± 2.85",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_4096 - B/op",
            "value": 368,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_4096 - allocs/op",
            "value": 4,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_4096",
            "value": 758.8,
            "range": "± 1.35",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_1048576 - B/op",
            "value": 4046980,
            "range": "± 3.5",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_1048576 - allocs/op",
            "value": 5,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_1048576",
            "value": 7924600.5,
            "range": "± 480223.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_4096 - B/op",
            "value": 11520,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_4096 - allocs/op",
            "value": 4,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_4096",
            "value": 27775.5,
            "range": "± 730",
            "unit": "ns/op",
            "extra": "10 samples, median"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Christopher Plieger",
            "username": "cplieger",
            "email": "917744+cplieger@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "f9577db6c2f2096d9cc325c89450a48686d66346",
          "message": "chore(deps): update cplieger/ci digest to aa0a018 (#649)",
          "timestamp": "2026-09-20T08:02:03Z",
          "url": "https://github.com/cplieger/ci/commit/f9577db6c2f2096d9cc325c89450a48686d66346"
        },
        "date": 1790123480789,
        "tool": "customSmallerIsBetter",
        "benches": [
          {
            "name": "BenchmarkAggregate/join_then_capped - B/op",
            "value": 532896,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/join_then_capped - allocs/op",
            "value": 3,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/join_then_capped",
            "value": 2689388.5,
            "range": "± 5655.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/streaming - B/op",
            "value": 624,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/streaming - allocs/op",
            "value": 3,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkAggregate/streaming",
            "value": 1719,
            "range": "± 53",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkCapBytes - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkCapBytes - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkCapBytes",
            "value": 4.375,
            "range": "± 0.009",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_256 - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_256 - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_256",
            "value": 811.15,
            "range": "± 7",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_4096 - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_4096 - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_4096",
            "value": 12858,
            "range": "± 143",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_65536 - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_65536 - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_ascii/bytes_65536",
            "value": 205073,
            "range": "± 1143",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_utf8/bytes_4096 - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_utf8/bytes_4096 - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/clean_utf8/bytes_4096",
            "value": 9784,
            "range": "± 65",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/invalid_utf8/bytes_4096 - B/op",
            "value": 11392,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/invalid_utf8/bytes_4096 - allocs/op",
            "value": 2,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/invalid_utf8/bytes_4096",
            "value": 29134,
            "range": "± 249.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_4096 - B/op",
            "value": 4864,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_4096 - allocs/op",
            "value": 1,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_4096",
            "value": 15973.5,
            "range": "± 529.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_65536 - B/op",
            "value": 73728,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_65536 - allocs/op",
            "value": 1,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_dense/bytes_65536",
            "value": 254680,
            "range": "± 730.5",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_tail/bytes_4096 - B/op",
            "value": 4864,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_tail/bytes_4096 - allocs/op",
            "value": 1,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitize/unsafe_tail/bytes_4096",
            "value": 14369,
            "range": "± 449",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitizeSingleLineBounded - B/op",
            "value": 208,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitizeSingleLineBounded - allocs/op",
            "value": 1,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkSanitizeSingleLineBounded",
            "value": 12930.5,
            "range": "± 20",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkUntrustedLogValue - B/op",
            "value": 0,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkUntrustedLogValue - allocs/op",
            "value": 0,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkUntrustedLogValue",
            "value": 12812,
            "range": "± 74",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_1048576 - B/op",
            "value": 368,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_1048576 - allocs/op",
            "value": 4,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_1048576",
            "value": 809.8,
            "range": "± 2.65",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_4096 - B/op",
            "value": 368,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_4096 - allocs/op",
            "value": 4,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/budgeted/bytes_4096",
            "value": 799.3,
            "range": "± 20.95",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_1048576 - B/op",
            "value": 4046981.5,
            "range": "± 2",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_1048576 - allocs/op",
            "value": 5,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_1048576",
            "value": 8328991,
            "range": "± 370554",
            "unit": "ns/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_4096 - B/op",
            "value": 11520,
            "range": "± 0",
            "unit": "B/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_4096 - allocs/op",
            "value": 4,
            "range": "± 0",
            "unit": "allocs/op",
            "extra": "10 samples, median"
          },
          {
            "name": "BenchmarkWorkBound/capped/bytes_4096",
            "value": 27889.5,
            "range": "± 258",
            "unit": "ns/op",
            "extra": "10 samples, median"
          }
        ]
      }
    ]
  }
}