package services

import (
	"context"
	"errors"
	"testing"
)

func TestJob_CallsAllDependencies(t *testing.T) {
	receiving := &fakeReceivingSvc{}
	notification := &fakeNotificationSvc{}
	user := &fakeUserSvc{}
	documents := &fakeDocumentSvc{}

	svc := NewSchedulerService(&SchedulerDeps{
		Notification: notification,
		User:         user,
		Receiving:    receiving,
		Documents:    documents,
	})

	svc.job()

	if receiving.forcedAllCalls != 1 {
		t.Fatalf("expected ForcedReceiptAll called once, got %d", receiving.forcedAllCalls)
	}
	if notification.usedCalls != 1 || notification.sentCalls != 1 || notification.verificationCalls != 1 {
		t.Fatalf("expected all notification checks called, got %+v", notification)
	}
	if user.calls != 1 {
		t.Fatalf("expected Sync called once, got %d", user.calls)
	}
	if documents.folderCalls != 1 {
		t.Fatalf("expected RemoveEmptyFolders called once, got %d", documents.folderCalls)
	}
}

func TestJob_StopsOnReceivingError(t *testing.T) {
	receiving := &fakeReceivingSvc{
		forcedReceiptAllFn: func(ctx context.Context) error { return errors.New("db error") },
	}
	notification := &fakeNotificationSvc{}
	user := &fakeUserSvc{}
	documents := &fakeDocumentSvc{}

	svc := NewSchedulerService(&SchedulerDeps{
		Notification: notification,
		User:         user,
		Receiving:    receiving,
		Documents:    documents,
	})

	svc.job()

	if notification.usedCalls != 0 || notification.sentCalls != 0 || notification.verificationCalls != 0 {
		t.Fatalf("expected no notification checks after receiving error, got %+v", notification)
	}
	if user.calls != 0 || documents.folderCalls != 0 {
		t.Fatal("expected remaining jobs to be skipped")
	}
}

func TestJob_ContinuesAfterNotificationError(t *testing.T) {
	receiving := &fakeReceivingSvc{}
	notification := &fakeNotificationSvc{
		checkUsedFn: func(ctx context.Context) error { return errors.New("check used error") },
	}
	user := &fakeUserSvc{}
	documents := &fakeDocumentSvc{}

	svc := NewSchedulerService(&SchedulerDeps{
		Notification: notification,
		User:         user,
		Receiving:    receiving,
		Documents:    documents,
	})

	svc.job()

	if notification.sentCalls != 1 || notification.verificationCalls != 1 {
		t.Fatalf("expected remaining checks to run, got %+v", notification)
	}
	if user.calls != 1 || documents.folderCalls != 1 {
		t.Fatal("expected remaining jobs to run")
	}
}
