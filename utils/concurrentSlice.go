package utils

/*
import "sync"

// ConcurrentSlice type that can be safely shared between goroutines
type ConcurrentSlice struct {
	sync.RWMutex
	items []interface{}
}

// ConcurrentSliceItem contains the index/value pair of an item in a
// concurrent slice
type ConcurrentSliceItem struct {
	Index int
	Value interface{}
}

// NewConcurrentSlice creates a new concurrent slice
func NewConcurrentSlice() *ConcurrentSlice {
	cs := &ConcurrentSlice{
		items: make([]interface{}, 0),
	}

	return cs
}

// Append adds an item to the concurrent slice
func (cs *ConcurrentSlice) Append(item interface{}) {
	cs.Lock()
	defer cs.Unlock()

	cs.items = append(cs.items, item)
}

func Contains(s *ConcurrentSlice, e interface{}) bool {
    for _, a := range s.Iter() {
        if a == e {
            return true
        }
    }
    return false
}


func removeElementFromSlice(series []interface{}, value interface{}) []interface{}{
	i := indexOf(value, series)
	if i != -1 {
		series[i] = series[len(series)-1] // Copy last element to index i.
		series[len(series)-1] = ""   // Erase last element (write zero value).
		series = series[:len(series)-1]   // Truncate slice.
	}
	return series
}

// Append adds an item to the concurrent slice
func (cs *ConcurrentSlice) Remove(item interface{}) {
	cs.Lock()
	defer cs.Unlock()
	cs.items = removeElementFromSlice(cs.items, item)
}

// Iter iterates over the items in the concurrent slice
// Each item is sent over a channel, so that
// we can iterate over the slice using the builtin range keyword
func (cs *ConcurrentSlice) Iter() <-chan ConcurrentSliceItem {
	c := make(chan ConcurrentSliceItem)

	f := func() {
		cs.Lock()
		defer cs.Unlock()
		for index, value := range cs.items {
			c <- ConcurrentSliceItem{index, value}
		}
		close(c)
	}
	go f()

	return c
}
*/