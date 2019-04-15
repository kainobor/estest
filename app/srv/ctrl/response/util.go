package response

import (
	"context"
	"encoding/json"
	"github.com/kainobor/estest/app/lib/logger"
	"net/http"
)

func WriteError(ctx context.Context, w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	if _, err := w.Write([]byte(msg)); err != nil {
		log := logger.New(ctx)
		log.Error(err.Error())
	}
}

func WriteJSON(ctx context.Context, w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)

	log := logger.New(ctx)
	dataJSON, err := json.Marshal(data)
	if err != nil {
		log.Errorw("failed to marshal response: %v", err)
	}

	if _, err := w.Write(dataJSON); err != nil {
		log.Error(err.Error())
	}
}
