package session

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/nasu/lifelog-aggregator/infrastructure/dynamodb"
)

type Repository struct {
	table        string
	partitionKey string
	db           *dynamodb.DB
}

type Session struct {
	SessionID string
	UserID    string
	Email     string
	CreatedAt time.Time
}

func NewRepository(db *dynamodb.DB) *Repository {
	return &Repository{
		table:        "lifelog",
		partitionKey: "session",
		db:           db,
	}
}

func (r *Repository) Save(ctx context.Context, sess *Session) error {
	item := make(map[string]types.AttributeValue)
	item["partition_key"] = &types.AttributeValueMemberS{Value: r.partitionKey}
	item["sort_key"] = &types.AttributeValueMemberS{Value: sess.SessionID}
	item["user_id"] = &types.AttributeValueMemberS{Value: sess.UserID}
	item["email"] = &types.AttributeValueMemberS{Value: sess.Email}
	item["created_at"] = &types.AttributeValueMemberS{Value: sess.CreatedAt.Format(time.RFC3339)}
	return r.db.PutItem(ctx, r.table, item)
}

func (r *Repository) Get(ctx context.Context, sessID string) (*Session, error) {
	values, err := r.db.GetItem(ctx, r.table, r.partitionKey, sessID)
	if err != nil {
		return nil, err
	}

	sess := &Session{}
	if v, ok := values["sort_key"].(*types.AttributeValueMemberS); ok {
		sess.SessionID = v.Value
	}
	if v, ok := values["user_id"].(*types.AttributeValueMemberS); ok {
		sess.UserID = v.Value
	}
	if v, ok := values["email"].(*types.AttributeValueMemberS); ok {
		sess.Email = v.Value
	}
	if v, ok := values["created_at"].(*types.AttributeValueMemberS); ok {
		sess.CreatedAt, _ = time.Parse(time.RFC3339, v.Value)
	}
	return sess, nil
}
