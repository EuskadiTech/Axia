package request

import (
	"axia4/ldap"
	"axia4/ldap/ldap_check"
	"axia4/types"
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
)

func LdapDel_tx(ctx context.Context, tx pgx.Tx, reqJson json.RawMessage) (any, error) {
	var req struct {
		Id int32 `json:"id"`
	}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}
	return nil, ldap.Del_tx(ctx, tx, req.Id)
}

func LdapGet_tx(ctx context.Context, tx pgx.Tx) (any, error) {
	return ldap.Get_tx(ctx, tx)
}

func LdapSet_tx(ctx context.Context, tx pgx.Tx, reqJson json.RawMessage) (any, error) {
	var req types.Ldap
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}
	return nil, ldap.Set_tx(ctx, tx, req)
}

func LdapCheck(reqJson json.RawMessage) (any, error) {
	var req struct {
		Id int32 `json:"id"`
	}
	if err := json.Unmarshal(reqJson, &req); err != nil {
		return nil, err
	}
	return nil, ldap_check.Run(req.Id)
}
