package gonertia

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// Option is an option parameter that modifies Inertia.
type Option func(i *Inertia) error

// WithVersion returns Option that will set Inertia's version.
func WithVersion(version string) Option {
	return func(i *Inertia) error {
		i.version = md5(version)
		return nil
	}
}

// WithVersionFromFile returns Option that will set Inertia's version based on file checksum.
func WithVersionFromFile(path string) Option {
	return func(i *Inertia) (err error) {
		i.version, err = md5File(path)
		if err != nil {
			return fmt.Errorf("calculating md5 hash of manifest file: %w", err)
		}
		return nil
	}
}

// WithJSONMarshaller returns Option that will set Inertia's JSON marshaller.
func WithJSONMarshaller(jsonMarshaller JSONMarshaller) Option {
	return func(i *Inertia) error {
		i.jsonMarshaller = jsonMarshaller
		return nil
	}
}

// WithLogger returns Option that will set Inertia's logger.
func WithLogger(logs ...Logger) Option {
	var l Logger
	if len(logs) > 0 {
		l = logs[0]
	} else {
		l = log.Default()
	}

	if l == nil {
		l = log.New(io.Discard, "", 0)
	}

	return func(i *Inertia) error {
		i.logger = l
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

// WithSSR returns Option that will enable server side rendering on Inertia.
func WithSSR(url ...string) Option {
	return func(i *Inertia) error {
		var u string
		if len(url) > 0 {
			u = url[0]
		} else {
			const defaultURL = "http://127.0.0.1:13714"
			u = defaultURL
		}

		i.ssrURL = u
		i.ssrHTTPClient = &http.Client{}
		return nil
	}
}

// WithFlashProvider returns Option that will set Inertia's flash data provider.
func WithFlashProvider(flashData FlashProvider) Option {
	return func(i *Inertia) error {
		i.flash = flashData
		return nil
	}
}
