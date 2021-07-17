package locationhistory

import (
	"reflect"
)

type SemanticLocationHistory struct {
	TimelineObjects  []TimelineObject `json:"timelineObjects"`
	placeVisits      []*PlaceVisit
	activitySegments []*ActivitySegment
}

func (hist *SemanticLocationHistory) GetPlaceVisits() []*PlaceVisit {
	if len(hist.placeVisits) == 0 {
		hist.initialize()
	}
	return hist.placeVisits
}

func (hist *SemanticLocationHistory) GetActivitySegment() []*ActivitySegment {
	if len(hist.placeVisits) == 0 {
		hist.initialize()
	}
	return hist.activitySegments
}

func (hist *SemanticLocationHistory) initialize() {
	for _, obj := range hist.TimelineObjects {
		switch obj.Type() {
		case TimelineObjectType_PLACE_VISIT:
			hist.placeVisits = append(hist.placeVisits, obj.PlaceVisit)
		case TimelineObjectType_ACTIVITY_SEGMENT:
			hist.activitySegments = append(hist.activitySegments, obj.ActivitySegment)
		default:
			break
		}
	}
}

type TimelineObject struct {
	PlaceVisit      *PlaceVisit      `json:"placeVisit"`
	ActivitySegment *ActivitySegment `json:"activitySegment"`
}

type TimelineObjectType reflect.Type

var (
	TimelineObjectType_PLACE_VISIT      TimelineObjectType = reflect.TypeOf(PlaceVisit{})
	TimelineObjectType_ACTIVITY_SEGMENT TimelineObjectType = reflect.TypeOf(ActivitySegment{})
)

func (obj TimelineObject) Type() TimelineObjectType {
	if obj.PlaceVisit != nil {
		return TimelineObjectType_PLACE_VISIT
	}
	if obj.ActivitySegment != nil {
		return TimelineObjectType_ACTIVITY_SEGMENT
	}
	return nil
}
