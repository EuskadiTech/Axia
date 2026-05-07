package request

import (
	"axia4/schema/pgTrigger"
	"axia4/types"
	"context"
	"encoding/json"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

func PgTriggerDel_tx(ctx context.Context, tx pgx.Tx, reqJson json.RawMessage) (any, error) {
	var req uuid.UUID
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}
	return nil, pgTrigger.Del_tx(ctx, tx, req)
}

func PgTriggerSet_tx(ctx context.Context, tx pgx.Tx, reqJson json.RawMessage) (any, error) {
	var req types.PgTrigger
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}
	return nil, pgTrigger.Set_tx(ctx, tx, req)
}
