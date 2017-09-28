package funcname

import (
	"regexp"
	"runtime"
	"sync"
)

var (
	fnNames         = NewFnNameCache()
	stripFnPreamble = regexp.MustCompile(`^.*\.(.*)$`)
)

type FnNameCache struct {
	mu    sync.Mutex
	items map[uintptr]string
}

func (c *FnNameCache) Get(key uintptr, factory func() string) (data string, found bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	data, found = c.items[key]
	if !found {
		data = factory()
		c.items[key] = data
	}
	return
}

func NewFnNameCache() FnNameCache {
	return FnNameCache{
		mu:    sync.Mutex{},
		items: make(map[uintptr]string),
	}
}

func FuncName(level int) string {
	fnName := "<unknown>"
	pc, _, _, ok := runtime.Caller(level)
	if ok {
		fnName, _ = fnNames.Get(pc, func() string {
			return stripFnPreamble.ReplaceAllString(runtime.FuncForPC(pc).Name(), "$1")
		})
	}
	return fnName
}
