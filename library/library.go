package library

import (
	"encoding/gob"
	"errors"
	"image"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/klippa-app/go-libheif/library/requests"
	"github.com/klippa-app/go-libheif/library/shared"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

type Config struct {
	Command Command
}

type Command struct {
	BinPath string
	Args    []string

	// StartTimeout is the timeout to wait for the plugin to say it
	// has started successfully.
	StartTimeout time.Duration
}

var client *plugin.Client
var gRPCClient plugin.ClientProtocol
var libheifplugin shared.Libheif
var currentConfig Config

func init() {
	// Needed to serialize the image interface.
	gob.Register(&image.YCbCr{})
	gob.Register(&image.RGBA64{})
	gob.Register(&image.RGBA{})
	gob.Register(&image.Gray{})
}

func Init(config Config) error {
	if client != nil {
		return nil
	}

	currentConfig = config

	return startPlugin()
}

func DeInit() {
	client.Kill()
	client = nil
	gRPCClient.Close()
	gRPCClient = nil
	libheifplugin = nil
}

func startPlugin() error {
	var handshakeConfig = plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "BASIC_PLUGIN",
		MagicCookieValue: "libheif",
	}

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"libheif": &shared.LibheifPlugin{},
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	client = plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command(currentConfig.Command.BinPath, currentConfig.Command.Args...),
		Logger:          logger,
		StartTimeout:    currentConfig.Command.StartTimeout,
	})

	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}

	gRPCClient = rpcClient

	raw, err := rpcClient.Dispense("libheif")
	if err != nil {
		log.Fatal(err)
	}

	pluginInstance := raw.(shared.Libheif)
	pong, err := pluginInstance.Ping()
	if err != nil {
		return err
	}

	if pong != "Pong" {
		return errors.New("Wrong ping/pong result")
	}

	libheifplugin = pluginInstance

	return nil
}

func checkPlugin() error {
	pong, err := libheifplugin.Ping()
	if err != nil {
		log.Printf("restarting libheif plugin due to wrong pong result: %s", err.Error())
		err = startPlugin()
		if err != nil {
			log.Printf("could not restart libheif plugin: %s", err.Error())
			return err
		}
	}

	if pong != "Pong" {
		log.Printf("restarting libheif plugin due to wrong pong result: %s", pong)
		err = startPlugin()
		if err != nil {
			log.Printf("could not restart libheif plugin: %s", err.Error())
			return err
		}
	}

	return nil
}

var NotInitializedError = errors.New("libheif was not initialized, you must call the Init() method")

type RenderFileOutputFormat string // The file format to render output as.

const (
	RenderFileOutputFormatJPG RenderFileOutputFormat = "jpg" // Render the file as a JPEG file.
	RenderFileOutputFormatPNG RenderFileOutputFormat = "png" // Render the file as a PNG file.
)

type RenderOptions struct {
	OutputFormat RenderFileOutputFormat // The format to output the image as
	MaxFileSize  int64                  // The maximum filesize, if jpg is chosen as output format, it will try to compress it until it fits
}

func RenderFile(data *[]byte, options RenderOptions) (*[]byte, error) {
	if libheifplugin == nil {
		return nil, NotInitializedError
	}

	err := checkPlugin()
	if err != nil {
		return nil, errors.New("could not check or start plugin")
	}

	resp, err := libheifplugin.RenderFile(&requests.RenderFile{Data: data, OutputFormat: requests.RenderFileOutputFormat(options.OutputFormat), MaxFileSize: options.MaxFileSize})
	if err != nil {
		return nil, err
	}
	return resp.Output, nil
}

func DecodeImage(r io.Reader) (image.Image, error) {
	if libheifplugin == nil {
		return nil, NotInitializedError
	}

	err := checkPlugin()
	if err != nil {
		return nil, errors.New("could not check or start plugin")
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	resp, err := libheifplugin.DecodeImage(&requests.DecodeImage{Data: &data})
	if err != nil {
		return nil, err
	}

	return resp.Image, nil
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	var config image.Config

	if libheifplugin == nil {
		return config, NotInitializedError
	}

	err := checkPlugin()
	if err != nil {
		return config, errors.New("could not check or start plugin")
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return config, err
	}

	resp, err := libheifplugin.DecodeConfig(&requests.DecodeConfig{Data: &data})
	if err != nil {
		return config, err
	}

	return resp.Config, nil
}
