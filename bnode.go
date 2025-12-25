package main

import (
	"encoding/binary"
)

const HEADER = 4

const (
	BNODE_NODE = 1 // internal node
	BNODE_LEAF = 2 // leaf node
)

type BNode []byte

func (node BNode) bType() uint16 {
	return binary.LittleEndian.Uint16(node[0:2])
}

func (node BNode) nKeys() uint16 {
	return binary.LittleEndian.Uint16(node[2:4])
}

func (node BNode) setHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node[0:2], btype)
	binary.LittleEndian.PutUint16(node[2:4], nkeys)
}

func (node BNode) getPtr(idx uint16) uint64 {
	if idx >= node.nKeys() {
		panic("btree: getPtr index out of range")
	}

	pos := HEADER + 8*idx
	return binary.LittleEndian.Uint64(node[pos:])
}

func (node BNode) setPtr(idx uint16, ptr uint64) {
	if idx > node.nKeys() {
		panic("btree: setPtr index out of range")
	}

	pos := HEADER + 8*idx
	binary.LittleEndian.PutUint64(node[pos:], ptr)
}
