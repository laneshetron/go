// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

type TicketLock struct {
	noCopy noCopy
	notify notifyList
}

// notifyListAdd adds the caller to a notify list such that it can receive
// notifications. The caller must eventually call notifyListWait to wait for
// such a notification, passing the returned ticket number.
func (lock *TicketLock) Add() (ticket uint32) {
	return runtime_notifyListAdd(&lock.notify)
}

// Done marks the given ticket as complete, and wakes the next goroutine waiting on it, if any.
func (lock *TicketLock) Done(ticket uint32) {
	runtime_notifyListNotifyTicket(&lock.notify, ticket)
}

// Wait atomically unlocks c.L and suspends execution
// of the calling goroutine. After later resuming execution,
// Wait locks c.L before returning. Unlike in other systems,
// Wait cannot return unless awoken by Broadcast or Signal.
//
// Because c.L is not locked when Wait first resumes, the caller
// typically cannot assume that the condition is true when
// Wait returns. Instead, the caller should Wait in a loop:
//
//    c.L.Lock()
//    for !condition() {
//        c.Wait()
//    }
//    ... make use of condition ...
//    c.L.Unlock()
//
func (lock *TicketLock) Wait(ticket uint32) {
	runtime_notifyListWait(&lock.notify, ticket)
}

func (lock *TicketLock) GoWithCallback(f, callback func()) {
	ticket := lock.Add()
	go func() {
		defer lock.Done(ticket)
		f()
		lock.Wait(ticket - 1)
		callback()
	}()
}
