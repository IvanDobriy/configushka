package configurator

import "io"

type Registry interface {
	Get(key string) (r io.Reader)
}
