package middlewares

import (
	"net/http"

	"github.com/gofrs/uuid"
	"playlistturbo.com/model"
	"playlistturbo.com/plterror"
)

type Middleware func(h http.Handler) http.Handler

func handleError(r *http.Request, w http.ResponseWriter, err error) {
	err = plterror.PropagateError(err, 2)

	appErr, ok := err.(*plterror.PLTError)
	if !ok {
		appErr = plterror.ErrServerError
	}

	id, ok := r.Context().Value(&model.CtxKeyID).(uuid.UUID)
	if ok {
		appErr.ID = id
	}

	appErr.Log(plterror.LogMessageErrorResponse)

	http.Error(w, appErr.Error(), appErr.Status())
}
