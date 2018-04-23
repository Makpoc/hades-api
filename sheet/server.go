package sheet

import (
	"net/http"

	"github.com/makpoc/hades-api/sheet/routes"
)

// Init ...
func Init() error {
	return routes.InitSheet()
}

// GetHandleFuncs returns a map with paths and handlers to attach to them
func GetHandleFuncs() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/timezones":            routes.TimeZonesHandler,
		"/timezones/{username}": routes.TimeZoneHandler,
		"/users":                routes.UsersHandler,
		"/users/{username}":     routes.UserHandler,
		"/respawns":             routes.RespawnsHandler,
		"/respawns/{username}":  routes.RespawnHandler,
	}
}
