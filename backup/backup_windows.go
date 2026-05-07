//go:build windows

package backup

import (
	"axia4/config"
	"axia4/db/embedded"
	"path/filepath"
)

func getPgDumpPath() string {
	if config.File.Db.Embedded {
		return filepath.Join(embedded.GetDbBinPath(), "pg_dump")
	}
	return "pg_dump"
}
