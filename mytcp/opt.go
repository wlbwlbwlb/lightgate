package mytcp

import (
	"fmt"
	"io"
	"strings"

	"github.com/DarthPestilane/easytcp"
)

type Options struct {
}

type Option func(*Options)

func Writer(w io.Writer) Option {
	return func(*Options) {
		easytcp.SetLogger(&WrapWriter{
			w: w,
		})
	}
}

type WrapWriter struct {
	w io.Writer
}

func (p *WrapWriter) Errorf(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(p.w, "[easytcp] "+format, args...)
}

func (p *WrapWriter) Tracef(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(p.w, "[easytcp] "+format, args...)
}
