package gsclient

import (
	"fmt"
	"log"
	"strings"

	"github.com/makpoc/hades-api/sheet/models"
)

const wsFleetSheet = "WS Fleet"

type valRange struct {
	start int
	end   int
}

var rowDesc = map[string]valRange{
	"name":      valRange{start: 0},  // B
	"bsRole":    valRange{start: 1},  // C
	"bsLvl":     valRange{start: 2},  // D
	"bsWeapon":  valRange{start: 3},  // E
	"bsShield":  valRange{start: 4},  // F
	"bsMods":    valRange{5, 13},     // G - N
	"utilType":  valRange{start: 13}, // O
	"tsLvl":     valRange{start: 14}, // P
	"tsCap":     valRange{start: 15}, // Q
	"tsMods":    valRange{16, 22},    // R - W
	"minerLvl":  valRange{start: 22}, // X
	"minerMods": valRange{23, 35},    // Y - AJ
}

// GetUsers returns all usernames from the WS Fleet sheet
func (s *Sheet) GetUsers() (models.Users, error) {
	const userColumnFrom = "B"
	const userColumnTo = "AJ"

	var sheetRange = fmt.Sprintf("%s!%s%d:%s%d", wsFleetSheet, userColumnFrom, minRowN, userColumnTo, maxRowN)
	users, err := s.service.Spreadsheets.Values.Get(s.id, sheetRange).Do()
	if err != nil {
		return nil, err
	}

	if len(users.Values) == 0 {
		return models.Users{}, nil
	}

	tz, err := s.GetTimeZones()
	if err != nil {
		log.Printf("failed to load time zone information: %v", err)
	}

	var result models.Users
	values := getDataSubset(users.Values)
	for _, u := range values {
		usr, err := buildUser(u)
		if err != nil {
			log.Printf("Failed to build user from row %#v\n", u)
			continue
		}
		usr.TZ = getUserTz(usr.Name, tz)
		result = append(result, *usr)
	}
	return result, nil
}

// GetUser returns the requested user from the WS sheet
func (s *Sheet) GetUser(username string) (*models.User, error) {
	const userColumnFrom = "B"
	const userColumnTo = "AJ"

	var sheetRange = fmt.Sprintf("%s!%s%d:%s%d", wsFleetSheet, userColumnFrom, minRowN, userColumnTo, maxRowN)
	users, err := s.service.Spreadsheets.Values.Get(s.id, sheetRange).Do()
	if err != nil {
		return nil, err
	}

	if len(users.Values) == 0 {
		return &models.User{}, nil
	}

	tz, err := s.GetTimeZone(username)
	if err != nil {
		log.Printf("failed to load time zone information for user %s: %v", username, err)
	}

	values := getDataSubset(users.Values)
	for _, u := range values {
		name, err := getSingleCellVal(u, rowDesc["name"])
		if err != nil {
			log.Printf("Failed to get username from sheet for user %#v, error was: %#v\n", u, err)
			continue
		}
		if strings.ToLower(name) == strings.ToLower(username) {
			usr, err := buildUser(u)
			if err != nil {
				return nil, err
			}
			usr.TZ = tz
			return usr, nil
		}
	}
	return nil, fmt.Errorf("no information about %s in sheet", username)
}

func buildUser(v []interface{}) (*models.User, error) {
	usr := &models.User{}
	if len(v) == 0 {
		// ignore empty row
		return usr, nil
	}

	name, err := getSingleCellVal(v, rowDesc["name"])
	if err != nil {
		log.Printf("Failed to parse name: %v", err)
		// without name there's nothing left
		return nil, err
	}
	usr.Name = name

	role, err := getSingleCellVal(v, rowDesc["bsRole"])
	if err != nil {
		log.Printf("Failed to parse BS role for user %s: %v", usr.Name, err)
	} else {
		usr.BsRole = role
	}

	weapon, err := getSingleCellVal(v, rowDesc["bsWeapon"])
	if err != nil {
		log.Printf("Failed to parse BS weapon for user %s: %v", usr.Name, err)
	} else {
		usr.BsWeapon = weapon
	}

	shield, err := getSingleCellVal(v, rowDesc["bsShield"])
	if err != nil {
		log.Printf("Failed to parse BS shield for user %s: %v", usr.Name, err)
	} else {
		usr.BsShield = shield
	}

	bsModules, err := getModules(v, rowDesc["bsMods"])
	if err != nil {
		log.Printf("Failed to parse BS Modules for user %s: %v", usr.Name, err)
	} else {
		usr.BsModules = bsModules
	}

	tsCap, err := getSingleCellVal(v, rowDesc["tsCap"])
	if err != nil {
		log.Printf("Failed to parse TS Capacity for user %s: %v", usr.Name, err)
	} else {
		usr.TsCapacity = tsCap
	}

	tsModules, err := getModules(v, rowDesc["tsMods"])
	if err != nil {
		log.Printf("Failed to parse TS Modules for user %s: %v", usr.Name, err)
	} else {
		usr.TsModules = tsModules
	}

	minerLvl, err := getSingleCellVal(v, rowDesc["minerLvl"])
	if err != nil {
		log.Printf("Failed to parse miner level for user %s: %v", usr.Name, err)
	} else {
		usr.MinerLevel = minerLvl
	}

	minerModules, err := getModules(v, rowDesc["minerMods"])
	if err != nil {
		log.Printf("Failed to parse MinerModules for user %s: %v", usr.Name, err)
	} else {
		usr.MinerModules = minerModules
	}

	return usr, nil
}

func getSingleCellVal(v []interface{}, vr valRange) (string, error) {
	index := vr.start
	if len(v) >= index {
		value, ok := v[index].(string)
		if !ok {
			return "", fmt.Errorf("value %v at %d not a string", v[index], index)
		}
		return value, nil
	}
	return "", fmt.Errorf("failed to parse value at index %d", index)
}

func getModules(v []interface{}, vr valRange) (models.Modules, error) {
	var modules = models.Modules{}
	var mStart = vr.start
	var mEnd = vr.end

	if len(v) < mStart {
		// no info about BS modules
		return modules, fmt.Errorf("no modules found")
	}

	if len(v) < mEnd {
		// load as much as we can
		mEnd = len(v)
	}

	vSub := v[mStart:mEnd]

	for i := 0; i < len(vSub); i = i + 2 {
		if vSub[i] == nil {
			// empty cell
			continue
		}
		name, ok := vSub[i].(string)
		if !ok {
			log.Println("Failed to parse module name from cell value: %v", vSub[i])
			continue
		}

		var lvl string
		if i+1 < len(vSub) {
			var ok bool
			lvl, ok = vSub[i+1].(string)
			if !ok {
				log.Println("Failed to parse module level from cell value: %v", vSub[i+1])
				continue
			}
		}
		if lvl == "" {
			lvl = getDefaultLevel(name)
		}

		if name != "" {
			module := models.Module{Name: name, Level: lvl}
			modules = append(modules, module)
		}
	}
	return modules, nil
}

func getUserTz(user string, tz models.UserTimes) models.UserTime {
	result := models.UserTime{}
	for _, t := range tz {
		if strings.ToLower(t.UserName) == strings.ToLower(user) {
			result = t
			break
		}
	}

	return result
}

// getDefaultLevel returns the default lvl for the given module
func getDefaultLevel(modName string) string {
	if strings.ToLower(modName) == "sanctuary" {
		return "1"
	}
	return "?"
}
