import type { EntryPoints } from './types/entrypoints';
import type { Overview } from './types/overview';
import type { Middleware } from './types/middlewares';
import type { Router, Service } from './types/router';

// Base types ------------------------------------------------------------------
export interface Profile {
	id: number;
	name: string;
	url: string;
	username?: string;
	password?: string;
	tls: boolean;
	created_at: string;
	updated_at: string;
}

export interface TraefikConfig {
	id: number;
	profile_id: number;
	source: TraefikSource;
	entrypoints: EntryPoints;
	overview: Overview;
	config: BaseTraefikConfig;
	created_at: string;
	updated_at: string;
}

export enum TraefikSource {
	API = 'api',
	LOCAL = 'local',
	AGENT = 'agent'
}

export interface BaseTraefikConfig {
	routers: Record<string, Router>;
	tcpRouters: Record<string, Router>;
	udpRouters: Record<string, Router>;
	services: Record<string, Service>;
	tcpServices: Record<string, Service>;
	udpServices: Record<string, Service>;
	middlewares: Record<string, Middleware>;
	tcpMiddlewares: Record<string, Middleware>;
}

export interface DNSProvider {
	id: number;
	name: string;
	type: string;
	config: Record<string, unknown>;
	is_active: boolean;
	created_at: string;
	updated_at: string;
}

export interface User {
	id: number;
	username: string;
	password?: string;
	email?: string;
	is_admin: boolean;
	last_login?: string;
	created_at: string;
	updated_at: string;
}

export interface Agent {
	id: string;
	profile_id: number;
	hostname: string;
	public_ip?: string;
	private_ips: string[];
	containers: Record<string, unknown>[];
	active_ip?: string;
	token: string;
	created_at: string;
	updated_at: string;
}

export interface Setting {
	key: string;
	value: string;
	updated_at: string;
}

export interface Plugin {
	id: string;
	name: string;
	displayName: string;
	author: string;
	type: string;
	import: string;
	summary: string;
	iconUrl: string;
	bannerUrl: string;
	readme: string;
	latestVersion: string;
	versions: string[];
	stars: number;
	snippet: Record<string, string>;
	createdAt: string;
}
