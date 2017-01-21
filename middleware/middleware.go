package main

import (
	"fmt"
)

type rqHandler func(req string, resp string) error

type middleware func(rqHandler) rqHandler

type Chain struct {
	middlewares []middleware
}

func NewChain(middlewares ...middleware) *Chain {
	c := new(Chain)
	c.middlewares = append(c.middlewares, middlewares...)
	return c
}

func (c *Chain) Then(h rqHandler) rqHandler {
	for i := range c.middlewares {
		reversedI := len(c.middlewares) - i - 1
		hInChain := c.middlewares[reversedI]
		h = hInChain(h)
	}
	return h
}

func logBefore(nextHandler rqHandler) rqHandler {
	h := func(req string, resp string) error {
		fmt.Println("before")
		err := nextHandler(req, resp)

		return err
	}
	return h
}

func logAfter(nextHandler rqHandler) rqHandler {
	return func(req string, resp string) error {
		err := nextHandler(req, resp)
		fmt.Println("after")
		return err
	}
}

func real(req string, resp string) error {
	fmt.Println("in real work")
	return nil
}

func main() {
	middlewares := NewChain(logBefore, logAfter)
	call := middlewares.Then(real)
	call("in", "out")
}
