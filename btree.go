package main

const BTREE_PAGE_SIZE = 4096
const BTREE_MAX_KEY_SIZE = 1000
const BTREE_MAX_VAL_SIZE = 3000

type BTree struct {
	root uint64              // page number pointer
	get  func(uint64) []byte // dereference pointer to page
	new  func([]byte) uint64 // allocate new page
	del  func(uint64)        // deallocate page
}
