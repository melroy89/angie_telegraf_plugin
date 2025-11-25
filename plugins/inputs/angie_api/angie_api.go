package angie_api

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/config"
	common_http "github.com/influxdata/telegraf/plugins/common/http"
	"github.com/influxdata/telegraf/plugins/inputs"
)

const sampleConfig = `
[[inputs.angie_api]]
  urls = ["http://localhost/status"]
  `

const (
	// Default settings
	defaultAPIVersion = 1

	// Paths
	processesPath   = "processes"
	connectionsPath = "connections"
	slabsPath       = "slabs"

	httpServerZonesPath   = "http/server_zones"
	httpLocationZonesPath = "http/location_zones"
	httpUpstreamsPath     = "http/upstreams"
	httpCachesPath        = "http/caches"
	httpLimitReqsPath     = "http/limit_reqs"
	httpLimitConnsPath    = "http/limit_conns"
	resolverZonesPath     = "resolvers"

	streamServerZonesPath = "stream/server_zones"
	streamUpstreamsPath   = "stream/upstreams"
	streamLimitConnsPath  = "stream/limit_conns"
)

type AngieAPI struct {
	Urls       []string        `toml:"urls"`
	APIVersion int64           `toml:"api_version"`
	Log        telegraf.Logger `toml:"-"`
	common_http.HTTPClientConfig

	client *http.Client
}

func (*AngieAPI) SampleConfig() string {
	return sampleConfig
}

func (n *AngieAPI) Gather(acc telegraf.Accumulator) error {
	var wg sync.WaitGroup

	// Only support API version 1 (currently APIVersion is not yet used)
	if n.APIVersion == 0 {
		n.APIVersion = defaultAPIVersion
	}

	// Create an HTTP client that is re-used for each
	// collection interval
	if n.client == nil {
		client, err := n.createHTTPClient()
		if err != nil {
			return err
		}
		n.client = client
	}

	for _, u := range n.Urls {
		addr, err := url.Parse(u)
		if err != nil {
			acc.AddError(fmt.Errorf("unable to parse address %q: %w", u, err))
			continue
		}

		wg.Add(1)
		go func(addr *url.URL) {
			defer wg.Done()
			n.gatherMetrics(addr, acc)
		}(addr)
	}

	wg.Wait()
	return nil
}

func (n *AngieAPI) createHTTPClient() (*http.Client, error) {
	if n.HTTPClientConfig.ResponseHeaderTimeout < config.Duration(time.Second) {
		n.HTTPClientConfig.ResponseHeaderTimeout = config.Duration(time.Second * 5)
	}

	n.Log.Debugf("Creating HTTP client with response timeout of %s", n.HTTPClientConfig.ResponseHeaderTimeout)

	// Create the client
	ctx := context.Background()
	client, err := n.HTTPClientConfig.CreateClient(ctx, n.Log)
	if err != nil {
		return nil, fmt.Errorf("creating client failed: %w", err)
	}

	return client, nil
}

func init() {
	inputs.Add("angie_api", func() telegraf.Input {
		return &AngieAPI{}
	})
}
