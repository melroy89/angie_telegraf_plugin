# Angie API Input Plugin

This plugin gathers metrics from the **free** and open-source
[Angie web server][angie] via the [REST API][api]. And see also my matching [Grafana dashboard](https://grafana.com/grafana/dashboards/24461-angie-metrics/).


⭐ Telegraf v1.9.0
🏷️ server, web
💻 all

[angie]: https://en.angie.software/
[api]: https://en.angie.software/angie/docs/configuration/modules/http/http_api/#a-api

## Building from source

Requirements:

- **Golang**: tested on version `go1.25.0 linux/amd64`

This plugin can be used an an external plugin to Telegraf. 

1. First clone and build the project:
```sh
$ git clone https://github.com/melroy89/angie_telegraf_plugin.git
$ cd angie_telegraf_plugin
$ make build
```
2. The standalone binary should be available in the root directory: `./angie_telegraf`
3. Copy both the binary at some location together with your `plugin.conf` file. If needed make changes to your configuration.
3. Add the plugin to your `telegraf.conf` file: 

```ini
[[inputs.execd]]
  command = ["/path/to/angie_telegraf", "-config", "/path/to/plugin.conf"]
  signal = "none"
```

## Global configuration options

In addition to the plugin-specific configuration settings, plugins support
additional global and plugin configuration settings. These settings are used to
modify metrics, tags, and field or create aliases and configure ordering, etc.
See the [CONFIGURATION.md](https://github.com/influxdata/telegraf/blob/master/docs/CONFIGURATION.md#plugins) for more details.

## Configuration

```toml @plugin.conf
# Read Angie API status information
[[inputs.angie_api]]
  ## An array of Angie API URIs to gather stats.
  urls = ["http://localhost/status"]
  # Angie API version, default: 1
  # api_version = 1

  # HTTP response timeout (default: 5s)
  response_timeout = "5s"

  ## Optional TLS Config
  # tls_ca = "/etc/telegraf/ca.pem"
  # tls_cert = "/etc/telegraf/cert.pem"
  # tls_key = "/etc/telegraf/key.pem"
  ## Use TLS but skip chain & host verification
  # insecure_skip_verify = false
```

## Grafana Dashboard

This Telegraf plugin _could_ be used together with [my matching Grafana Angie Metrics Dashboard](https://grafana.com/grafana/dashboards/24461-angie-metrics/).

## Developers

For developers (thank you for contributing), first make a copy of the `plugin.conf` and call it `dev.conf`:

```sh
cp plugin.conf dev.conf
```

Make any changes to the configuration if you wish.

Finally, you can use following command to build + run (it will use the `dev.conf` file in the root of directory of this project):

```sh
make rundev
```

(Press enter to trigger a fetch).

## Measurements by API version

| Measurement                     | API version (api_version) |
|---------------------------------|---------------------------|
| angie_api_connections           | >= 1                      |
| angie_api_slabs_pages           | >= 1                      |
| angie_api_slabs_slots           | >= 1                      |
| angie_api_http_requests         | >= 1                      |
| angie_api_http_server_zones     | >= 1                      |
| angie_api_http_upstreams        | >= 1                      |
| angie_api_http_upstream_peers   | >= 1                      |
| angie_api_http_caches           | >= 1                      |
| angie_api_http_location_zones   | >= 1                      |
| angie_api_resolver_zones        | >= 1                      |
| angie_api_http_limit_reqs       | >= 1                      |
| angie_api_http_limit_conns      | >= 1                      |
| angie_api_stream_server_zones   | >= 1                      |
| angie_api_stream_upstream_peers | >= 1                      |
| angie_api_stream_limit_conns    | >= 1                      |

## Metrics

- angie_api_connections
  - accepted
  - dropped
  - active
  - idle
- angie_api_slabs_pages
  - used
  - free
- angie_api_slabs_slots
  - used
  - free
  - reqs
  - fails
- angie_api_http_server_zones
  - requests_total
  - requests_processing
  - requests_discarded
  - received
  - sent
  - responses_xxx
     - Where `xxx` is the status code (100-599) 
  - ssl_handhaked (in case of SSL)
  - ssl_reuses (in case of SSL)
  - ssl_timedout (in case of SSL)
  - ssl_failed (in case of SSL)
- angie_api_http_upstreams
  - keepalive
- angie_api_http_upstream_peers
  - backup
  - weight
  - state
  - selected_current
  - selected_total
  - selected_last (if present)
  - sent
  - received
  - health_fails
  - health_unavailable
  - health_downtime
  - health_downstart (if present)
  - responses_xxx
     - Where `xxx` is the status code (100-599) 
  - service (if configured)
  - max_conns (if present)
- angie_api_http_caches
  - size
  - max_size
  - cold
  - hit_responses
  - hit_bytes
  - stale_responses
  - stale_bytes
  - updating_responses
  - updating_bytes
  - revalidated_responses
  - revalidated_bytes
  - miss_responses
  - miss_bytes
  - miss_responses_written
  - miss_bytes_written
  - expired_responses
  - expired_bytes
  - expired_responses_written
  - expired_bytes_written
  - bypass_responses
  - bypass_bytes
  - bypass_responses_written
  - bypass_bytes_written
- angie_api_http_location_zones
  - requests_total
  - requests_processing
  - requests_discarded
  - received
  - sent
  - responses_xxx
     - Where `xxx` is the status code (100-599) 
- angie_api_resolver_zones
	- queries_name
	- queries_srv
	- queries_addr
	- sent_a
	- sent_aaaa
	- sent_srv
	- sent_ptr
	- success
	- timedout
	- format_error
	- server_failure
	- not_found
	- unimplemented
	- refused
	- other
- angie_api_http_limit_reqs
  - passed
  - skipped
  - delayed
  - rejected
  - exhausted
- angie_api_http_limit_conns
  - passed
  - skipped
  - rejected
  - exhausted
- angie_api_stream_server_zones
  - connections_total
  - connections_processing
  - connections_discarded
  - sessions_success
  - sessions_invalid
  - sessions_forbidden
  - sessions_internal_error
  - sessions_bad_gateway
  - sessions_service_unavailable
  - received
  - sent
  - ssl_handhaked (in case of SSL)
  - ssl_reuses (in case of SSL)
  - ssl_timedout (in case of SSL)
  - ssl_failed (in case of SSL)
- angie_api_stream_upstream_peers
  - backup
  - weight
  - state
  - selected_current
  - selected_total
	- sent
	- received      
  - health_fails
  - health_unavailable
  - health_downtime
  - health_downstart (if present)
  - service (if configured)
  - max_conns (if present)
- angie_api_stream_limit_conns
  - passed
  - skipped
  - rejected
  - exhausted

### Tags

- angie_api_connections, angie_api_http_requests
  - source
  - port

- angie_api_http_upstreams
  - upstream
  - source
  - port

- angie_api_http_server_zones, angie_api_http_location_zones, 
  angie_api_resolver_zones, angie_api_slabs_pages, angie_api_stream_server_zones
  - source
  - port
  - zone

- angie_api_slabs_slots
  - source
  - port
  - zone
  - slot

- angie_api_http_caches
  - source
  - port

- angie_api_http_limit_reqs, angie_api_http_limit_conns, angie_api_stream_limit_conns
  - source
  - port
  - limit

- angie_api_http_upstream_peers
  - peer
  - sid

- angie_api_stream_upstream_peers
  - peer

## Example Output

Using this configuration:

```toml
[[inputs.angie_api]]
  ## An array of Angie API URIs to gather stats.
  urls = ["http://localhost/api"]
```

When run with:

```sh
telegraf --config telegraf.conf --input-filter angie_api --test
```

*Note:* Use `--sample-config` if you wish to use the included example configuration instead.

It produces (example output):

```text
angie_api_connections,port=80,source=angie.host.tld accepted=3108i,dropped=0i,active=103i,idle=75i 1764031139994287336
angie_api_slabs_pages,port=80,source=angie.host.tld,zone=addr used=2i,free=2542i 1764031139998323977
angie_api_slabs_slots,port=80,slot=32,source=angie.host.tld,zone=addr fails=0i,used=1i,free=126i,reqs=1i 1764031139998334023
angie_api_slabs_slots,port=80,slot=64,source=angie.host.tld,zone=addr reqs=2i,fails=0i,used=0i,free=0i 1764031139998337548
angie_api_slabs_slots,port=80,slot=128,source=angie.host.tld,zone=addr used=1i,free=31i,reqs=1i,fails=0i 1764031139998392180
angie_api_slabs_pages,port=80,source=angie.host.tld,zone=SSL used=152i,free=2392i 1764031139998394774
angie_api_slabs_slots,port=80,slot=128,source=angie.host.tld,zone=SSL used=1592i,free=8i,reqs=1592i,fails=0i 1764031139998414794
angie_api_slabs_slots,port=80,slot=256,source=angie.host.tld,zone=SSL fails=0i,used=1592i,free=8i,reqs=1592i 1764031139998445841
angie_api_slabs_slots,port=80,slot=512,source=angie.host.tld,zone=SSL used=1i,free=7i,reqs=1i,fails=0i 1764031139998450908
angie_api_slabs_slots,port=80,slot=64,source=angie.host.tld,zone=SSL used=1i,free=63i,reqs=1i,fails=0i 1764031139998453272
angie_api_slabs_pages,port=80,source=angie.host.tld,zone=upstream used=6i,free=57i 1764031139998483648
angie_api_slabs_slots,port=80,slot=8,source=angie.host.tld,zone=upstream used=2i,free=502i,reqs=2i,fails=0i 1764031139998486632
angie_api_slabs_slots,port=80,slot=16,source=angie.host.tld,zone=upstream used=2i,free=252i,reqs=2i,fails=0i 1764031139998489236
angie_api_slabs_slots,port=80,slot=32,source=angie.host.tld,zone=upstream reqs=2i,fails=0i,used=2i,free=125i 1764031139998511249
angie_api_slabs_slots,port=80,slot=128,source=angie.host.tld,zone=upstream reqs=4i,fails=0i,used=4i,free=28i 1764031139998516417
angie_api_slabs_slots,port=80,slot=512,source=angie.host.tld,zone=upstream used=1i,free=7i,reqs=1i,fails=0i 1764031139998518640
angie_api_slabs_pages,port=80,source=angie.host.tld,zone=ip used=2i,free=2542i 1764031139998536307
angie_api_slabs_slots,port=80,slot=32,source=angie.host.tld,zone=ip fails=0i,used=1i,free=126i,reqs=1i 1764031139998540213
angie_api_slabs_slots,port=80,slot=128,source=angie.host.tld,zone=ip used=11i,free=21i,reqs=89i,fails=0i 1764031139998566312
angie_api_slabs_pages,port=80,source=angie.host.tld,zone=matrix_limit used=2i,free=2542i 1764031139998569978
angie_api_slabs_slots,port=80,slot=64,source=angie.host.tld,zone=matrix_limit free=63i,reqs=1i,fails=0i,used=1i 1764031139998572662
angie_api_slabs_slots,port=80,slot=128,source=angie.host.tld,zone=matrix_limit free=21i,reqs=75i,fails=0i,used=11i 1764031139998599322
angie_api_slabs_pages,port=80,source=angie.host.tld,zone=some_limit used=2i,free=2542i 1764031139998602256
angie_api_slabs_slots,port=80,slot=32,source=angie.host.tld,zone=some_limit used=1i,free=126i,reqs=1i,fails=0i 1764031139998604760
angie_api_slabs_slots,port=80,slot=128,source=angie.host.tld,zone=some_limit used=13i,free=19i,reqs=104i,fails=0i 1764031139998626202
angie_api_slabs_pages,port=80,source=angie.host.tld,zone=another_zone used=6i,free=248i 1764031139998628175
angie_api_slabs_slots,port=80,slot=32,source=angie.host.tld,zone=another_zone used=2i,free=125i,reqs=2i,fails=0i 1764031139998629958
angie_api_slabs_slots,port=80,slot=128,source=angie.host.tld,zone=another_zone used=4i,free=28i,reqs=4i,fails=0i 1764031139998672662
angie_api_slabs_slots,port=80,slot=512,source=angie.host.tld,zone=another_zone used=1i,free=7i,reqs=1i,fails=0i 1764031139998676478
angie_api_slabs_slots,port=80,slot=8,source=angie.host.tld,zone=another_zone used=2i,free=502i,reqs=2i,fails=0i 1764031139998679412
angie_api_slabs_slots,port=80,slot=16,source=angie.host.tld,zone=another_zone free=252i,reqs=2i,fails=0i,used=2i 1764031139998708085
angie_api_slabs_pages,port=80,source=angie.host.tld,zone=CACHE used=38i,free=2506i 1764031139998712151
angie_api_slabs_slots,port=80,slot=32,source=angie.host.tld,zone=CACHE fails=0i,used=1i,free=126i,reqs=1i 1764031139998715426
angie_api_slabs_slots,port=80,slot=128,source=angie.host.tld,zone=CACHE used=1143i,free=9i,reqs=1145i,fails=0i 1764031139998738311
angie_api_slabs_slots,port=80,slot=512,source=angie.host.tld,zone=CACHE used=1i,free=7i,reqs=1i,fails=0i 1764031139998740254
angie_api_slabs_pages,port=80,source=angie.host.tld,zone=limit_req_zone used=2i,free=2542i 1764031139998741946
angie_api_slabs_slots,port=80,slot=64,source=angie.host.tld,zone=limit_req_zone used=2i,free=62i,reqs=2i,fails=0i 1764031139998762217
angie_api_slabs_slots,port=80,slot=128,source=angie.host.tld,zone=limit_req_zone used=1i,free=31i,reqs=1i,fails=0i 1764031139998764951
angie_api_slabs_pages,port=80,source=angie.host.tld,zone=dns_upstream used=6i,free=248i 1764031139998766804
angie_api_slabs_slots,port=80,slot=8,source=angie.host.tld,zone=dns_upstream free=502i,reqs=2i,fails=0i,used=2i 1764031139998789347
angie_api_slabs_slots,port=80,slot=16,source=angie.host.tld,zone=dns_upstream reqs=3i,fails=0i,used=3i,free=251i 1764031139998804811
angie_api_slabs_slots,port=80,slot=32,source=angie.host.tld,zone=dns_upstream used=2i,free=125i,reqs=2i,fails=0i 1764031139998807305
angie_api_slabs_slots,port=80,slot=64,source=angie.host.tld,zone=dns_upstream used=1i,free=63i,reqs=1i,fails=0i 1764031139998826543
angie_api_slabs_slots,port=80,slot=128,source=angie.host.tld,zone=dns_upstream used=6i,free=26i,reqs=6i,fails=0i 1764031139998829177
angie_api_slabs_slots,port=80,slot=512,source=angie.host.tld,zone=dns_upstream used=2i,free=6i,reqs=2i,fails=0i 1764031139998832633
angie_api_http_server_zones,port=80,source=angie.host.tld,zone=server_zone requests_processing=0i,sent=0i,ssl_reuses=0i,ssl_timedout=0i,requests_discarded=0i,received=0i,ssl_handhaked=0i,ssl_failed=0i,requests_total=0i 1763846935470146020
angie_api_http_server_zones,port=80,source=angie.host.tld,zone=example.zone.tld requests_total=849i,requests_processing=17i,requests_discarded=0i,sent=14214953i,responses_101=54i,responses_200=639i,responses_304=139i,ssl_handhaked=664i,received=267614i,ssl_reuses=424i,ssl_timedout=0i,ssl_failed=0i 1763846935470152089
angie_api_http_upstreams,port=80,source=angie.host.tld,upstream=some_upstream_name keepalive=0i 1763846935470758731
angie_api_http_upstream_peers,peer=127.0.0.1:3005,port=80,sid=0349acf60535cd8bdf89fb53de0f959e,source=angie.host.tld,upstream=some_upstream_name state="up",sent=0i,weight=1i,selected_current=0i,selected_total=0i,reveived=0i,health_fails=0i,health_unavailable=0i,health_downtime=0i,backup=false 1763846935470766403
angie_api_http_upstreams,port=80,source=angie.host.tld,upstream=another_upstream keepalive=0i 1763846935470769057
angie_api_http_upstream_peers,peer=127.0.0.1:8999,port=80,sid=adbdc4c737eef0c63976e2f697c8c8b3,source=angie.host.tld,upstream=another_upstream backup=false,weight=1i,selected_total=674i,sent=398003i,reveived=14527538i,health_fails=0i,selected_last="2025-11-22T21:28:00Z",responses_200=557i,state="up",selected_current=17i,health_unavailable=0i,health_downtime=0i,responses_101=54i,responses_304=46i 1763846935470800063
angie_api_http_caches,cache=CACHE,port=80,source=angie.host.tld miss_responses=67i,miss_bytes_written=834493i,expired_bytes_written=1351616i,bypass_bytes=0i,bypass_bytes_written=0i,cold=false,stale_responses=0i,stale_bytes=0i,bypass_responses=0i,bypass_responses_written=0i,updating_responses=0i,updating_bytes=0i,miss_responses_written=36i,expired_responses=89i,max_size=1073741824i,revalidated_responses=0i,revalidated_bytes=0i,miss_bytes=843208i,expired_bytes=1351616i,expired_responses_written=89i,size=7794688i,hit_responses=64i,hit_bytes=421058i 1763846935471333465
angie_api_resolver_zones,port=80,source=angie.host.tld,zone=resolver_zone server_failure=0i,not_found=0i,unimplemented=0i,refused=0i,sent_aaaa=0i,queries_srv=0i,sent_a=0i,success=0i,format_error=0i,queries_name=0i,sent_srv=0i,sent_ptr=0i,other=0i,queries_addr=0i,timedout=0i 1763929363672772193
angie_api_http_limit_reqs,limit=ip,port=80,source=angie.host.tld passed=102772i,skipped=0i,delayed=208i,rejected=4i,exhausted=0i 1763922956010183369
angie_api_http_limit_reqs,limit=some_limit,port=80,source=angie.host.tld passed=16223i,skipped=0i,delayed=0i,rejected=0i,exhausted=0i 1763922956010187055
angie_api_http_limit_reqs,limit=another_limit,port=80,source=angie.host.tld passed=47085i,skipped=0i,delayed=0i,rejected=0i,exhausted=0i 1763922956010191722
angie_api_http_limit_conns,limit=addr,port=80,source=angie.host.tld skipped=0i,rejected=1i,exhausted=0i,passed=355i 1763922956011417885
angie_api_stream_server_zones,port=80,source=angie.host.tld,zone=dns connections_total=140i,sessions_success=78i,sessions_invalid=0i,sessions_forbidden=0i,sessions_bad_gateway=0i,received=4118i,sent=10755i,connections_processing=62i,connections_discarded=0i,sessions_internal_error=0i,sessions_service_unavailable=0i 1764030545409478780
angie_api_stream_upstream_peers,peer=192.168.1.1:53,port=80,source=angie.host.tld,upstream=dns health_downtime=0i,backup=false,state="up",sent=2035i,received=5302i,health_fails=0i,weight=1i,selected_current=31i,selected_total=70i,health_unavailable=0i 1764030545409918130
angie_api_stream_upstream_peers,peer=8.8.8.8:53,port=80,source=angie.host.tld,upstream=dns state="up",selected_current=31i,received=5453i,health_fails=0i,health_downtime=0i,backup=false,weight=1i,selected_total=70i,sent=2083i,health_unavailable=0i 1764030545409923959
angie_api_stream_limit_conns,limit=limit_req_zone,port=80,source=angie.host.tld rejected=0i,exhausted=0i,passed=140i,skipped=0i 1764030545410299754
```

### Reference material

- [Angie API documentation](https://en.angie.software/angie/docs/configuration/modules/http/http_api/)
