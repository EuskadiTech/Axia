package request

import (
	"axia4/schema/searchBar"
	"axia4/types"
	"context"
	"encoding/json"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

func SearchBarDel_tx(ctx context.Context, tx pgx.Tx, reqJson json.RawMessage) (any, error) {
	var req uuid.UUID
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}
	return nil, searchBar.Del_tx(ctx, tx, req)
}

func SearchBarSet_tx(ctx context.Context, tx pgx.Tx, reqJson json.RawMessage) (any, error) {
	var req types.SearchBar
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}
	return nil, searchBar.Set_tx(ctx, tx, req)
}
