package request

import (
	"axia4/schema/pgIndex"
	"axia4/types"
	"context"
	"encoding/json"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

func PgIndexDel_tx(ctx context.Context, tx pgx.Tx, reqJson json.RawMessage) (any, error) {
	var req uuid.UUID
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}
	return nil, pgIndex.Del_tx(ctx, tx, req)
}

func PgIndexSet_tx(ctx context.Context, tx pgx.Tx, reqJson json.RawMessage) (any, error) {
	var req types.PgIndex
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}
	// overwrite values that can only be set on the backend
	req.AutoFki = false
	req.PrimaryKey = false
	return nil, pgIndex.Set_tx(ctx, tx, req)
}
