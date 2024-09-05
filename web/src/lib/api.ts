import { goto } from '$app/navigation';
import { toast } from 'svelte-sonner';
import type { Config, Profile, DNSProvider, User } from './types/base';
import { newMiddleware, type Middleware } from './types/middlewares';
import { derived, get, writable, type Writable } from 'svelte/store';
import { newRouter, newService, type Router, type Service } from './types/config';
import type { Selected } from 'bits-ui';

// Global state variables
export const loggedIn = writable(false);
export const profiles: Writable<Profile[]> = writable();
export const profile: Writable<Profile> = writable();
export const config: Writable<Config> = writable();
export const users: Writable<User[]> = writable();
export const provider: Writable<DNSProvider[]> = writable();
export const API_URL = import.meta.env.PROD ? '/api' : 'http://localhost:3000/api';

// Derived stores
export const routers = derived(config, ($config) => Object.values($config?.routers ?? []));
export const services = derived(config, ($config) => Object.values($config?.services ?? []));
export const middlewares = derived(config, ($config) => Object.values($config?.middlewares ?? []));
export const version = derived(config, ($config) => $config?.version ?? '');
export const entrypoints = derived(config, ($config) => $config?.entrypoints ?? []);

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
		let data = await response.json();
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

export async function getProfile(id: number) {
	const respProfile = await handleRequest(`/profile/${id}`, 'GET');
	if (respProfile) {
		let data = await respProfile.json();
		profile.set(data);
		localStorage.setItem('profile', data.id.toString());
	} else {
		localStorage.removeItem('profile');
		return;
	}

	const respConfig = await handleRequest(`/config/${id}`, 'GET');
	if (respConfig) {
		let data = await respConfig.json();
		config.set(data);
	}
}

export async function createProfile(p: Profile): Promise<void> {
	const response = await handleRequest('/profile', 'POST', p);
	if (response) {
		let data = await response.json();
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
		let data = await response.json();
		profile.set(data);
		profiles.update((items) => items.map((i) => (i.id === p.id ? data : i)));
		toast.success(`Profile ${data.name} updated`);

		if (get(profile).id === data.id) {
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
		let data = await response.json();
		users.set(data);
	}
}

export async function getUser(id: number): Promise<User> {
	const response = await handleRequest(`/user/${id}`, 'GET');
	if (response) {
		let data = await response.json();
		return data;
	}
	return {} as User;
}

export async function createUser(u: User): Promise<void> {
	const response = await handleRequest('/user', 'POST', u);
	if (response) {
		let data = await response.json();
		users.update((items) => [...(items ?? []), data]);
		toast.success(`User ${data.username} created`);
	}
}

export async function updateUser(u: User): Promise<void> {
	const response = await handleRequest(`/user`, 'PUT', u);
	if (response) {
		let data = await response.json();
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
		let data = await response.json();
		provider.set(data);
	}
}

export async function getProvider(id: number): Promise<DNSProvider> {
	const response = await handleRequest(`/provider/${id}`, 'GET');
	if (response) {
		let data = await response.json();
		return data;
	}
	return {} as DNSProvider;
}

export async function createProvider(p: DNSProvider): Promise<void> {
	const response = await handleRequest('/provider', 'POST', p);
	if (response) {
		let data = await response.json();
		provider.update((items) => [...(items ?? []), data]);
		toast.success(`Provider ${data.name} created`);
	}
	await getProviders();
}

export async function updateProvider(p: DNSProvider): Promise<void> {
	const response = await handleRequest(`/provider`, 'PUT', p);
	if (response) {
		let data = await response.json();
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
export async function updateConfig(c: Config): Promise<void> {
	const response = await handleRequest(`/config/${get(profile).id}`, 'PUT', c);
	if (response) {
		let data = await response.json();
		config.set(data);
		toast.success(`Config updated`);
	}
}

export async function deleteRouterDNS(r: Router): Promise<void> {
	const response = await handleRequest(`/dns`, 'POST', r);
	if (response) {
		let data = await response.json();
		config.set(data);
		toast.success(`DNS record of router ${r.name} deleted`);
	}
}

// Backup ---------------------------------------------------------------------
export async function downloadBackup() {
	const response = await handleRequest('/backup', 'GET');
	if (response?.ok) {
		const blob = await response.blob();
		const url = URL.createObjectURL(blob);
		const link = document.createElement('a');
		link.href = url;
		link.download = `backup-${new Date().toISOString().split('T')[0]}.tar.gz`;
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
}

// Helper functions -----------------------------------------------------------
// Create or update a router and its service
export async function upsertRouter(name: string, router: Router, service: Service): Promise<void> {
	let data = get(config);
	if (!data.routers) data.routers = {};
	if (!data.services) data.services = {};

	if (router.name !== name) {
		delete data.routers[name];
		delete data.services[name];
	}
	// Ensure the service name is the same as the router name
	router.service = router.name;
	service.name = router.name;

	data.routers[router.name] = router;
	data.services[router.name] = service;
	await updateConfig(data);
}

// Create or update a middleware
export async function upsertMiddleware(name: string, middleware: Middleware): Promise<void> {
	let data = get(config);
	if (!data.middlewares) data.middlewares = {};
	if (middleware.name !== name) {
		delete data.middlewares[name];
	}

	data.middlewares[middleware.name] = middleware;
	await updateConfig(data);
}

// Delete a router with its service by name
export async function deleteRouter(name: string): Promise<void> {
	let data = get(config);
	if (!data.routers || !data.services) return;
	await deleteRouterDNS(data.routers[name]);
	delete data.routers[name];
	delete data.services[name];
	await updateConfig(data);
}

// Delete a middleware by name
export async function deleteMiddleware(name: string): Promise<void> {
	let data = get(config);
	if (!data.middlewares) return;
	delete data.middlewares[name];
	await updateConfig(data);
}

// Toggle functions -----------------------------------------------------------
export async function toggleEntrypoint(router: Router, item: Selected<unknown>[] | undefined) {
	if (item === undefined) return;
	router.entrypoints = item.map((i) => i.value) as string[];
	let service = get(config)?.services?.[router.service + '@' + router.provider];
	if (service === undefined) return;
	upsertRouter(router.name, router, service);
}

export async function toggleMiddleware(router: Router, item: Selected<unknown>[] | undefined) {
	if (item === undefined) return;
	router.middlewares = item.map((i) => i.value) as string[];
	let service = get(config)?.services?.[router.service + '@' + router.provider];
	if (service === undefined) return;
	upsertRouter(router.name, router, service);
}

export async function toggleDNSProvider(router: Router, item: Selected<unknown> | undefined) {
	let newProvider = (item?.value as string) ?? '';
	if (newProvider === '' && router.dnsProvider !== '') {
		deleteRouterDNS(router);
	}
	router.dnsProvider = newProvider;
	let service = get(config)?.services?.[router.service + '@' + router.provider];
	if (service === undefined) return;
	upsertRouter(router.name, router, service);
}
