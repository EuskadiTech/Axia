package request_login

import (
	"axia4/login/login_widget"
	"axia4/types"
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
)

func WidgetGroupsGet_tx(ctx context.Context, tx pgx.Tx, loginId int64) (any, error) {
	return login_widget.Get_tx(ctx, tx, loginId)
}
func WidgetGroupsSet_tx(ctx context.Context, tx pgx.Tx, reqJson json.RawMessage, loginId int64) (any, error) {
	var req []types.LoginWidgetGroup
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}
	return nil, login_widget.Set_tx(ctx, tx, loginId, req)
}
