package policies

import "github.com/meshery/meshkit/errors"

// Please reference the following before contributing an error code:
// https://docs.meshery.io/project/contributing/contributing-error
// https://github.com/meshery/meshkit/blob/master/errors/errors.go
const (
	ErrApplyPatchCode       = "meshery-server-1376"
	ErrParsePayloadCode     = "meshery-server-1377"
	ErrFlappingDetectedCode = "meshery-server-1378"
)

func ErrApplyPatch(err error) error {
	return errors.New(ErrApplyPatchCode, errors.Alert, []string{"failed to apply configuration patches"}, []string{err.Error()}, []string{"component configuration may be incompatible with patch"}, []string{"verify the relationship definition's mutatorRef and mutatedRef paths"})
}

func ErrParsePayload(err error) error {
	return errors.New(ErrParsePayloadCode, errors.Alert, []string{"failed to parse action payload"}, []string{err.Error()}, []string{"action value may be malformed"}, []string{"check the policy engine action output"})
}

func ErrFlappingDetected(designID string) error {
	return errors.New(
		ErrFlappingDetectedCode,
		errors.Alert,
		[]string{"policy evaluation flapping detected"},
		[]string{"design " + designID + " produced alternating results across consecutive evaluations"},
		[]string{"two or more relationship policies are emitting patches that invert each other"},
		[]string{"review conflicting mutator/mutated refs; consider assigning distinct policy priorities"},
	)
}
