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

export interface DNSProviderConfig {
	apiKey: string;
	apiUrl: string;
	traefikIp: string;
	proxied: boolean;
	zoneType: string;
}

export interface DNSProvider {
	id: number;
	name: string;
	type: string;
	config: DNSProviderConfig;
	isActive: boolean;
	createdAt?: string;
	updatedAt?: string;
}

export enum DNSProviderTypes {
	CLOUDFLARE = 'cloudflare',
	POWERDNS = 'powerdns',
	TECHNITIUM = 'technitium'
}

export interface User {
	id: number;
	username: string;
	password?: string;
	email?: string;
	isAdmin: boolean;
	lastLogin?: string;
	createdAt: string;
	updatedAt: string;
}

export interface Agent {
	id: string;
	profileId: number;
	hostname: string;
	publicIp?: string;
	private_ips: string[];
	containers: Record<string, unknown>[];
	activeIp?: string;
	token: string;
	createdAt: string;
	updatedAt: string;
}

export interface Settings {
	key: Record<string, Setting>;
}

export interface Setting {
	value: string | number | boolean;
	description: string;
}

export interface UpsertSettingsParams {
	key: string;
	value: string;
	description: string;
}

export interface RouterDNSProvider {
	traefikId: number;
	providerId: number;
	routerName: string;
	providerName: string;
	providerType: string;
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
