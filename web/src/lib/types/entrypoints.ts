export interface EntryPoints {
	address: string;
	asDefault: boolean;
	transport: Transport;
	// forwardedHeaders: ForwardedHeaders;
	http: HTTP;
	http2: HTTP2;
	udp: UDP;
	observability: Observability;
	name: string;
}

export interface Transport {
	lifeCycle: LifeCycle;
	respondingTimeouts: RespondingTimeouts;
}

export interface LifeCycle {
	graceTimeOut: string;
}

export interface RespondingTimeouts {
	readTimeout: string;
	idleTimeout: string;
}

// export interface ForwardedHeaders {}

export interface HTTP {
	tls: TLS;
	maxHeaderBytes: number;
}

export interface TLS {
	certResolver: string;
}

export interface HTTP2 {
	maxConcurrentStreams: number;
}

export interface UDP {
	timeout: string;
}

export interface Observability {
	accessLogs: boolean;
	tracing: boolean;
	metrics: boolean;
}
