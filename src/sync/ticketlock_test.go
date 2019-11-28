package sync_test

import (
	"runtime"
	. "sync"
	"testing"
)

func TestTicketLockGoWithCallbackRace(t *testing.T) {
	var (
		wg   WaitGroup
		lock TicketLock
	)

	out := []uint32{}
	for i := 0; i < 1000; i++ {
		i := uint32(i)
		wg.Add(1)

		lock.GoWithCallback(func() {
			runtime.Gosched()
		}, func() {
			defer wg.Done()
			out = append(out, i)
		})
	}
	wg.Wait()

	for want, seen := range out {
		if uint32(want) != seen {
			t.Errorf("Callbacks arrived out of order: expected %d, got %d.", want, seen)
		}
	}
}
