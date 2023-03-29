package main

type OcapData struct {
	Markers      [][]interface{} `json:"Markers"`
	AddonVersion string          `json:"addonVersion"`
	CaptureDelay float64         `json:"captureDelay"`
	EndFrame     int             `json:"endFrame"`
	Entities     []struct {
		FramesFired   [][]interface{} `json:"framesFired"`
		Group         string          `json:"group,omitempty"`
		ID            int             `json:"id"`
		IsPlayer      int             `json:"isPlayer,omitempty"`
		Name          string          `json:"name"`
		Positions     [][]interface{} `json:"positions"`
		Role          string          `json:"role,omitempty"`
		Side          string          `json:"side,omitempty"`
		StartFrameNum int             `json:"startFrameNum"`
		Type          string          `json:"type"`
		Class         string          `json:"class,omitempty"`
	} `json:"entities"`
	Events           [][]interface{} `json:"events"`
	ExtensionBuild   string          `json:"extensionBuild"`
	ExtensionVersion string          `json:"extensionVersion"`
	MissionAuthor    string          `json:"missionAuthor"`
	MissionName      string          `json:"missionName"`
	Tags             string          `json:"tags"`
	Times            []struct {
		Date           string  `json:"date"`
		FrameNum       int     `json:"frameNum"`
		SystemTimeUTC  string  `json:"systemTimeUTC"`
		Time           float64 `json:"time"`
		TimeMultiplier float64 `json:"timeMultiplier"`
	} `json:"times"`
	WorldName string `json:"worldName"`
}

type PlayerStats struct {
	Name            string
	ID              int
	Group           string
	PrimaryWeapon   string
	TotalShotsFired int
	TotalKills      int
	Weapons         map[string]int
}
