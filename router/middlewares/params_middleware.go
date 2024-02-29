package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
	"playlistturbo.com/model"
	"playlistturbo.com/plterror"
)

var validate = validator.New() // it should be used as a singleton, see docs for details

// Param define a single param for validation.
type Param struct {
	Name       string // param name used to read it from request
	Validation string // go-validator like string, i.e. "required,uuid"
	IsQuery    bool
}

// Validate checks if the given value is valid for the param.
func (p Param) Validate(v string) error {
	errs := validate.Var(v, p.Validation)
	if errs != nil {
		return plterror.ErrBadSyntax
	}
	return nil
}

// Params define query and path params.
//
// Query params to be parsed like ?email=mrossi@email.com
//
// Path params to be parsed like /student/{id}
type Params []Param

type FormValues map[string]string

func handleValidationError(err error) error {
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return err
	}

	miiErr := plterror.ErrBadSyntax
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	for _, err := range errs {
		miiErr.AddStackTraceItem(fmt.Sprintf("%s %s %s %s", err.Field(), err.Type(), err.Value(), err.Tag()))
	}

	return miiErr
}

// ParseBody parse the body in the given request.
func ParseBody(r *http.Request, body interface{}) (interface{}, error) {
	if body == nil {
		return nil, nil
	}

	if r.ContentLength == 0 {
		return nil, plterror.ErrBadSyntax
	}

	// read original type
	originalType := reflect.TypeOf(body)

	// create a new instance of the original struct type
	// otherwise json will decode a map inside the interface
	v := reflect.New(originalType).Elem()

	if err := json.NewDecoder(r.Body).Decode(v.Addr().Interface()); err != nil {
		return nil, plterror.ErrBadSyntax
	}

	if err := validate.Struct(v.Addr().Interface()); err != nil {
		return nil, handleValidationError(err)
	}

	return v.Interface(), nil
}

// ParseParams parse the params in the given request.
func (params Params) ParseParams(r *http.Request) (map[string]string, error) {
	values := make(map[string]string)
	if len(params) == 0 {
		return values, nil // avoid nil map
	}

	for _, p := range params {
		value := chi.URLParam(r, p.Name)
		if query := r.URL.Query().Get(p.Name); query != "" {
			value = query
		}

		err := p.Validate(value)
		if err != nil {
			return nil, handleValidationError(err)
		}

		values[p.Name] = value
	}

	return values, nil
}

// ParamsMiddleware looks if the request is well made by checking
// the query params, the path params, and the request body defined.
func ParamsMiddleware(params Params, body interface{}) Middleware {
	f := func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// read actual context
			ctx := r.Context()

			// parse body
			body, err := ParseBody(r, body)
			if err != nil {
				log.Println("error parsing body:", err)
				handleError(r, w, err)
				return
			}
			ctx = context.WithValue(ctx, &model.CtxKeyBody, body)

			// parse params
			values, err := params.ParseParams(r)
			if err != nil {
				log.Println("error parsing parameters:", err)
				handleError(r, w, err)
				return
			}

			ctx = context.WithValue(ctx, &model.CtxKeyParams, values)

			h.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}

	return f
}
