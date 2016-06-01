package main

// DO NOT CHANGE THIS CACHE SIZE VALUE
const CACHE_SIZE int = 3

var lcache = &Lrucache{m: make(map[int]*Node), head: nil, end: nil}

func Set(key int, value int) {
	temp := &Node{value: value, key: key}
	if val, ok := lcache.m[key]; ok {
		temp = val
		temp.value = value
		remove(temp)
		setHead(temp)
	} else {
		temp.key = key
		temp.value = value
		temp.pre = nil
		temp.next = nil
		if len(lcache.m) >= CACHE_SIZE {
			delete(lcache.m, lcache.end.key)
			remove(lcache.end)
		}
		setHead(temp)
		lcache.m[key] = temp
	}
}

func Get(key int) int {
	if val, ok := lcache.m[key]; ok {
		var temp *Node = val
		remove(temp)
		setHead(temp)
		return temp.value
	}
	return -1
}

type Lrucache struct {
	head *Node
	end  *Node
	m    map[int]*Node
}
type Node struct {
	pre   *Node
	next  *Node
	key   int
	value int
}

func remove(temp *Node) {
	var a *Node = temp.pre
	if a != nil {
		a.next = temp.next
	} else {
		lcache.head = temp.next
	}
	if temp.next != nil {
		temp.next.pre = a
	} else {
		lcache.end = temp.pre
	}
}
func setHead(temp *Node) {
	if lcache.head == nil {
		lcache.head = temp
	} else {
		temp.next = lcache.head
		lcache.head.pre = temp
		lcache.head = temp
	}

	if lcache.end == nil {
		lcache.end = lcache.head
	}
}
