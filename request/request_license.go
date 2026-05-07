package request

import (
	"axia4/cluster"
	"axia4/config"
	"context"

	"github.com/jackc/pgx/v5"
)

func LicenseDel_tx(ctx context.Context, tx pgx.Tx) (any, error) {
	if err := config.SetString_tx(ctx, tx, "licenseFile", ""); err != nil {
		return nil, err
	}
	return nil, cluster.ConfigChanged_tx(ctx, tx, true, false, false)
}
