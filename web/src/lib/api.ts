import {
	TraefikSource,
	type Agent,
	type BackupFile,
	type DNSProvider,
	type Plugin,
	type Profile,
	type RouterDNSProvider,
	type Settings,
	type Stats,
	type TraefikConfig,
	type UpdateAgentIPParams,
	type UpsertSettingsParams,
	type User
} from './types';
import type { EntryPoints } from './types/entrypoints';
import {
	flattenMiddlewareData,
	type Middleware,
	type UpsertMiddlewareParams
} from './types/middlewares';
import {
	flattenRouterData,
	flattenServiceData,
	type Router,
	type Service,
	type UpsertRouterParams
} from './types/router';
import YAML from 'yaml';
import { get, writable, type Writable } from 'svelte/store';
import { toast } from 'svelte-sonner';
import { goto } from '$app/navigation';
import type { Overview } from './types/overview';
import { profile } from './stores/profile';
import { user } from './stores/user';
import { source } from './stores/source';
import { token } from './stores/common';

export type RouterWithService = { router: Router; service: Service };

// Global state variables
export const BACKEND_PORT = import.meta.env.PORT || 3000;
export const BASE_URL = import.meta.env.PROD ? '/api' : `http://127.0.0.1:${BACKEND_PORT}/api`;

// DB Stores
export const profiles: Writable<Profile[]> = writable();
export const traefik: Writable<TraefikConfig[]> = writable();
export const entrypoints: Writable<EntryPoints[]> = writable([]);
export const overview: Writable<Overview> = writable({} as Overview);
export const version: Writable<string> = writable('');
export const routers: Writable<Router[]> = writable([]);
export const services: Writable<Service[]> = writable([]);
export const routerServiceMerge: Writable<RouterWithService[]> = writable([]);
export const middlewares: Writable<Middleware[]> = writable([]);
export const users: Writable<User[]> = writable([]);
export const rdps: Writable<RouterDNSProvider[]> = writable([]);
export const dnsProviders: Writable<DNSProvider[]> = writable([]);
export const agents: Writable<Agent[]> = writable([]);
export const settings: Writable<Settings> = writable({} as Settings);
export const plugins: Writable<Plugin[]> = writable([]);
export const backups: Writable<BackupFile[]> = writable([]);

// App state
export const stats: Writable<Stats> = writable({} as Stats);
export const mwNames: Writable<string[]> = writable([]);
export const dynamicJSON: Writable<string> = writable('');
export const dynamicYAML: Writable<string> = writable('');

// Loading and error states
export const loading = writable<boolean>(false);
export const error = writable<string | null>(null);

interface APIOptions {
	method?: string;
	// eslint-disable-next-line
	body?: any;
	headers?: Record<string, string>;
}

async function send(endpoint: string, options: APIOptions = {}, fetch?: typeof window.fetch) {
	// Custom fetch function that adds the Authorization header
	const customFetch: typeof window.fetch = async (url, options) => {
		const headers = new Headers(options?.headers); // Get existing headers
		if (token.value) {
			headers.set('Authorization', 'Bearer ' + token.value); // Add the Authorization header
		}
		// Don't set Content-Type for FormData
		const isFormData = options?.body instanceof FormData;
		if (!isFormData) {
			headers.set('Content-Type', 'application/json');
		}
		const customOptions = {
			...options,
			headers,
			body: isFormData ? options?.body : options?.body ? JSON.stringify(options.body) : undefined
		};
		return fetch ? fetch(url, customOptions) : window.fetch(url, customOptions); // Use custom fetch or default
	};

	try {
		loading.set(true);
		const response = await customFetch(`${BASE_URL}${endpoint}`, {
			method: options.method || 'GET',
			body: options.body,
			headers: options.headers
		});

		if (!response.ok) {
			throw new Error(await response.text());
		}
		if (response.headers.get('Content-Type') === 'application/json') {
			return await response.json();
		} else {
			return response;
		}
	} catch (err: unknown) {
		error.set(err instanceof Error ? err.message : String(err));
		throw err;
	} finally {
		loading.set(false);
	}
}

export const api = {
	// Auth ----------------------------------------------------------------------
	async login(username: string, password: string, remember: boolean = false) {
		const endpoint = remember ? `login?remember=true` : `/login`;
		const data = await send(endpoint, {
			method: 'POST',
			body: { username, password }
		});
		if (data.token) {
			token.value = data.token;
			goto('/');
		}

		await api.load();
	},

	async verify(fetch: typeof window.fetch = window.fetch) {
		try {
			const data = await send(
				'/verify',
				{
					method: 'POST',
					body: token.value
				},
				fetch
			);

			if (data) {
				user.value = data;
			}
		} catch (err: unknown) {
			const error = err instanceof Error ? err.message : String(err);
			toast.error('Session expired', { description: error });
			api.logout();
			return;
		}
	},

	async sendResetEmail(username: string) {
		await send(`/reset/${username}`, { method: 'POST' });
	},

	async verifyOTP(username: string, otp: string) {
		const data = await send('/verify/otp', {
			method: 'POST',
			body: { username, token: otp }
		});

		if (data.token) {
			token.value = data.token;
			goto('/');
		}
		await api.load();
	},

	async load() {
		if (!user.isLoggedIn()) return;

		// Load profiles
		await api.listProfiles();
		if (get(profiles) && !profile.id) {
			profile.value = get(profiles)[0];
		}

		// Load Traefik Config
		if (profile.value && source.value) {
			await api.getTraefikConfig(source.value);
		}
	},

	logout() {
		token.value = null;
		user.clear();
		goto('/login');
	},

	async loadStats() {
		await api.listDNSProviders();
		await api.listUsers();
		await api.listProfiles();
		const agents = await api.listAgents();

		stats.set({
			profiles: get(profiles)?.length ?? 0,
			agents: agents?.length ?? 0,
			users: get(users)?.length ?? 0,
			dnsProviders: get(dnsProviders)?.length ?? 0,
			activeDNS: get(dnsProviders)?.find((item) => item.isActive)?.name ?? 'None'
		});
	},

	// Profiles ------------------------------------------------------------------
	async listProfiles() {
		const data: Profile[] = await send('/profile');
		profiles.set(data);
		if (profile.isValid()) {
			profile.value = data.find((item) => item.id === profile.value?.id) ?? data[0];
		}
	},

	async getProfile(id: number) {
		return await send(`/profile/${id}`);
	},

	async createProfile(profile: Omit<Profile, 'id' | 'created_at' | 'updated_at'>) {
		await send('/profile', {
			method: 'POST',
			body: profile
		});
		await api.listProfiles(); // Refresh the list
	},

	async updateProfile(p: Profile) {
		await send('/profile', {
			method: 'PUT',
			body: p
		});
		await api.listProfiles(); // Refresh the list
		if (p.id === profile.id) {
			await api.getTraefikConfig(TraefikSource.API);
		}
	},

	async deleteProfile(id: number) {
		await send(`/profile/${id}`, { method: 'DELETE' });
		if (id === profile.value?.id) {
			profile.value = {} as Profile;
		}
		await api.listProfiles(); // Refresh the list
	},

	// Traefik -------------------------------------------------------------------
	async getTraefikConfig(source: TraefikSource) {
		await fetchTraefikMetadata();
		await fetchTraefikConfig(source);
	},

	async getTraefikConfigLocal() {
		if (!profile.isValid()) {
			toast.error('No valid profile selected');
			return;
		}
		// Get the local config without mutating the stores
		const res = await send(`/traefik/${profile.id}/${TraefikSource.LOCAL}`);
		if (!res) {
			return;
		}
		const traefik = res as TraefikConfig;
		const routers = flattenRouterData(res);
		const services = flattenServiceData(res);
		const middlewares = flattenMiddlewareData(res);
		return { traefik, routers, services, middlewares };
	},

	async getDynamicConfig() {
		if (!profile.hasValidName()) {
			toast.error('Profile name is required');
			return;
		}
		const data = await send(`/${profile.name}`);
		if (!data || (typeof data === 'object' && Object.keys(data).length === 0)) {
			return;
		}
		dynamicJSON.set(JSON.stringify(data, null, 2));
		dynamicYAML.set(YAML.stringify(data));
	},

	// Routers -------------------------------------------------------------------
	async upsertRouter(data: UpsertRouterParams) {
		if (!profile.hasValidId()) {
			toast.error('Invalid profile ID');
			return;
		}

		await send(`/router/${profile.id}`, {
			method: 'POST',
			body: data
		});
		await api.getTraefikConfig(TraefikSource.LOCAL);
	},

	async deleteRouter(data: Router) {
		if (!profile.hasValidId()) {
			toast.error('Invalid profile ID');
			return;
		}

		await send(`/router/${profile.id}/${data.name}/${data.protocol}`, {
			method: 'DELETE'
		});
		await api.getTraefikConfig(TraefikSource.LOCAL);
	},

	// Middlewares ---------------------------------------------------------------
	async upsertMiddleware(data: UpsertMiddlewareParams) {
		if (!profile.hasValidId()) {
			toast.error('Invalid profile ID');
			return;
		}
		await send(`/middleware/${profile.id}`, {
			method: 'POST',
			body: data
		});
		await api.getTraefikConfig(TraefikSource.LOCAL);
	},

	async deleteMiddleware(data: Middleware) {
		if (!profile.hasValidId()) {
			toast.error('Invalid profile ID');
			return;
		}
		await send(`/middleware/${profile.id}/${data.name}/${data.protocol}`, {
			method: 'DELETE'
		});
		await api.getTraefikConfig(TraefikSource.LOCAL);
	},

	// DNS Providers -------------------------------------------------------------
	async listDNSProviders() {
		const data = await send('/dns');
		dnsProviders.set(data);
	},

	async getDNSProvider(id: number) {
		return await send(`/dns/${id}`);
	},

	async createDNSProvider(provider: Omit<DNSProvider, 'id' | 'created_at' | 'updated_at'>) {
		await send('/dns', {
			method: 'POST',
			body: provider
		});
		await api.listDNSProviders();
	},

	async updateDNSProvider(provider: Omit<DNSProvider, 'created_at' | 'updated_at'>) {
		await send(`/dns`, {
			method: 'PUT',
			body: provider
		});
		await api.listDNSProviders();
	},

	async deleteDNSProvider(id: number) {
		if (!id) return;
		await send(`/dns/${id}`, {
			method: 'DELETE'
		});
		await api.listDNSProviders();
	},

	async getRouterDNSProvider(traefikId: number, routerName: string) {
		if (!traefikId || !routerName) return;
		return await send(`/dns/router`, {
			method: 'GET',
			body: { traefikId, routerName }
		});
	},

	async listRouterDNSProviders(traefikId: number) {
		if (!traefikId) return;
		const data = await send(`/dns/router/${traefikId}`, {
			method: 'GET'
		});
		rdps.set(data);
	},

	async setRouterDNSProvider(traefikId: number, providerId: number, routerName: string) {
		await send(`/dns/router`, {
			method: 'POST',
			body: { traefikId, providerId, routerName }
		});
		await api.listRouterDNSProviders(traefikId);
	},

	async deleteRouterDNSProvider(traefikId: number, routerName: string) {
		await send(`/dns/router`, {
			method: 'DELETE',
			body: { traefikId, routerName }
		});
		await api.listRouterDNSProviders(traefikId);
	},

	// Users ---------------------------------------------------------------------
	async listUsers() {
		const data = await send('/user');
		users.set(data);
	},

	async getUser(id: number) {
		const data = await send(`/user/${id}`);
		return data;
	},

	async createUser(user: Omit<User, 'id' | 'created_at' | 'updated_at'>) {
		await send('/user', {
			method: 'POST',
			body: user
		});
		await api.listUsers();
	},

	async updateUser(user: Omit<User, 'created_at' | 'updated_at'>) {
		await send(`/user`, {
			method: 'PUT',
			body: user
		});
		await api.listUsers();
	},

	async deleteUser(id: number) {
		await send(`/user/${id}`, {
			method: 'DELETE'
		});
		await api.listUsers();
	},

	// Agents
	async listAgents(): Promise<Agent[]> {
		const data = await send(`/agent`);
		return data;
	},

	async listAgentsByProfile(): Promise<Agent[]> {
		if (!profile.hasValidId()) {
			toast.error('Invalid profile ID');
			return [];
		}
		const data = await send(`/agent/list/${profile.id}`);
		agents.set(data);
		return data;
	},

	async getAgent(id: string) {
		return await send(`/agent/${id}`);
	},

	async createAgent() {
		if (!profile.hasValidId()) {
			toast.error('Invalid profile ID');
			return;
		}
		await send(`/agent/${profile.id}`, {
			method: 'POST'
		});
		await api.listAgentsByProfile();
	},

	async updateAgentIP(params: UpdateAgentIPParams) {
		await send(`/agent`, {
			method: 'PUT',
			body: params
		});
		await api.listAgentsByProfile();
	},

	async deleteAgent(id: string) {
		await send(`/agent/${id}`, {
			method: 'DELETE'
		});
		await api.listAgentsByProfile();
	},

	// Settings ------------------------------------------------------------------
	async listSettings() {
		const data = await send('/settings');
		settings.set(data);
	},

	async getSetting(id: number) {
		return await send(`/settings/${id}`);
	},

	async upsertSetting(setting: UpsertSettingsParams) {
		await send(`/settings`, {
			method: 'POST',
			body: setting
		});
		await api.listSettings();
	},

	// Backups -------------------------------------------------------------------
	async listBackups() {
		const data = await send('/backups');
		backups.set(data);
	},

	async createBackup() {
		await send('/backups', { method: 'POST' });
		await api.listBackups();
	},

	async downloadBackup() {
		try {
			const response = await send('/backups/download', { method: 'GET' });

			const blob = await response.blob();
			const filename =
				response.headers.get('content-disposition')?.split('filename=')[1] || 'backup.db';
			const url = window.URL.createObjectURL(blob);
			const a = document.createElement('a');
			a.href = url;
			a.download = filename;
			document.body.appendChild(a);
			a.click();
			window.URL.revokeObjectURL(url);
			document.body.removeChild(a);
		} catch (error) {
			throw new Error(`Failed to download backup: ${error}`);
		}
	},

	async downloadBackupByName(name: string) {
		try {
			const response = await send(`/backups/download/${name}`, { method: 'GET' });
			const blob = await response.blob();
			const filename = response.headers.get('content-disposition')?.split('filename=')[1] || name;
			const url = window.URL.createObjectURL(blob);
			const a = document.createElement('a');
			a.href = url;
			a.download = filename;
			document.body.appendChild(a);
			a.click();
			window.URL.revokeObjectURL(url);
			document.body.removeChild(a);
		} catch (error) {
			throw new Error(`Failed to download backup: ${error}`);
		}
	},

	async restoreBackup(files: FileList | null) {
		if (!files?.length) return;
		const formData = new FormData();
		formData.append('file', files[0]);

		await send(`/backups/restore`, {
			method: 'POST',
			body: formData
		});
		toast.success('Backup restored successfully');
	},

	async deleteBackup(name: string) {
		await send(`/backups/${name}`, {
			method: 'DELETE'
		});
		await api.listBackups();
	},

	// Plugins
	async getMiddlewarePlugins() {
		const data = await send('/middleware/plugins');
		plugins.set(data);
	},

	async getVersion() {
		return await send('/version');
	}
};

// Helper
async function fetchTraefikMetadata() {
	if (!profile.isValid()) {
		toast.error('No valid profile selected');
		return;
	}
	const res = await send(`/traefik/${profile.id}/${TraefikSource.API}`);
	if (!res) {
		// Reset metadata stores
		overview.set({} as Overview);
		entrypoints.set([]);
		version.set('');
		return false;
	}

	// Set metadata stores
	const meta = res[0];
	overview.set(meta.overview);
	entrypoints.set(meta.entrypoints);
	version.set(meta.version);

	// Set middleware names (used for chain)
	const middlewares = flattenMiddlewareData(res);
	mwNames.set(middlewares.map((mw) => mw.name));
	return true;
}

async function fetchTraefikConfig(src: TraefikSource) {
	if (!profile.isValid() || !source.isValid(src)) {
		toast.error('No valid profile selected');
		return;
	}

	const res = await send(`/traefik/${profile.id}/${src}`);
	if (!res) {
		// Reset stores
		traefik.set([]);
		routers.set([]);
		services.set([]);
		middlewares.set([]);
		routerServiceMerge.set([]);
		return;
	}

	// Update stores with proper diffing
	traefik.update((current) => {
		return JSON.stringify(current) === JSON.stringify(res) ? current : res;
	});

	const newRouters = flattenRouterData(res);
	const newServices = flattenServiceData(res);
	const newMiddlewares = flattenMiddlewareData(res);
	const newMerge = newRouters.map((router) => {
		let service = newServices.find((service) => service.name === router.service);
		if (!service) {
			service = newServices.find((service) => service.name === router.name);
		}
		return { router, service: service || ({} as Service) };
	});

	routers.update((current) => {
		if (!current || current.length === 0) return newRouters;
		if (JSON.stringify(current) === JSON.stringify(newRouters)) return current;
		return newRouters;
	});

	services.update((current) => {
		if (!current || current.length === 0) return newServices;
		if (JSON.stringify(current) === JSON.stringify(newServices)) return current;
		return newServices;
	});

	middlewares.update((current) => {
		if (!current || current.length === 0) return newMiddlewares;
		if (JSON.stringify(current) === JSON.stringify(newMiddlewares)) return current;
		return newMiddlewares;
	});

	routerServiceMerge.update((current) => {
		if (!current || current.length === 0) return newMerge;
		if (JSON.stringify(current) === JSON.stringify(newMerge)) return current;
		return newMerge;
	});

	// Fetch the router dns relations
	await api.listRouterDNSProviders(res.id);

	// Fetch dynamic config
	await api.getDynamicConfig();
}
