package request_login

import (
	"axia4/login/login_setting"
	"axia4/types"
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func SettingsGet_tx(ctx context.Context, tx pgx.Tx, loginId int64) (any, error) {
	return login_setting.Get_tx(ctx, tx,
		pgtype.Int8{Int64: loginId, Valid: true},
		pgtype.Int8{})
}
func SettingsSet_tx(ctx context.Context, tx pgx.Tx, reqJson json.RawMessage, loginId int64) (any, error) {
	var req types.Settings
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}
	return nil, login_setting.Set_tx(ctx, tx,
		pgtype.Int8{Int64: loginId, Valid: true},
		pgtype.Int8{},
		req, false)
}
