package helpers

import (
	stderrs "errors"
	"fmt"
	"testing"

	"github.com/meshery/meshkit/errors"
)

func TestErrInvalidK8SConfig(t *testing.T) {
	mockErr := fmt.Errorf("mock error")
	err := ErrInvalidK8SConfig(mockErr)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedCode := ErrInvalidK8SConfigCode
	expectedShortDesc := "No valid kubernetes config found"
	expectedLongDesc := "mock error"
	expectedProbableCause := "Kubernetes config is not accessible to meshery or not valid"
	expectedSuggestedRemediation := "Upload your kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"

	var e *errors.Error
	if stderrs.As(err, &e) {
		if e.Code != expectedCode {
			t.Errorf("Expected code %s, got %s", expectedCode, e.Code)
		}
		if len(e.ShortDescription) == 0 || e.ShortDescription[0] != expectedShortDesc {
			t.Errorf("Expected ShortDescription to be '%s', got '%v'", expectedShortDesc, e.ShortDescription)
		}
		if len(e.LongDescription) == 0 || e.LongDescription[0] != expectedLongDesc {
			t.Errorf("Expected LongDescription to be '%s', got '%v'", expectedLongDesc, e.LongDescription)
		}
		if len(e.ProbableCause) == 0 || e.ProbableCause[0] != expectedProbableCause {
			t.Errorf("Expected ProbableCause to be '%s', got '%v'", expectedProbableCause, e.ProbableCause)
		}
		if len(e.SuggestedRemediation) == 0 || e.SuggestedRemediation[0] != expectedSuggestedRemediation {
			t.Errorf("Expected SuggestedRemediation to be '%s', got '%v'", expectedSuggestedRemediation, e.SuggestedRemediation)
		}
		if e.Severity != errors.Alert {
			t.Errorf("Expected Severity to be %v, got %v", errors.Alert, e.Severity)
		}
	} else {
		t.Errorf("Expected error to be of type *errors.Error, got %T", err)
	}
}
