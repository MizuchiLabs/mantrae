import { goto } from '$app/navigation';
import { toast } from 'svelte-sonner';
import type { Profile } from './types/dynamic';
import type { Middleware } from './types/middlewares';
import { derived, get, writable, type Writable } from 'svelte/store';
import { newRouter, newService, type Router, type Service } from './types/config';
import type { Provider } from './types/provider';

export const loggedIn = writable(false);
export const profile: Writable<Profile> = writable();
export const profiles: Writable<Record<string, Profile>> = writable();
export const provider: Writable<Record<string, Provider>> = writable();
export const API_URL = import.meta.env.PROD ? '/api' : 'http://localhost:3000/api';

export const routers = derived(profile, ($profile) =>
	Object.values($profile?.dynamic?.routers ?? [])
);
export const services = derived(profile, ($profile) =>
	Object.values($profile?.dynamic?.services ?? [])
);
export const middlewares = derived(profile, ($profile) =>
	Object.values($profile?.dynamic?.middlewares ?? [])
);
export const version = derived(profile, ($profile) => $profile?.dynamic?.version ?? '');
export const entrypoints = derived(profile, ($profile) => $profile?.dynamic?.entrypoints ?? []);

async function handleRequest(endpoint: string, method: string, body?: any): Promise<any> {
	if (!get(loggedIn)) return;

	const token = localStorage.getItem('token');
	try {
		const response = await fetch(`${API_URL}${endpoint}`, {
			method: method,
			body: body ? JSON.stringify(body) : undefined,
			headers: { Authorization: `Bearer ${token}` }
		});
		return await response.json();
	} catch (e: any) {
		toast.error('Request failed', {
			description: e.message,
			duration: 3000
		});
	}
}

// Login ----------------------------------------------------------------------
export async function login(username: string, password: string) {
	try {
		const response = await fetch(`${API_URL}/login`, {
			method: 'POST',
			body: JSON.stringify({ username, password })
		});
		const { token } = await response.json();
		localStorage.setItem('token', token);
		loggedIn.set(true);
		await getProfiles();
		goto('/');
		toast.success('Login successful');
	} catch (e: any) {
		toast.error('Login failed', {
			description: e.message,
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
	const response = await handleRequest('/profiles', 'GET');
	profiles.set(response);

	// Get saved profile
	const savedProfile = localStorage.getItem('profile');
	if (savedProfile !== null) {
		getProfile(savedProfile);
	}
	if (!get(profile) && Object.keys(response).length > 0) {
		getProfile(Object.keys(response)[0]);
	}
}

export async function getProfile(name: string) {
	const response = await handleRequest('/profile/' + name, 'GET');
	profile.set(response);
	localStorage.setItem('profile', response.name);
}

export async function createProfile(profile: Profile): Promise<void> {
	const response = await handleRequest('/profiles', 'POST', profile);
	if (response) {
		profiles.set(response);
		toast.success(`Profile ${profile.name} created`);
	}
}

export async function updateProfile(name: string, p: Profile): Promise<void> {
	if (!get(profile)) return;
	const response = await handleRequest(`/profiles/${get(profile).name}`, 'PUT', p);
	if (response) {
		profile.set(response);
		toast.success(`Profile ${name} updated`);
		if (get(profile).name === name) {
			localStorage.setItem('profile', response.name);
		}
	}
}

export async function deleteProfile(name: string): Promise<void> {
	const response = await handleRequest(`/profiles/${name}`, 'DELETE');
	if (response) {
		profiles.set(response);
		toast.success(`Profile ${name} deleted`);
		if (get(profile).name === name) {
			profile.set({} as Profile);
			localStorage.removeItem('profile');
		}
	}
}

// Providers ------------------------------------------------------------------
export async function getProviders() {
	const response = await handleRequest('/providers', 'GET');
	provider.set(response);
}

export async function updateProvider(oldName: string, p: Provider): Promise<void> {
	const response = await handleRequest(`/providers/${oldName}`, 'PUT', p);
	provider.set(response);
	toast.success(`Provider ${p.name} updated`);
}

export async function deleteProvider(name: string): Promise<void> {
	const response = await handleRequest(`/providers/${name}`, 'DELETE');
	provider.set(response);
	toast.success(`Provider ${name} deleted`);
}

// Routers --------------------------------------------------------------------
export async function updateRouter(
	oldName: string,
	router: Router,
	service: Service
): Promise<void> {
	if (!get(profile)) return;

	const resRouter = await handleRequest(`/routers/${get(profile).name}/${oldName}`, 'PUT', router);
	if (resRouter) {
		const resService = await handleRequest(
			`/services/${get(profile).name}/${oldName}`,
			'PUT',
			service
		);
		if (resService) {
			profile.set(resService);
			toast.success(`Router ${router.name} updated`);
		}
	}
}

export async function deleteRouter(name: string): Promise<void> {
	if (!get(profile)) return;
	await handleRequest(`/routers/${get(profile).name}/${name}`, 'DELETE');
	const response = await handleRequest(`/services/${get(profile).name}/${name}`, 'DELETE');
	profile.set(response);
	toast.success(`Router ${name} deleted`);
}

// Middlewares ----------------------------------------------------------------
export async function updateMiddleware(
	middleware: Middleware,
	oldMiddleware: string
): Promise<void> {
	if (!get(profile)) return;
	const response = await handleRequest(
		`/middlewares/${get(profile).name}/${oldMiddleware}`,
		'PUT',
		middleware
	);

	profile.set(response);
	toast.success(`Middleware ${middleware.name} updated`);
}

export async function deleteMiddleware(name: string): Promise<void> {
	if (!get(profile)) return;
	const response = await handleRequest(`/middlewares/${get(profile).name}/${name}`, 'DELETE');
	profile.set(response);
	toast.success(`Middleware ${name} deleted`);
}

// Helper functions -----------------------------------------------------------
export function getRouter(routerName: string): Router {
	const router = get(profile)?.dynamic?.routers?.[routerName];
	return router ?? newRouter();
}
export function getService(serviceName: string): Service {
	const service = get(profile)?.dynamic?.services?.[serviceName];
	return service ?? newService();
}
