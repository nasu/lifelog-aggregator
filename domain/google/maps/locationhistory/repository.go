package locationhistory

import (
	"context"
	"strconv"

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
	sortKey := userID + "-" + start

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
	sortKey := userID + "-" + start

	item := make(map[string]types.AttributeValue)
	item["partition_key"] = &types.AttributeValueMemberS{Value: r.partitionKey}
	item["sort_key"] = &types.AttributeValueMemberS{Value: sortKey}
	item["start"] = &types.AttributeValueMemberN{Value: start}
	item["end"] = &types.AttributeValueMemberN{Value: end}
	item["name"] = &types.AttributeValueMemberS{Value: string(act.ActivityType)}
	item["distance"] = &types.AttributeValueMemberN{Value: strconv.Itoa(int(act.Distance))}
	return r.db.PutItem(ctx, r.table, item)
}
