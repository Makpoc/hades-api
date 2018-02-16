package gsclient

const rangeStartMarker = "BOT_RANGE_START"
const rangeEndMarker = "BOT_RANGE_END"

// getDataSubset cleans the list from empty rows and gets the sublist between rangeStartMarker and rangeEndMarker
func getDataSubset(values [][]interface{}) [][]interface{} {
	var result [][]interface{}
	var started bool
	for _, v := range values {
		//fmt.Printf("Checking %d: %v\n", i, v)
		if !started {
			//fmt.Println("Not started yet")
			if len(v) > 0 && v[0] == rangeStartMarker {
				//	fmt.Println("Start marker found")
				started = true
			}
			continue
		}

		// if we are here - we've started
		if len(v) > 0 && v[0] == rangeEndMarker {
			//fmt.Println("End marker found")
			break
		}

		if len(v) > 0 {
			//fmt.Printf("Added %v\n", v)
			result = append(result, v)
		}
	}
	return result
}
