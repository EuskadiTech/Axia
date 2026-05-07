package transfer_export

import (
	"axia4/cache"
	"axia4/config"
	"axia4/db"
	"axia4/handler"
	"axia4/log"
	"axia4/login/login_auth"
	"axia4/tools"
	"axia4/transfer"
	"context"
	"errors"
	"net/http"
	"os"
)

var genErr = "could not finish module export"

func Handler(w http.ResponseWriter, r *http.Request) {

	// get authentication token
	token, err := handler.ReadGetterFromUrl(r, "token")
	if err != nil {
		log.Error(log.ContextServer, genErr, err)
		return
	}

	ctx, ctxCanc := context.WithTimeout(context.Background(), db.CtxDefTimeoutTransfer)
	defer ctxCanc()

	// authenticate via token
	login, err := login_auth.Token(ctx, token)
	if err != nil {
		log.Error(log.ContextServer, genErr, err)
		return
	}

	if !login.Admin {
		log.Error(log.ContextServer, genErr, errors.New(handler.ErrUnauthorized))
		return
	}

	exportKey, err := cache.GetExportKey(login.Id)
	if err != nil {
		log.Error(log.ContextServer, genErr, err)
		return
	}

	// get module ID
	moduleId, err := handler.ReadUuidGetterFromUrl(r, "module_id")
	if err != nil {
		log.Error(log.ContextServer, genErr, err)
		return
	}

	filePath, err := tools.GetUniqueFilePath(config.File.Paths.Temp, 8999999, 9999999)
	if err != nil {
		log.Error(log.ContextServer, genErr, err)
		return
	}

	if err := transfer.ExportToFile(ctx, moduleId, exportKey, filePath); err != nil {
		log.Error(log.ContextServer, genErr, err)
		return
	}
	http.ServeFile(w, r, filePath)
	if err := os.Remove(filePath); err != nil {
		log.Warning(log.ContextServer, "could not delete temporary export file", err)
	}
}
