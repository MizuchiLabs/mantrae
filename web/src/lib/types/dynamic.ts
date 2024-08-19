import type { Router, Service } from './config';
import type { Middleware } from './middlewares';
import type { CertAndStores, Options, Store } from './tls';

export interface Profile {
	name: string;
	url: string;
	username: string;
	password: string;
	dynamic?: Dynamic;
}

export interface Dynamic {
	entrypoints?: Entrypoint[];
	routers?: Record<string, Router>;
	services?: Record<string, Service>;
	middlewares?: Record<string, Middleware>;
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

export const newProfile = (): Profile => {
	return {
		name: '',
		url: '',
		username: '',
		password: '',
		dynamic: {}
	};
};
