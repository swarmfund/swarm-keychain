package problem

import (
	"encoding/json"
	"fmt"
	"net/http"

	"time"

	"bullioncoin.githost.io/development/keychain/context/requestid"
	"bullioncoin.githost.io/development/keychain/log"
	"bullioncoin.githost.io/development/keychain/utils"
	"github.com/go-errors/errors"
	"golang.org/x/net/context"
)

var (
	errToProblemMap = map[error]P{}
)

// RegisterError records an error -> P mapping, allowing the app to register
// specific errors that may occur in other packages to be rendered as a specific
// P instance.
//
// For example, you might want to render any sql.ErrNoRows errors as a
// problem.NotFound, and you would do so by calling:
//
// problem.RegisterError(sql.ErrNoRows, problem.NotFound) in you application
// initialization sequence
func RegisterError(err error, p P) {
	errToProblemMap[err] = p
}

// HasProblem types can be transformed into a problem.
// Implement it for custom errors.
type HasProblem interface {
	Problem() P
}

// P is a struct that represents an error response to be rendered to a connected
// client.
type P struct {
	Type     string                 `json:"type"`
	Title    string                 `json:"title"`
	Status   int                    `json:"status"`
	Detail   string                 `json:"detail,omitempty"`
	Instance string                 `json:"instance,omitempty"`
	Extras   map[string]interface{} `json:"extras,omitempty"`
}

func (p *P) Error() string {
	return fmt.Sprintf("problem: %s", p.Type)
}

// Inflate expands a problem with contextal information.
// At present it adds the request's id as the problem's Instance, if available.
func Inflate(ctx context.Context, p *P) {
	p.Instance = requestid.FromContext(ctx)
}

// Render writes a http response to `w`, compliant with the "Problem
// Details for HTTP APIs" RFC:
//   https://tools.ietf.org/html/draft-ietf-appsawg-http-problem-00
//
// `p` is the problem, which may be either a concrete P struct, an implementor
// of the `HasProblem` interface, or an error.  Any other value for `p` will
// panic.
func Render(ctx context.Context, w http.ResponseWriter, p interface{}) {
	switch p := p.(type) {
	case P:
		render(ctx, w, p)
	case *P:
		render(ctx, w, *p)
	case HasProblem:
		render(ctx, w, p.Problem())
	case error:
		renderErr(ctx, w, p)
	default:
		panic(fmt.Sprintf("Invalid problem: %v+", p))
	}
}

func render(ctx context.Context, w http.ResponseWriter, p P) {

	Inflate(ctx, &p)

	w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
	js, err := json.MarshalIndent(p, "", "  ")

	if err != nil {
		err := errors.Wrap(err, 1)
		log.Ctx(ctx).WithStack(err).Error(err)
		http.Error(w, "error rendering problem", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(p.Status)
	w.Write(js)
}

func renderErr(ctx context.Context, w http.ResponseWriter, err error) {
	origErr := err

	if err, ok := err.(*errors.Error); ok {
		origErr = err.Err
	}

	p, ok := errToProblemMap[origErr]

	// If this error is not a registered error
	// log it and replace it with a 500 error
	if !ok {
		log.Ctx(ctx).WithStack(err).Error(err)
		p = ServerError
	}

	render(ctx, w, p)
}

// Well-known and reused problems below:
var (
	// NotFound is a well-known problem type.  Use it as a shortcut
	// in your actions.
	NotFound = P{
		Type:   "not_found",
		Title:  "Resource Missing",
		Status: http.StatusNotFound,
		Detail: "The resource at the url requested was not found.  This is usually " +
			"occurs for one of two reasons:  The url requested is not valid, or no " +
			"data in our database could be found with the parameters provided.",
	}

	Success = P{
		Type:   "success",
		Title:  "Success",
		Status: http.StatusOK,
	}

	// ServerError is a well-known problem type.  Use it as a shortcut
	// in your actions.
	ServerError = P{
		Type:   "server_error",
		Title:  "Internal Server Error",
		Status: http.StatusInternalServerError,
		Detail: "An error occurred while processing this request.  This is usually due " +
			"to a bug within the server software.  Trying this request again may " +
			"succeed if the bug is transient, otherwise please report this issue " +
			"to the issue tracker at: https://bullioncoin.githost.io/development/horizon/issues." +
			" Please include this response in your issue.",
	}

	// RateLimitExceeded is a well-known problem type.  Use it as a shortcut
	// in your actions.
	RateLimitExceeded = P{
		Type:   "rate_limit_exceeded",
		Title:  "Rate limit exceeded",
		Status: 429,
		Detail: "The rate limit for the requesting IP address is over its alloted " +
			"limit.  The allowed limit and requests left per time period are " +
			"communicated to clients via the http response headers 'X-RateLimit-*' " +
			"headers.",
	}

	// NotImplemented is a well-known problem type.  Use it as a shortcut
	// in your actions.
	NotImplemented = P{
		Type:   "not_implemented",
		Title:  "Resource Not Yet Implemented",
		Status: http.StatusNotFound,
		Detail: "While the requested URL is expected to eventually point to a " +
			"valid resource, the work to implement the resource has not yet " +
			"been completed.",
	}

	// NotAcceptable is a well-known problem type.  Use it as a shortcut
	// in your actions.
	NotAcceptable = P{
		Type: "not_acceptable",
		Title: "An acceptable response content-type could not be provided for " +
			"this request",
		Status: http.StatusNotAcceptable,
	}

	// BadRequest is a well-known problem type.  Use it as a shortcut
	// in your actions.
	BadRequest = P{
		Type:   "bad_request",
		Title:  "Bad Request",
		Status: http.StatusBadRequest,
		Detail: "The request you sent was invalid in some way",
	}

	// ServerOverCapacity is a well-known problem type.  Use it as a shortcut
	// in your actions.
	ServerOverCapacity = P{
		Type:   "server_over_capacity",
		Title:  "Server Over Capacity",
		Status: http.StatusServiceUnavailable,
		Detail: "This horizon server is currently overloaded.  Please wait for " +
			"several minutes before trying your request again.",
	}

	// Timeout is a well-known problem type.  Use it as a shortcut
	// in your actions.
	Timeout = P{
		Type:   "timeout",
		Title:  "Timeout",
		Status: http.StatusGatewayTimeout,
		Detail: "Your request timed out before completing.  Please try your " +
			"request again.",
	}

	// UnsupportedMediaType is a well-known problem type.  Use it as a shortcut
	// in your actions.
	UnsupportedMediaType = P{
		Type:   "unsupported_media_type",
		Title:  "Unsupported Media Type",
		Status: http.StatusUnsupportedMediaType,
		Detail: "The request has an unsupported content type. Presently, the " +
			"only supported content type is application/x-www-form-urlencoded.",
	}

	// BeforeHistory is a well-known problem type.  Use it as a shortcut
	// in your actions.
	BeforeHistory = P{
		Type:   "before_history",
		Title:  "Data Requested Is Before Recorded History",
		Status: http.StatusGone,
		Detail: "This horizon instance is configured to only track a " +
			"portion of the stellar network's latest history. This request " +
			"is asking for results prior to the recorded history known to " +
			"this horizon instance.",
	}

	// StaleHistory is a well-known problem type.  Use it as a shortcut
	// in your actions.
	StaleHistory = P{
		Type:   "stale_history",
		Title:  "Historical DB Is Too Stale",
		Status: http.StatusServiceUnavailable,
		Detail: "This horizon instance is configured to reject client requests " +
			"when it can determine that the history database is lagging too far " +
			"behind the connected instance of stellar-core.  If you operate this " +
			"server, please ensure that the ingestion system is properly running.",
	}
	NotAllowed = P{
		Type:   "not_allowed",
		Title:  "It is not allowed to access this data uthing provided account id",
		Status: http.StatusUnauthorized,
		Detail: "Provided account id does not have rights to access requsted data",
	}
	SignNotVerified = P{
		Type:   "not_allowed",
		Title:  "Failed to verify your signature",
		Status: http.StatusUnauthorized,
		Detail: "Your signature is incorrect",
	}

	// Forbidden is a well-known problem type.  Use it as a shortcut
	// in your actions.
	Forbidden = P{
		Type:   "forbidden",
		Title:  "Forbidden",
		Status: http.StatusForbidden,
	}

	Conflict = P{
		Type:   "conflict",
		Title:  "Conflict",
		Status: http.StatusConflict,
		Detail: "Already exists",
	}

	// Only one operation is allowed for adminOperations.
	AdminOperationsRestrictionViolated = P{
		Type:   "not_acceptable",
		Title:  "Admin Operations Restriction Violated",
		Status: http.StatusBadRequest,
		Detail: "Only one ADMIN operation per admin tx allowed",
	}

	AdminOperationAlreadyExist = P{
		Type:   "already_submitted",
		Title:  "Admin Operations Already Submitted",
		Status: http.StatusBadRequest,
		Detail: "You can not add an admin operation, the same already exists",
	}

	AdminOperationAlreadySubmitted = P{
		Type:   "already_submitted",
		Title:  "Admin Operations Already Submitted",
		Status: http.StatusBadRequest,
		Detail: "You can not submit admin operation with enough signatures twice",
	}

	AdminOperationRejectRestrictionViolated = P{
		Type:   "not_acceptable",
		Title:  "Admin Operation Reject Restriction Violated",
		Status: http.StatusBadRequest,
		Detail: "Only pending operation can be rejected",
	}

	ExpiredTFA = P{
		Type:   "gone",
		Title:  "Gone",
		Status: http.StatusGone,
		Detail: "Your verification code has expired. Please request a new code",
	}

	AttemptsLimitExceeded = P{
		Type:   "too_many_requests",
		Title:  "Limit Exceeded",
		Status: http.StatusTooManyRequests,
		Detail: "You have exceeded the number of attempts to verify the code. Request a new code",
	}

	DayLimitExceeded = P{
		Type:   "too_many_requests",
		Title:  "Limit Exceeded",
		Status: http.StatusTooManyRequests,
		Detail: "You exceeded the daily limit of code requests. Request a new code tomorrow",
	}

	AdminOperationDeleteRestrictionViolated = P{
		Type:   "not_acceptable",
		Title:  "Admin Operation Delete Restriction Violated",
		Status: http.StatusBadRequest,
		Detail: "Only pending operation can be deleted",
	}

	RecoveryRequestLimitExceeded = P{
		Type:   "too_many_requests",
		Title:  "Limit Exceeded",
		Status: http.StatusTooManyRequests,
	}

	TFARequired = func(token, phoneMask string, retryIn *time.Duration) *P {
		var seconds int64
		if retryIn != nil {
			seconds = retryIn.Nanoseconds() / 1000000000
		}
		return &P{
			Type:   "tfa_required",
			Title:  "Two factor verification is required to access this resource",
			Status: http.StatusForbidden,
			Extras: map[string]interface{}{
				"token":      token,
				"phone_mask": utils.MaskPhoneNumber(phoneMask),
				"retry_in":   seconds,
			},
		}
	}
)
