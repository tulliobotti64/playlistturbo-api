package utils

import (
	"log"
	"net/http"

	"playlistturbo.com/plterror"
)

type Handler interface {
	// HandleError handles the given error by propagating it if different than nil, otherwise does nothing.
	HandleError(err error) error

	// HandleError handles the given error and writes the response accordingly.
	HandleControllerError(r *http.Request, w http.ResponseWriter, err error)
}

func (utils) HandleError(err error) error {
	if err != nil {
		return plterror.PropagateError(err, 2)
	}

	return nil
}

func (utils) HandleControllerError(r *http.Request, w http.ResponseWriter, err error) {
	err = plterror.PropagateError(err, 3)

	// id, ok := r.Context().Value(&model.CtxKeyID).(uuid.UUID)
	// if !ok {
	// 	plterror.Logger.Info(plterror.ErrServerError) // TODO
	// }

	appErr, ok := err.(*plterror.PLTError)
	if !ok {
		log.Println("HandleControllerError!ok:", appErr)
		appErr = plterror.ErrServerError
		appErr.Message = err.Error()
	}

	// appErr.ID = id

	http.Error(w, appErr.Error(), appErr.Status()) // TODO
	log.Println("HandleControllerError:", appErr)
	// appErr.Log(plterror.LogMessageErrorResponse)
}
