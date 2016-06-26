package libwebawk

import (
	"strconv"
)

type Context struct {
	name     string
	value    string
	children map[string]*Context
	stats    map[string]int64
}

func NewContext(n string, v string) *Context {
	c := new(Context)
	c.name = n
	c.value = v
	c.children = make(map[string]*Context)
	c.stats = make(map[string]int64)
	return c
}

func (ctxt *Context) CreateChild(k string, v string) *Context {
	id, ok := ctxt.stats[k]
	if !ok {
		id = 0
	}
	nk := k + "[" + strconv.FormatInt(id, 10) + "]"
	ctxt.stats[k] = id + 1
	c := NewContext(nk, v)
	ctxt.children[nk] = c
	return c
}

func (ctxt *Context) AppendText(v string) {
	ctxt.value += v
}

func (ctxt *Context) GetValue(addr Address) string {
	if ctxt.name != addr.name[0] {
		return "NIL"
	}
	for i := 1; i < len(addr.name); i++ {
		c, ok := ctxt.children[addr.name[i]]
		if !ok {
			return "NIL"
		}
		ctxt = c
	}
	return ctxt.value
}
