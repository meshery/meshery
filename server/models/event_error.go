package models

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	stderrors "errors"
	"net"
	"syscall"

	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// An error event carries one occurrence of a MeshKit error definition as its
// metadata.error, in the schemas EventError shape: the code, the descriptions
// rendered for this occurrence, the causes nested with their remediations, and
// matchedCauses naming the causes that apply this time
// (docs/error-definition-contract.md in meshery/schemas).

// matchedCauseError is a MeshKit error together with the id of the probable
// cause the failure it reports points at. It is transparent to every reader of
// the MeshKit error: Error, errors.As (Unwrap) and JSON (MarshalJSON) all see
// the wrapped error. Build one with withMatchedCause.
type matchedCauseError struct {
	err     error
	causeID string
}

func (e *matchedCauseError) Error() string { return e.err.Error() }

func (e *matchedCauseError) Unwrap() error { return e.err }

// MarshalJSON keeps the serialized form of the MeshKit error, which events
// that still carry the legacy metadata.error shape store.
func (e *matchedCauseError) MarshalJSON() ([]byte, error) { return json.Marshal(e.err) }

// withMatchedCause records that err's failure points at the cause with id
// causeID. An empty causeID (the cause could not be told) returns err as is.
func withMatchedCause(err error, causeID string) error {
	if err == nil || causeID == "" {
		return err
	}
	return &matchedCauseError{err: err, causeID: causeID}
}

// MatchedCauses returns the probable causes of the MeshKit error in err that
// Meshery Server found apply to this occurrence, most likely first, for an
// event's metadata.error.matchedCauses. It is empty when the server could not
// tell which cause applies.
func MatchedCauses(err error) []events.EventMatchedCause {
	var matched *matchedCauseError
	if !stderrors.As(err, &matched) {
		return nil
	}
	return []events.EventMatchedCause{events.MatchCause(matched.causeID)}
}

// AttachEventError sets builder's metadata.error to err as an EventError,
// naming the causes MatchedCauses finds, and keeps every other metadata key.
// Call it after WithMetadata, which replaces the whole map.
//
// When err carries no MeshKit error, or its EventError would break the schemas
// contract, metadata.error keeps the legacy form (err itself, which serializes
// with MeshKit's Go field names) and ErrAttachEventError is logged, so the
// event still reports the failure. A nil log is allowed.
func AttachEventError(builder *events.EventBuilder, err error, log logger.Handler) *events.EventBuilder {
	if err == nil {
		return builder
	}
	if _, attachErr := builder.WithError(err, MatchedCauses(err)...); attachErr != nil {
		if log != nil {
			log.Warn(ErrAttachEventError(attachErr))
		}
		// WithError leaves the event unchanged on failure, so the legacy key is
		// set here, keeping the rest of the metadata.
		event := builder.Build()
		if event.Metadata == nil {
			event.Metadata = map[string]interface{}{}
		}
		event.Metadata[events.MetadataErrorKey] = err
	}
	return builder
}

// classifyKubeAPIError tells which probable cause of ErrUnreachableKubeAPI a
// Kubernetes client error points at, or "" when it cannot tell. Certificate
// failures are checked first: the client reports them inside a *url.Error,
// which is also a net.Error.
func classifyKubeAPIError(err error) string {
	if err == nil {
		return ""
	}
	var (
		unknownAuthority x509.UnknownAuthorityError
		certInvalid      x509.CertificateInvalidError
		hostname         x509.HostnameError
		verification     *tls.CertificateVerificationError
	)
	if stderrors.As(err, &unknownAuthority) || stderrors.As(err, &certInvalid) ||
		stderrors.As(err, &hostname) || stderrors.As(err, &verification) {
		return KubeAPICauseCertificateUntrusted
	}
	if apierrors.IsUnauthorized(err) {
		return KubeAPICauseCredentialsRejected
	}
	var netErr net.Error
	if stderrors.As(err, &netErr) || stderrors.Is(err, context.DeadlineExceeded) ||
		stderrors.Is(err, syscall.ECONNREFUSED) || stderrors.Is(err, syscall.EHOSTUNREACH) {
		return KubeAPICauseUnreachable
	}
	return ""
}
