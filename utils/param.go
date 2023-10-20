package utils

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gofrs/uuid"
	"playlistturbo.com/model"
	"playlistturbo.com/plterror"
)

type Param interface {
	// GetParam reads the query or path param already validated from the request context.
	// The param must be defined in the relative route.
	GetParam(r *http.Request, name string) string

	// // ctrl.GetClaims get claims from a request.
	// GetClaims(r *http.Request) model.KcClaims

	// GetParamUUID like getParam plus the string is converted to uuid.
	// The validate tag "uuid" should be set in the relative route.
	GetParamUUID(r *http.Request, name string) uuid.UUID

	// GetParamInt like getParam plus the string is converted to an integer.
	// The validate tag "numeric" should be set in the relative route.
	GetParamInt(r *http.Request, name string) int

	// GetBody reads the body already validated from the request context.
	// The body type must be defined in the relative route.
	GetBody(r *http.Request) interface{}
}

func (u utils) GetParamUUID(r *http.Request, name string) uuid.UUID {
	param := u.GetParam(r, name)

	if param == "" { // only if not required
		return uuid.Nil // default value
	}

	id, err := uuid.FromString(param)
	if err != nil {
		panic(plterror.ErrBadSyntax) // it should never happen
	}

	return id
}

func (u utils) GetBody(r *http.Request) interface{} {
	body := r.Context().Value(&model.CtxKeyBody)
	if body == nil {
		panic(plterror.ErrBadSyntax) // it should never happen
	}
	if err := validateBody(body); err != nil {
		panic(err)
	}

	return body
}

func validateBody(body interface{}) error {
	v := validator.New()
	var err error
	if err = v.Struct(body); err != nil {
		//nolint:forcetypeassert
		for range err.(validator.ValidationErrors) {
			return plterror.ErrBadSyntax
		}
	}
	return nil
}

func (u utils) GetParam(r *http.Request, name string) string {
	params := r.Context().Value(&model.CtxKeyParams)
	if params == nil {
		plterror.Logger.Info(plterror.ErrBadSyntax) // it should never happen
	}

	paramsParsed, ok := params.(map[string]string)
	if !ok {
		plterror.Logger.Info(plterror.ErrBadSyntax) // it should never happen
	}

	return paramsParsed[name]
}

func (u utils) GetParamInt(r *http.Request, name string) int {
	param := u.GetParam(r, name)

	number, err := strconv.ParseInt(param, 10, 0)
	if err != nil {
		panic(plterror.ErrBadSyntax) // it should never happen
	}

	return int(number)
}

// func (u utils) ValidatePath()
