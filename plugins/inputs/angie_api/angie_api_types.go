package angie_api

type processes struct {
	Respawned int `json:"respawned"`
}

type connections struct {
	Accepted int64 `json:"accepted"`
	Dropped  int64 `json:"dropped"`
	Active   int64 `json:"active"`
	Idle     int64 `json:"idle"`
}

type slabs map[string]struct {
	Pages struct {
		Used int64 `json:"used"`
		Free int64 `json:"free"`
	} `json:"pages"`
	Slots map[string]struct {
		Used  int64 `json:"used"`
		Free  int64 `json:"free"`
		Reqs  int64 `json:"reqs"`
		Fails int64 `json:"fails"`
	} `json:"slots"`
}

type ssl struct {
	Handshaked int64 `json:"handshaked"`
	Reuses     int64 `json:"reuses"`
	TimedOut   int64 `json:"timedout"`
	Failed     int64 `json:"failed"`
}

type resolverZones map[string]struct {
	Requests struct {
		Name int64 `json:"name"`
		Srv  int64 `json:"srv"`
		Addr int64 `json:"addr"`
	} `json:"requests"`
	Responses struct {
		Noerror  int64 `json:"noerror"`
		Formerr  int64 `json:"formerr"`
		Servfail int64 `json:"servfail"`
		Nxdomain int64 `json:"nxdomain"`
		Notimp   int64 `json:"notimp"`
		Refused  int64 `json:"refused"`
		Timedout int64 `json:"timedout"`
		Unknown  int64 `json:"unknown"`
	} `json:"responses"`
}

type httpRequests struct {
	Total   int64 `json:"total"`
	Current int64 `json:"current"`
}

type requestsStats struct {
	Total      int64 `json:"total"`
	Processing int64 `json:"processing"`
	Discarded  int64 `json:"discarded"`
}

type responseStats struct {
	// Table of response codes (100-599)
	Response100 *int64 `json:"100"`
	Response101 *int64 `json:"101"`
	Response102 *int64 `json:"102"`
	Response103 *int64 `json:"103"`
	Response200 *int64 `json:"200"`
	Response201 *int64 `json:"201"`
	Response202 *int64 `json:"202"`
	Response203 *int64 `json:"203"`
	Response204 *int64 `json:"204"`
	Response205 *int64 `json:"205"`
	Response206 *int64 `json:"206"`
	Response207 *int64 `json:"207"`
	Response208 *int64 `json:"208"`
	Response226 *int64 `json:"226"`
	Response300 *int64 `json:"300"`
	Response301 *int64 `json:"301"`
	Response302 *int64 `json:"302"`
	Response303 *int64 `json:"303"`
	Response304 *int64 `json:"304"`
	Response305 *int64 `json:"305"`
	Response307 *int64 `json:"307"`
	Response308 *int64 `json:"308"`
	Response400 *int64 `json:"400"`
	Response401 *int64 `json:"401"`
	Response402 *int64 `json:"402"`
	Response403 *int64 `json:"403"`
	Response404 *int64 `json:"404"`
	Response405 *int64 `json:"405"`
	Response406 *int64 `json:"406"`
	Response407 *int64 `json:"407"`
	Response408 *int64 `json:"408"`
	Response409 *int64 `json:"409"`
	Response410 *int64 `json:"410"`
	Response411 *int64 `json:"411"`
	Response412 *int64 `json:"412"`
	Response413 *int64 `json:"413"`
	Response414 *int64 `json:"414"`
	Response415 *int64 `json:"415"`
	Response416 *int64 `json:"416"`
	Response417 *int64 `json:"417"`
	Response418 *int64 `json:"418"`
	Response421 *int64 `json:"421"`
	Response422 *int64 `json:"422"`
	Response423 *int64 `json:"423"`
	Response424 *int64 `json:"424"`
	Response425 *int64 `json:"425"`
	Response426 *int64 `json:"426"`
	Response428 *int64 `json:"428"`
	Response429 *int64 `json:"429"`
	Response431 *int64 `json:"431"`
	Response451 *int64 `json:"451"`
	Response500 *int64 `json:"500"`
	Response501 *int64 `json:"501"`
	Response502 *int64 `json:"502"`
	Response503 *int64 `json:"503"`
	Response504 *int64 `json:"504"`
	Response505 *int64 `json:"505"`
	Response506 *int64 `json:"506"`
	Response507 *int64 `json:"507"`
	Response508 *int64 `json:"508"`
	Response510 *int64 `json:"510"`
	Response511 *int64 `json:"511"`
}

type httpServerZones map[string]struct {
	Ssl       *ssl          `json:"ssl"`
	Requests  requestsStats `json:"requests"`
	Responses responseStats `json:"responses"`
	Data      data          `json:"data"`
}

type httpLocationZones map[string]struct {
	Requests  requestsStats `json:"requests"`
	Responses responseStats `json:"responses"`
	Data      data          `json:"data"`
}

type healthCheckStats struct {
	Fails     int64 `json:"fails"`
	Unavaible int64 `json:"unavailable"`
	Downtime  int64 `json:"downtime"`
	// downstart?
}

type selected struct {
	Current int64   `json:"current"`
	Total   int64   `json:"total"`
	Last    *string `json:"last"`
}

type data struct {
	Sent     int64 `json:"sent"`
	Received int64 `json:"received"`
}

type httpUpstreams map[string]struct {
	Peers map[string]struct {
		Service   *string          `json:"service"`
		Backup    bool             `json:"backup"`
		Weight    int              `json:"weight"`
		State     string           `json:"state"`
		Selected  selected         `json:"selected"`
		MaxConns  *int             `json:"max_conns"`
		Responses responseStats    `json:"responses"`
		Data      data             `json:"data"`
		Health    healthCheckStats `json:"health"`
		SID       string           `json:"sid"`
	} `json:"peers"`
	Keepalive int `json:"keepalive"`
	// Zombies   int `json:"zombies"`
	// Queue     *struct {
	// 	Size      int   `json:"size"`
	// 	MaxSize   int   `json:"max_size"`
	// 	Overflows int64 `json:"overflows"`
	// } `json:"queue"`
	// backup_switch is also only in Pro version
}

type streamServerZones map[string]struct {
	Processing  int            `json:"processing"`
	Connections int            `json:"connections"`
	Sessions    *responseStats `json:"sessions"`
	Discarded   *int64         `json:"discarded"`
	Received    int64          `json:"received"`
	Sent        int64          `json:"sent"`
}

type streamUpstreams map[string]struct {
	Peers []struct {
		ID            int              `json:"id"`
		Server        string           `json:"server"`
		Backup        bool             `json:"backup"`
		Weight        int              `json:"weight"`
		State         string           `json:"state"`
		Active        int              `json:"active"`
		Connections   int64            `json:"connections"`
		ConnectTime   *int             `json:"connect_time"`
		FirstByteTime *int             `json:"first_byte_time"`
		ResponseTime  *int             `json:"response_time"`
		Sent          int64            `json:"sent"`
		Received      int64            `json:"received"`
		Fails         int64            `json:"fails"`
		Unavail       int64            `json:"unavail"`
		HealthChecks  healthCheckStats `json:"health_checks"`
		Downtime      int64            `json:"downtime"`
	} `json:"peers"`
	Zombies int `json:"zombies"`
}

type basicHitStats struct {
	Responses int64 `json:"responses"`
	Bytes     int64 `json:"bytes"`
}

type extendedHitStats struct {
	basicHitStats
	ResponsesWritten int64 `json:"responses_written"`
	BytesWritten     int64 `json:"bytes_written"`
}

type httpCaches map[string]struct {
	Size        int64            `json:"size"`
	MaxSize     int64            `json:"max_size"`
	Cold        bool             `json:"cold"`
	Hit         basicHitStats    `json:"hit"`
	Stale       basicHitStats    `json:"stale"`
	Updating    basicHitStats    `json:"updating"`
	Revalidated *basicHitStats   `json:"revalidated"`
	Miss        extendedHitStats `json:"miss"`
	Expired     extendedHitStats `json:"expired"`
	Bypass      extendedHitStats `json:"bypass"`
}

type httpLimitReqs map[string]struct {
	Passed         int64 `json:"passed"`
	Delayed        int64 `json:"delayed"`
	Rejected       int64 `json:"rejected"`
	DelayedDryRun  int64 `json:"delayed_dry_run"`
	RejectedDryRun int64 `json:"rejected_dry_run"`
}
