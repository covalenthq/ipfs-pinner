// contains client to register and interact with a IPFS remote pinning service
// + with the custom endpoints for file/car upload for supported service
// TODO: might want to separate out the custom endpoints
package pinclient

import (
	"context"
	"fmt"
	"net/http"
	"os"

	core "github.com/covalenthq/ipfs-pinner/core"
	ihttp "github.com/covalenthq/ipfs-pinner/http"
	"github.com/covalenthq/ipfs-pinner/openapi"
	"github.com/covalenthq/ipfs-pinner/w3up"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multibase"
	"github.com/pkg/errors"

	logging "github.com/ipfs/go-log/v2"
)

var logger = logging.Logger("ipfs-pinner")

const UserAgent = "ipfs-pinner"

type Client struct {
	client     *openapi.APIClient
	ps         core.PinningService
	cidVersion int
	w3up       *w3up.W3up
}

func NewClient(request ClientCreateRequest, cidVersion int, w3up *w3up.W3up) PinServiceAPI {
	// assuming we are getting a supported pinning service request
	config := openapi.NewConfiguration()
	if request.httpClient == nil {
		request.httpClient = ihttp.NewHttpClient(nil)
	}
	config.UserAgent = UserAgent
	//bearer := fmt.Sprintf("Bearer %s", request.bearerToken)
	//config.AddDefaultHeader("Authorization", bearer)
	config.Servers = openapi.ServerConfigurations{
		openapi.ServerConfiguration{
			URL:         request.pinningServiceBaseUrl,
			Description: "IPFS pinning service API base",
		},
		openapi.ServerConfiguration{
			URL:         request.filePinBaseUrl,
			Description: "pinning service url to upload files to",
		},
	}
	config.HTTPClient = request.httpClient

	return &Client{client: openapi.NewAPIClient(config), ps: request.ps, cidVersion: cidVersion, w3up: w3up}
}

func (c *Client) IsIPFSSupportedFor(ps core.PinningService) bool {
	return ps != core.Other
}

func (c *Client) Add(ctx context.Context, cid cid.Cid, opts ...AddOption) (core.PinStatusGetter, error) {
	settings := new(addSettings)
	for _, o := range opts {
		if err := o(settings); err != nil {
			return nil, err
		}
	}

	adder := c.client.PinsApi.PinsPost(ctx)
	p := openapi.Pin{
		Cid: cid.Encode(getCIDEncoder()),
	}

	if len(settings.origins) > 0 {
		p.SetOrigins(settings.origins)
	}
	if settings.meta != nil {
		p.SetMeta(settings.meta)
	}
	if len(settings.name) > 0 {
		p.SetName(settings.name)
	}

	result, httpresp, err := adder.Pin(p).Execute()
	if err != nil {
		err := httperr(httpresp, err)
		return nil, err
	}

	return core.NewPinStatusGetter(*result), nil
}

func (c *Client) UploadFile(ctx context.Context, file *os.File) (cid.Cid, error) {
	fcid, err := c.uploadFileViaWeb3Storage(ctx, file)
	if err != nil {
		return cid.Undef, fmt.Errorf("%w", err)
	}

	return fcid, nil
}

func (c *Client) ServiceType() core.PinningService {
	return c.ps
}

func getCIDEncoder() multibase.Encoder {
	enc, err := multibase.NewEncoder(multibase.Base32)
	if err != nil {
		panic(err)
	}
	return enc
}

func httperr(resp *http.Response, e error) error {
	oerr, ok := e.(openapi.GenericOpenAPIError)
	if ok {
		ferr, ok := oerr.Model().(openapi.Failure)
		if ok {
			return errors.Wrapf(e, "reason: %q, details: %q", ferr.Error.GetReason(), ferr.Error.GetDetails())
		}
	}

	if resp == nil {
		return errors.Wrapf(e, "empty response from remote pinning service")
	}

	return errors.Wrapf(e, "remote pinning service returned http error %d", resp.StatusCode)
}

func (c *Client) uploadFileViaPinata(ctx context.Context, file *os.File) (core.PinataResponseGetter, error) {
	//ctx = context.WithValue(ctx, openapi.ContextServerIndex, 1) // index = 1 is the file pin url

	poster := c.client.FilepinApi.PinataFileUpload(ctx)
	opt := openapi.NewPinataOptions()
	opt.SetCidVersion(string(rune(c.cidVersion)))

	result, httpresp, err := poster.PinataOptions(*opt).File(file).Execute()
	if err != nil {
		err := httperr(httpresp, err)
		return nil, err
	}

	return core.NewPinataResponseGetter(*result), nil
}

func (c *Client) uploadFileViaWeb3Storage(ctx context.Context, file *os.File) (cid.Cid, error) {
	return c.w3up.UploadCarFile(file)

	// if err != nil {
	// 	return nil, err
	// }

	// return

	// poster := c.client.FilepinApi.Web3StorageCarUpload(ctx)
	// result, httpresp, err := poster.Body(file).Execute()
	// if err != nil {
	// 	err := httperr(httpresp, err)
	// 	return nil, err
	// }

	// return core.NewWeb3StorageResponseGetter(*result), nil
}
