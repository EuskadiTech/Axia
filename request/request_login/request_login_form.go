package request_login

import (
	"axia4/schema/loginForm"
	"axia4/types"
	"context"
	"encoding/json"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

func FormDel_tx(ctx context.Context, tx pgx.Tx, reqJson json.RawMessage) (any, error) {
	var req uuid.UUID
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}
	return nil, loginForm.Del_tx(ctx, tx, req)
}

func FormSet_tx(ctx context.Context, tx pgx.Tx, reqJson json.RawMessage) (any, error) {
	var req types.LoginForm
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}
	return nil, loginForm.Set_tx(ctx, tx, req.ModuleId, req.Id, req.AttributeIdLogin,
		req.AttributeIdLookup, req.FormId, req.Name, req.Captions)
}
