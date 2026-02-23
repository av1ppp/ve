package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "060102"
const startMicro = 1

var getNow = func() time.Time {
	return time.Now().UTC()
}

type Version struct {
	date  time.Time
	micro int
}

func NewVersion() Version {
	return Version{
		date:  normalizeDate(getNow()),
		micro: startMicro,
	}
}

func ParseVersion(s string) (Version, error) {
	parts := strings.Split(s, ".")
	if len(parts) != 2 {
		return Version{}, errors.New("invalid version file")
	}

	rDate := parts[0]
	rMicro := parts[1]

	date, err := time.Parse(dateLayout, rDate)
	if err != nil {
		return Version{}, fmt.Errorf("invalid date: %w", err)
	}

	micro, err := strconv.Atoi(rMicro)
	if err != nil {
		return Version{}, fmt.Errorf("invalid micro: %w", err)
	}
	if micro < startMicro {
		return Version{}, fmt.Errorf("invalid micro: must be >= %d", startMicro)
	}

	return Version{
		date:  normalizeDate(date),
		micro: micro,
	}, nil
}

func ParseVersionFile(name string) (Version, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return Version{}, err
	}

	return ParseVersion(strings.TrimSpace(string(data)))
}

func (self Version) String() string {
	return self.date.Format(dateLayout) + "." + strconv.Itoa(self.micro)
}

func (self Version) Increment() Version {
	now := normalizeDate(getNow())

	y1, m1, d1 := self.date.Date()
	y2, m2, d2 := now.Date()

	if y1 == y2 && m1 == m2 && d1 == d2 {
		return Version{date: self.date, micro: self.micro + 1}
	}
	return Version{date: now, micro: startMicro}
}

func (self Version) WriteFile(name string) error {
	return os.WriteFile(name, []byte(self.String()), 0644)
}

func normalizeDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
