package gsclient

import (
	"fmt"
	"log"

	"github.com/makpoc/hades-api/sheet/models"
)

const wsFleetSheet = "WS Fleet"

// GetUsers returns all usernames from the WS Fleet sheet
func (s *Sheet) GetUsers() (models.Users, error) {
	const userColumn = "B"

	users, err := s.service.Spreadsheets.Values.Get(s.id, fmt.Sprintf("%s!%s%d:%s%d", wsFleetSheet, userColumn, minRowN, userColumn, maxRowN)).Do()
	if err != nil {
		return nil, err
	}

	if len(users.Values) == 0 {
		return models.Users{}, nil
	}

	var result models.Users
	values := getDataSubset(users.Values)
	for _, u := range values {
		usr, ok := u[0].(string)
		if !ok {
			log.Printf("Value not of type string: %v. It's %t\n", u[0], u[0])
			continue
		}
		result = append(result, models.User(usr))
	}
	return result, nil
}
