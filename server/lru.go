package server


type ListNode struct {
	Left *ListNode
	Right *ListNode
	Key string
}


type Element struct {
	P *ListNode
	Val int
	Expire int64
}

type LRUCache struct {
	Cap int
	Len int
	Data map[string]Element
	Head *ListNode
	Tail *ListNode
}


func (lru *LRUCache) Get(key string) int {
	if v, ok := lru.Data[key]; !ok {
		return -1
	} else {
		left, right := v.P.Left, v.P.Right
		thenode := v.P
		left.Right, right.Left = right, left
		first := lru.Head.Right
		thenode.Left, thenode.Right = lru.Head, first
		lru.Head.Right, first.Left = thenode, thenode
		return v.Val
	}

}


func (lru *LRUCache) Put(key string, value int, expire int64) {
	if v, ok := lru.Data[key]; ok {
		thenode := v.P
		left, right := thenode.Left, thenode.Right
		left.Right, right.Left = left, left
		delete(lru.Data, key)
		lru.Len --
	}

	if lru.Cap == lru.Len {
		last := lru.Tail.Left
		left := last.Left
		left.Right = lru.Tail
		lru.Tail.Left = left
		delete(lru.Data, last.Key)
		lru.Len --
	}

	pnode := &ListNode{Key: key}
	first := lru.Head.Right
	pnode.Left = lru.Head
	pnode.Right = first
	lru.Head.Right, first.Left = pnode, pnode
	element := Element{P: pnode, Val: value, Expire: expire}
	lru.Data[key] = element
	lru.Len ++
}


func NewLRUCache(Cap int) *LRUCache {
	if Cap <= 0 {
		return nil
	}
	cache := &LRUCache{}
	cache.Cap = Cap
	cache.Len = 0
	cache.Data = make(map[string]Element)
	head, tail := &ListNode{}, &ListNode{}
	head.Right, tail.Left = tail, head
	cache.Head = head
	cache.Tail = tail
	return cache
}