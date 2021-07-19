package locationhistory

import (
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/nasu/lifelog-aggregator/infrastructure/dynamodb"
)

type VisitRepository struct {
	table        string
	partitionKey string
	db           *dynamodb.DB
}

func NewVisitRepository(db *dynamodb.DB) *VisitRepository {
	return &VisitRepository{
		table:        "lifelog-metrics",
		partitionKey: "visit",
		db:           db,
	}
}

func (r *VisitRepository) Save(ctx context.Context, userID string, vis *PlaceVisit) error {
	start := strconv.FormatUint(vis.Duration.StartTimestampMs, 10)
	end := strconv.FormatUint(vis.Duration.EndTimestampMs, 10)
	sortKey := userID + "#" + start

	item := make(map[string]types.AttributeValue)
	item["partition_key"] = &types.AttributeValueMemberS{Value: r.partitionKey}
	item["sort_key"] = &types.AttributeValueMemberS{Value: sortKey}
	item["start"] = &types.AttributeValueMemberN{Value: start}
	item["end"] = &types.AttributeValueMemberN{Value: end}
	item["name"] = &types.AttributeValueMemberS{Value: vis.Location.Name}
	item["address"] = &types.AttributeValueMemberS{Value: vis.Location.Address}
	item["latitude"] = &types.AttributeValueMemberN{Value: strconv.Itoa(int(vis.Location.LatitudeE7))}
	item["longitude"] = &types.AttributeValueMemberN{Value: strconv.Itoa(int(vis.Location.LatitudeE7))}
	return r.db.PutItem(ctx, r.table, item)
}

type MoveRepository struct {
	table        string
	partitionKey string
	db           *dynamodb.DB
}

func NewMoveRepository(db *dynamodb.DB) *MoveRepository {
	return &MoveRepository{
		table:        "lifelog-metrics",
		partitionKey: "move",
		db:           db,
	}
}

func (r *MoveRepository) Save(ctx context.Context, userID string, act *ActivitySegment) error {
	start := strconv.FormatUint(act.Duration.StartTimestampMs, 10)
	end := strconv.FormatUint(act.Duration.EndTimestampMs, 10)
	sortKey := userID + "#" + start

	item := make(map[string]types.AttributeValue)
	item["partition_key"] = &types.AttributeValueMemberS{Value: r.partitionKey}
	item["sort_key"] = &types.AttributeValueMemberS{Value: sortKey}
	item["start"] = &types.AttributeValueMemberN{Value: start}
	item["end"] = &types.AttributeValueMemberN{Value: end}
	item["name"] = &types.AttributeValueMemberS{Value: string(act.ActivityType)}
	item["distance"] = &types.AttributeValueMemberN{Value: strconv.Itoa(int(act.Distance))}
	return r.db.PutItem(ctx, r.table, item)
}

func (r *MoveRepository) GetWithSinceAndUntil(ctx context.Context, userID, since, until string) ([]*ActivitySegment, error) {
	start, err := time.ParseInLocation("2006-01-02", since, time.Local)
	if err != nil {
		return nil, err
	}
	end, err := time.ParseInLocation("2006-01-02", until, time.Local)
	if err != nil {
		return nil, err
	}
	end = end.Add(time.Hour * 24).Add(time.Nanosecond * -1)

	items, err := r.db.GetItemsWithSortKeyRange(ctx, r.table, r.partitionKey,
		userID+"#"+strconv.FormatInt(start.UnixNano()/1000000, 10),
		userID+"#"+strconv.FormatInt(end.UnixNano()/1000000, 10))
	if err != nil {
		return nil, err
	}

	acts := make([]*ActivitySegment, len(items))
	for i, item := range items {
		act := &ActivitySegment{}
		if v, ok := item["start"].(*types.AttributeValueMemberN); ok {
			act.Duration.StartTimestampMs, _ = strconv.ParseUint(v.Value, 10, 64)
		}
		if v, ok := item["end"].(*types.AttributeValueMemberN); ok {
			act.Duration.EndTimestampMs, _ = strconv.ParseUint(v.Value, 10, 64)
		}
		if v, ok := item["name"].(*types.AttributeValueMemberS); ok {
			act.ActivityType = ActivityType(v.Value)
		}
		if v, ok := item["distance"].(*types.AttributeValueMemberN); ok {
			d, _ := strconv.ParseInt(v.Value, 10, 32)
			act.Distance = int32(d)
		}
		acts[i] = act
	}
	return acts, nil
}
