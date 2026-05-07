package request

import (
	"axia4/schema/variable"
	"axia4/types"
	"context"
	"encoding/json"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

func VariableDel_tx(ctx context.Context, tx pgx.Tx, reqJson json.RawMessage) (any, error) {
	var req uuid.UUID
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}
	return nil, variable.Del_tx(ctx, tx, req)
}

func VariableSet_tx(ctx context.Context, tx pgx.Tx, reqJson json.RawMessage) (any, error) {
	var req types.Variable
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}
	return nil, variable.Set_tx(ctx, tx, req)
}
