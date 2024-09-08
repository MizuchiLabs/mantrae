export interface Overview {
	http: HTTPOverview;
	tcp: TCPOverview;
	udp: UDPOverview;
	features: {
		tracing: string;
		metrics: string;
		accessLog: boolean;
	};
	providers: string[];
}

export interface BasicOverview {
	total: number;
	warnings: number;
	errors: number;
}

export interface HTTPOverview {
	routers: BasicOverview;
	services: BasicOverview;
	middlewares: BasicOverview;
}

export interface TCPOverview {
	routers: BasicOverview;
	services: BasicOverview;
	middlewares: BasicOverview;
}

export interface UDPOverview {
	routers: BasicOverview;
	services: BasicOverview;
}
