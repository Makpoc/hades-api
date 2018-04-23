package gsclient

import (
	"fmt"
	"sort"
	"time"

	"github.com/makpoc/hades-api/sheet/models"
	sheets "google.golang.org/api/sheets/v4"
)

const respawnSheet = "RespawnTimer"

// GetTimeZones returns the list with users and their corresponding offset and currentTime
func (s *Sheet) GetRespawnTimes() (models.RespawnTimes, error) {
	const startColumn = "A" // Affiliation
	const endColumn = "C"   // RespawnTime
	respawns, err := s.service.Spreadsheets.Values.Get(s.id, fmt.Sprintf("%s!%s%d:%s%d", respawnSheet, startColumn, minRowN, endColumn, 10)).Do()
	fmt.Printf("%#v\n", respawns)
	if err != nil {
		return nil, err
	}

	if len(respawns.Values) == 0 {
		return models.RespawnTimes{}, nil
	}

	values := getDataSubset(respawns.Values)
	var result models.RespawnTimes
	for _, v := range values {
		if len(v) == 0 {
			// empty row, skip
			continue
		}

		rt, err := buildRespawnTime(v)
		if err != nil {
			fmt.Println(err)
			continue
		}
		result = append(result, rt)
	}

	sort.Sort(result)
	return result, nil
}

func (s *Sheet) GetRespawnTime(user string) (models.RespawnTime, error) {
	if user == "" {
		return models.RespawnTime{}, fmt.Errorf("empty user")
	}
	rTimes, err := s.GetRespawnTimes()
	if err != nil {
		return models.RespawnTime{}, err
	}

	for _, rt := range rTimes {
		if rt.User == user {
			return rt, nil
		}
	}

	return models.RespawnTime{}, fmt.Errorf("User %s not found in respawn table", user)
}

func (s *Sheet) AddRespawnTime(affiliation, user string, respawnTime time.Time) error {
	// get the existing data
	const startColumn = "A" // Affiliation
	existingRespawns, err := s.service.Spreadsheets.Values.Get(s.id, fmt.Sprintf("%s!%s1:%s%d", respawnSheet, startColumn, startColumn, maxRowN)).Do()
	fmt.Printf("%v\n", existingRespawns)
	if err != nil {
		return err
	}

	rangeEndRow := getRangeEndRow(existingRespawns)
	if rangeEndRow == -1 {
		return fmt.Errorf("Failed to find BOT_RANGE_END in sheet!")
	}

	entry := &sheets.ValueRange{}
	entry.Values = append(entry.Values, []interface{}{affiliation, user, respawnTime.Format(time.Stamp)})
	_, err = s.service.Spreadsheets.Values.Append(s.id, fmt.Sprintf("%s!%s%d", respawnSheet, startColumn, rangeEndRow), entry).InsertDataOption("INSERT_ROWS").ValueInputOption("RAW").Do()

	return err
}

func getRangeEndRow(values *sheets.ValueRange) int {
	for i, row := range values.Values {
		fmt.Println(row[0])
		if len(row) > 0 && fmt.Sprintf("%s", row[0]) == rangeEndMarker {
			return i - 1
		}
	}

	return -1
}

// buildUserTime builds UserTime from sheet cell values
func buildRespawnTime(v []interface{}) (models.RespawnTime, error) {
	entry := models.RespawnTime{}
	if len(v) != 3 {
		return entry, fmt.Errorf("Invalid Entry: %v", v)
	}

	entry.Affiliation = fmt.Sprintf("%s", v[0])
	entry.User = fmt.Sprintf("%s", v[1])
	rt, err := time.Parse(time.Stamp, fmt.Sprintf("%s", v[2]))
	if err != nil {
		return models.RespawnTime{}, fmt.Errorf("Failed to parse RespawnTime: %v", err)
	}
	entry.RespawnTime = rt

	return entry, nil
}
