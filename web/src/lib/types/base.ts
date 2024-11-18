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

export interface Entrypoint {
	name: string;
	address: string;
	asDefault: boolean;
	http?: {
		middlewares: string[];
		tls?: TLSConfiguration;
	};
}

export interface User {
	id: number;
	username: string;
	password: string;
	email: string;
	isAdmin: boolean;
}

export function newUser(): User {
	return {
		id: 0,
		username: '',
		password: '',
		email: '',
		isAdmin: false
	};
}

export interface Agent {
	id: number;
	profileId: number;
	hostname: string;
	publicIp: string;
	privateIps: string;
	containers: any;
	activeIp: string;
	deleted: boolean;
	lastSeen: string;
}

export interface DNSProvider {
	id: number;
	name: string;
	type: string;
	externalIp: string;
	apiKey?: string;
	apiUrl?: string;
	zoneType?: string;
	proxied: boolean;
	isActive: boolean;
}

export function newProvider(): DNSProvider {
	return {
		id: 0,
		name: '',
		type: 'cloudflare',
		externalIp: '',
		apiKey: '',
		apiUrl: '',
		proxied: false,
		isActive: false
	};
}

export interface Setting {
	key: string;
	value: string;
}

export interface TLSConfiguration {
	certificates?: (CertAndStores | undefined)[];
	options?: Record<string, Options>;
	stores?: Record<string, Store>;
}
