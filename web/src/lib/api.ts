import { goto } from '$app/navigation';
import { toast } from 'svelte-sonner';
import type { Config, Profile, DNSProvider, User, Setting } from './types/base';
import { type Middleware } from './types/middlewares';
import { derived, get, writable, type Writable } from 'svelte/store';
import { newService, type Router, type Service } from './types/config';
import type { Selected } from 'bits-ui';

// Global state variables
export const API_URL = import.meta.env.PROD ? '/api' : 'http://localhost:3000/api';
export const loggedIn = writable(false);
export const profiles: Writable<Profile[]> = writable();
export const profile: Writable<Profile> = writable();
export const config: Writable<Config> = writable();
export const users: Writable<User[]> = writable();
export const provider: Writable<DNSProvider[]> = writable();
export const settings: Writable<Setting[]> = writable();
export const dynamic = writable('');
export const version = writable('');

// Derived stores
export const entrypoints = derived(config, ($config) => $config?.entrypoints ?? []);
export const routers = derived(config, ($config) => Object.values($config?.routers ?? []));
export const services = derived(config, ($config) => Object.values($config?.services ?? []));
export const middlewares = derived(config, ($config) => Object.values($config?.middlewares ?? []));

async function handleRequest(
	endpoint: string,
	method: string,
	body?: any
): Promise<Response | undefined> {
	if (!get(loggedIn)) return;

	const token = localStorage.getItem('token');
	const response = await fetch(`${API_URL}${endpoint}`, {
		method: method,
		body: body ? JSON.stringify(body) : undefined,
		headers: { Authorization: `Bearer ${token}` }
	});
	if (response.status === 200) {
		return response;
	} else {
		toast.error('Request failed', {
			description: await response.text(),
			duration: 3000
		});
	}
}

// Login ----------------------------------------------------------------------
export async function login(username: string, password: string) {
	const response = await fetch(`${API_URL}/login`, {
		method: 'POST',
		body: JSON.stringify({ username, password })
	});
	if (response.status === 200) {
		const { token } = await response.json();
		localStorage.setItem('token', token);
		loggedIn.set(true);
		await getProfiles();
		goto('/');
		toast.success('Login successful');
	} else {
		toast.error('Login failed', {
			description: await response.text(),
			duration: 3000
		});
		return;
	}
}

export async function logout() {
	localStorage.removeItem('token');
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
			const profileID = parseInt(localStorage.getItem('profile') ?? '');
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
		localStorage.setItem('profile', data.id.toString());
	} else {
		localStorage.removeItem('profile');
		return;
	}
	await getConfig();
}

export async function createProfile(p: Profile): Promise<void> {
	const response = await handleRequest('/profile', 'POST', p);
	if (response) {
		const data = await response.json();
		profiles.update((items) => [...(items ?? []), data]);
		toast.success(`Profile ${data.name} created`);

		const profileID = parseInt(localStorage.getItem('profile') ?? '');
		if (!profileID) {
			localStorage.setItem('profile', data.id.toString());
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
			localStorage.setItem('profile', data.id.toString());
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
			localStorage.removeItem('profile');
		}
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

export async function createUser(u: User): Promise<void> {
	const response = await handleRequest('/user', 'POST', u);
	if (response) {
		const data = await response.json();
		users.update((items) => [...(items ?? []), data]);
		toast.success(`User ${data.username} created`);
	}
}

export async function updateUser(u: User): Promise<void> {
	const response = await handleRequest(`/user`, 'PUT', u);
	if (response) {
		const data = await response.json();
		users.update((items) => items.map((i) => (i.id === u.id ? data : i)));
		toast.success(`User ${data.username} updated`);
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

export async function getProvider(id: number): Promise<DNSProvider> {
	const response = await handleRequest(`/provider/${id}`, 'GET');
	if (response) {
		const data = await response.json();
		return data;
	}
	return {} as DNSProvider;
}

export async function createProvider(p: DNSProvider): Promise<void> {
	const response = await handleRequest('/provider', 'POST', p);
	if (response) {
		const data = await response.json();
		provider.update((items) => [...(items ?? []), data]);
		toast.success(`Provider ${data.name} created`);
	}
	await getProviders();
}

export async function updateProvider(p: DNSProvider): Promise<void> {
	const response = await handleRequest(`/provider`, 'PUT', p);
	if (response) {
		const data = await response.json();
		provider.update((items) => items.map((i) => (i.id === p.id ? data : i)));
		toast.success(`Provider ${data.name} updated`);
	}
	await getProviders();
}

export async function deleteProvider(id: number): Promise<void> {
	const response = await handleRequest(`/provider/${id}`, 'DELETE');
	if (response) {
		provider.update((items) => items.filter((i) => i.id !== id));
		toast.success(`Provider deleted`);
	}
}

// Config ---------------------------------------------------------------------
export async function getConfig() {
	const response = await handleRequest(`/config/${get(profile).id}`, 'GET');
	if (response) {
		const data = await response.json();
		config.set(data);
	}
}

export async function updateRouter(r: Router): Promise<void> {
	const response = await handleRequest(`/router/${get(profile).id}`, 'PUT', r);
	if (response) {
		const data = await response.json();
		config.set(data);
		toast.success(`Router ${r.name} updated`);
	}
}

export async function updateService(s: Service): Promise<void> {
	const response = await handleRequest(`/service/${get(profile).id}`, 'PUT', s);
	if (response) {
		const data = await response.json();
		config.set(data);
		toast.success(`Service ${s.name} updated`);
	}
}

export async function updateMiddleware(m: Middleware): Promise<void> {
	const response = await handleRequest(`/middleware/${get(profile).id}`, 'PUT', m);
	if (response) {
		const data = await response.json();
		config.set(data);
		toast.success(`Middleware ${m.name} updated`);
	}
}

export async function deleteRouter(r: Router): Promise<void> {
	const response = await handleRequest(`/router/${get(profile).id}/${r.name}`, 'DELETE');
	if (response) {
		const data = await response.json();
		config.set(data);
		toast.success(`Router ${r.name} deleted`);
	}
}

export async function deleteMiddleware(m: Middleware): Promise<void> {
	const response = await handleRequest(`/middleware/${get(profile).id}/${m.name}`, 'DELETE');
	if (response) {
		const data = await response.json();
		config.set(data);
		toast.success(`Middleware ${m.name} deleted`);
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
export async function getSettings() {
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
	const token = localStorage.getItem('token');
	await fetch(`${API_URL}/restore`, {
		method: 'POST',
		body: formData,
		headers: { Authorization: `Bearer ${token}` }
	});
	toast.success('Backup restored!');
	await getProfiles();
	await getUsers();
	await getProviders();
	await getConfig();
	await getSettings();
}

export async function getVersion() {
	const response = await handleRequest('/version', 'GET');
	if (response) {
		const data = await response.text();
		version.set(data);
	}
}

export async function getTraefikConfig() {
	if (!get(profile)) return '';

	const response = await handleRequest(`/${get(profile)?.name}`, 'GET');
	if (response) {
		const data = await response.text();
		dynamic.set(data);
	}
}

// Helper functions -----------------------------------------------------------
// Create or update a router and its service
function nameCheck(router: Router) {
	const name = router.name?.trim().toLowerCase();
	const provider = router.provider?.trim().toLowerCase() ?? 'http';
	const parts = name.split('@');
	return parts[0] + '@' + provider;
}
export async function upsertRouter(
	name: string,
	router: Router,
	service: Service | undefined
): Promise<void> {
	if (name === '' || router.name === '') return;

	// Ensure the service name is the same as the router name
	if (service === undefined) {
		service = getService(router);
	}
	router.name = nameCheck(router);
	router.service = nameCheck(router);
	service.name = nameCheck(router);
	service.serviceType = router.routerType;

	await updateRouter(router);
	await updateService(service);
}

// TODO: Handle this differently
export const getService = (router: Router): Service => {
	if (router === undefined) return newService();

	const baseName = router.service.split('@')[0];
	let service = get(config)?.services?.[baseName + '@' + router.provider];
	if (service === undefined) {
		service = get(config)?.services?.[router.service];
		if (service === undefined) return newService();
	}
	return service;
};

// Toggle functions -----------------------------------------------------------
export async function toggleEntrypoint(
	router: Router,
	item: Selected<unknown>[] | undefined,
	update: boolean
) {
	if (item === undefined) return;
	router.entrypoints = item.map((i) => i.value) as string[];

	if (update) {
		upsertRouter(router.name, router, undefined);
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
		upsertRouter(router.name, router, undefined);
	}
}

export async function toggleDNSProvider(
	router: Router,
	item: Selected<unknown> | undefined,
	update: boolean
) {
	const newProvider = (item?.value as string) ?? '';
	if (newProvider === '' && router.dnsProvider !== '') {
		deleteRouterDNS(router);
	}
	router.dnsProvider = newProvider;

	if (update) {
		upsertRouter(router.name, router, undefined);
	}
}
