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

// UserTime contains information about users and their time
type UserTime struct {
	UserName    string
	CurrentTime time.Time
	Offset      Offset
}

// UserTimes is a list of UserTimes
type UserTimes []UserTime

func (ut UserTimes) Len() int      { return len(ut) }
func (ut UserTimes) Swap(i, j int) { ut[i], ut[j] = ut[j], ut[i] }
func (ut UserTimes) Less(i, j int) bool {
	return ut[i].Offset.Less(ut[j].Offset)
}
