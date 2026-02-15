package policies

import "github.com/meshery/meshkit/errors"

// Please reference the following before contributing an error code:
// https://docs.meshery.io/project/contributing/contributing-error
// https://github.com/meshery/meshkit/blob/master/errors/errors.go
const (
	ErrConvertDesignCode  = "meshery-server-1378"
	ErrConvertResultCode  = "meshery-server-1375"
	ErrApplyPatchCode     = "meshery-server-1376"
	ErrParsePayloadCode   = "meshery-server-1377"
)

func ErrConvertDesign(err error) error {
	return errors.New(ErrConvertDesignCode, errors.Alert, []string{"failed to convert design for policy evaluation"}, []string{err.Error()}, []string{"design file may be malformed"}, []string{"check the design file structure"})
}

func ErrConvertResult(err error) error {
	return errors.New(ErrConvertResultCode, errors.Alert, []string{"failed to convert evaluation result"}, []string{err.Error()}, []string{"evaluation may have produced invalid output"}, []string{"check the policy engine output"})
}

func ErrApplyPatch(err error) error {
	return errors.New(ErrApplyPatchCode, errors.Alert, []string{"failed to apply configuration patches"}, []string{err.Error()}, []string{"component configuration may be incompatible with patch"}, []string{"verify the relationship definition's mutatorRef and mutatedRef paths"})
}

func ErrParsePayload(err error) error {
	return errors.New(ErrParsePayloadCode, errors.Alert, []string{"failed to parse action payload"}, []string{err.Error()}, []string{"action value may be malformed"}, []string{"check the policy engine action output"})
}
