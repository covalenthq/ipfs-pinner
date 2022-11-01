//nolint:unused
package pinclient

import (
	"fmt"

	"github.com/multiformats/go-multiaddr"
)

const maxNameSize = 255

// TODO: We should probably make sure there are no duplicates sent
type addSettings struct {
	name    string
	origins []string
	meta    map[string]string
}

type AddOption func(options *addSettings) error

type pinAddOpts struct{}

func (pinAddOpts) WithName(name string) AddOption {
	return func(options *addSettings) error {
		if len(name) > maxNameSize {
			return fmt.Errorf("name cannot be longer than %d", maxNameSize)
		}
		options.name = name
		return nil
	}
}

func (pinAddOpts) WithOrigins(origins ...multiaddr.Multiaddr) AddOption {
	return func(options *addSettings) error {
		for _, o := range origins {
			options.origins = append(options.origins, o.String())
		}
		return nil
	}
}

func (pinAddOpts) AddMeta(meta map[string]string) AddOption {
	return func(options *addSettings) error {
		options.meta = meta
		return nil
	}
}
