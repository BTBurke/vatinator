package vat

import "github.com/BTBurke/vatinator/svc"

type options struct {
	zipOutputDir string
	zipFilename  string
	hooks        svc.Hooks
	formData     FormData
}

type Option func(c *options)

func WithZipOutput(dir string, fname string) Option {
	return func(c *options) {
		c.zipOutputDir = dir
		c.zipFilename = fname
	}
}

func WithHooks(h svc.Hooks) Option {
	return func(c *options) {
		c.hooks = h
	}
}

func WithFormData(fd FormData) Option {
	return func(c *options) {
		c.formData = fd
	}
}
