import { goto } from '$app/navigation';
import { toast } from 'svelte-sonner';
import type { Config, Profile } from './types/dynamic';
import { newMiddleware, type Middleware } from './types/middlewares';
import { derived, get, writable, type Writable } from 'svelte/store';
import { newRouter, newService, type Router, type Service } from './types/config';
import type { Provider } from './types/provider';

export const loggedIn = writable(false);
export const profile: Writable<Profile> = writable();
export const config: Writable<Config> = writable();
export const profiles: Writable<Profile[]> = writable();
export const provider: Writable<Provider[]> = writable();
export const API_URL = import.meta.env.PROD ? '/api' : 'http://localhost:3000/api';

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

// Provider -------------------------------------------------------------------
export async function getProviders() {
	const response = await handleRequest('/provider', 'GET');
	if (response) {
		let data = await response.json();
		provider.set(data);
	}
}

export async function getProvider(id: number): Promise<Provider> {
	const response = await handleRequest(`/provider/${id}`, 'GET');
	if (response) {
		let data = await response.json();
		return data;
	}
	return {} as Provider;
}

export async function createProvider(p: Provider): Promise<void> {
	const response = await handleRequest('/provider', 'POST', p);
	if (response) {
		let data = await response.json();
		provider.update((items) => [...(items ?? []), data]);
		toast.success(`Provider ${data.name} created`);
	}
}

export async function updateProvider(p: Provider): Promise<void> {
	const response = await handleRequest(`/provider`, 'PUT', p);
	if (response) {
		let data = await response.json();
		provider.update((items) => items.map((i) => (i.id === p.id ? data : i)));
		toast.success(`Provider ${data.name} updated`);
	}
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

// Helper functions -----------------------------------------------------------
export function getRouter(routerName: string): Router {
	const router = get(config)?.routers?.[routerName];
	return router ?? newRouter();
}
export function getService(serviceName: string): Service {
	const service = get(config)?.services?.[serviceName];
	return service ?? newService();
}
export function getMiddleware(middlewareName: string): Middleware {
	const middleware = get(config)?.middlewares?.[middlewareName];
	return middleware ?? newMiddleware();
}

function nameCheck(name: string): string {
	return name.split('@')[0].toLowerCase() + '@http';
}

// Create or update a router
export async function upsertRouter(name: string, router: Router, service: Service): Promise<void> {
	let data = get(config);
	if (!data.routers) data.routers = {};
	if (!data.services) data.services = {};

	if (router.name !== name) {
		delete data.routers[name];
		delete data.services[name];
	}
	router.name = nameCheck(router.name);
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
	middleware.name = nameCheck(middleware.name);
	data.middlewares[middleware.name] = middleware;
	await updateConfig(data);
}

// Delete a router by name
export async function deleteRouter(name: string): Promise<void> {
	let data = get(config);
	if (!data.routers || !data.services) return;
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
