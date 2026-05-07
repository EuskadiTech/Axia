package request

import (
	"axia4/schema"
	"axia4/schema/article"
	"axia4/types"
	"context"
	"encoding/json"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

func ArticleAssign_tx(ctx context.Context, tx pgx.Tx, reqJson json.RawMessage) (any, error) {
	var req struct {
		Target     schema.DbEntity `json:"target`
		TargetId   uuid.UUID       `json:"targetId"`
		ArticleIds []uuid.UUID     `json:"articleIds"`
	}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}
	return nil, article.Assign_tx(ctx, tx, req.Target, req.TargetId, req.ArticleIds)
}

func ArticleDel_tx(ctx context.Context, tx pgx.Tx, reqJson json.RawMessage) (any, error) {
	var req uuid.UUID
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}
	return nil, article.Del_tx(ctx, tx, req)
}

func ArticleSet_tx(ctx context.Context, tx pgx.Tx, reqJson json.RawMessage) (any, error) {
	var req types.Article
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}
	return nil, article.Set_tx(ctx, tx, req.ModuleId, req.Id, req.Name, req.Captions)
}
