package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"playlistturbo.com/model"
)

type Encoder interface {
	// EncodeDataResponse with status 200 and given data as json if err is nil, otherwise an error response is sent.
	EncodeDataResponse(r *http.Request, w http.ResponseWriter, resp interface{}, err error)

	// EncodeMarshaledJSON with status 200 and given json if err is nil, otherwise an error response is sent.
	EncodeMarshaledJSON(r *http.Request, w http.ResponseWriter, json string, err error)

	// ctrl.EncodeEmptyResponse with status 204 No Content and empty body if err is nil, otherwise an error response is sent.
	EncodeEmptyResponse(r *http.Request, w http.ResponseWriter, err error)

	// EncodeBinaryResponse with status 200 and given binary file as body if err is nil, otherwise an error response is sent.
	EncodeBinaryResponse(r *http.Request, w http.ResponseWriter, content *bytes.Buffer, mimeType string, err error)

	// EncodeTextResponse with status 200 and string as body if err is nil, otherwise an error response is sent.
	EncodeTextResponse(r *http.Request, w http.ResponseWriter, text string, err error)

	EncodeDataResponseFiles(r *http.Request, w http.ResponseWriter, content []byte, err error)
}

func (u utils) EncodeDataResponse(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	if err != nil {
		u.HandleControllerError(r, w, err)
		return
	}

	marshaled, err := json.Marshal(resp)
	if err != nil {
		u.HandleControllerError(r, w, err)
		return
	}

	w.Header().Set(model.HeaderContentType, model.MimeTypeJSON)

	if _, err = w.Write(marshaled); err != nil {
		u.HandleControllerError(r, w, err)
		return
	}
}

func (u utils) EncodeMarshaledJSON(r *http.Request, w http.ResponseWriter, json string, err error) {
	if err != nil {
		u.HandleControllerError(r, w, err)
		return
	}

	w.Header().Set(model.HeaderContentType, model.MimeTypeJSON)

	if _, err = w.Write([]byte(json)); err != nil {
		u.HandleControllerError(r, w, err)
		return
	}
}

func (u utils) EncodeEmptyResponse(r *http.Request, w http.ResponseWriter, err error) {
	if err != nil {
		u.HandleControllerError(r, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (u utils) EncodeBinaryResponse(r *http.Request, w http.ResponseWriter, content *bytes.Buffer, mimeType string, err error) {
	if err != nil {
		w.Header().Set("Content-Type", model.MimeTypeText)
		u.HandleControllerError(r, w, err)
		return
	}

	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Content-Disposition", "attachment; filename=export."+mimeType)
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	_, err = content.WriteTo(w)
	if err != nil {
		log.Print("couldn't write buffer to response")
	}
}

func (u utils) EncodeTextResponse(r *http.Request, w http.ResponseWriter, text string, err error) {
	if err != nil {
		u.HandleControllerError(r, w, err)
		return
	}

	w.Header().Set(model.HeaderContentType, model.MimeTypeText)
	_, err = w.Write([]byte(text))
	if err != nil {
		u.HandleControllerError(r, w, err)
		return
	}
}

func (u utils) EncodeDataResponseFiles(r *http.Request, w http.ResponseWriter, content []byte, err error) {
	if err != nil {
		w.Header().Set("Content-Type", model.MimeTypeText)
		u.HandleControllerError(r, w, err)
		return
	}

	contentType := http.DetectContentType(content)

	w.Header().Set("Content-Type", contentType)
	if _, err := w.Write(content); err != nil {
		log.Print("unable to write file.")
	}
}
