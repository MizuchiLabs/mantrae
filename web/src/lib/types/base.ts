import type { Router, Service } from './config';
import type { Middleware } from './middlewares';
import type { Overview } from './overview';
import type { CertAndStores, Options, Store } from './tls';

export interface Profile {
	id?: number;
	name: string;
	url: string;
	username: string;
	password: string;
	tls: boolean;
}

export const newProfile = (): Profile => {
	return {
		name: '',
		url: '',
		username: '',
		password: '',
		tls: false
	};
};

export interface Config {
	profile_id: number;
	overview?: Overview;
	entrypoints?: Entrypoint[];
	routers?: Record<string, Router>;
	services?: Record<string, Service>;
	middlewares?: Record<string, Middleware>;
	tls?: TLSConfiguration;
	version?: string;
}

export interface User {
	id: number;
	username: string;
	password: string;
	email: string;
	type: string;
}

export function newUser(): User {
	return {
		id: 0,
		username: '',
		password: '',
		email: '',
		type: 'user'
	};
}

export interface DNSProvider {
	id: number;
	name: string;
	type: string;
	external_ip: string;
	api_key?: string;
	api_url?: string;
	proxied: boolean;
	is_active: boolean;
}

export function newProvider(): DNSProvider {
	return {
		id: 0,
		name: '',
		type: 'cloudflare',
		external_ip: '',
		api_key: '',
		api_url: '',
		proxied: false,
		is_active: false
	};
}

export interface Setting {
	id: number;
	key: string;
	value: string;
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
