package main

import "fmt"

//create a type that:
/*
	holds data
	holds pointer tp next node
*/
type node struct {
	data int
	next *node
}

/*
	 contains the many nodes
	 requires at lest a head node
*/
type LinkedList struct {
	head *node   //needs to hold the head at least
	length int
}

//creates a method for linked list ds
/*
	we describe the type it belongs to before name of method ie: methofd receiver
*/
func (l *LinkedList) prepend(n *node) {
	/*
		first we save the current head node
		we then change LL head to n our passed node
		then we set/point our saved original head node as next
		then add the length of our LL
	*/
	second := l.head
	l.head = n
	l.head.next = second
	l.length++
}

func (l LinkedList) printListData() {
	toPrint := l.head
	for range l.length{
		fmt.Printf("%d\n", toPrint.data)
		toPrint =toPrint.next
	}
}
func main(){
	mylist := LinkedList {}
	node1 := &node{data:48}
	node2 := &node{data:24}
	node3 := &node{data:16}
	node4 := &node{data:12}
	node5 := &node{data:8}
	node6 := &node{data:5}


	mylist.prepend(node1)
	mylist.prepend(node2)
	mylist.prepend(node3)
	mylist.prepend(node4)
	mylist.prepend(node5)
	mylist.prepend(node6)

	mylist.printListData()
}