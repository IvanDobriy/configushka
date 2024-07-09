package configurator

import "io"

type Agent interface {
	Update(r io.Reader)
}
