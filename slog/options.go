package slog

import (
	"strconv"
	"strings"
)

type CallerFormater func(file, function string, line int) any

func ShortFormater() CallerFormater {
	return func(file, function string, line int) any {
		idx := strings.LastIndexByte(file, '/')
		if idx == -1 {
			return file
		}
		idx = strings.LastIndexByte(file[:idx], '/')
		if idx == -1 {
			return file
		}
		return file[idx+1:] + ":" + strconv.Itoa(line)
	}
}

type HandlerOptions struct {
	mainPKG        string
	excludePKGList []string
	addCaller      bool
	callerFormater CallerFormater
}

type OptionFunc func(*HandlerOptions)

func WithCaller(mainPKG string, excludePKGList ...string) OptionFunc {
	return func(opts *HandlerOptions) {
		opts.mainPKG = mainPKG
		opts.excludePKGList = excludePKGList
		opts.addCaller = true
	}
}

func WithCallerFormater(formater CallerFormater) OptionFunc {
	return func(opts *HandlerOptions) {
		opts.callerFormater = formater
	}
}

func (opts HandlerOptions) Match(pkg string) bool {
	if pkg == "main" {
		return true
	}
	if !strings.HasPrefix(pkg, opts.mainPKG) {
		return false
	}
	if len(opts.excludePKGList) > 0 {
		for _, exclude := range opts.excludePKGList {
			if strings.HasPrefix(pkg, exclude) {
				return false
			}
		}
	}
	return true
}
