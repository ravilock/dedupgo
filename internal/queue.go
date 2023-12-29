package internal

type node[T any] struct {
	value T
	depth int
	next  *node[T]
}

func (n *node[T]) setNext(node *node[T]) {
	n.next = node
}

type Queue[T any] struct {
	first  *node[T]
	last   *node[T]
	length int
}

func (q *Queue[T]) Length() int {
	return q.length
}

func (q *Queue[T]) Enqueue(data T, depth int) {
	node := &node[T]{data, depth, nil}
	if q.length == 0 {
		q.first = node
		q.last = node
		q.length += 1
		return
	}
	q.length += 1
	q.last.setNext(node)
	q.last = node
}

func (q *Queue[T]) Dequeue() (T, int) {
	var x T
	if q.length == 0 {
		return x, 0
	}
	previousFirst := q.first
	if q.length == 1 {
		q.first = nil
		q.last = nil
		q.length = 0
		return previousFirst.value, previousFirst.depth
	}
	q.length -= 1
	q.first = q.first.next
	if q.length == 1 {
		q.last = q.first
	}
	return previousFirst.value, previousFirst.depth
}
