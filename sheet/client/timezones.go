package gsclient

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/makpoc/hades-api/sheet/models"
)

const tzSheet = "Timezones"

// GetTimeZones returns the list with users and their corresponding offset and currentTime
func (s *Sheet) GetTimeZones() ([]models.UserTime, error) {
	const userColumn = "A"
	const offsetColumn = "C"
	users, err := s.service.Spreadsheets.Values.Get(s.id, fmt.Sprintf("%s!%s%d:%s%d", tzSheet, userColumn, minRowN, offsetColumn, maxRowN)).Do()
	if err != nil {
		return nil, err
	}

	if len(users.Values) == 0 {
		return nil, fmt.Errorf("no data found")
	}

	values := getDataSubset(users.Values)
	var result models.UserTimes
	for _, v := range values {
		if len(v) == 0 {
			// empty row, skip
			continue
		}

		result = append(result, buildUserTime(v))
	}

	sort.Sort(result)
	return result, nil
}

// GetTimeZone returns the time zone information for the provided user
func (s *Sheet) GetTimeZone(user string) (models.UserTime, error) {
	if user == "" {
		return models.UserTime{}, fmt.Errorf("empty user provided")
	}

	allTz, err := s.GetTimeZones()
	if err != nil {
		return models.UserTime{}, err
	}

	for _, tz := range allTz {
		if strings.ToLower(user) == strings.ToLower(tz.UserName) {
			return tz, nil
		}
	}

	return models.UserTime{}, fmt.Errorf("no time zone information found for user %s", user)
}

// buildUserTime builds UserTime from sheet cell values
func buildUserTime(v []interface{}) models.UserTime {
	entry := models.UserTime{}
	if len(v) >= 1 {
		entry.UserName = fmt.Sprintf("%s", v[0])
	}
	if len(v) == 3 {
		offsetStr, ok := v[2].(string)
		if !ok {
			fmt.Println("Cell value is not a string!")
		}

		var err error
		entry.Offset, err = models.ParseOffset(offsetStr)
		if err != nil {
			fmt.Println(err)
		}
	}

	currTime, err := getCurrentTime(entry.Offset)
	if err != nil {
		log.Printf("Failed to calculate current time for user %s: %v", entry.UserName, err)
	}
	entry.CurrentTime = currTime

	return entry
}

// getCurrentTime calculates the time based on given offset
func getCurrentTime(offset models.Offset) (time.Time, error) {
	now := time.Now().UTC()
	userTime := now.Add(time.Hour*time.Duration(offset.H) + time.Minute*time.Duration(offset.M))

	return userTime, nil
}
