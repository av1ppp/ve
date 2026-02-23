package main

import (
	"path/filepath"
	"testing"
	"time"
)

func withNow(t *testing.T, fake time.Time, fn func()) {
	orig := getNow
	getNow = func() time.Time { return fake }
	t.Cleanup(func() { getNow = orig })
	fn()
}

func TestVersion_StringParse_Roundtrip(t *testing.T) {
	v := Version{
		date:  time.Date(2026, 2, 23, 0, 0, 0, 0, time.UTC),
		micro: 3,
	}

	s := v.String()

	parsed, err := ParseVersion(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if parsed.date != v.date {
		t.Fatalf("date mismatch: %v vs %v", parsed.date, v.date)
	}

	if parsed.micro != v.micro {
		t.Fatalf("micro mismatch: %d vs %d", parsed.micro, v.micro)
	}
}

func TestParseVersion_InvalidFormat(t *testing.T) {
	_, err := ParseVersion("invalid")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseVersion_InvalidDate(t *testing.T) {
	_, err := ParseVersion("bad123.1")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseVersion_InvalidMicro(t *testing.T) {
	_, err := ParseVersion("260223.abc")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseVersion_MicroTooSmall(t *testing.T) {
	_, err := ParseVersion("260223.0")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewVersion(t *testing.T) {
	fakeNow := time.Date(2026, 2, 23, 15, 30, 0, 0, time.UTC)

	withNow(t, fakeNow, func() {
		v := NewVersion()

		expectedDate := time.Date(2026, 2, 23, 0, 0, 0, 0, time.UTC)

		if v.date != expectedDate {
			t.Fatalf("expected %v, got %v", expectedDate, v.date)
		}

		if v.micro != startMicro {
			t.Fatalf("expected micro %d, got %d", startMicro, v.micro)
		}
	})
}

func TestIncrement_SameDay(t *testing.T) {
	date := time.Date(2026, 2, 23, 0, 0, 0, 0, time.UTC)
	fakeNow := time.Date(2026, 2, 23, 18, 0, 0, 0, time.UTC)

	withNow(t, fakeNow, func() {
		v := Version{date: date, micro: 5}
		next := v.Increment()

		if next.micro != 6 {
			t.Fatalf("expected 6, got %d", next.micro)
		}

		if next.date != date {
			t.Fatalf("date should not change")
		}
	})
}

func TestIncrement_NewDay(t *testing.T) {
	oldDate := time.Date(2026, 2, 22, 0, 0, 0, 0, time.UTC)
	newDate := time.Date(2026, 2, 23, 10, 0, 0, 0, time.UTC)

	withNow(t, newDate, func() {
		v := Version{date: oldDate, micro: 7}
		next := v.Increment()

		expectedDate := time.Date(2026, 2, 23, 0, 0, 0, 0, time.UTC)

		if next.date != expectedDate {
			t.Fatalf("expected %v, got %v", expectedDate, next.date)
		}

		if next.micro != startMicro {
			t.Fatalf("expected %d, got %d", startMicro, next.micro)
		}
	})
}

func TestWriteAndReadFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "VERSION")

	v := Version{
		date:  time.Date(2026, 2, 23, 0, 0, 0, 0, time.UTC),
		micro: 4,
	}

	if err := v.WriteFile(path); err != nil {
		t.Fatalf("write error: %v", err)
	}

	read, err := ParseVersionFile(path)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}

	if read.String() != v.String() {
		t.Fatalf("expected %s, got %s", v.String(), read.String())
	}
}
