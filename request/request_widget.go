package request

import (
	"axia4/schema/widget"
	"axia4/types"
	"context"
	"encoding/json"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

func WidgetDel_tx(ctx context.Context, tx pgx.Tx, reqJson json.RawMessage) (any, error) {
	var req uuid.UUID
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}
	return nil, widget.Del_tx(ctx, tx, req)
}

func WidgetSet_tx(ctx context.Context, tx pgx.Tx, reqJson json.RawMessage) (any, error) {
	var req types.Widget
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}
	return nil, widget.Set_tx(ctx, tx, req)
}
