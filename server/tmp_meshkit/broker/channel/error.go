package channel

import (
	"fmt"
	"strings"

	"github.com/meshery/meshkit/errors"
)

// TODO run error-util
const (
	ErrPublishCode = "replace_me"
)

func ErrPublish(err error) error {
	return errors.New(ErrPublishCode, errors.Alert, []string{"Publish failed"}, []string{err.Error()}, []string{"Publish to channel failed", "Subject channel buffer is full"}, []string{"Make sure there is a consumer from the subject"})
}

type ErrChannelPublishType struct {
	Err              error
	SuccessQueueList []string
	FailedQueueList  []string
}

func ErrChannelPublish(
	err error,
	successQueueList []string,
	failedQueueList []string,
) *ErrChannelPublishType {
	return &ErrChannelPublishType{
		Err:              err,
		SuccessQueueList: successQueueList,
		FailedQueueList:  failedQueueList,
	}
}

func (e *ErrChannelPublishType) Error() string {
	return fmt.Sprintf(
		"%s, success queue list: [%s], failed queue list [%s]",
		e.Err.Error(),
		strings.Join(e.SuccessQueueList, ","),
		strings.Join(e.FailedQueueList, ","),
	)
}
