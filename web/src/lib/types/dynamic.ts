import type { Router, Service } from './config';
import type { HttpMiddleware, TCPMiddleware } from './middlewares';
import type { CertAndStores, Options, Store } from './tls';

export interface Profile {
	name: string;
	instance: Instance;
}

export interface Instance {
	url: string;
	username: string;
	password: string;
	dynamic?: Dynamic;
}

export interface Dynamic {
	entrypoints?: Entrypoint[];
	routers?: Router[];
	services?: Service[];
	httpmiddlewares?: HttpMiddleware[];
	tcpmiddlewares?: TCPMiddleware[];
	version?: string;
}

export interface Entrypoint {
	name: string;
	address: string;
	asDefault: boolean;
	http?: {
		middlewares: string[];
		tls?: TLSConfiguration;
	};
}

export interface TLSConfiguration {
	certificates?: (CertAndStores | undefined)[];
	options?: Record<string, Options>;
	stores?: Record<string, Store>;
}
