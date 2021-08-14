package datamodels

import (
	"time"
)

type Shard struct {
	tableName struct{} `sql:"shard_info"`

	ShardId   int64  `json:"shardId" pg:"shard_id"`
	ShardName string `json:"shardName" pg:"shard_name"`

	//  standard DB fields
	CreatedAt  time.Time  `json:"createdAt, omitempty" pg:"created_at, default:now()"`
	UpdatedAt  *time.Time `json:"updatedAt, omitempty" pg:"updated_at, default:now()"`
	ArchivedAt *time.Time `json:"archivedAt, omitempty" pg:"archived_at"`
}
