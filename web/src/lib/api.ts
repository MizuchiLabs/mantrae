import { get, writable, type Writable } from 'svelte/store';
import type { Agent, DNSProvider, Plugin, Profile, Setting, TraefikConfig, User } from './types';
import { PROFILE_SK, TOKEN_SK } from './store';
import { flattenRouterData, flattenServiceData, type Router, type Service } from './types/router';
import type { EntryPoints } from './types/entrypoints';
import { flattenMiddlewareData, type Middleware } from './types/middlewares';

// Global state variables
export const BACKEND_PORT = import.meta.env.PORT || 3000;
export const BASE_URL = import.meta.env.PROD ? '/api' : `http://127.0.0.1:${BACKEND_PORT}/api`;

// Stores
export const profiles: Writable<Profile[]> = writable();
export const traefik: Writable<TraefikConfig> = writable();
export const entrypoints: Writable<EntryPoints[]> = writable([]);
export const routers: Writable<Router[]> = writable([]);
export const services: Writable<Service[]> = writable([]);
export const middlewares: Writable<Middleware[]> = writable([]);
export const users: Writable<User[]> = writable([]);
export const dnsProviders: Writable<DNSProvider[]> = writable([]);
export const agents: Writable<Agent[]> = writable([]);
export const settings: Writable<Setting[]> = writable([]);
export const plugins: Writable<Plugin[]> = writable([]);

// App state
export const profile: Writable<Profile> = writable({} as Profile);
export const user: Writable<User | null> = writable({} as User);

// Loading and error states
export const loading = writable<boolean>(false);
export const error = writable<string | null>(null);

interface APIOptions {
	method?: string;
	// eslint-disable-next-line
	body?: any;
	headers?: Record<string, string>;
}

async function send(endpoint: string, options: APIOptions = {}) {
	const token = localStorage.getItem(TOKEN_SK);
	const headers = {
		'Content-Type': 'application/json',
		...(token && { Authorization: `Bearer ${token}` }),
		...options.headers
	};
	try {
		loading.set(true);
		const response = await fetch(`${BASE_URL}${endpoint}`, {
			method: options.method || 'GET',
			body: options.body ? JSON.stringify(options.body) : undefined,
			headers
		});

		if (!response.ok) {
			throw new Error(`${await response.text()}`);
		}
		if (response.status !== 204) {
			return await response.json();
		}
	} catch (e) {
		error.set(e.message);
		throw e;
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
		}
		return data;
	},

	async verify() {
		const token = localStorage.getItem(TOKEN_SK);
		try {
			const data = await send('/verify', {
				method: 'POST',
				body: token
			});
			return data;
		} catch (_) {
			localStorage.removeItem(TOKEN_SK);
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
	},

	// Profiles ------------------------------------------------------------------
	async listProfiles() {
		const data = await send('/profile');
		profiles.set(data);
		return data;
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
	async getTraefikConfig(id: number) {
		const res = await send(`/traefik/${id}`);
		traefik.set(res);
		entrypoints.set(res.entrypoints);
		routers.set(flattenRouterData(res.config));
		services.set(flattenServiceData(res.config));
		middlewares.set(flattenMiddlewareData(res.config));
	},

	async upsertRouter(id: number, data: unknown) {
		return await send(`/router/${id}`, {
			method: 'POST',
			body: data
		});
	},

	async deleteRouter(id: number, name: string, type: string) {
		return await send(`/router/${id}/${name}/${type}`, {
			method: 'DELETE'
		});
	},

	async upsertMiddleware(id: number, data: unknown) {
		return await send(`/middleware/${id}`, {
			method: 'POST',
			body: data
		});
	},

	async deleteMiddleware(id: number, name: string, type: string) {
		return await send(`/middleware/${id}/${name}/${type}`, {
			method: 'DELETE'
		});
	},

	// DNS Providers -------------------------------------------------------------
	async listDNSProviders() {
		const data = await send('/provider');
		dnsProviders.set(data);
	},

	async getDNSProvider(id: number) {
		const data = await send(`/provider/${id}`);
		return data;
	},

	async createDNSProvider(provider: Omit<DNSProvider, 'id' | 'created_at' | 'updated_at'>) {
		const data = await send('/provider', {
			method: 'POST',
			body: provider
		});
		await api.listDNSProviders();
		return data;
	},

	async updateDNSProvider(
		id: number,
		provider: Omit<DNSProvider, 'id' | 'created_at' | 'updated_at'>
	) {
		const data = await send(`/provider/${id}`, {
			method: 'PUT',
			body: provider
		});
		await api.listDNSProviders();
		return data;
	},

	async deleteDNSProvider(id: number) {
		const data = await send(`/provider/${id}`, {
			method: 'DELETE'
		});
		await api.listDNSProviders();
		return data;
	},

	// Users ---------------------------------------------------------------------
	async listUsers() {
		const data = await send('/user');
		users.set(data);
		return data;
	},

	async getUser(id: number) {
		const data = await send(`/user/${id}`);
		return data;
	},

	async createUser(user: Omit<User, 'id' | 'created_at' | 'updated_at'>) {
		const data = await send('/user', {
			method: 'POST',
			body: user
		});
		await api.listUsers();
		return data;
	},

	async updateUser(id: number, user: Omit<User, 'id' | 'created_at' | 'updated_at'>) {
		const data = await send(`/user/${id}`, {
			method: 'PUT',
			body: user
		});
		await api.listUsers();
		return data;
	},

	async deleteUser(id: number) {
		const data = await send(`/user/${id}`, {
			method: 'DELETE'
		});
		await api.listUsers();
		return data;
	},

	// Agents
	async listAgents() {
		const data = await send('/agent');
		agents.set(data);
		return data;
	},

	// Settings ------------------------------------------------------------------
	async listSettings() {
		const data = await send('/settings');
		settings.set(data);
		return data;
	},

	async getSetting(id: number) {
		const data = await send(`/settings/${id}`);
		return data;
	},

	async upsertSetting(id: number, setting: Omit<Setting, 'id' | 'created_at' | 'updated_at'>) {
		const data = await send(`/settings/${id}`, {
			method: 'POST',
			body: setting
		});
		await api.listSettings();
		return data;
	},

	// Plugins
	async getMiddlewarePlugins() {
		const data = await send('/middleware/plugins');
		plugins.set(data);
		return data;
	},

	async getVersion() {
		const data = await send('/version');
		return data;
	}
};

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

// export async function logout() {
// 	localStorage.removeItem(TOKEN_SK);
// 	loggedIn.set(false);
// }

// // Profiles -------------------------------------------------------------------
// export async function getProfiles() {
// 	const response = await handleRequest('/profile', 'GET');
// 	if (response) {
// 		const data = await response.json();
// 		if (data) {
// 			profiles.set(data);

// 			// Get saved profile
// 			const profileID = parseInt(localStorage.getItem(PROFILE_SK) ?? '');
// 			if (profileID) {
// 				getProfile(profileID);
// 				return;
// 			}
// 			if (data === undefined) return;
// 			if (!get(profile) && data.length > 0) {
// 				getProfile(data[0].id);
// 			}
// 		}
// 	}
// }

// export async function getProfile(id: number) {
// 	const response = await handleRequest(`/profile/${id}`, 'GET');
// 	if (response) {
// 		const data = await response.json();
// 		profile.set(data);
// 		localStorage.setItem(PROFILE_SK, data.id.toString());
// 	} else {
// 		localStorage.removeItem(PROFILE_SK);
// 		return;
// 	}
// 	await getEntrypoints();
// }

// export async function upsertProfile(p: Profile): Promise<void> {
// 	const response = await handleRequest(`/profile`, 'PUT', p);
// 	if (response) {
// 		const data = await response.json();
// 		if (data && get(profiles)) {
// 			profiles.update((items) => {
// 				const index = items.findIndex((item) => item.id === p.id);
// 				if (index !== -1) {
// 					items[index] = data;
// 					return [...items];
// 				} else {
// 					return [...items, data];
// 				}
// 			});
// 		}
// 		toast.success(`Profile ${data.name} updated`);

// 		if (get(profile) && get(profile).id === data.id) {
// 			profile.set(data);
// 			localStorage.setItem(PROFILE_SK, data.id.toString());
// 		}
// 	}
// }

// export async function deleteProfile(p: Profile): Promise<void> {
// 	const response = await handleRequest(`/profile/${p.id}`, 'DELETE', p);
// 	if (response) {
// 		profiles.update((items) => items.filter((i) => i.id !== p.id));
// 		toast.success(`Profile deleted`);

// 		if (get(profile).id === p.id) {
// 			profile.set({} as Profile);
// 			localStorage.removeItem(PROFILE_SK);
// 		}
// 	}
// }

// // Routers ----------------------------------------------------------------------
// export async function getRouters() {
// 	const profileID = get(profile)?.id;
// 	if (!profileID) return;
// 	const response = await handleRequest(`/router/${profileID}`, 'GET');
// 	if (response) {
// 		const data = await response.json();
// 		if (data) routers.set(data);
// 	}
// }

// export async function upsertRouter(r: Router): Promise<void> {
// 	const response = await handleRequest(`/router`, 'POST', r);
// 	if (response) {
// 		const data = await response.json();
// 		if (data && get(routers)) {
// 			routers.update((items) => {
// 				const index = items.findIndex((item) => item.id === r.id);
// 				if (index !== -1) {
// 					items[index] = data;
// 					return [...items];
// 				} else {
// 					return [...items, data];
// 				}
// 			});
// 			toast.success(`Router ${r.name} updated`);
// 		}
// 	}
// }

// export async function deleteRouter(r: Router): Promise<void> {
// 	const response = await handleRequest(`/router/${r.id}`, 'DELETE');
// 	if (response) {
// 		routers.update((items) => items.filter((i) => i.id !== r.id));
// 		toast.success(`Router ${r.name} deleted`);
// 	}
// }

// // Services ----------------------------------------------------------------------
// export async function getServices() {
// 	const profileID = get(profile)?.id;
// 	if (!profileID) return;
// 	const response = await handleRequest(`/service/${profileID}`, 'GET');
// 	if (response) {
// 		const data = await response.json();
// 		if (data) services.set(data);
// 	}
// }

// export async function upsertService(s: Service): Promise<void> {
// 	const response = await handleRequest(`/service`, 'POST', s);
// 	if (response) {
// 		const data = await response.json();
// 		if (data && get(services)) {
// 			services.update((items) => {
// 				const index = items.findIndex((item) => item.id === s.id);
// 				if (index !== -1) {
// 					items[index] = data;
// 					return [...items];
// 				} else {
// 					return [...items, data];
// 				}
// 			});
// 			//toast.success(`Service ${data.name} created`);
// 		}
// 	}
// }

// export async function deleteService(s: Service): Promise<void> {
// 	const response = await handleRequest(`/service/${s.id}`, 'DELETE');
// 	if (response) {
// 		services.update((items) => items.filter((i) => i.id !== s.id));
// 	}
// }

// // Middleware -----------------------------------------------------------------
// export async function getMiddlewares() {
// 	const profileID = get(profile)?.id;
// 	if (!profileID) return;
// 	const response = await handleRequest(`/middleware/${profileID}`, 'GET');
// 	if (response) {
// 		const data = await response.json();
// 		if (data) middlewares.set(data);
// 	}
// }

// export async function upsertMiddleware(m: Middleware): Promise<void> {
// 	const response = await handleRequest(`/middleware`, 'POST', m);
// 	if (response) {
// 		const data = await response.json();
// 		if (data && get(middlewares)) {
// 			middlewares.update((items) => {
// 				const index = items.findIndex((item) => item.id === m.id);
// 				if (index !== -1) {
// 					items[index] = data;
// 					return [...items];
// 				} else {
// 					return [...items, data];
// 				}
// 			});
// 			toast.success(`Middleware ${data.name} updated`);
// 		}
// 	}
// }

// export async function deleteMiddleware(m: Middleware): Promise<void> {
// 	const response = await handleRequest(`/middleware/${m.id}`, 'DELETE');
// 	if (response) {
// 		middlewares.update((items) => items.filter((i) => i.id !== m.id));
// 		toast.success(`Middleware ${m.name} deleted`);
// 	}
// }

// // Users ----------------------------------------------------------------------
// export async function getUsers() {
// 	const response = await handleRequest('/user', 'GET');
// 	if (response) {
// 		const data = await response.json();
// 		if (data) users.set(data);
// 	}
// }

// export async function upsertUser(u: User): Promise<void> {
// 	const response = await handleRequest(`/user`, 'POST', u);
// 	if (response) {
// 		const data = await response.json();
// 		if (data && get(users)) {
// 			users.update((items) => {
// 				const index = items.findIndex((item) => item.id === u.id);
// 				if (index !== -1) {
// 					items[index] = data;
// 					return [...items];
// 				} else {
// 					return [...items, data];
// 				}
// 			});
// 			toast.success(`User ${data.username} updated`);
// 		}
// 	}
// }

// export async function deleteUser(id: number): Promise<void> {
// 	const response = await handleRequest(`/user/${id}`, 'DELETE');
// 	if (response) {
// 		users.update((items) => items.filter((i) => i.id !== id));
// 		toast.success(`User deleted`);
// 	}
// }

// // Provider -------------------------------------------------------------------
// export async function getProviders() {
// 	const response = await handleRequest('/provider', 'GET');
// 	if (response) {
// 		const data = await response.json();
// 		if (data) provider.set(data);
// 	}
// }

// export async function upsertProvider(p: DNSProvider): Promise<void> {
// 	const response = await handleRequest(`/provider`, 'POST', p);
// 	if (response) {
// 		const data = await response.json();
// 		if (data && get(provider)) {
// 			provider.update((items) => {
// 				const index = items.findIndex((item) => item.id === p.id);
// 				if (index !== -1) {
// 					items[index] = data;
// 					return [...items];
// 				} else {
// 					return [...items, data];
// 				}
// 			});
// 			toast.success(`Provider ${data.name} updated`);
// 		}
// 	}
// }

// export async function deleteProvider(id: number): Promise<void> {
// 	const response = await handleRequest(`/provider/${id}`, 'DELETE');
// 	if (response) {
// 		provider.update((items) => items.filter((i) => i.id !== id));
// 		toast.success(`Provider deleted`);
// 	}
// }

// // Entrypoints ----------------------------------------------------------------
// export async function getEntrypoints() {
// 	const profileID = get(profile)?.id;
// 	if (!profileID) return;
// 	const response = await handleRequest(`/entrypoint/${profileID}`, 'GET');
// 	if (response) {
// 		const data = await response.json();
// 		entrypoints.set(data);
// 	}
// }

// export async function deleteRouterDNS(r: Router): Promise<void> {
// 	if (!r || !r.dnsProvider) return;

// 	const response = await handleRequest(`/dns`, 'POST', r);
// 	if (response) {
// 		toast.success(`DNS record of router ${r.name} deleted`);
// 	}
// }

// // Settings -------------------------------------------------------------------
// export async function getSettings(): Promise<void> {
// 	const response = await handleRequest('/settings', 'GET');
// 	if (response) {
// 		const data = await response.json();
// 		settings.set(data);
// 	}
// }

// export async function getSetting(key: string) {
// 	const response = await handleRequest(`/settings/${key}`, 'GET');
// 	if (response) {
// 		const data = await response.json();
// 		return data;
// 	}
// 	return {} as Setting;
// }

// export async function updateSetting(s: Setting): Promise<void> {
// 	const response = await handleRequest(`/settings`, 'PUT', s);
// 	if (response) {
// 		const data = await response.json();
// 		settings.update((items) => items.map((i) => (i.key === s.key ? data : i)));
// 		toast.success(`Setting ${s.key} updated`);
// 	}
// }

// // Agents ---------------------------------------------------------------------
// export async function getAgents(id: number) {
// 	const response = await handleRequest(`/agent/${id}`, 'GET');
// 	if (response) {
// 		const data = await response.json();
// 		if (data) agents.set(data);
// 	}
// }

// export async function upsertAgent(agent: Agent) {
// 	const response = await handleRequest(`/agent`, 'PUT', agent);
// 	if (response) {
// 		const data = await response.json();
// 		if (data && get(agents)) {
// 			agents.update((items) => {
// 				const index = items.findIndex((item) => item.id === agent.id);
// 				if (index !== -1) {
// 					items[index] = data;
// 					return [...items];
// 				} else {
// 					return [...items, data];
// 				}
// 			});
// 		}
// 	}
// }

// export async function deleteAgent(id: string) {
// 	const response = await handleRequest(`/agent/${id}`, 'DELETE');
// 	if (response) {
// 		agents.update((items) => items.filter((i) => i.id !== id));
// 		toast.success(`Agent deleted`);
// 	}
// }

// // Backup ---------------------------------------------------------------------
// export async function downloadBackup() {
// 	const response = await handleRequest('/backup', 'GET');
// 	if (response) {
// 		const data = await response.json();
// 		const jsonString = JSON.stringify(data, null, 2);
// 		const blob = new Blob([jsonString], { type: 'application/json' });
// 		const url = URL.createObjectURL(blob);
// 		const link = document.createElement('a');
// 		link.href = url;
// 		link.download = `backup-${new Date().toISOString().split('T')[0]}.json`;
// 		document.body.appendChild(link);
// 		link.click();
// 		URL.revokeObjectURL(url);
// 		document.body.removeChild(link);
// 	}
// }

// export async function uploadBackup(file: File) {
// 	const formData = new FormData();
// 	formData.append('file', file);
// 	const token = localStorage.getItem(TOKEN_SK);
// 	await fetch(`${BASE_URL}/restore`, {
// 		method: 'POST',
// 		body: formData,
// 		headers: { Authorization: `Bearer ${token}` }
// 	});
// 	toast.success('Backup restored!');
// 	await getProfiles();
// 	await getUsers();
// 	await getProviders();
// 	await getEntrypoints();
// 	await getSettings();
// }

// // Plugins --------------------------------------------------------------------
// export async function getPlugins() {
// 	const response = await handleRequest('/middleware/plugins', 'GET');
// 	if (response) {
// 		const data = await response.json();
// 		plugins.set(data);
// 	}
// }

// // Extras ---------------------------------------------------------------------
// export async function getVersion() {
// 	const response = await handleRequest('/version', 'GET');
// 	if (response) {
// 		const data = await response.text();
// 		version.set(data);
// 	}
// }

// export async function getPublicIP() {
// 	const response = await handleRequest(`/ip/${get(profile).id}`, 'GET');
// 	if (response) {
// 		const data = await response.json();
// 		return data.ip;
// 	}
// 	return '';
// }

// export async function getTraefikOverview() {
// 	const profileID = get(profile)?.id;
// 	if (!profileID) return;
// 	const response = await handleRequest(`/traefik/${profileID}`, 'GET');
// 	if (response) {
// 		const data = await response.json();
// 		traefikError.set('');
// 		return data;
// 	} else {
// 		traefikError.set('No connection to Traefik');
// 	}
// 	return '';
// }

// export async function getTraefikConfig() {
// 	const profileName = get(profile)?.name;
// 	if (!profileName) return;

// 	const response = await handleRequest(`/${profileName}?yaml=true`, 'GET');
// 	if (response) {
// 		const data = await response.text();
// 		if (!data.includes('{}')) {
// 			dynamic.set(data);
// 		}
// 	} else {
// 		dynamic.set('');
// 	}
// }

// // Toggle functions -----------------------------------------------------------
// export async function toggleEntrypoint(
// 	router: Router,
// 	item: Selected<unknown>[] | undefined,
// 	update: boolean
// ) {
// 	if (item === undefined) return;
// 	router.entryPoints = item.map((i) => i.value) as string[];

// 	if (update) {
// 		upsertRouter(router);
// 	}
// }

// export async function toggleMiddleware(
// 	router: Router,
// 	item: Selected<unknown>[] | undefined,
// 	update: boolean
// ) {
// 	if (item === undefined) return;
// 	router.middlewares = item.map((i) => i.value) as string[];

// 	if (update) {
// 		upsertRouter(router);
// 	}
// }

// export async function toggleDNSProvider(
// 	router: Router,
// 	item: Selected<unknown> | undefined,
// 	update: boolean
// ) {
// 	const providerID = (item?.value as number) ?? 0;
// 	if (providerID === 0 && router.dnsProvider !== 0) {
// 		deleteRouterDNS(router);
// 	}
// 	router.dnsProvider = providerID;

// 	if (update) {
// 		upsertRouter(router);
// 	}
// }
