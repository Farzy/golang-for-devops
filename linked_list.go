package main

import "fmt"

type Node struct {
	data int
	next *Node
}

type LinkedList struct {
	head *Node
}

func (ll *LinkedList) append(data int) {
	newNode := &Node{
		data: data,
		next: nil,
	}

	if ll.head == nil {
		ll.head = newNode
		return
	}

	current := ll.head
	for current.next != nil {
		current = current.next
	}
	current.next = newNode
}

func (ll *LinkedList) display() {
	current := ll.head
	for current != nil {
		fmt.Printf("%d -> ", current.data)
		current = current.next
	}
	fmt.Println("nil")
}

func testLinkedList() {
	ll := LinkedList{}

	ll.append(43)
	ll.append(-10)
	ll.append(123)

	ll.display()
}
