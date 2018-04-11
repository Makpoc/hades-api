package gsclient

import (
	"fmt"
	"strings"
)

const metaSheet = "bot_meta"

// getUserMappings ...
func (s *Sheet) getUserDiscordIDs() (map[string]string, error) {
	const userCol = "A"
	const discordIDCol = "C"
	userData, err := s.service.Spreadsheets.Values.Get(s.id, fmt.Sprintf("%s!%s%d:%s%d", metaSheet, userCol, minRowN, discordIDCol, maxRowN)).Do()
	if err != nil {
		return nil, err
	}

	if len(userData.Values) == 0 {
		return nil, fmt.Errorf("no data found")
	}

	values := getDataSubset(userData.Values)
	var result = make(map[string]string)
	for _, v := range values {
		if len(v) == 0 {
			// empty row, skip
			continue
		}

		if len(v) != 2 {
			// empty or incomplete row
			continue
		}

		userName, ok := v[0].(string)
		if !ok {
			return nil, fmt.Errorf("username value not a string")
		}

		discordId, ok := v[1].(string)
		if !ok {
			return nil, fmt.Errorf("discordId value not a string")
		}

		result[strings.ToLower(userName)] = discordId
	}

	return result, nil
}

func (s *Sheet) getUserDiscordID(username string) (string, error) {
	if username == "" {
		return "", fmt.Errorf("Empty username")
	}

	username = strings.ToLower(username)
	discordIDs, err := s.getUserDiscordIDs()
	if err != nil {
		return "", err
	}

	if discordID, ok := discordIDs[username]; ok {
		return discordID, nil
	}

	return "", fmt.Errorf("User %s doesn't have assigned discordId", username)
}
