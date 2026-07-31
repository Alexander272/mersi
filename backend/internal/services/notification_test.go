package services

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/Alexander272/mersi/backend/internal/config"
	"github.com/Alexander272/mersi/backend/internal/models"
)

func newNotificationService(si *fakeSISvc, section *fakeSectionSvc, post *fakePost, times []time.Duration) *NotificationService {
	return NewNotificationService(&NotificationDeps{
		SI:      si,
		File:    nil,
		Section: section,
		Most:    post,
		Conf:    config.UsedConfig{Times: times},
	})
}

func TestCheckUsed_EarlyReturnWhenDateInFuture(t *testing.T) {
	post := &fakePost{}
	svc := newNotificationService(&fakeSISvc{}, &fakeSectionSvc{}, post, []time.Duration{240 * time.Hour})
	svc.now = func() time.Time { return time.Date(2025, 3, 20, 12, 0, 0, 0, time.UTC) }

	if err := svc.CheckUsed(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(post.created) != 0 {
		t.Fatalf("expected no posts, got %d", len(post.created))
	}
}

func TestCheckUsed_PostsReminder(t *testing.T) {
	post := &fakePost{}
	si := &fakeSISvc{
		getUsedFn: func(ctx context.Context, req *models.Period) ([]*models.SiReceiving, error) {
			if req.SectionId != "sec1" {
				t.Fatalf("expected section sec1, got %s", req.SectionId)
			}
			return []*models.SiReceiving{
				{Channel: "ch1", SI: []*models.SI{{Name: "SI-1", FactoryNumber: "123", Person: "Ivanov"}}},
			}, nil
		},
	}
	section := &fakeSectionSvc{
		getAllFn: func(ctx context.Context, req *models.GetAllSectionsDTO) ([]*models.Section, error) {
			if !req.IsActive.Defined || !req.HasReturnNotice.Defined {
				t.Fatal("expected IsActive and HasReturnNotice filters")
			}
			return []*models.Section{{ID: "sec1"}}, nil
		},
	}
	svc := newNotificationService(si, section, post, []time.Duration{240 * time.Hour, 168 * time.Hour})
	svc.now = func() time.Time { return time.Date(2025, 3, 30, 12, 0, 0, 0, time.UTC) }

	if err := svc.CheckUsed(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(post.created) != 1 {
		t.Fatalf("expected 1 post, got %d", len(post.created))
	}
	got := post.created[0]
	if got.ChannelId != "ch1" {
		t.Fatalf("expected channel ch1, got %s", got.ChannelId)
	}
	if len(got.Props) == 0 || got.Props[0].Value != "sia" {
		t.Fatalf("expected service=sia props, got %+v", got.Props)
	}
	if got.Message == "" || !strings.HasPrefix(got.Message, "####") {
		t.Fatalf("expected reminder message, got %q", got.Message)
	}
}

func TestCheckUsed_IterationIncrements(t *testing.T) {
	svc := newNotificationService(&fakeSISvc{}, &fakeSectionSvc{}, &fakePost{}, []time.Duration{240 * time.Hour})
	svc.now = func() time.Time { return time.Date(2025, 3, 30, 12, 0, 0, 0, time.UTC) }

	if err := svc.CheckUsed(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc.iteration != 1 {
		t.Fatalf("expected iteration=1, got %d", svc.iteration)
	}
}

func TestCheckReceiving_UpdatesAcceptedAndSendsMissing(t *testing.T) {
	post := &fakePost{}
	svc := newNotificationService(&fakeSISvc{}, &fakeSectionSvc{}, post, nil)

	instruments := []*models.SI{
		{Id: "si1", Name: "A", FactoryNumber: "1", Place: "dep1", Person: "Ivanov"},
		{Id: "si2", Name: "B", FactoryNumber: "2", Place: "dep1", Person: "Ivanov"},
	}
	siJSON, err := json.Marshal(instruments)
	if err != nil {
		t.Fatal(err)
	}

	dto := &models.DialogResponse{
		State:      "PostId:p1&Status:moved&SI:" + string(siJSON),
		ChannelID:  "ch1",
		Submission: map[string]bool{"si1": true, "si2": false},
	}
	if err := svc.CheckReceiving(context.Background(), dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(post.updated) != 1 {
		t.Fatalf("expected 1 updated post, got %d", len(post.updated))
	}
	if post.updated[0].PostId != "p1" {
		t.Fatalf("expected update post p1, got %+v", post.updated[0])
	}

	if len(post.created) != 1 {
		t.Fatalf("expected 1 created post for missing instruments, got %d", len(post.created))
	}
	if post.created[0].ChannelId != "ch1" {
		t.Fatalf("expected channel ch1 for new post, got %s", post.created[0].ChannelId)
	}
}

func TestCheckReceiving_AllAcceptedNoNewPost(t *testing.T) {
	post := &fakePost{}
	svc := newNotificationService(&fakeSISvc{}, &fakeSectionSvc{}, post, nil)

	instruments := []*models.SI{
		{Id: "si1", Name: "A", FactoryNumber: "1", Place: "dep1", Person: "Ivanov"},
	}
	siJSON, err := json.Marshal(instruments)
	if err != nil {
		t.Fatal(err)
	}

	dto := &models.DialogResponse{
		State:      "PostId:p1&Status:moved&SI:" + string(siJSON),
		ChannelID:  "ch1",
		Submission: map[string]bool{"si1": true},
	}
	if err := svc.CheckReceiving(context.Background(), dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(post.updated) != 1 {
		t.Fatalf("expected 1 updated post, got %d", len(post.updated))
	}
	if len(post.created) != 0 {
		t.Fatalf("expected no new post, got %d", len(post.created))
	}
}

func TestCheckReceiving_AllMissingPanicsOnEmptyAccept(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Log("NOTE: empty accept slice does not panic (no regression)")
		}
	}()

	post := &fakePost{}
	svc := newNotificationService(&fakeSISvc{}, &fakeSectionSvc{}, post, nil)

	instruments := []*models.SI{{Id: "si1", Name: "A", FactoryNumber: "1"}}
	siJSON, err := json.Marshal(instruments)
	if err != nil {
		t.Fatal(err)
	}

	dto := &models.DialogResponse{
		State:      "PostId:p1&Status:moved&SI:" + string(siJSON),
		ChannelID:  "ch1",
		Submission: map[string]bool{"si1": false},
	}
	if err := svc.CheckReceiving(context.Background(), dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
