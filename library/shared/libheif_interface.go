package shared

import (
	"fmt"
	"net/rpc"

	"github.com/klippa-app/go-libheif/library/requests"
	"github.com/klippa-app/go-libheif/library/responses"

	"github.com/hashicorp/go-plugin"
)

type Libheif interface {
	Ping() (string, error)
	DecodeImage(*requests.DecodeImage) (*responses.DecodeImage, error)
	DecodeConfig(*requests.DecodeConfig) (*responses.DecodeConfig, error)
	RenderFile(*requests.RenderFile) (*responses.RenderFile, error)
}

type LibheifRPC struct{ client *rpc.Client }

func (g *LibheifRPC) Ping() (string, error) {
	var resp string
	err := g.client.Call("Plugin.Ping", new(interface{}), &resp)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func (g *LibheifRPC) DecodeImage(request *requests.DecodeImage) (*responses.DecodeImage, error) {
	resp := &responses.DecodeImage{}
	err := g.client.Call("Plugin.DecodeImage", request, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (g *LibheifRPC) DecodeConfig(request *requests.DecodeConfig) (*responses.DecodeConfig, error) {
	resp := &responses.DecodeConfig{}
	err := g.client.Call("Plugin.DecodeConfig", request, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (g *LibheifRPC) RenderFile(request *requests.RenderFile) (*responses.RenderFile, error) {
	resp := &responses.RenderFile{}
	err := g.client.Call("Plugin.RenderFile", request, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type LibheifRPCServer struct {
	Impl Libheif
}

func (s *LibheifRPCServer) Ping(args interface{}, resp *string) error {
	var err error
	*resp, err = s.Impl.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (s *LibheifRPCServer) DecodeImage(request *requests.DecodeImage, resp *responses.DecodeImage) (err error) {
	defer func() {
		if panicError := recover(); panicError != nil {
			err = fmt.Errorf("panic occurred in %s: %v", "DecodeImage", panicError)
		}
	}()

	implResp, err := s.Impl.DecodeImage(request)
	if err != nil {
		return err
	}

	// Overwrite the target address of resp to the target address of implResp.
	*resp = *implResp

	return nil
}

func (s *LibheifRPCServer) DecodeConfig(request *requests.DecodeConfig, resp *responses.DecodeConfig) (err error) {
	defer func() {
		if panicError := recover(); panicError != nil {
			err = fmt.Errorf("panic occurred in %s: %v", "DecodeConfig", panicError)
		}
	}()

	implResp, err := s.Impl.DecodeConfig(request)
	if err != nil {
		return err
	}

	// Overwrite the target address of resp to the target address of implResp.
	*resp = *implResp

	return nil
}

func (s *LibheifRPCServer) RenderFile(request *requests.RenderFile, resp *responses.RenderFile) (err error) {
	defer func() {
		if panicError := recover(); panicError != nil {
			err = fmt.Errorf("panic occurred in %s: %v", "RenderFile", panicError)
		}
	}()

	implResp, err := s.Impl.RenderFile(request)
	if err != nil {
		return err
	}

	// Overwrite the target address of resp to the target address of implResp.
	*resp = *implResp

	return nil
}

type LibheifPlugin struct {
	Impl Libheif
}

func (p *LibheifPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &LibheifRPCServer{Impl: p.Impl}, nil
}

func (LibheifPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &LibheifRPC{client: c}, nil
}
