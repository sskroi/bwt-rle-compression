## Implementation of the compression algorithm using Burrows–Wheeler transform and Run-length encoding.
### Tests
You can check the effectiveness of the algorithm by placing the test files for compression and decompression in the `./internal/compression/test_files` directory and running `go test ./internal/compression -v`

Some examples:
```
$ go test ./internal/compression -v
=== RUN   TestCompression
"a.out" - test file
15.61 KiB - size before compress
7 ms - time to compress
4.18 KiB - size after compress
1 ms - time to decompress

"telegram-desktop-5.1.7-3-x86_64.pkg.tar" - test file
99.32 MiB - size before compress
3958 ms - time to compress
65.71 MiB - size after compress
2691 ms - time to decompress

"neovim-0.10.0-5-x86_64.pkg.tar" - test file
30.03 MiB - size before compress
1198 ms - time to compress
18.58 MiB - size after compress
825 ms - time to decompress

"python-3.12.4-1-x86_64.pkg.tar" - test file
72.73 MiB - size before compress
3201 ms - time to compress
42.48 MiB - size after compress
1871 ms - time to decompress

"qtcreator-13.0.2-2-x86_64.pkg.tar" - test file
119.78 MiB - size before compress
4659 ms - time to compress
78.41 MiB - size after compress
3161 ms - time to decompress
```
