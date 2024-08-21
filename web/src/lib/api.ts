import { goto } from '$app/navigation';
import { toast } from 'svelte-sonner';
import type { Profile } from './types/dynamic';
import type { Middleware } from './types/middlewares';
import { derived, get, writable, type Writable } from 'svelte/store';
import { newRouter, newService, type Router, type Service } from './types/config';

export const loggedIn = writable(false);
export const profile = writable('');
export const profiles: Writable<Record<string, Profile>> = writable({});
export const API_URL = import.meta.env.PROD ? '/api' : 'http://localhost:3000/api';

export const entrypoints = derived(
	[profiles, profile],
	([$profiles, $profile]) => $profiles?.[$profile]?.dynamic?.entrypoints ?? []
);
export const routers = derived([profiles, profile], ([$profiles, $profile]) =>
	Object.values($profiles?.[$profile]?.dynamic?.routers ?? [])
);
export const services = derived([profiles, profile], ([$profiles, $profile]) =>
	Object.values($profiles?.[$profile]?.dynamic?.services ?? [])
);
export const middlewares = derived([profiles, profile], ([$profiles, $profile]) =>
	Object.values($profiles?.[$profile]?.dynamic?.middlewares ?? [])
);
export const version = derived(
	[profiles, profile],
	([$profiles, $profile]) => $profiles?.[$profile]?.dynamic?.version ?? ''
);

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
		profile.set(savedProfile);
	}
	if (get(profile) === '' && Object.keys(response).length > 0) {
		profile.set(Object.keys(response)[0]);
		localStorage.setItem('profile', Object.keys(response)[0]);
	}
}

export async function createProfile(profile: Profile): Promise<void> {
	const response = await handleRequest('/profiles', 'POST', profile);
	if (response) {
		profiles.set(response);
		toast.success(`Profile ${profile.name} created`);
	}
}

export async function updateProfile(name: string, p: Profile): Promise<void> {
	const response = await handleRequest(`/profiles/${name}`, 'PUT', p);
	if (response) {
		profiles.set(response);
		toast.success(`Profile ${name} updated`);
		if (get(profile) === name) {
			profile.set(p.name);
			localStorage.setItem('profile', p.name);
		}
	}
}

export async function deleteProfile(name: string): Promise<void> {
	const response = await handleRequest(`/profiles/${name}`, 'DELETE');
	if (response) {
		profiles.set(response);
		toast.success(`Profile ${name} deleted`);
	}
}

// Routers --------------------------------------------------------------------
export async function updateRouter(
	profileName: string,
	router: Router,
	oldRouter: string
): Promise<void> {
	const response = await handleRequest(`/routers/${profileName}/${oldRouter}`, 'PUT', router);
	profiles.set(response);
	toast.success(`Router ${router.name} updated`);
}

export async function deleteRouter(profileName: string, routerName: string): Promise<void> {
	const response = await handleRequest(`/routers/${profileName}/${routerName}`, 'DELETE');
	profiles.set(response);
	toast.success(`Router ${routerName} deleted`);
}

// Services -------------------------------------------------------------------
export async function updateService(
	profileName: string,
	service: Service,
	oldService: string
): Promise<void> {
	const response = await handleRequest(`/services/${profileName}/${oldService}`, 'PUT', service);
	profiles.set(response);
	toast.success(`Service ${service.name} updated`);
}

export async function deleteService(profileName: string, serviceName: string): Promise<void> {
	const response = await handleRequest(`/services/${profileName}/${serviceName}`, 'DELETE');
	profiles.set(response);
	toast.success(`Service ${serviceName} deleted`);
}

// Middlewares ----------------------------------------------------------------
export async function updateMiddleware(
	profileName: string,
	middleware: Middleware,
	oldMiddleware: string
): Promise<void> {
	const response = await handleRequest(
		`/middlewares/${profileName}/${oldMiddleware}`,
		'PUT',
		middleware
	);

	profiles.set(response);
	toast.success(`Middleware ${middleware.name} updated`);
}

export async function deleteMiddleware(profileName: string, name: string): Promise<void> {
	const response = await handleRequest(`/middlewares/${profileName}/${name}`, 'DELETE');
	profiles.set(response);
	toast.success(`Middleware ${name} deleted`);
}

// Helper functions -----------------------------------------------------------
export function getRouter(profileName: string, routerName: string): Router {
	const router = get(profiles)?.[profileName].dynamic?.routers?.[routerName];
	return router ?? newRouter();
}
export function getService(profileName: string, serviceName: string): Service {
	const service = get(profiles)?.[profileName].dynamic?.services?.[serviceName];
	return service ?? newService();
}
