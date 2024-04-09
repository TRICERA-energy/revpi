export CFLAGS="-I$(pwd)/interface"
go tool cgo -godefs=true -- $CFLAGS ctypes_linux.go >types_linux.go
