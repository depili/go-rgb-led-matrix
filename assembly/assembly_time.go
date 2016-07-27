package assembly

import (
	"time"
)

type AssemblyTime struct {
	time.Time
}

const ctLayout = "2006-01-02T15:04-0700"

func (ct *AssemblyTime) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	ct.Time, err = time.Parse(ctLayout, string(b))
	return
}

func (ct *AssemblyTime) MarshalJSON() ([]byte, error) {
	return []byte(ct.Time.Format(ctLayout)), nil
}

var nilTime = (time.Time{}).UnixNano()

func (ct *AssemblyTime) IsSet() bool {
	return ct.UnixNano() != nilTime
}
