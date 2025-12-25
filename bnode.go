package main

import (
	"bytes"
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

func offsetPos(node BNode, idx uint16) uint16 {
	if 1 > idx || idx > node.nKeys() {
		panic("btree: offsetPos index out of range")
	}

	return HEADER + 8*node.nKeys() + 2*(idx-1)
}

func (node BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}

	pos := offsetPos(node, idx)
	return binary.LittleEndian.Uint16(node[pos:])
}

func (node BNode) setOffset(idx uint16, offset uint16) {
	if idx == 0 {
		panic("btree: setOffset index out of range")
	}

	pos := offsetPos(node, idx)
	binary.LittleEndian.PutUint16(node[pos:], offset)
}

func (node BNode) kvPos(idx uint16) uint16 {
	if idx > node.nKeys() {
		panic("btree: kvPos index out of range")
	}

	return HEADER + 8*node.nKeys() + 2*node.nKeys() + node.getOffset(idx)
}

func (node BNode) getKey(idx uint16) []byte {
	if idx >= node.nKeys() {
		panic("btree: getKey index out of range")
	}

	pos := node.kvPos(idx)
	keyLength := binary.LittleEndian.Uint16(node[pos:])
	return node[pos+4:][:keyLength]
}

func (node BNode) getVal(idx uint16) []byte {
	if idx >= node.nKeys() {
		panic("btree: getVal index out of range")
	}

	pos := node.kvPos(idx)
	keyLength := binary.LittleEndian.Uint16(node[pos:])
	valLength := binary.LittleEndian.Uint16(node[pos+2:])
	return node[pos+4+keyLength:][:valLength]
}

func (node BNode) nbytes() uint16 {
	return node.kvPos(node.nKeys())
}

// TODO: implement binary search
func nodeLookupLE(node BNode, key []byte) uint16 {
	nkeys := node.nKeys()
	found := uint16(0)

	for i := uint16(1); i < nkeys; i++ {
		cmp := bytes.Compare(node.getKey(i), key)

		if cmp <= 0 {
			found = i
		}

		if cmp >= 0 {
			break
		}
	}

	return found
}

func leafInsert(new BNode, old BNode, idx uint16, key []byte, val []byte) {
	new.setHeader(BNODE_LEAF, old.nKeys()+1) // Incrementing key count
	nodeAppendRange(new, old, 0, 0, idx)     // FIrst append the
	nodeAppendKV(new, idx, 0, key, val)
	nodeAppendRange(new, old, idx+1, idx, old.nKeys()-idx)
}

func nodeAppendKV(new BNode, idx uint16, ptr uint64, key []byte, val []byte) {
	new.setPtr(idx, ptr)

	pos := new.kvPos(idx)
	binary.LittleEndian.PutUint16(new[pos:], uint16(len(key)))
	binary.LittleEndian.PutUint16(new[pos+2:], uint16(len(val)))
	copy(new[pos+4:], key)
	copy(new[pos+4+uint16(len(key)):], val)

	new.setOffset(idx+1, new.getOffset(idx)+4+uint16(len(key)+len(val)))
}

func nodeAppendRange(new BNode, old BNode, dstNew uint16, srcOld uint16, n uint16) {
	for i := uint16(0); i < n; i++ {
		ptr := old.getPtr(srcOld + i)
		new.setPtr(dstNew+i, ptr)
	}
}
