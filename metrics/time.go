package metrics

import "time"

var TimeFunc = func() time.Time {
	return time.Now().UTC()
}
