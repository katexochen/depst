package dependonstd

import (
	"bytes"
	"errors"
	"fmt"
	"time"
)

var (
	err    = errors.New("failed")
	errTwo = fmt.Errorf("failed")
	buf    = bytes.Buffer{}
	T      = time.Time{}
)
