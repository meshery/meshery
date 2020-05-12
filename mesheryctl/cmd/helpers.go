package cmd

import (
	"io"

	log "github.com/sirupsen/logrus"
)

// SafeClose is a helper function help to close the io
func SafeClose(co io.Closer) {
	if cerr := co.Close(); cerr != nil {
		log.Error(cerr)
	}
}
