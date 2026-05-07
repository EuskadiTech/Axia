package login_clusterEvent

import (
	"axia4/cluster"
	"axia4/log"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func Kick_tx(ctx context.Context, tx pgx.Tx, loginId int64, loginName string) {
	log.Info(log.ContextServer, fmt.Sprintf("user account '%s' is locked, kicking active sessions", loginName))

	if err := cluster.LoginDisabled_tx(ctx, tx, true, loginId); err != nil {
		log.Warning(log.ContextServer, fmt.Sprintf("could not kick active sessions for '%s'", loginName), err)
	}
}

func Reauth_tx(ctx context.Context, tx pgx.Tx, loginId int64, loginName string) {
	log.Info(log.ContextServer, fmt.Sprintf("user account '%s' received new roles, renewing access permissions", loginName))

	if err := cluster.LoginReauthorized_tx(ctx, tx, true, loginId); err != nil {
		log.Warning(log.ContextServer, fmt.Sprintf("could not renew access permissions for '%s'", loginName), err)
	}
}
