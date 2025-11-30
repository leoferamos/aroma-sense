package email

import (
	"log"
	"sync"

	"github.com/leoferamos/aroma-sense/internal/model"
)

// AsyncEmailService is a wrapper that enqueues email send operations and processes them in background workers.
type AsyncEmailService struct {
	svc   EmailService
	queue chan func()
	wg    sync.WaitGroup
	stop  chan struct{}
}

// NewAsyncEmailService wraps a concrete EmailService with an asynchronous worker pool.
func NewAsyncEmailService(svc EmailService, workerCount int, queueSize int) *AsyncEmailService {
	if workerCount <= 0 {
		workerCount = 2
	}
	if queueSize <= 0 {
		queueSize = 100
	}

	a := &AsyncEmailService{
		svc:   svc,
		queue: make(chan func(), queueSize),
		stop:  make(chan struct{}),
	}

	for i := 0; i < workerCount; i++ {
		a.wg.Add(1)
		go func(workerID int) {
			defer a.wg.Done()
			for {
				select {
				case job := <-a.queue:
					safeExecute(job)
				case <-a.stop:
					return
				}
			}
		}(i)
	}

	return a
}

func safeExecute(job func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in async email job: %v", r)
		}
	}()
	job()
}

// enqueue tries to enqueue the job; if the queue is full it will spawn a goroutine to execute the job immediately to avoid blocking the caller.
func (a *AsyncEmailService) enqueue(job func()) {
	select {
	case a.queue <- job:
		// enqueued
	default:
		// queue full — execute in separate goroutine to avoid blocking
		go safeExecute(job)
	}
}

// Stop gracefully stops workers.
func (a *AsyncEmailService) Stop() {
	close(a.stop)
	a.wg.Wait()
}

// The following methods implement EmailService by enqueuing the actual send operations.
func (a *AsyncEmailService) SendPasswordResetCode(to, code string) error {
	a.enqueue(func() { _ = a.svc.SendPasswordResetCode(to, code) })
	return nil
}
func (a *AsyncEmailService) SendOrderConfirmation(to string, order *model.Order) error {
	a.enqueue(func() { _ = a.svc.SendOrderConfirmation(to, order) })
	return nil
}
func (a *AsyncEmailService) SendWelcomeEmail(to, name string) error {
	a.enqueue(func() { _ = a.svc.SendWelcomeEmail(to, name) })
	return nil
}
func (a *AsyncEmailService) SendPromotional(to, subject, htmlBody string) error {
	a.enqueue(func() { _ = a.svc.SendPromotional(to, subject, htmlBody) })
	return nil
}
func (a *AsyncEmailService) SendAccountDeactivated(to, reason string, contestationDeadline string) error {
	a.enqueue(func() { _ = a.svc.SendAccountDeactivated(to, reason, contestationDeadline) })
	return nil
}
func (a *AsyncEmailService) SendContestationReceived(to string) error {
	a.enqueue(func() { _ = a.svc.SendContestationReceived(to) })
	return nil
}
func (a *AsyncEmailService) SendContestationResult(to string, approved bool, reason string) error {
	a.enqueue(func() { _ = a.svc.SendContestationResult(to, approved, reason) })
	return nil
}
func (a *AsyncEmailService) SendDeletionRequested(to string, cancelLink string) error {
	a.enqueue(func() { _ = a.svc.SendDeletionRequested(to, cancelLink) })
	return nil
}
func (a *AsyncEmailService) SendDeletionAutoConfirmed(to string) error {
	a.enqueue(func() { _ = a.svc.SendDeletionAutoConfirmed(to) })
	return nil
}
func (a *AsyncEmailService) SendDataAnonymized(to string) error {
	a.enqueue(func() { _ = a.svc.SendDataAnonymized(to) })
	return nil
}
