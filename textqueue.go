package main

import "container/list"

// TextQueue is a specialisation of a List for implementing an elastic sliding
// text window with all the casting gunk hidden away
type TextQueue struct {
	Data *list.List
}

// NewTextQueue constructs an empty new TextQueue
func NewTextQueue() *TextQueue {
	return &TextQueue{Data: list.New()}
}

// AddFront adds an element to the front of the queue
func (tq *TextQueue) AddFront(s string) {
	tq.Data.PushFront(s)
}

// RemoveBack discards the rearmost element, if any
func (tq *TextQueue) RemoveBack() {
	if back := tq.Data.Back(); back != nil {
		tq.Data.Remove(back)
	}
}

// Len returns the number of items currently in the TextQueue
func (tq *TextQueue) Len() int {
	return tq.Data.Len()
}

// StringSlice returns all the items in the TextQueue as a string slice
func (tq *TextQueue) StringSlice() []string {
	var a []string
	for elem := tq.Data.Back(); elem != nil; elem = elem.Prev() {
		a = append(a, elem.Value.(string))
	}
	return a
}

// Purge removes all items from the queue
func (tq *TextQueue) Purge() {
	tq.Data.Init()
}
