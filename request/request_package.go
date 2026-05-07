package request

import (
	"axia4/cache"
	"axia4/config"
	"axia4/tools"
	"axia4/transfer"
	"context"
	"os"
)

func PackageInstall(ctx context.Context) (any, error) {

	// store package file from embedded binary data to temp folder
	filePath, err := tools.GetUniqueFilePath(config.File.Paths.Temp, 8999999, 9999999)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(filePath, cache.Package_CoreCompany, 0644); err != nil {
		return nil, err
	}
	return nil, transfer.ImportFromFiles(ctx, []string{filePath})
}
