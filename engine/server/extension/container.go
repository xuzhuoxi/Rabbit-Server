// Package protox
// Created by xuzhuoxi
// on 2019-02-26.
// @author xuzhuoxi
//
package extension

import (
	"errors"
	"github.com/xuzhuoxi/Rabbit-Server/engine/server"
)

func NewIProtoExtensionContainer() server.IProtoExtensionContainer {
	return &ProtocolContainer{ExtensionContainer: NewExtensionContainer()}
}

func NewProtoExtensionContainer() ProtocolContainer {
	return ProtocolContainer{ExtensionContainer: NewExtensionContainer()}
}

func NewIExtensionContainer() server.IExtensionContainer {
	return &ExtensionContainer{extensionMap: make(map[string]server.IExtension)}
}

func NewExtensionContainer() ExtensionContainer {
	return ExtensionContainer{extensionMap: make(map[string]server.IExtension)}
}

type ExtensionContainer struct {
	extensions   []server.IExtension
	extensionMap map[string]server.IExtension
}

func (m *ExtensionContainer) AppendExtension(extension server.IExtension) {
	name := extension.ExtensionName()
	if m.checkMap(name) {
		panic("Repeat Name In Map: " + name)
	}
	m.extensionMap[name] = extension
	m.extensions = append(m.extensions, extension)
}

func (m *ExtensionContainer) CheckExtension(name string) bool {
	_, ok := m.extensionMap[name]
	return ok
}

func (m *ExtensionContainer) GetExtension(name string) server.IExtension {
	if !m.checkMap(name) {
		return nil
	}
	rs, _ := m.extensionMap[name]
	return rs
}

func (m *ExtensionContainer) Len() int {
	return len(m.extensions)
}

func (m *ExtensionContainer) Extensions() []server.IExtension {
	ln := len(m.extensions)
	if 0 == ln {
		return nil
	}
	cp := make([]server.IExtension, ln)
	copy(cp, m.extensions)
	return cp
}

func (m *ExtensionContainer) ExtensionsReversed() []server.IExtension {
	ln := len(m.extensions)
	if 0 == ln {
		return nil
	}
	cp := make([]server.IExtension, ln)
	for i, j := 0, ln-1; i < j; i, j = i+1, j-1 {
		cp[i], cp[j] = m.extensions[j], m.extensions[i]
	}
	return cp
}

func (m *ExtensionContainer) Range(handler func(index int, extension server.IExtension)) {
	for index, extension := range m.extensions {
		handler(index, extension)
	}
}

func (m *ExtensionContainer) RangeReverse(handler func(index int, extension server.IExtension)) {
	ln := len(m.extensions)
	for index := ln - 1; index >= 0; index-- {
		handler(index, m.extensions[index])
	}
}

func (m *ExtensionContainer) HandleAt(index int, handler func(index int, extension server.IExtension)) error {
	if index < 0 || index >= len(m.extensions) {
		return errors.New("HandleAt Error : Out of index! ")
	}
	handler(index, m.extensions[index])
	return nil
}

func (m *ExtensionContainer) HandleAtName(name string, handler func(name string, extension server.IExtension)) error {
	if !m.CheckExtension(name) {
		return errors.New("HandleAtName Error : No such name [" + name + "]")
	}
	handler(name, m.extensionMap[name])
	return nil
}

func (m *ExtensionContainer) checkMap(name string) bool {
	_, ok := m.extensionMap[name]
	return ok
}

type ProtocolContainer struct {
	ExtensionContainer
}

func (c *ProtocolContainer) InitExtensions() []error {
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

func (c *ProtocolContainer) DestroyExtensions() []error {
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

func (c *ProtocolContainer) SaveExtensions() []error {
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

func (c *ProtocolContainer) SaveExtension(name string) error {
	var err error
	c.HandleAtName(name, func(_ string, extension server.IExtension) {
		if e, ok := extension.(server.ISaveExtension); ok {
			err = e.SaveExtension()
		}
	})
	return err
}

func (c *ProtocolContainer) EnableExtensions(enable bool) []error {
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

func (c *ProtocolContainer) EnableExtension(name string, enable bool) error {
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
