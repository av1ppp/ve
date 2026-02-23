package main

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type VersionStruct struct {
	year      int
	dayOfYear int
	micro     int
}

func NewVersionStruct(micro int) *VersionStruct {
	now := time.Now().UTC()
	version := VersionStruct{}
	version.year = now.Year() - 2000
	version.dayOfYear = int(now.YearDay())
	version.micro = micro
	return &version
}

func ParseVersionFile(name string) (*VersionStruct, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(data), ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid version file")
	}

	rawYear := parts[0]
	rawDayOfYear := parts[1]
	rawMicro := parts[2]

	year, err := strconv.Atoi(rawYear)
	if err != nil {
		return nil, errors.New("invalid year in version file")
	}

	dayOfYear, err := strconv.Atoi(rawDayOfYear)
	if err != nil {
		return nil, errors.New("invalid day of year in version file")
	}

	micro, err := strconv.Atoi(rawMicro)
	if err != nil {
		return nil, errors.New("invalid micro in version file")
	}

	return &VersionStruct{
		year:      year,
		dayOfYear: dayOfYear,
		micro:     micro,
	}, nil
}

func (self *VersionStruct) String() string {
	return strconv.Itoa(self.year) + "." + strconv.Itoa(self.dayOfYear) + "." + strconv.Itoa(self.micro)
}

func (self *VersionStruct) WriteFile(name string) error {
	return os.WriteFile(name, []byte(self.String()), 0644)
}
