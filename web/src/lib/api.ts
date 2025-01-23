import {
	TraefikSource,
	type Agent,
	type DNSProvider,
	type Plugin,
	type Profile,
	type Settings,
	type TraefikConfig,
	type UpsertSettingsParams,
	type User
} from './types';
import type { EntryPoints } from './types/entrypoints';
import { PROFILE_SK, TOKEN_SK } from './store';
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
import { get, writable, type Writable } from 'svelte/store';
import { toast } from 'svelte-sonner';
import { goto } from '$app/navigation';
import type { Overview } from './types/overview';

// Global state variables
export const BACKEND_PORT = import.meta.env.PORT || 3000;
export const BASE_URL = import.meta.env.PROD ? '/api' : `http://127.0.0.1:${BACKEND_PORT}/api`;

// Stores
export const profiles: Writable<Profile[]> = writable();
export const traefik: Writable<TraefikConfig> = writable();
export const entrypoints: Writable<EntryPoints[]> = writable([]);
export const overview: Writable<Overview> = writable({} as Overview);
export const version: Writable<string> = writable('');
export const routers: Writable<Router[]> = writable([]);
export const services: Writable<Service[]> = writable([]);
export const middlewares: Writable<Middleware[]> = writable([]);
export const users: Writable<User[]> = writable([]);
export const dnsProviders: Writable<DNSProvider[]> = writable([]);
export const agents: Writable<Agent[]> = writable([]);
export const settings: Writable<Settings> = writable({} as Settings);
export const plugins: Writable<Plugin[]> = writable([]);

// App state
export const profile: Writable<Profile> = writable({} as Profile);
export const user: Writable<User | null> = writable({} as User);
export const source: Writable<TraefikSource> = writable({} as TraefikSource);

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
	const token = localStorage.getItem(TOKEN_SK);

	// Custom fetch function that adds the Authorization header
	const customFetch: typeof window.fetch = async (url, options) => {
		const headers = new Headers(options?.headers); // Get existing headers
		if (token) {
			headers.set('Authorization', 'Bearer ' + token); // Add the Authorization header
		}
		const customOptions = {
			'Content-Type': 'application/json',
			...options,
			headers
		};
		return fetch ? fetch(url, customOptions) : window.fetch(url, customOptions); // Use custom fetch or default
	};

	try {
		loading.set(true);
		const response = await customFetch(`${BASE_URL}${endpoint}`, {
			method: options.method || 'GET',
			body: options.body ? JSON.stringify(options.body) : undefined
		});

		if (!response.ok) {
			throw new Error(`${await response.text()}`);
		}
		if (response.status !== 204) {
			return await response.json();
		}
	} catch (err: unknown) {
		error.set(err instanceof Error ? err.message : String(err));
		loading.set(false);
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
			localStorage.setItem(TOKEN_SK, data.token);
			goto('/');
		}
		return data;
	},

	async verify(fetch: typeof window.fetch = window.fetch) {
		const token = localStorage.getItem(TOKEN_SK);
		try {
			const data = await send(
				'/verify',
				{
					method: 'POST',
					body: token
				},
				fetch
			);
			return data;
		} catch (err: unknown) {
			const error = err instanceof Error ? err.message : String(err);
			toast.error('Session expired', { description: error });
			api.logout();
			return;
		}
	},

	async resetPassword(token: string, password: string) {
		const data = await send('/reset', {
			method: 'POST',
			body: { token, password }
		});
		return data;
	},

	logout() {
		localStorage.removeItem(TOKEN_SK);
		goto('/login');
	},

	// Profiles ------------------------------------------------------------------
	async listProfiles() {
		const data = await send('/profile');
		profiles.set(data);
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

	async updateProfile(profile: Profile) {
		await send('/profile', {
			method: 'PUT',
			body: profile
		});
		await api.listProfiles(); // Refresh the list
	},

	async deleteProfile(id: number) {
		await send(`/profile/${id}`, { method: 'DELETE' });
		if (id === get(profile).id) {
			localStorage.removeItem(PROFILE_SK);
			profile.set({} as Profile);
		}
		await api.listProfiles(); // Refresh the list
	},

	// Traefik -------------------------------------------------------------------
	async getTraefikConfig(id: number, source: TraefikSource) {
		if (source === TraefikSource.API) {
			// For API source, fetch everything from Traefik
			const hasMetadata = await fetchTraefikMetadata(id);
			if (hasMetadata) {
				await fetchTraefikConfig(id, source);
			}
		} else {
			// For everything else, fetch metadata from API and config from source
			await fetchTraefikMetadata(id);
			await fetchTraefikConfig(id, source);
		}
	},

	async getDynamicConfig(profileName: string) {
		return await send(`/${profileName}`);
	},

	// Routers -------------------------------------------------------------------
	async upsertRouter(id: number, data: UpsertRouterParams) {
		await send(`/router/${id}`, {
			method: 'POST',
			body: data
		});
		await api.getTraefikConfig(id, TraefikSource.LOCAL);
	},

	async deleteRouter(id: number, data: Router) {
		await send(`/router/${id}/${data.name}/${data.protocol}`, {
			method: 'DELETE'
		});
		await api.getTraefikConfig(id, TraefikSource.LOCAL);
	},

	// Middlewares ---------------------------------------------------------------
	async upsertMiddleware(id: number, data: UpsertMiddlewareParams) {
		await send(`/middleware/${id}`, {
			method: 'POST',
			body: data
		});
		await api.getTraefikConfig(id, TraefikSource.LOCAL);
	},

	async deleteMiddleware(id: number, data: Middleware) {
		await send(`/middleware/${id}/${data.name}/${data.protocol}`, {
			method: 'DELETE'
		});
		await api.getTraefikConfig(id, TraefikSource.LOCAL);
	},

	// DNS Providers -------------------------------------------------------------
	async listDNSProviders() {
		const data = await send('/provider');
		dnsProviders.set(data);
	},

	async getDNSProvider(id: number) {
		return await send(`/provider/${id}`);
	},

	async createDNSProvider(provider: Omit<DNSProvider, 'id' | 'created_at' | 'updated_at'>) {
		await send('/provider', {
			method: 'POST',
			body: provider
		});
		await api.listDNSProviders();
	},

	async updateDNSProvider(provider: Omit<DNSProvider, 'created_at' | 'updated_at'>) {
		await send(`/provider/${provider.id}`, {
			method: 'PUT',
			body: provider
		});
		await api.listDNSProviders();
	},

	async deleteDNSProvider(id: number) {
		await send(`/provider/${id}`, {
			method: 'DELETE'
		});
		await api.listDNSProviders();
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
	async listAgents() {
		const data = await send('/agent');
		agents.set(data);
	},

	async getAgent(id: string) {
		return await send(`/agent/${id}`);
	},

	async createAgent(profileID: number) {
		await send(`/agent/${profileID}`, { method: 'POST' });
		await api.listAgents();
	},

	async updateAgent(agent: Omit<Agent, 'created_at' | 'updated_at'>) {
		await send(`/agent/${agent.id}`, {
			method: 'PUT',
			body: agent
		});
		await api.listAgents();
	},

	async deleteAgent(id: string) {
		await send(`/agent/${id}`, {
			method: 'DELETE'
		});
		await api.listAgents();
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
async function fetchTraefikMetadata(id: number) {
	const res = await send(`/traefik/${id}/${TraefikSource.API}`);
	if (!res) {
		// Reset metadata stores
		traefik.set({} as TraefikConfig);
		overview.set({} as Overview);
		entrypoints.set([]);
		version.set('');
		return false;
	}

	// Set metadata stores
	traefik.set(res);
	overview.set(res.overview);
	entrypoints.set(res.entrypoints);
	version.set(res.version);
	return true;
}

async function fetchTraefikConfig(id: number, source: TraefikSource) {
	const res = await send(`/traefik/${id}/${source}`);
	if (!res) {
		// Reset routing stores
		routers.set([]);
		services.set([]);
		middlewares.set([]);
		return;
	}

	// Set routing stores
	routers.set(flattenRouterData(res.config));
	services.set(flattenServiceData(res.config));
	middlewares.set(flattenMiddlewareData(res.config));
}

// Login ----------------------------------------------------------------------
// export async function login(username: string, password: string, remember: boolean) {
// 	const loginURL = remember ? `${BASE_URL}/login?remember=true` : `${BASE_URL}/login`;
// 	const response = await fetch(loginURL, {
// 		method: 'POST',
// 		body: JSON.stringify({ username, password })
// 	});
// 	if (response.ok) {
// 		const { token } = await response.json();
// 		localStorage.setItem(TOKEN_SK, token);
// 		loggedIn.set(true);
// 		goto('/');
// 		toast.success('Login successful');
// 		await getProfiles();
// 		await getProviders();
// 	} else {
// 		toast.error('Login failed', {
// 			description: await response.text(),
// 			duration: 3000
// 		});
// 		return;
// 	}
// }

// export async function sendResetEmail(username: string) {
// 	const response = await fetch(`${BASE_URL}/reset/${username}`, {
// 		method: 'POST'
// 	});
// 	if (response.ok) {
// 		toast.success('Password reset email sent!');
// 	} else {
// 		toast.error('Request failed', {
// 			description: await response.text(),
// 			duration: 3000
// 		});
// 	}
// }

// export async function resetPassword(token: string, password: string) {
// 	const response = await fetch(`${BASE_URL}/reset`, {
// 		method: 'POST',
// 		body: JSON.stringify({ password }),
// 		headers: { Authorization: `Bearer ${token}` }
// 	});
// 	if (response.ok) {
// 		toast.success('Password reset successful!');
// 		goto('/login');
// 	} else {
// 		toast.error('Request failed', {
// 			description: await response.text(),
// 			duration: 3000
// 		});
// 	}
// }
