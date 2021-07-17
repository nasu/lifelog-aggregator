package locationhistory

type ActivitySegment struct {
	StartLocation ActivityLocation `json:"startLocation"`
	EndLocation   ActivityLocation `json:"endLocation"`
	Duration      Duration         `json:"duration"`
	Distance      int32            `json:"distance"`
	ActivityType  ActivityType     `json:"activityType"`
	Confidence    ConfidenceType   `json:"confidence"`
	Activities    []Activity       `json:"activities"`
	WaypointPath  Waypoints        `json:"waypointPath"`
}

type ActivityLocation struct {
	LatitudeE7  LL_E7      `json:"latitudeE7"`
	LongitudeE7 LL_E7      `json:"longitudeE7"`
	SourceInfo  SourceInfo `json:"sourceInfo"`
}

type ConfidenceType string

const (
	ConfidenceType_LOW    ConfidenceType = "LOW"
	ConfidenceType_MEDIUM ConfidenceType = "MEDIUM"
	ConfidenceType_HIGH   ConfidenceType = "HIGH"
)

type ActivityType string

const (
	ActivityType_WALKING              ActivityType = "WALKING"
	ActivityType_STILL                ActivityType = "STILL"
	ActivityType_IN_PASSENGER_VEHICLE ActivityType = "IN_PASSENGER_VEHICLE"
	ActivityType_CYCLING              ActivityType = "CYCLING"
	ActivityType_RUNNING              ActivityType = "RUNNING"
	ActivityType_MOTORCYCLING         ActivityType = "MOTORCYCLING"
	ActivityType_IN_BUS               ActivityType = "IN_BUS"
	ActivityType_IN_FERRY             ActivityType = "IN_FERRY"
	ActivityType_IN_SUBWAY            ActivityType = "IN_SUBWAY"
	ActivityType_FLYING               ActivityType = "FLYING"
	ActivityType_IN_TRAIN             ActivityType = "IN_TRAIN"
	ActivityType_IN_TRAM              ActivityType = "IN_TRAM"
	ActivityType_SKIING               ActivityType = "SKIING"
	ActivityType_SAILING              ActivityType = "SAILING"
	ActivityType_IN_VEHICLE           ActivityType = "IN_VEHICLE"
)

type Activity struct {
	ActivityType ActivityType `json:"activityType"`
	Probability  float64      `json:"probability"`
}

type Waypoints struct {
	Waypoints []Waypoint `json:"waypoints"`
}

type Waypoint struct {
	LatE7 LL_E7 `json:"latE7"`
	LngE7 LL_E7 `json:"lngE7"`
}
