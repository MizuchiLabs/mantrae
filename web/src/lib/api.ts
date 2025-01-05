import { goto } from '$app/navigation';
import { toast } from 'svelte-sonner';
import type { Profile, DNSProvider, User, Setting, Agent, Entrypoint } from './types/base';
import type { Plugin } from './types/plugins';
import type { Middleware } from './types/middlewares';
import { get, writable, type Writable } from 'svelte/store';
import { type Router, type Service } from './types/config';
import type { Selected } from 'bits-ui';
import { PROFILE_SK, TOKEN_SK } from './store';

// Global state variables
export const BACKEND_PORT = import.meta.env.PORT || 3000;
export const API_URL = import.meta.env.PROD ? '/api' : `http://127.0.0.1:${BACKEND_PORT}/api`;
export const loggedIn = writable(false);
export const profiles: Writable<Profile[]> = writable();
export const profile: Writable<Profile> = writable();
export const entrypoints: Writable<Entrypoint[]> = writable();

export const routers: Writable<Router[]> = writable();
export const services: Writable<Service[]> = writable();
export const middlewares: Writable<Middleware[]> = writable();
export const users: Writable<User[]> = writable();
export const agents: Writable<Agent[]> = writable();
export const provider: Writable<DNSProvider[]> = writable();
export const settings: Writable<Setting[]> = writable();
export const plugins: Writable<Plugin[]> = writable();

export const dynamic = writable('');
export const version = writable('');
export const configError = writable('');
export const traefikError = writable('');
export const agentToken = writable('');

async function handleRequest(
	endpoint: string,
	method: string,
	body?: object
): Promise<Response | undefined> {
	if (!get(loggedIn)) return;

	const token = localStorage.getItem(TOKEN_SK);
	const response = await fetch(`${API_URL}${endpoint}`, {
		method: method,
		body: body ? JSON.stringify(body) : undefined,
		headers: { Authorization: `Bearer ${token}` }
	});
	if (response.ok) {
		return response;
	} else {
		toast.error('Request failed', {
			description: await response.text(),
			duration: 3000
		});
	}
}

// Login ----------------------------------------------------------------------
export async function login(username: string, password: string, remember: boolean) {
	const loginURL = remember ? `${API_URL}/login?remember=true` : `${API_URL}/login`;
	const response = await fetch(loginURL, {
		method: 'POST',
		body: JSON.stringify({ username, password })
	});
	if (response.ok) {
		const { token } = await response.json();
		localStorage.setItem(TOKEN_SK, token);
		loggedIn.set(true);
		goto('/');
		toast.success('Login successful');
		await getProfiles();
		await getProviders();
	} else {
		toast.error('Login failed', {
			description: await response.text(),
			duration: 3000
		});
		return;
	}
}

export async function sendResetEmail(username: string) {
	const response = await fetch(`${API_URL}/reset/${username}`, {
		method: 'POST'
	});
	if (response.ok) {
		toast.success('Password reset email sent!');
	} else {
		toast.error('Request failed', {
			description: await response.text(),
			duration: 3000
		});
	}
}

export async function resetPassword(token: string, password: string) {
	const response = await fetch(`${API_URL}/reset`, {
		method: 'POST',
		body: JSON.stringify({ password }),
		headers: { Authorization: `Bearer ${token}` }
	});
	if (response.ok) {
		toast.success('Password reset successful!');
		goto('/login');
	} else {
		toast.error('Request failed', {
			description: await response.text(),
			duration: 3000
		});
	}
}

export async function logout() {
	localStorage.removeItem(TOKEN_SK);
	loggedIn.set(false);
}

// Profiles -------------------------------------------------------------------
export async function getProfiles() {
	const response = await handleRequest('/profile', 'GET');
	if (response) {
		const data = await response.json();
		if (data) {
			profiles.set(data);

			// Get saved profile
			const profileID = parseInt(localStorage.getItem(PROFILE_SK) ?? '');
			if (profileID) {
				getProfile(profileID);
				return;
			}
			if (data === undefined) return;
			if (!get(profile) && data.length > 0) {
				getProfile(data[0].id);
			}
		}
	}
}

export async function getProfile(id: number) {
	const response = await handleRequest(`/profile/${id}`, 'GET');
	if (response) {
		const data = await response.json();
		profile.set(data);
		localStorage.setItem(PROFILE_SK, data.id.toString());
	} else {
		localStorage.removeItem(PROFILE_SK);
		return;
	}
	await getEntrypoints();
}

export async function createProfile(p: Profile): Promise<void> {
	const response = await handleRequest('/profile', 'POST', p);
	if (response) {
		const data = await response.json();
		profiles.update((items) => [...(items ?? []), data]);
		toast.success(`Profile ${data.name} created`);

		const profileID = parseInt(localStorage.getItem(PROFILE_SK) ?? '');
		if (!profileID) {
			localStorage.setItem(PROFILE_SK, data.id.toString());
			profile.set(data);
		}
	}
}

export async function updateProfile(p: Profile): Promise<void> {
	const response = await handleRequest(`/profile`, 'PUT', p);
	if (response) {
		const data = await response.json();
		profiles.update((items) => items.map((i) => (i.id === p.id ? data : i)));
		toast.success(`Profile ${data.name} updated`);

		if (get(profile) && get(profile).id === data.id) {
			profile.set(data);
			localStorage.setItem(PROFILE_SK, data.id.toString());
		}
	}
}

export async function deleteProfile(p: Profile): Promise<void> {
	const response = await handleRequest(`/profile/${p.id}`, 'DELETE', p);
	if (response) {
		profiles.update((items) => items.filter((i) => i.id !== p.id));
		toast.success(`Profile deleted`);

		if (get(profile).id === p.id) {
			profile.set({} as Profile);
			localStorage.removeItem(PROFILE_SK);
		}
	}
}

// Routers ----------------------------------------------------------------------
export async function getRouters() {
	const profileID = get(profile)?.id;
	if (!profileID) return;
	const response = await handleRequest(`/router/${profileID}`, 'GET');
	if (response) {
		const data = await response.json();
		routers.set(data);
	}
}

export async function upsertRouter(r: Router): Promise<void> {
	const response = await handleRequest(`/router`, 'POST', r);
	if (response) {
		const data = await response.json();
		if (data && get(routers)) {
			routers.update((items) => {
				const existingIndex = items.findIndex((item) => item.id === r.id);

				if (existingIndex !== -1) {
					// Update existing item
					const updatedItems = [...items];
					updatedItems[existingIndex] = data;
					return updatedItems;
				} else {
					// Add new item
					return [...items, data];
				}
			});
			toast.success(`Router ${r.name} updated`);
		}
	}
}

export async function deleteRouter(r: Router): Promise<void> {
	const response = await handleRequest(`/router/${r.id}`, 'DELETE');
	if (response) {
		routers.update((items) => items.filter((i) => i.id !== r.id));
		toast.success(`Router ${r.name} deleted`);
	}
}

// Services ----------------------------------------------------------------------
export async function getServices() {
	const profileID = get(profile)?.id;
	if (!profileID) return;
	const response = await handleRequest(`/service/${profileID}`, 'GET');
	if (response) {
		const data = await response.json();
		services.set(data);
	}
}

export async function upsertService(s: Service): Promise<void> {
	const response = await handleRequest(`/service`, 'POST', s);
	if (response) {
		const data = await response.json();
		if (data && get(services)) {
			services.update((items) => {
				const existingIndex = items.findIndex((item) => item.id === s.id);

				if (existingIndex !== -1) {
					// Update existing item
					const updatedItems = [...items];
					updatedItems[existingIndex] = data;
					return updatedItems;
				} else {
					// Add new item
					return [...items, data];
				}
			});
			//toast.success(`Service ${data.name} created`);
		}
	}
}

export async function deleteService(s: Service): Promise<void> {
	const response = await handleRequest(`/service/${s.id}`, 'DELETE');
	if (response) {
		services.update((items) => items.filter((i) => i.id !== s.id));
	}
}

// Middleware -----------------------------------------------------------------
export async function getMiddlewares() {
	const profileID = get(profile)?.id;
	if (!profileID) return;
	const response = await handleRequest(`/middleware/${profileID}`, 'GET');
	if (response) {
		const data = await response.json();
		middlewares.set(data);
	}
}

export async function upsertMiddleware(m: Middleware): Promise<void> {
	const response = await handleRequest(`/middleware`, 'POST', m);
	if (response) {
		const data = await response.json();
		if (data && get(middlewares)) {
			middlewares.update((items) => {
				const existingIndex = items.findIndex((item) => item.id === m.id);

				if (existingIndex !== -1) {
					// Update existing item
					const updatedItems = [...items];
					updatedItems[existingIndex] = data;
					return updatedItems;
				} else {
					// Add new item
					return [...items, data];
				}
			});
			toast.success(`Middleware ${data.name} updated`);
		}
	}
}

export async function deleteMiddleware(m: Middleware): Promise<void> {
	const response = await handleRequest(`/middleware/${m.id}`, 'DELETE');
	if (response) {
		middlewares.update((items) => items.filter((i) => i.id !== m.id));
		toast.success(`Middleware ${m.name} deleted`);
	}
}

// Users ----------------------------------------------------------------------
export async function getUsers() {
	const response = await handleRequest('/user', 'GET');
	if (response) {
		const data = await response.json();
		users.set(data);
	}
}

export async function getUser(id: number): Promise<User> {
	const response = await handleRequest(`/user/${id}`, 'GET');
	if (response) {
		const data = await response.json();
		return data;
	}
	return {} as User;
}

export async function upsertUser(u: User): Promise<void> {
	const response = await handleRequest(`/user`, 'POST', u);
	if (response) {
		const data = await response.json();
		if (data && get(users)) {
			users.update((items) => {
				const existingIndex = items.findIndex((item) => item.id === u.id);

				if (existingIndex !== -1) {
					// Update existing item
					const updatedItems = [...items];
					updatedItems[existingIndex] = data;
					return updatedItems;
				} else {
					// Add new item
					return [...items, data];
				}
			});
			toast.success(`User ${data.username} updated`);
		}
	}
}

export async function deleteUser(id: number): Promise<void> {
	const response = await handleRequest(`/user/${id}`, 'DELETE');
	if (response) {
		users.update((items) => items.filter((i) => i.id !== id));
		toast.success(`User deleted`);
	}
}

// Provider -------------------------------------------------------------------
export async function getProviders() {
	const response = await handleRequest('/provider', 'GET');
	if (response) {
		const data = await response.json();
		provider.set(data);
	}
}

export async function upsertProvider(p: DNSProvider): Promise<void> {
	const response = await handleRequest(`/provider`, 'POST', p);
	if (response) {
		const data = await response.json();
		if (data && get(provider)) {
			provider.update((items) => {
				const existingIndex = items.findIndex((item) => item.id === p.id);

				if (existingIndex !== -1) {
					// Update existing item
					const updatedItems = [...items];
					updatedItems[existingIndex] = data;
					return updatedItems;
				} else {
					// Add new item
					return [...items, data];
				}
			});
			toast.success(`Provider ${data.name} updated`);
		}
	}
}

export async function deleteProvider(id: number): Promise<void> {
	const response = await handleRequest(`/provider/${id}`, 'DELETE');
	if (response) {
		provider.update((items) => items.filter((i) => i.id !== id));
		toast.success(`Provider deleted`);
	}
}

// Entrypoints ----------------------------------------------------------------
export async function getEntrypoints() {
	const profileID = get(profile)?.id;
	if (!profileID) return;
	const response = await handleRequest(`/entrypoint/${profileID}`, 'GET');
	if (response) {
		const data = await response.json();
		entrypoints.set(data);
	}
}

export async function deleteRouterDNS(r: Router): Promise<void> {
	if (!r || !r.dnsProvider) return;

	const response = await handleRequest(`/dns`, 'POST', r);
	if (response) {
		toast.success(`DNS record of router ${r.name} deleted`);
	}
}

// Settings -------------------------------------------------------------------
export async function getSettings(): Promise<void> {
	const response = await handleRequest('/settings', 'GET');
	if (response) {
		const data = await response.json();
		settings.set(data);
	}
}

export async function getSetting(key: string) {
	const response = await handleRequest(`/settings/${key}`, 'GET');
	if (response) {
		const data = await response.json();
		return data;
	}
	return {} as Setting;
}

export async function updateSetting(s: Setting): Promise<void> {
	const response = await handleRequest(`/settings`, 'PUT', s);
	if (response) {
		const data = await response.json();
		settings.update((items) => items.map((i) => (i.key === s.key ? data : i)));
		toast.success(`Setting ${s.key} updated`);
	}
}

// Agents ---------------------------------------------------------------------
export async function getAgents() {
	const profileID = get(profile)?.id;
	if (!profileID) return;
	const response = await handleRequest(`/agent/${profileID}`, 'GET');
	if (response) {
		const data = await response.json();
		agents.set(data);
	}
}

export async function upsertAgent(a: Agent) {
	const response = await handleRequest(`/agent`, 'PUT', a);
	if (response) {
		const data = await response.json();
		if (data && get(agents)) {
			agents.update((items) => {
				const existingIndex = items.findIndex((item) => item.id === a.id);

				if (existingIndex !== -1) {
					// Update existing item
					const updatedItems = [...items];
					updatedItems[existingIndex] = data;
					return updatedItems;
				} else {
					// Add new item
					return [...items, data];
				}
			});
		}
	}
}

export async function deleteAgent(id: string) {
	const response = await handleRequest(`/agent/${id}`, 'DELETE');
	if (response) {
		agents.update((items) => items.filter((i) => i.id !== id));
		toast.success(`Agent deleted`);
	}
}

// Agent Token derived store
export async function getAgentToken() {
	const profileID = get(profile)?.id;
	if (!profileID) return;
	const response = await handleRequest(`/agent/token/${profileID}`, 'GET');
	if (response) {
		const data = await response.json();
		agentToken.set(data.token || '');
	}
}

// Backup ---------------------------------------------------------------------
export async function downloadBackup() {
	const response = await handleRequest('/backup', 'GET');
	if (response) {
		const data = await response.json();
		const jsonString = JSON.stringify(data, null, 2);
		const blob = new Blob([jsonString], { type: 'application/json' });
		const url = URL.createObjectURL(blob);
		const link = document.createElement('a');
		link.href = url;
		link.download = `backup-${new Date().toISOString().split('T')[0]}.json`;
		document.body.appendChild(link);
		link.click();
		URL.revokeObjectURL(url);
		document.body.removeChild(link);
	}
}

export async function uploadBackup(file: File) {
	const formData = new FormData();
	formData.append('file', file);
	const token = localStorage.getItem(TOKEN_SK);
	await fetch(`${API_URL}/restore`, {
		method: 'POST',
		body: formData,
		headers: { Authorization: `Bearer ${token}` }
	});
	toast.success('Backup restored!');
	await getProfiles();
	await getUsers();
	await getProviders();
	await getEntrypoints();
	await getSettings();
}

// Plugins --------------------------------------------------------------------
export async function getPlugins() {
	const response = await handleRequest('/middleware/plugins', 'GET');
	if (response) {
		const data = await response.json();
		plugins.set(data);
	}
}

// Extras ---------------------------------------------------------------------
export async function getVersion() {
	const response = await handleRequest('/version', 'GET');
	if (response) {
		const data = await response.text();
		version.set(data);
	}
}

export async function getPublicIP() {
	const response = await handleRequest(`/ip/${get(profile).id}`, 'GET');
	if (response) {
		const data = await response.json();
		return data.ip;
	}
	return '';
}

export async function getTraefikOverview() {
	const profileID = get(profile)?.id;
	if (!profileID) return;
	const response = await handleRequest(`/traefik/${profileID}`, 'GET');
	if (response) {
		const data = await response.json();
		traefikError.set('');
		return data;
	} else {
		traefikError.set('No connection to Traefik');
	}
	return '';
}

export async function getTraefikConfig() {
	const profileName = get(profile)?.name;
	if (!profileName) return;

	const response = await handleRequest(`/${profileName}?yaml=true`, 'GET');
	if (response) {
		const data = await response.text();
		if (!data.includes('{}')) {
			dynamic.set(data);
		}
	} else {
		dynamic.set('');
	}
}

// Toggle functions -----------------------------------------------------------
export async function toggleEntrypoint(
	router: Router,
	item: Selected<unknown>[] | undefined,
	update: boolean
) {
	if (item === undefined) return;
	router.entryPoints = item.map((i) => i.value) as string[];

	if (update) {
		upsertRouter(router);
	}
}

export async function toggleMiddleware(
	router: Router,
	item: Selected<unknown>[] | undefined,
	update: boolean
) {
	if (item === undefined) return;
	router.middlewares = item.map((i) => i.value) as string[];

	if (update) {
		upsertRouter(router);
	}
}

export async function toggleDNSProvider(
	router: Router,
	item: Selected<unknown> | undefined,
	update: boolean
) {
	const providerID = (item?.value as number) ?? 0;
	if (providerID === 0 && router.dnsProvider !== 0) {
		deleteRouterDNS(router);
	}
	router.dnsProvider = providerID;

	if (update) {
		upsertRouter(router);
	}
}
