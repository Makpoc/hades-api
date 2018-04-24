package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Offset is time offset in hours and minutes
type Offset struct {
	H int
	M int
}

// Less checks if an offset is smaller than the provided one
func (o Offset) Less(other Offset) bool {
	if o == other {
		return true
	}

	if o.H != other.H {
		return o.H < other.H
	}

	return o.M < other.M
}

// String implements Stringer
func (o Offset) String() string {
	var hoursPrefix string
	if o.H > 0 {
		hoursPrefix = "+"
	}
	return fmt.Sprintf("%s%d.%d", hoursPrefix, o.H, o.M)
}

// ParseOffset parses a string into Hour:Minutes offset
func ParseOffset(offset string) (Offset, error) {
	// test if it's a number at all
	_, err := strconv.ParseFloat(offset, 64)
	if err != nil {
		return Offset{}, err
	}

	var offsetH, offsetM = "0", "0"
	if strings.Contains(offset, ".") {
		// float offset in form "h.m"
		parts := strings.Split(offset, ".")
		if len(parts) != 2 {
			return Offset{}, fmt.Errorf("invalid offset format")
		}
		offsetH, offsetM = parts[0], parts[1]
	} else {
		offsetH = offset
	}

	offsetHInt, err := strconv.Atoi(offsetH)
	if err != nil {
		return Offset{}, fmt.Errorf("failed to parse hours offset: %v", err)
	}
	offsetMInt, err := strconv.Atoi(offsetM)
	if err != nil {
		return Offset{}, fmt.Errorf("failed to parse minutes offset: %v", err)
	}

	return Offset{H: offsetHInt, M: offsetMInt}, nil
}

type Availability struct {
	From time.Duration
	To   time.Duration
}

func (a Availability) String() string {
	return fmt.Sprintf("%02d:%02d-%02d:%02d", int(a.From.Hours()), int(a.From.Minutes())-int(a.From.Hours())*60, int(a.To.Hours()), int(a.To.Minutes())-int(a.To.Hours())*60)
}

func ParseAvailability(rawValue string) ([]Availability, error) {
	// format:
	// 8-10;11:00-13:00;15:30-20:45
	if rawValue == "" {
		return []Availability{}, nil
	}
	allAvails := strings.Split(rawValue, ";")

	var result []Availability
	for _, avail := range allAvails {
		fromTo := strings.Split(avail, "-")
		if len(fromTo) != 2 {
			return nil, fmt.Errorf("Failed to parse availability. Unknown range: %s", avail)
		}

		from := strings.TrimSpace(fromTo[0])
		to := strings.TrimSpace(fromTo[1])

		var currAvailability Availability
		var err error
		currAvailability.From, err = parseTime(from)
		if err != nil {
			return nil, err
		}
		currAvailability.To, err = parseTime(to)
		if err != nil {
			return nil, err
		}

		result = append(result, currAvailability)
	}

	return result, nil
}

func parseTime(avail string) (time.Duration, error) {
	if avail == "" {
		return -1, fmt.Errorf("empty time")
	}

	var hour, min int
	var err error
	if strings.Contains(avail, ":") {
		hourMin := strings.SplitN(avail, ":", 2)
		hour, err = strconv.Atoi(hourMin[0])
		if err != nil {
			return -1, fmt.Errorf("failed to parse hours: %v", err)
		}
		min, err = strconv.Atoi(hourMin[1])
		if err != nil {
			return -1, fmt.Errorf("failed to parse minutes: %v", err)
		}
	} else {
		hour, err = strconv.Atoi(avail)
		if err != nil {
			return -1, fmt.Errorf("failed to parse hours: %v", err)
		}
	}

	return time.ParseDuration(fmt.Sprintf("%dh%dm", hour, min))
}

// UserTime contains information about users and their time
type UserTime struct {
	UserName     string
	CurrentTime  time.Time
	Offset       Offset
	Availability []Availability `json:",omitempty"`
}

// UserTimes is a list of UserTimes
type UserTimes []UserTime

func (ut UserTimes) Len() int      { return len(ut) }
func (ut UserTimes) Swap(i, j int) { ut[i], ut[j] = ut[j], ut[i] }
func (ut UserTimes) Less(i, j int) bool {
	return ut[i].Offset.Less(ut[j].Offset)
}
