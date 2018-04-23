package models

import "time"

type RespawnTime struct {
	Affiliation string
	User        string
	RespawnTime time.Time
}

type RespawnTimes []RespawnTime

func (rt RespawnTimes) Len() int      { return len(rt) }
func (rt RespawnTimes) Swap(i, j int) { rt[i], rt[j] = rt[j], rt[i] }
func (rt RespawnTimes) Less(i, j int) bool {
	return rt[i].RespawnTime.Before(rt[j].RespawnTime)
}
