package tcpx

import "github.com/DarthPestilane/easytcp"

func Init(opts ...Option) (serve *easytcp.Server, e error) {
	custom := Options{}

	for _, opt := range opts {
		opt(&custom)
	}

	opt := easytcp.ServerOption{
		Packer: easytcp.NewDefaultPacker(),
		Codec:  &easytcp.ProtobufCodec{},
	}

	serve = easytcp.NewServer(&opt)

	serve.AddRoute(1, func(c easytcp.Context) {
		c.SetResponseMessage(easytcp.NewMessage(2, []byte("pong")))
	})

	return
}
