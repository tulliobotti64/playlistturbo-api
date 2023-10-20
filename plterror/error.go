package plterror

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/gofrs/uuid"
)

type PLTError struct {
	message    string
	Message    string
	status     int // HTTP status codes as registered with IANA.
	parameters map[string]interface{}
	stackTrace []string
	ID         uuid.UUID // request uuid
}

func (e PLTError) Error() string {
	return e.message
}

func (e PLTError) Status() int {
	return e.status
}

func (e PLTError) Params() map[string]interface{} {
	return e.parameters
}

func (e *PLTError) AddStackTraceItem(item string) {
	e.stackTrace = append(e.stackTrace, item)
}

func (e PLTError) PrintStackTrace() string {
	res := fmt.Sprintf("%s:", e.message)
	for i := len(e.stackTrace) - 1; i >= 0; i-- {
		res = fmt.Sprintf("%s\n\t%s", res, e.stackTrace[i])
	}
	return res
}

func LogError(message string) {
	Logger.Errorln(message)
}

// Log displays an error with the right format using the given message.
func (e *PLTError) Log(message string) { //nolint
	status := e.status // TODO verify

	Logger.Warnf("%s %s %s\n%s", e.ID.String(), e.Error(), status, e.GenerateStackTrace())

	e.ResetStackTrace()
}

// ResetStackTrace reset stack trace of an HGErr.
func (e *PLTError) ResetStackTrace() {
	e.stackTrace = []string{}
}

// GenerateStackTrace joins the string to create the new stack trace message
func (e PLTError) GenerateStackTrace() string {
	if len(e.stackTrace) == 0 { // empty stacktrace
		return "[]"
	}

	// create stack trace
	return strings.Join(e.stackTrace, "\n")
}

func LogFatalError(message string) {
	LogError(message)
	os.Exit(1)
}

func LogErrorWithCode(message string, code int) error {
	return &PLTError{message: "ERROR: " + message, status: code}
}

// LogErrorsResp logs an error response (no Internal Server Error) with the appropriate format.
func LogErrorsResp(method string, url string, errorMsg string) {
	colorYellow := "\033[1;33m"
	noColor := "\033[0m"
	log.Printf("%s[ERROR RESPONSE] %s %s %s %s\n", colorYellow, method, url, errorMsg, noColor)
}

// LogWarning logs a warning message in the right format.
func LogWarning(message string) {
	colorYellow := "\033[1;33m"
	noColor := "\033[0m"
	log.Println(colorYellow + "[WARNING] " + message + noColor)
}

func PropagateError(err error, skips int) error {
	if err == nil {
		return nil
	}

	appErr, ok := err.(*PLTError)
	if !ok {
		appErr = ErrServerError
	}

	pc, file, line, _ := runtime.Caller(skips)
	funcName := runtime.FuncForPC(pc).Name()

	appErr.AddStackTraceItem(fmt.Sprintf("[%s:%v:%s %s]", file, line, funcName, err.Error()))

	return appErr
}

// Log message types are used to define the type of each log.
var (
	LogMessageErrorResponse   = "ERROR RESPONSE"   // Handled errors that could happen.
	LogMessageUnexpectedError = "UNEXPECTED ERROR" // Panic errors that should never happen.
)
