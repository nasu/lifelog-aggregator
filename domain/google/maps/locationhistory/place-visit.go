package locationhistory

type PlaceVisit struct {
	Location                PlaceLocation          `json:"location"`
	Duration                Duration               `json:"duration"`
	PlaceConfidence         PlaceConfidence        `json:"placeConfidence"`
	CenterLatE7             LL_E7                  `json:"centerLatE7"`
	CenterLngE7             LL_E7                  `json:"centerLngE7"`
	VisitConfidence         int                    `json:"visitConfidence"`
	EditConfirmationStatus  EditConfirmationStatus `json:"editConfirmationStatus"`
	OtherCandidateLocations []PlaceLocation        `json:"otherCandidateLocations"`
}

type PlaceLocation struct {
	LatitudeE7         LL_E7        `json:"latitudeE7"`
	LongitudeE7        LL_E7        `json:"longitudeE7"`
	PlaceID            string       `json:"placeId"`
	Address            string       `json:"address"`
	Name               string       `json:"name"`
	SemanticType       SemanticType `json:"semanticType"`
	SourceInfo         SourceInfo   `json:"sourceInfo"`
	LocationConfidence float64      `json:"locationConfidence"`
}

type SemanticType string

const (
	SemanticType_HOME             SemanticType = "TYPE_HOME"
	SemanticType_WORK             SemanticType = "TYPE_WORK"
	SemanticType_SEARCHED_ADDRESS SemanticType = "TYPE_SEARCHED_ADDRESS"
)

type PlaceConfidence string

const (
	PlaceConfidence_HIGH   PlaceConfidence = "HIGH_CONFIDENCE"
	PlaceConfidence_MEDIUM PlaceConfidence = "MEDIUM_CONFIDENCE"
	PlaceConfidence_LOW    PlaceConfidence = "LOW_CONFIDENCE"
	PlaceConfidence_USER   PlaceConfidence = "USER_CONFIRMED"
)

type EditConfirmationStatus string

const (
	EditConfirmationStatus_DONE EditConfirmationStatus = "CONFIRMED"
	EditConfirmationStatus_NOT  EditConfirmationStatus = "NOT_CONFIRMED"
)
