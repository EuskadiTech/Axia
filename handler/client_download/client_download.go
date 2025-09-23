package client_download

import (
	"context"
	"errors"
	"net/http"
	"r3/bruteforce"
	"r3/config"
	"r3/handler"
	"r3/login/login_auth"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	if blocked := bruteforce.Check(r); blocked {
		handler.AbortRequestNoLog(w, handler.ErrBruteforceBlock)
		return
	}

	// get authentication token
	token, err := handler.ReadGetterFromUrl(r, "token")
	if err != nil {
		handler.AbortRequest(w, handler.ContextClientDownload, err, handler.ErrGeneral)
		return
	}

	ctx, ctxCanc := context.WithTimeout(context.Background(),
		time.Duration(int64(config.GetUint64("dbTimeoutDataWs")))*time.Second)

	defer ctxCanc()

	// authenticate via token
	if _, err := login_auth.Token(ctx, token); err != nil {
		handler.AbortRequest(w, handler.ContextClientDownload, err, handler.ErrAuthFailed)
		bruteforce.BadAttempt(r)
		return
	}

	// parse getters
	requestedOs, err := handler.ReadGetterFromUrl(r, "os")
	if err != nil {
		handler.AbortRequest(w, handler.ContextClientDownload, err, handler.ErrGeneral)
		return
	}

	// Client distribution files have been removed to reduce repository size
	// Return appropriate error message
	switch requestedOs {
	case "amd64_windows", "amd64_linux", "arm64_linux", "amd64_mac":
		// Valid OS types but files not available
	default:
		handler.AbortRequest(w, handler.ContextClientDownload, errors.New("unsupported OS"), handler.ErrGeneral)
		return
	}

	// Return error indicating client files are not available
	err = errors.New("client distribution files have been removed to reduce repository size - please build clients separately or contact administrator")
	handler.AbortRequest(w, handler.ContextClientDownload, err, handler.ErrGeneral)
}
