package internal

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueue(t *testing.T) {
	t.Run("Should be able to enqueue when empty", func(t *testing.T) {
		// Arrange
		expectedLength := 1
		expecedValue := "first"
		queue := new(Queue[string])

		// Act
		queue.Enqueue(expecedValue)

		// Assert
		firstNode := queue.first
		lastNode := queue.last
		require.Equal(t, expectedLength, queue.length)
		require.Equal(t, expecedValue, firstNode.value)
		require.Equal(t, expecedValue, lastNode.value)
		require.Equal(t, firstNode, lastNode)
		require.Nil(t, firstNode.next)
		require.Nil(t, lastNode.next)
	})

	t.Run("Should be able to enqueue when not empty", func(t *testing.T) {
		// Arrange
		expectedLength := 2
		expecedValue := "second"
		queue := new(Queue[string])

		// Act
		queue.Enqueue("first")
		queue.Enqueue(expecedValue)

		// Assert
		lastNode := queue.last
		firstNode := queue.first
		require.Equal(t, expectedLength, queue.length)
		require.Equal(t, expecedValue, lastNode.value)
		require.NotEqual(t, firstNode, lastNode)
		require.Equal(t, firstNode.next, lastNode)
		require.Nil(t, lastNode.next)
	})

	t.Run("Should return queue's type zero value and not ok when dequeuing", func(t *testing.T) {
		// Arrange
		expectedLength := 0
		queue := new(Queue[string])

		// Act
		value, ok := queue.Dequeue()

		// Assert
		firstNode := queue.first
		lastNode := queue.last
		require.Zero(t, value)
		require.False(t, ok)
		require.Equal(t, expectedLength, queue.length)
		require.Nil(t, firstNode)
		require.Nil(t, lastNode)
	})

	t.Run("Should properly dequeue and clean queue if only one value is queued", func(t *testing.T) {
		// Arrange
		expectedLength := 0
		expecedValue := "first"
		queue := new(Queue[string])

		// Act
		queue.Enqueue(expecedValue)
		value, ok := queue.Dequeue()

		// Assert
		firstNode := queue.first
		lastNode := queue.last
		require.Equal(t, expectedLength, queue.length)
		require.Equal(t, expecedValue, value)
		require.True(t, ok)
		require.Nil(t, firstNode)
		require.Nil(t, lastNode)
	})

	t.Run("Should be able to dequeue when length is bigger than 1", func(t *testing.T) {
		// Arrange
		expectedLength := 1
		expectedDequeued := "first"
		exepctedRemainder := "second"
		queue := new(Queue[string])

		// Act
		queue.Enqueue(expectedDequeued)
		queue.Enqueue(exepctedRemainder)
		dequeued, ok := queue.Dequeue()

		// Assert
		firstNode := queue.first
		lastNode := queue.last
		require.Equal(t, expectedLength, queue.length)
		require.Equal(t, expectedDequeued, dequeued)
		require.True(t, ok)
		require.Equal(t, exepctedRemainder, firstNode.value)
		require.Equal(t, exepctedRemainder, lastNode.value)
		require.Equal(t, firstNode, lastNode)
		require.Nil(t, firstNode.next)
		require.Nil(t, lastNode.next)
	})

	t.Run("Should be FIFO", func(t *testing.T) {
		// Arrange
		expectedLength := 0
		expectedFirst := "first"
		exepctedSecond := "second"
		queue := new(Queue[string])

		// Act
		queue.Enqueue(expectedFirst)
		queue.Enqueue(exepctedSecond)
		firstDequeued, ok1 := queue.Dequeue()
		secondDequeued, ok2 := queue.Dequeue()

		// Assert
		firstNode := queue.first
		lastNode := queue.last
		require.Equal(t, expectedLength, queue.length)
		require.Equal(t, expectedFirst, firstDequeued)
		require.Equal(t, exepctedSecond, secondDequeued)
		require.True(t, ok1)
		require.True(t, ok2)
		require.Nil(t, firstNode)
		require.Nil(t, lastNode)
	})
}
