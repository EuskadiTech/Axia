package request

import (
	"axia4/backup"
	"axia4/config"
	"axia4/types"
)

func BackupGet() (interface{}, error) {

	// no backup directory set, return empty value
	if config.GetString("backupDir") == "" {
		return types.BackupTocFile{Backups: make([]types.BackupDef, 0)}, nil
	}
	return backup.TocFileReadCreate()
}
