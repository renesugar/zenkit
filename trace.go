package zenkit

import (
	"regexp"
	"runtime"
	"sync"
)

var (
	fnNames = fnNameCache{
		mu:    sync.Mutex{},
		items: make(map[uintptr]string),
	}
	stripFnPreamble = regexp.MustCompile(`^.*\.(.*)$`)
)

type fnNameCache struct {
	mu    sync.Mutex
	items map[uintptr]string
}

func (c *fnNameCache) Get(key uintptr, factory func() string) (data string, found bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	data, found = c.items[key]
	if !found {
		data = factory()
		c.items[key] = data
	}
	return
}

func funcName() string {
	fnName := "<unknown>"
	pc, _, _, ok := runtime.Caller(3)
	if ok {
		fnName, _ = fnNames.Get(pc, func() string {
			return stripFnPreamble.ReplaceAllString(runtime.FuncForPC(pc).Name(), "$1")
		})
	}
	return fnName
}
