package gonertia

import (
	"fmt"
	"io"
	"io/fs"
	"log"
)

// Option is an option parameter that modifies Inertia.
type Option func(i *Inertia) error

// WithTemplateFS returns Option that will set Inertia's templateFS.
func WithTemplateFS(templateFS fs.FS) Option {
	return func(i *Inertia) error {
		i.templateFS = templateFS
		return nil
	}
}

// WithVersion returns Option that will set Inertia's version.
func WithVersion(version string) Option {
	return func(i *Inertia) error {
		i.version = version
		return nil
	}
}

// WithAssetURL returns Option that will set Inertia's version based on asset url.
func WithAssetURL(url string) Option {
	return WithVersion(md5(url))
}

// WithManifestFile returns Option that will set Inertia's version based on manifest file.
func WithManifestFile(path string) Option {
	version, err := md5File(path)
	if err == nil {
		return WithVersion(version)
	}

	return func(i *Inertia) error {
		return fmt.Errorf("calculating md5 hash of manifest file: %w", err)
	}
}

// WithMarshalJSON returns Option that will set Inertia's marshallJSON func.
func WithMarshalJSON(f marshallJSON) Option {
	return func(i *Inertia) error {
		i.marshallJSON = f
		return nil
	}
}

// WithLogger returns Option that will set Inertia's logger.
func WithLogger(log logger) Option {
	if log == nil {
		return WithoutLogger()
	}

	return func(i *Inertia) error {
		i.logger = log
		return nil
	}
}

// WithoutLogger returns Option that will unset Inertia's logger.
// Actually set a logger with io.Discard output.
func WithoutLogger() Option {
	return func(i *Inertia) error {
		i.logger = log.New(io.Discard, "", 0)
		return nil
	}
}

// WithContainerID returns Option that will set Inertia's container id.
func WithContainerID(id string) Option {
	return func(i *Inertia) error {
		i.containerID = id
		return nil
	}
}
