package angie_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/influxdata/telegraf"
)

var (
	// errNotFound signals that the Angie API routes does not exist.
	errNotFound = errors.New("not found")
)

func (n *AngieAPI) gatherMetrics(addr *url.URL, acc telegraf.Accumulator) {
	addError(acc, n.gatherProcessesMetrics(addr, acc))
	addError(acc, n.gatherConnectionsMetrics(addr, acc))
	addError(acc, n.gatherSlabsMetrics(addr, acc))
	// addError(acc, n.gatherSslMetrics(addr, acc))
	addError(acc, n.gatherHTTPServerZonesMetrics(addr, acc))
	addError(acc, n.gatherHTTPUpstreamsMetrics(addr, acc))
	addError(acc, n.gatherHTTPCachesMetrics(addr, acc))
	addError(acc, n.gatherHTTPLocationZonesMetrics(addr, acc))
	addError(acc, n.gatherResolverZonesMetrics(addr, acc))
	addError(acc, n.gatherHTTPLimitReqsMetrics(addr, acc))
	addError(acc, n.gatherHTTPLimitConnsMetrics(addr, acc))
	addError(acc, n.gatherStreamServerZonesMetrics(addr, acc))
	addError(acc, n.gatherStreamUpstreamsMetrics(addr, acc))
	addError(acc, n.gatherStreamLimitConnsMetrics(addr, acc))
}

func addError(acc telegraf.Accumulator, err error) {
	// This plugin has hardcoded API resource paths it checks that may not
	// be in the angie.conf.  Currently, this is to prevent logging of
	// paths that are not configured.
	//
	// The correct solution is to do a GET to /api to get the available paths
	// on the server rather than simply ignore.
	if !errors.Is(err, errNotFound) {
		acc.AddError(err)
	}
}

func (n *AngieAPI) gatherURL(addr *url.URL, path string) ([]byte, error) {
	// Turn off pretty output to safe bandwidth
	address := fmt.Sprintf("%s/%s?pretty=off", addr.String(), path)
	resp, err := n.client.Get(address)

	if err != nil {
		return nil, fmt.Errorf("error making HTTP request to %q: %w", address, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		// format as special error to catch and ignore as some Angie API
		// features are either optional, or only available in some versions
		return nil, errNotFound
	default:
		return nil, fmt.Errorf("%s returned HTTP status %s", address, resp.Status)
	}

	contentType := strings.Split(resp.Header.Get("Content-Type"), ";")[0]
	switch contentType {
	case "application/json":
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	default:
		return nil, fmt.Errorf("%s returned unexpected content type %s", address, contentType)
	}
}

func (n *AngieAPI) gatherProcessesMetrics(addr *url.URL, acc telegraf.Accumulator) error {
	body, err := n.gatherURL(addr, processesPath)
	if err != nil {
		return err
	}

	var processes = &processes{}

	if err := json.Unmarshal(body, processes); err != nil {
		return err
	}

	acc.AddFields(
		"angie_api_processes",
		map[string]interface{}{
			"respawned": processes.Respawned,
		},
		getTags(addr),
	)

	return nil
}

func (n *AngieAPI) gatherConnectionsMetrics(addr *url.URL, acc telegraf.Accumulator) error {
	body, err := n.gatherURL(addr, connectionsPath)
	if err != nil {
		return err
	}

	var connections = &connections{}

	if err := json.Unmarshal(body, connections); err != nil {
		return err
	}

	acc.AddFields(
		"angie_api_connections",
		map[string]interface{}{
			"accepted": connections.Accepted,
			"dropped":  connections.Dropped,
			"active":   connections.Active,
			"idle":     connections.Idle,
		},
		getTags(addr),
	)

	return nil
}

func (n *AngieAPI) gatherSlabsMetrics(addr *url.URL, acc telegraf.Accumulator) error {
	body, err := n.gatherURL(addr, slabsPath)
	if err != nil {
		return err
	}

	var slabs slabs

	if err := json.Unmarshal(body, &slabs); err != nil {
		return err
	}

	tags := getTags(addr)

	for zoneName, slab := range slabs {
		slabTags := make(map[string]string, len(tags)+1)
		for k, v := range tags {
			slabTags[k] = v
		}
		slabTags["zone"] = zoneName

		acc.AddFields(
			"angie_api_slabs_pages",
			map[string]interface{}{
				"used": slab.Pages.Used,
				"free": slab.Pages.Free,
			},
			slabTags,
		)

		for slotID, slot := range slab.Slots {
			slotTags := make(map[string]string, len(slabTags)+1)
			for k, v := range slabTags {
				slotTags[k] = v
			}
			slotTags["slot"] = slotID

			acc.AddFields(
				"angie_api_slabs_slots",
				map[string]interface{}{
					"used":  slot.Used,
					"free":  slot.Free,
					"reqs":  slot.Reqs,
					"fails": slot.Fails,
				},
				slotTags,
			)
		}
	}

	return nil
}

// func (n *AngieAPI) gatherSslMetrics(addr *url.URL, acc telegraf.Accumulator) error {
// 	body, err := n.gatherURL(addr, sslPath)
// 	if err != nil {
// 		return err
// 	}

// 	var ssl = &ssl{}

// 	if err := json.Unmarshal(body, ssl); err != nil {
// 		return err
// 	}

// 	acc.AddFields(
// 		"angie_api_ssl",
// 		map[string]interface{}{
// 			"handshakes":        ssl.Handshakes,
// 			"handshakes_failed": ssl.HandshakesFailed,
// 			"session_reuses":    ssl.SessionReuses,
// 		},
// 		getTags(addr),
// 	)

// 	return nil
// }

func (n *AngieAPI) gatherHTTPServerZonesMetrics(addr *url.URL, acc telegraf.Accumulator) error {
	body, err := n.gatherURL(addr, httpServerZonesPath)
	if err != nil {
		return err
	}

	var httpServerZones httpServerZones

	if err := json.Unmarshal(body, &httpServerZones); err != nil {
		return err
	}

	tags := getTags(addr)
	for zoneName, zone := range httpServerZones {
		zoneTags := make(map[string]string, len(tags)+1)
		for k, v := range tags {
			zoneTags[k] = v
		}
		zoneTags["zone"] = zoneName
		acc.AddFields(
			"angie_api_http_server_zones",
			func() map[string]interface{} {
				result := map[string]interface{}{
					"requests_total":      zone.Requests.Total,
					"requests_processing": zone.Requests.Processing,
					"requests_discarded":  zone.Requests.Discarded,
					"received":            zone.Data.Received,
					"sent":                zone.Data.Sent,
				}

				// Response codes fields (only include those that are present)
				if zone.Responses.Response100 != nil {
					result["responses_100"] = zone.Responses.Response100
				}
				if zone.Responses.Response101 != nil {
					result["responses_101"] = zone.Responses.Response101
				}
				if zone.Responses.Response102 != nil {
					result["responses_102"] = zone.Responses.Response102
				}
				if zone.Responses.Response200 != nil {
					result["responses_200"] = zone.Responses.Response200
				}
				if zone.Responses.Response201 != nil {
					result["responses_201"] = zone.Responses.Response201
				}
				if zone.Responses.Response202 != nil {
					result["responses_202"] = zone.Responses.Response202
				}
				if zone.Responses.Response203 != nil {
					result["responses_203"] = zone.Responses.Response203
				}
				if zone.Responses.Response204 != nil {
					result["responses_204"] = zone.Responses.Response204
				}
				if zone.Responses.Response205 != nil {
					result["responses_205"] = zone.Responses.Response205
				}
				if zone.Responses.Response206 != nil {
					result["responses_206"] = zone.Responses.Response206
				}
				if zone.Responses.Response300 != nil {
					result["responses_300"] = zone.Responses.Response300
				}
				if zone.Responses.Response301 != nil {
					result["responses_301"] = zone.Responses.Response301
				}
				if zone.Responses.Response302 != nil {
					result["responses_302"] = zone.Responses.Response302
				}
				if zone.Responses.Response303 != nil {
					result["responses_303"] = zone.Responses.Response303
				}
				if zone.Responses.Response304 != nil {
					result["responses_304"] = zone.Responses.Response304
				}
				if zone.Responses.Response305 != nil {
					result["responses_305"] = zone.Responses.Response305
				}
				if zone.Responses.Response307 != nil {
					result["responses_307"] = zone.Responses.Response307
				}
				if zone.Responses.Response308 != nil {
					result["responses_308"] = zone.Responses.Response308
				}
				if zone.Responses.Response400 != nil {
					result["responses_400"] = zone.Responses.Response400
				}
				if zone.Responses.Response401 != nil {
					result["responses_401"] = zone.Responses.Response401
				}
				if zone.Responses.Response402 != nil {
					result["responses_402"] = zone.Responses.Response402
				}
				if zone.Responses.Response403 != nil {
					result["responses_403"] = zone.Responses.Response403
				}
				if zone.Responses.Response404 != nil {
					result["responses_404"] = zone.Responses.Response404
				}
				if zone.Responses.Response405 != nil {
					result["responses_405"] = zone.Responses.Response405
				}
				if zone.Responses.Response406 != nil {
					result["responses_406"] = zone.Responses.Response406
				}
				if zone.Responses.Response407 != nil {
					result["responses_407"] = zone.Responses.Response407
				}
				if zone.Responses.Response408 != nil {
					result["responses_408"] = zone.Responses.Response408
				}
				if zone.Responses.Response409 != nil {
					result["responses_409"] = zone.Responses.Response409
				}
				if zone.Responses.Response410 != nil {
					result["responses_410"] = zone.Responses.Response410
				}
				if zone.Responses.Response411 != nil {
					result["responses_411"] = zone.Responses.Response411
				}
				if zone.Responses.Response412 != nil {
					result["responses_412"] = zone.Responses.Response412
				}
				if zone.Responses.Response413 != nil {
					result["responses_413"] = zone.Responses.Response413
				}
				if zone.Responses.Response421 != nil {
					result["responses_421"] = zone.Responses.Response421
				}
				if zone.Responses.Response422 != nil {
					result["responses_422"] = zone.Responses.Response422
				}
				if zone.Responses.Response423 != nil {
					result["responses_423"] = zone.Responses.Response423
				}
				if zone.Responses.Response424 != nil {
					result["responses_424"] = zone.Responses.Response424
				}
				if zone.Responses.Response425 != nil {
					result["responses_425"] = zone.Responses.Response425
				}
				if zone.Responses.Response426 != nil {
					result["responses_426"] = zone.Responses.Response426
				}
				if zone.Responses.Response428 != nil {
					result["responses_428"] = zone.Responses.Response428
				}
				if zone.Responses.Response429 != nil {
					result["responses_429"] = zone.Responses.Response429
				}
				if zone.Responses.Response431 != nil {
					result["responses_431"] = zone.Responses.Response431
				}
				if zone.Responses.Response500 != nil {
					result["responses_500"] = zone.Responses.Response500
				}
				if zone.Responses.Response501 != nil {
					result["responses_501"] = zone.Responses.Response501
				}
				if zone.Responses.Response502 != nil {
					result["responses_502"] = zone.Responses.Response502
				}
				if zone.Responses.Response503 != nil {
					result["responses_503"] = zone.Responses.Response503
				}
				if zone.Responses.Response504 != nil {
					result["responses_504"] = zone.Responses.Response504
				}
				if zone.Responses.Response505 != nil {
					result["responses_505"] = zone.Responses.Response505
				}
				if zone.Responses.Response511 != nil {
					result["responses_511"] = zone.Responses.Response511
				}

				// SSL (if present)
				if zone.Ssl != nil {
					result["ssl_handhaked"] = zone.Ssl.Handshaked
					result["ssl_reuses"] = zone.Ssl.Reuses
					result["ssl_timedout"] = zone.Ssl.TimedOut
					result["ssl_failed"] = zone.Ssl.Failed
				}
				return result
			}(),
			zoneTags,
		)
	}

	return nil
}

func (n *AngieAPI) gatherHTTPLocationZonesMetrics(addr *url.URL, acc telegraf.Accumulator) error {
	body, err := n.gatherURL(addr, httpLocationZonesPath)
	if err != nil {
		return err
	}

	var httpLocationZones httpLocationZones

	if err := json.Unmarshal(body, &httpLocationZones); err != nil {
		return err
	}

	tags := getTags(addr)

	for zoneName, zone := range httpLocationZones {
		zoneTags := make(map[string]string, len(tags)+1)
		for k, v := range tags {
			zoneTags[k] = v
		}
		zoneTags["zone"] = zoneName
		acc.AddFields(
			"angie_api_http_location_zones",
			func() map[string]interface{} {
				result := map[string]interface{}{
					"requests_total":      zone.Requests.Total,
					"requests_processing": zone.Requests.Processing,
					"requests_discarded":  zone.Requests.Discarded,
					"received":            zone.Data.Received,
					"sent":                zone.Data.Sent,
				}

				// Response codes fields (only include those that are present)
				if zone.Responses.Response100 != nil {
					result["responses_100"] = zone.Responses.Response100
				}
				if zone.Responses.Response101 != nil {
					result["responses_101"] = zone.Responses.Response101
				}
				if zone.Responses.Response102 != nil {
					result["responses_102"] = zone.Responses.Response102
				}
				if zone.Responses.Response200 != nil {
					result["responses_200"] = zone.Responses.Response200
				}
				if zone.Responses.Response201 != nil {
					result["responses_201"] = zone.Responses.Response201
				}
				if zone.Responses.Response202 != nil {
					result["responses_202"] = zone.Responses.Response202
				}
				if zone.Responses.Response203 != nil {
					result["responses_203"] = zone.Responses.Response203
				}
				if zone.Responses.Response204 != nil {
					result["responses_204"] = zone.Responses.Response204
				}
				if zone.Responses.Response205 != nil {
					result["responses_205"] = zone.Responses.Response205
				}
				if zone.Responses.Response206 != nil {
					result["responses_206"] = zone.Responses.Response206
				}
				if zone.Responses.Response300 != nil {
					result["responses_300"] = zone.Responses.Response300
				}
				if zone.Responses.Response301 != nil {
					result["responses_301"] = zone.Responses.Response301
				}
				if zone.Responses.Response302 != nil {
					result["responses_302"] = zone.Responses.Response302
				}
				if zone.Responses.Response303 != nil {
					result["responses_303"] = zone.Responses.Response303
				}
				if zone.Responses.Response304 != nil {
					result["responses_304"] = zone.Responses.Response304
				}
				if zone.Responses.Response305 != nil {
					result["responses_305"] = zone.Responses.Response305
				}
				if zone.Responses.Response307 != nil {
					result["responses_307"] = zone.Responses.Response307
				}
				if zone.Responses.Response308 != nil {
					result["responses_308"] = zone.Responses.Response308
				}
				if zone.Responses.Response400 != nil {
					result["responses_400"] = zone.Responses.Response400
				}
				if zone.Responses.Response401 != nil {
					result["responses_401"] = zone.Responses.Response401
				}
				if zone.Responses.Response402 != nil {
					result["responses_402"] = zone.Responses.Response402
				}
				if zone.Responses.Response403 != nil {
					result["responses_403"] = zone.Responses.Response403
				}
				if zone.Responses.Response404 != nil {
					result["responses_404"] = zone.Responses.Response404
				}
				if zone.Responses.Response405 != nil {
					result["responses_405"] = zone.Responses.Response405
				}
				if zone.Responses.Response406 != nil {
					result["responses_406"] = zone.Responses.Response406
				}
				if zone.Responses.Response407 != nil {
					result["responses_407"] = zone.Responses.Response407
				}
				if zone.Responses.Response408 != nil {
					result["responses_408"] = zone.Responses.Response408
				}
				if zone.Responses.Response409 != nil {
					result["responses_409"] = zone.Responses.Response409
				}
				if zone.Responses.Response410 != nil {
					result["responses_410"] = zone.Responses.Response410
				}
				if zone.Responses.Response411 != nil {
					result["responses_411"] = zone.Responses.Response411
				}
				if zone.Responses.Response412 != nil {
					result["responses_412"] = zone.Responses.Response412
				}
				if zone.Responses.Response413 != nil {
					result["responses_413"] = zone.Responses.Response413
				}
				if zone.Responses.Response421 != nil {
					result["responses_421"] = zone.Responses.Response421
				}
				if zone.Responses.Response422 != nil {
					result["responses_422"] = zone.Responses.Response422
				}
				if zone.Responses.Response423 != nil {
					result["responses_423"] = zone.Responses.Response423
				}
				if zone.Responses.Response424 != nil {
					result["responses_424"] = zone.Responses.Response424
				}
				if zone.Responses.Response425 != nil {
					result["responses_425"] = zone.Responses.Response425
				}
				if zone.Responses.Response426 != nil {
					result["responses_426"] = zone.Responses.Response426
				}
				if zone.Responses.Response428 != nil {
					result["responses_428"] = zone.Responses.Response428
				}
				if zone.Responses.Response429 != nil {
					result["responses_429"] = zone.Responses.Response429
				}
				if zone.Responses.Response431 != nil {
					result["responses_431"] = zone.Responses.Response431
				}
				if zone.Responses.Response500 != nil {
					result["responses_500"] = zone.Responses.Response500
				}
				if zone.Responses.Response501 != nil {
					result["responses_501"] = zone.Responses.Response501
				}
				if zone.Responses.Response502 != nil {
					result["responses_502"] = zone.Responses.Response502
				}
				if zone.Responses.Response503 != nil {
					result["responses_503"] = zone.Responses.Response503
				}
				if zone.Responses.Response504 != nil {
					result["responses_504"] = zone.Responses.Response504
				}
				if zone.Responses.Response505 != nil {
					result["responses_505"] = zone.Responses.Response505
				}
				if zone.Responses.Response511 != nil {
					result["responses_511"] = zone.Responses.Response511
				}
				return result
			}(),
			zoneTags,
		)
	}

	return nil
}

func (n *AngieAPI) gatherHTTPUpstreamsMetrics(addr *url.URL, acc telegraf.Accumulator) error {
	body, err := n.gatherURL(addr, httpUpstreamsPath)
	if err != nil {
		return err
	}

	var httpUpstreams httpUpstreams

	if err := json.Unmarshal(body, &httpUpstreams); err != nil {
		return err
	}

	tags := getTags(addr)

	for upstreamName, upstream := range httpUpstreams {
		upstreamTags := make(map[string]string, len(tags)+1)
		for k, v := range tags {
			upstreamTags[k] = v
		}
		upstreamTags["upstream"] = upstreamName
		upstreamFields := map[string]interface{}{
			"keepalive": upstream.Keepalive,
		}
		acc.AddFields(
			"angie_api_http_upstreams",
			upstreamFields,
			upstreamTags,
		)
		for peerName, peer := range upstream.Peers {
			peerFields := map[string]interface{}{
				"backup":             peer.Backup,
				"weight":             peer.Weight,
				"state":              peer.State,
				"selected_current":   peer.Selected.Current,
				"selected_total":     peer.Selected.Total,
				"sent":               peer.Data.Sent,
				"received":           peer.Data.Received,
				"health_fails":       peer.Health.Fails,
				"health_unavailable": peer.Health.Unavailable,
				"health_downtime":    peer.Health.Downtime,
			}
			// Optional selected last data field
			if peer.Selected.Last != nil {
				peerFields["selected_last"] = peer.Selected.Last
			}
			// Optional downstart field (only present if peer became unavailable)
			if peer.Health.Downstart != nil {
				peerFields["health_downstart"] = peer.Health.Downstart
			}

			// Only include codes fields that are present (and only the important ones)
			if peer.Responses.Response100 != nil {
				peerFields["responses_100"] = peer.Responses.Response100
			}
			if peer.Responses.Response101 != nil {
				peerFields["responses_101"] = peer.Responses.Response101
			}
			if peer.Responses.Response102 != nil {
				peerFields["responses_102"] = peer.Responses.Response102
			}
			if peer.Responses.Response200 != nil {
				peerFields["responses_200"] = peer.Responses.Response200
			}
			if peer.Responses.Response201 != nil {
				peerFields["responses_201"] = peer.Responses.Response201
			}
			if peer.Responses.Response202 != nil {
				peerFields["responses_202"] = peer.Responses.Response202
			}
			if peer.Responses.Response203 != nil {
				peerFields["responses_203"] = peer.Responses.Response203
			}
			if peer.Responses.Response204 != nil {
				peerFields["responses_204"] = peer.Responses.Response204
			}
			if peer.Responses.Response205 != nil {
				peerFields["responses_205"] = peer.Responses.Response205
			}
			if peer.Responses.Response206 != nil {
				peerFields["responses_206"] = peer.Responses.Response206
			}
			if peer.Responses.Response300 != nil {
				peerFields["responses_300"] = peer.Responses.Response300
			}
			if peer.Responses.Response301 != nil {
				peerFields["responses_301"] = peer.Responses.Response301
			}
			if peer.Responses.Response302 != nil {
				peerFields["responses_302"] = peer.Responses.Response302
			}
			if peer.Responses.Response303 != nil {
				peerFields["responses_303"] = peer.Responses.Response303
			}
			if peer.Responses.Response304 != nil {
				peerFields["responses_304"] = peer.Responses.Response304
			}
			if peer.Responses.Response305 != nil {
				peerFields["responses_305"] = peer.Responses.Response305
			}
			if peer.Responses.Response307 != nil {
				peerFields["responses_307"] = peer.Responses.Response307
			}
			if peer.Responses.Response308 != nil {
				peerFields["responses_308"] = peer.Responses.Response308
			}
			if peer.Responses.Response400 != nil {
				peerFields["responses_400"] = peer.Responses.Response400
			}
			if peer.Responses.Response401 != nil {
				peerFields["responses_401"] = peer.Responses.Response401
			}
			if peer.Responses.Response402 != nil {
				peerFields["responses_402"] = peer.Responses.Response402
			}
			if peer.Responses.Response403 != nil {
				peerFields["responses_403"] = peer.Responses.Response403
			}
			if peer.Responses.Response404 != nil {
				peerFields["responses_404"] = peer.Responses.Response404
			}
			if peer.Responses.Response405 != nil {
				peerFields["responses_405"] = peer.Responses.Response405
			}
			if peer.Responses.Response406 != nil {
				peerFields["responses_406"] = peer.Responses.Response406
			}
			if peer.Responses.Response407 != nil {
				peerFields["responses_407"] = peer.Responses.Response407
			}
			if peer.Responses.Response408 != nil {
				peerFields["responses_408"] = peer.Responses.Response408
			}
			if peer.Responses.Response409 != nil {
				peerFields["responses_409"] = peer.Responses.Response409
			}
			if peer.Responses.Response410 != nil {
				peerFields["responses_410"] = peer.Responses.Response410
			}
			if peer.Responses.Response411 != nil {
				peerFields["responses_411"] = peer.Responses.Response411
			}
			if peer.Responses.Response412 != nil {
				peerFields["responses_412"] = peer.Responses.Response412
			}
			if peer.Responses.Response413 != nil {
				peerFields["responses_413"] = peer.Responses.Response413
			}
			if peer.Responses.Response421 != nil {
				peerFields["responses_421"] = peer.Responses.Response421
			}
			if peer.Responses.Response422 != nil {
				peerFields["responses_422"] = peer.Responses.Response422
			}
			if peer.Responses.Response423 != nil {
				peerFields["responses_423"] = peer.Responses.Response423
			}
			if peer.Responses.Response424 != nil {
				peerFields["responses_424"] = peer.Responses.Response424
			}
			if peer.Responses.Response425 != nil {
				peerFields["responses_425"] = peer.Responses.Response425
			}
			if peer.Responses.Response426 != nil {
				peerFields["responses_426"] = peer.Responses.Response426
			}
			if peer.Responses.Response428 != nil {
				peerFields["responses_428"] = peer.Responses.Response428
			}
			if peer.Responses.Response429 != nil {
				peerFields["responses_429"] = peer.Responses.Response429
			}
			if peer.Responses.Response431 != nil {
				peerFields["responses_431"] = peer.Responses.Response431
			}
			if peer.Responses.Response500 != nil {
				peerFields["responses_500"] = peer.Responses.Response500
			}
			if peer.Responses.Response501 != nil {
				peerFields["responses_501"] = peer.Responses.Response501
			}
			if peer.Responses.Response502 != nil {
				peerFields["responses_502"] = peer.Responses.Response502
			}
			if peer.Responses.Response503 != nil {
				peerFields["responses_503"] = peer.Responses.Response503
			}
			if peer.Responses.Response504 != nil {
				peerFields["responses_504"] = peer.Responses.Response504
			}
			if peer.Responses.Response505 != nil {
				peerFields["responses_505"] = peer.Responses.Response505
			}
			if peer.Responses.Response511 != nil {
				peerFields["responses_511"] = peer.Responses.Response511
			}

			// Other optional fields
			if peer.Service != nil {
				peerFields["service"] = *peer.Service
			}
			if peer.MaxConns != nil {
				peerFields["max_conns"] = *peer.MaxConns
			}
			peerTags := make(map[string]string, len(upstreamTags)+2)
			for k, v := range upstreamTags {
				peerTags[k] = v
			}
			peerTags["peer"] = peerName
			peerTags["sid"] = peer.SID
			acc.AddFields("angie_api_http_upstream_peers", peerFields, peerTags)
		}
	}
	return nil
}

func (n *AngieAPI) gatherHTTPCachesMetrics(addr *url.URL, acc telegraf.Accumulator) error {
	body, err := n.gatherURL(addr, httpCachesPath)
	if err != nil {
		return err
	}

	var httpCaches httpCaches

	if err := json.Unmarshal(body, &httpCaches); err != nil {
		return err
	}

	tags := getTags(addr)

	for cacheName, cache := range httpCaches {
		cacheTags := make(map[string]string, len(tags)+1)
		for k, v := range tags {
			cacheTags[k] = v
		}
		cacheTags["cache"] = cacheName
		acc.AddFields(
			"angie_api_http_caches",
			map[string]interface{}{
				"size":                      cache.Size,
				"max_size":                  cache.MaxSize,
				"cold":                      cache.Cold,
				"hit_responses":             cache.Hit.Responses,
				"hit_bytes":                 cache.Hit.Bytes,
				"stale_responses":           cache.Stale.Responses,
				"stale_bytes":               cache.Stale.Bytes,
				"updating_responses":        cache.Updating.Responses,
				"updating_bytes":            cache.Updating.Bytes,
				"revalidated_responses":     cache.Revalidated.Responses,
				"revalidated_bytes":         cache.Revalidated.Bytes,
				"miss_responses":            cache.Miss.Responses,
				"miss_bytes":                cache.Miss.Bytes,
				"miss_responses_written":    cache.Miss.ResponsesWritten,
				"miss_bytes_written":        cache.Miss.BytesWritten,
				"expired_responses":         cache.Expired.Responses,
				"expired_bytes":             cache.Expired.Bytes,
				"expired_responses_written": cache.Expired.ResponsesWritten,
				"expired_bytes_written":     cache.Expired.BytesWritten,
				"bypass_responses":          cache.Bypass.Responses,
				"bypass_bytes":              cache.Bypass.Bytes,
				"bypass_responses_written":  cache.Bypass.ResponsesWritten,
				"bypass_bytes_written":      cache.Bypass.BytesWritten,
			},
			cacheTags,
		)
	}

	return nil
}

func (n *AngieAPI) gatherResolverZonesMetrics(addr *url.URL, acc telegraf.Accumulator) error {
	body, err := n.gatherURL(addr, resolverZonesPath)
	if err != nil {
		return err
	}

	var resolverZones resolverZones

	if err := json.Unmarshal(body, &resolverZones); err != nil {
		return err
	}

	tags := getTags(addr)

	for zoneName, resolver := range resolverZones {
		zoneTags := make(map[string]string, len(tags)+1)
		for k, v := range tags {
			zoneTags[k] = v
		}
		zoneTags["zone"] = zoneName
		acc.AddFields(
			"angie_api_resolver_zones",
			map[string]interface{}{
				"queries_name":   resolver.Queries.Name,
				"queries_srv":    resolver.Queries.Srv,
				"queries_addr":   resolver.Queries.Addr,
				"sent_a":         resolver.Sent.A,
				"sent_aaaa":      resolver.Sent.AAAA,
				"sent_srv":       resolver.Sent.Srv,
				"sent_ptr":       resolver.Sent.Ptr,
				"success":        resolver.Responses.Success,
				"timedout":       resolver.Responses.TimedOut,
				"format_error":   resolver.Responses.FormatError,
				"server_failure": resolver.Responses.ServerFailure,
				"not_found":      resolver.Responses.NotFound,
				"unimplemented":  resolver.Responses.Unimplemented,
				"refused":        resolver.Responses.Refused,
				"other":          resolver.Responses.Other,
			},
			zoneTags,
		)
	}

	return nil
}

func (n *AngieAPI) gatherHTTPLimitReqsMetrics(addr *url.URL, acc telegraf.Accumulator) error {
	body, err := n.gatherURL(addr, httpLimitReqsPath)
	if err != nil {
		return err
	}

	var httpLimitReqs httpLimitReqs

	if err := json.Unmarshal(body, &httpLimitReqs); err != nil {
		return err
	}

	tags := getTags(addr)

	for limitReqName, limit := range httpLimitReqs {
		limitReqsTags := make(map[string]string, len(tags)+1)
		for k, v := range tags {
			limitReqsTags[k] = v
		}
		limitReqsTags["limit"] = limitReqName
		acc.AddFields(
			"angie_api_http_limit_reqs",
			map[string]interface{}{
				"passed":    limit.Passed,
				"skipped":   limit.Skipped,
				"delayed":   limit.Delayed,
				"rejected":  limit.Rejected,
				"exhausted": limit.Exhausted,
			},
			limitReqsTags,
		)
	}

	return nil
}

func (n *AngieAPI) gatherHTTPLimitConnsMetrics(addr *url.URL, acc telegraf.Accumulator) error {
	body, err := n.gatherURL(addr, httpLimitConnsPath)
	if err != nil {
		return err
	}

	var httpLimitConns limitConns

	if err := json.Unmarshal(body, &httpLimitConns); err != nil {
		return err
	}

	tags := getTags(addr)

	for limitConnName, limit := range httpLimitConns {
		limitConnsTags := make(map[string]string, len(tags)+1)
		for k, v := range tags {
			limitConnsTags[k] = v
		}
		limitConnsTags["limit"] = limitConnName
		acc.AddFields(
			"angie_api_http_limit_conns",
			map[string]interface{}{
				"passed":    limit.Passed,
				"skipped":   limit.Skipped,
				"rejected":  limit.Rejected,
				"exhausted": limit.Exhausted,
			},
			limitConnsTags,
		)
	}

	return nil
}

func (n *AngieAPI) gatherStreamServerZonesMetrics(addr *url.URL, acc telegraf.Accumulator) error {
	body, err := n.gatherURL(addr, streamServerZonesPath)
	if err != nil {
		return err
	}

	var streamServerZones streamServerZones

	if err := json.Unmarshal(body, &streamServerZones); err != nil {
		return err
	}

	tags := getTags(addr)

	for zoneName, zone := range streamServerZones {
		zoneTags := make(map[string]string, len(tags)+1)
		for k, v := range tags {
			zoneTags[k] = v
		}
		zoneTags["zone"] = zoneName

		acc.AddFields(
			"angie_api_stream_server_zones",
			func() map[string]interface{} {
				result := map[string]interface{}{
					"connections_total":            zone.Connections.Total,
					"connections_processing":       zone.Connections.Processing,
					"connections_discarded":        zone.Connections.Discarded,
					"sessions_success":             zone.Sessions.Success,
					"sessions_invalid":             zone.Sessions.Invalid,
					"sessions_forbidden":           zone.Sessions.Forbidden,
					"sessions_internal_error":      zone.Sessions.InternalError,
					"sessions_bad_gateway":         zone.Sessions.BadGateway,
					"sessions_service_unavailable": zone.Sessions.ServiceUnavailable,
					"received":                     zone.Data.Received,
					"sent":                         zone.Data.Sent,
				}
				// SSL (if present)
				if zone.Ssl != nil {
					result["ssl_handhaked"] = zone.Ssl.Handshaked
					result["ssl_reuses"] = zone.Ssl.Reuses
					result["ssl_timedout"] = zone.Ssl.TimedOut
					result["ssl_failed"] = zone.Ssl.Failed
				}
				return result
			}(),
			zoneTags,
		)
	}

	return nil
}

func (n *AngieAPI) gatherStreamUpstreamsMetrics(addr *url.URL, acc telegraf.Accumulator) error {
	body, err := n.gatherURL(addr, streamUpstreamsPath)
	if err != nil {
		return err
	}

	var streamUpstreams streamUpstreams

	if err := json.Unmarshal(body, &streamUpstreams); err != nil {
		return err
	}

	tags := getTags(addr)

	for upstreamName, upstream := range streamUpstreams {
		upstreamTags := make(map[string]string, len(tags)+1)
		for k, v := range tags {
			upstreamTags[k] = v
		}
		upstreamTags["upstream"] = upstreamName
		for peerName, peer := range upstream.Peers {
			peerFields := map[string]interface{}{
				"backup":             peer.Backup,
				"weight":             peer.Weight,
				"state":              peer.State,
				"selected_current":   peer.Selected.Current,
				"selected_total":     peer.Selected.Total,
				"sent":               peer.Data.Sent,
				"received":           peer.Data.Received,
				"health_fails":       peer.Health.Fails,
				"health_unavailable": peer.Health.Unavailable,
				"health_downtime":    peer.Health.Downtime,
			}
			// Optional fields
			if peer.Service != nil {
				peerFields["service"] = *peer.Service
			}
			if peer.MaxConns != nil {
				peerFields["max_conns"] = *peer.MaxConns
			}
			if peer.Health.Downstart != nil {
				peerFields["health_downstart"] = *peer.Health.Downstart
			}

			peerTags := make(map[string]string, len(upstreamTags)+1)
			for k, v := range upstreamTags {
				peerTags[k] = v
			}
			peerTags["peer"] = peerName
			acc.AddFields("angie_api_stream_upstream_peers", peerFields, peerTags)
		}
	}

	return nil
}

func (n *AngieAPI) gatherStreamLimitConnsMetrics(addr *url.URL, acc telegraf.Accumulator) error {
	body, err := n.gatherURL(addr, streamLimitConnsPath)
	if err != nil {
		return err
	}

	var streamLimitConns limitConns

	if err := json.Unmarshal(body, &streamLimitConns); err != nil {
		return err
	}

	tags := getTags(addr)

	for limitConnName, limit := range streamLimitConns {
		limitConnsTags := make(map[string]string, len(tags)+1)
		for k, v := range tags {
			limitConnsTags[k] = v
		}
		limitConnsTags["limit"] = limitConnName
		acc.AddFields(
			"angie_api_stream_limit_conns",
			map[string]interface{}{
				"passed":    limit.Passed,
				"skipped":   limit.Skipped,
				"rejected":  limit.Rejected,
				"exhausted": limit.Exhausted,
			},
			limitConnsTags,
		)
	}

	return nil
}

func getTags(addr *url.URL) map[string]string {
	h := addr.Host
	host, port, err := net.SplitHostPort(h)
	if err != nil {
		host = addr.Host
		switch addr.Scheme {
		case "http":
			port = "80"
		case "https":
			port = "443"
		default:
			port = ""
		}
	}
	return map[string]string{"source": host, "port": port}
}
