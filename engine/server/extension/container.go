// Package extension
// Created by xuzhuoxi
// on 2019-02-26.
// @author xuzhuoxi
//
package extension

import (
	"errors"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
)

func NewIExtensionContainer() server.IExtensionContainer {
	return NewExtensionContainer()
}

func NewExtensionContainer() *ExtensionContainer {
	return &ExtensionContainer{extensionMap: make(map[string]server.IExtension)}
}

type ExtensionContainer struct {
	extensions   []server.IExtension
	extensionMap map[string]server.IExtension
}

func (c *ExtensionContainer) AppendExtension(extension server.IExtension) {
	extName := extension.ExtensionName()
	if c.checkMap(extName) {
		panic("Repeat Name In Map: " + extName)
	}
	c.extensionMap[extName] = extension
	c.extensions = append(c.extensions, extension)
}

func (c *ExtensionContainer) CheckExtension(extName string) bool {
	_, ok := c.extensionMap[extName]
	return ok
}

func (c *ExtensionContainer) GetExtension(extName string) server.IExtension {
	if !c.checkMap(extName) {
		return nil
	}
	rs, _ := c.extensionMap[extName]
	return rs
}

func (c *ExtensionContainer) Len() int {
	return len(c.extensions)
}

func (c *ExtensionContainer) Extensions() []server.IExtension {
	ln := len(c.extensions)
	if 0 == ln {
		return nil
	}
	cp := make([]server.IExtension, ln)
	copy(cp, c.extensions)
	return cp
}

func (c *ExtensionContainer) ExtensionsReversed() []server.IExtension {
	ln := len(c.extensions)
	if 0 == ln {
		return nil
	}
	cp := make([]server.IExtension, ln)
	for i, j := 0, ln-1; i < j; i, j = i+1, j-1 {
		cp[i], cp[j] = c.extensions[j], c.extensions[i]
	}
	return cp
}

func (c *ExtensionContainer) Range(handler func(index int, extension server.IExtension)) {
	for index, extension := range c.extensions {
		handler(index, extension)
	}
}

func (c *ExtensionContainer) RangeReverse(handler func(index int, extension server.IExtension)) {
	ln := len(c.extensions)
	for index := ln - 1; index >= 0; index-- {
		handler(index, c.extensions[index])
	}
}

func (c *ExtensionContainer) HandleAt(index int, handler func(index int, extension server.IExtension)) error {
	if index < 0 || index >= len(c.extensions) {
		return errors.New("HandleAt Error : Out of index! ")
	}
	handler(index, c.extensions[index])
	return nil
}

func (c *ExtensionContainer) HandleAtName(name string, handler func(name string, extension server.IExtension)) error {
	if !c.CheckExtension(name) {
		return errors.New("HandleAtName Error : No such name [" + name + "]")
	}
	handler(name, c.extensionMap[name])
	return nil
}

func (c *ExtensionContainer) checkMap(name string) bool {
	_, ok := c.extensionMap[name]
	return ok
}

func NewIRabbitExtensionContainer() server.IRabbitExtensionContainer {
	return NewRabbitExtensionContainer()
}

func NewRabbitExtensionContainer() *RabbitExtensionContainer {
	return &RabbitExtensionContainer{ExtensionContainer: *NewExtensionContainer()}
}

type RabbitExtensionContainer struct {
	ExtensionContainer
}

func (c *RabbitExtensionContainer) InitExtensions() []error {
	ln := c.Len()
	if ln == 0 {
		return nil
	}
	var rs []error
	c.Range(func(_ int, extension server.IExtension) {
		if e, ok := extension.(server.IInitExtension); ok {
			err := e.InitExtension()
			rs = appendError(rs, err)
		}
	})
	return rs
}

func (c *RabbitExtensionContainer) DestroyExtensions() []error {
	ln := c.Len()
	if ln == 0 {
		return nil
	}
	var rs []error
	c.RangeReverse(func(_ int, extension server.IExtension) {
		if e, ok := extension.(server.IInitExtension); ok {
			err := e.DestroyExtension()
			rs = appendError(rs, err)
		}
	})
	return rs
}

func (c *RabbitExtensionContainer) SaveExtensions() []error {
	ln := c.Len()
	if ln == 0 {
		return nil
	}
	var rs []error
	c.Range(func(_ int, extension server.IExtension) {
		if e, ok := extension.(server.ISaveExtension); ok {
			err := e.SaveExtension()
			if nil != err {
				rs = append(rs, err)
			}
		}
	})
	return rs
}

func (c *RabbitExtensionContainer) SaveExtension(name string) error {
	var err error
	c.HandleAtName(name, func(_ string, extension server.IExtension) {
		if e, ok := extension.(server.ISaveExtension); ok {
			err = e.SaveExtension()
		}
	})
	return err
}

func (c *RabbitExtensionContainer) EnableExtensions(enable bool) []error {
	ln := c.Len()
	if ln == 0 {
		return nil
	}
	var rs []error
	if enable {
		c.Range(func(_ int, extension server.IExtension) {
			if e, ok := extension.(server.IEnableExtension); ok && !e.Enable() {
				err := e.EnableExtension()
				rs = appendError(rs, err)
			}
		})
	} else {
		c.RangeReverse(func(_ int, extension server.IExtension) {
			if e, ok := extension.(server.IEnableExtension); ok && e.Enable() {
				err := e.DisableExtension()
				rs = appendError(rs, err)
			}
		})
	}
	return rs
}

func (c *RabbitExtensionContainer) EnableExtension(name string, enable bool) error {
	var err error
	c.HandleAtName(name, func(_ string, extension server.IExtension) {
		if e, ok := extension.(server.IEnableExtension); ok {
			if e.Enable() != enable {
				if enable {
					err = e.EnableExtension()
				} else {
					err = e.DisableExtension()
				}
			}
		}
	})
	return err
}

func appendError(errs []error, err error) []error {
	if nil != err {
		return append(errs, err)
	} else {
		return errs
	}
}
