package types

import "os"

type FileDataBase struct {
	File *os.File
	Header FileHeader
	Keys []Key
	Values []Value
}

type FileHeader struct {
	MagicNubmer uint32
	Version uint16
	KeysCount uint32
	ValuesCount uint32
}

type Key struct {
	Vector Vector
	Index uint32
	Wehgit float64
	Load uint32
}

type Value struct {
	stringLen uint32
	Value string
}