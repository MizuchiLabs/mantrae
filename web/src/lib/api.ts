import { derived, get, writable, type Writable } from 'svelte/store';
import { toast } from 'svelte-sonner';
import type { Profile } from './types/dynamic';
import type { Router, Service } from './types/config';
import type { Middleware } from './types/middlewares';
import { goto } from '$app/navigation';

export const loggedIn = writable(false);
export const profiles: Writable<Profile[]> = writable([]);
export const activeProfile: Writable<Profile> = writable({} as Profile);
export const API_URL = import.meta.env.PROD ? '/api' : 'http://localhost:3000/api';

export const entrypoints = derived(
	activeProfile,
	($activeProfile) => $activeProfile?.client?.dynamic?.entrypoints ?? []
);
export const routers = derived(activeProfile, ($activeProfile) =>
	Object.values($activeProfile?.client?.dynamic?.routers ?? [])
);
export const services = derived(activeProfile, ($activeProfile) =>
	Object.values($activeProfile?.client?.dynamic?.services ?? [])
);
export const middlewares = derived(activeProfile, ($activeProfile) =>
	Object.values($activeProfile?.client?.dynamic?.middlewares ?? [])
);

async function handleError(response: Response) {
	if (!response.ok) {
		toast.error('Request failed', {
			description: await response.text(),
			duration: 3000
		});
		throw new Error(`Failed to fetch: ${response}`);
	}
}

async function handleRequest(endpoint: string, method: string, body?: any): Promise<any> {
	if (!get(loggedIn)) return;

	const token = localStorage.getItem('token');
	const response = await fetch(`${API_URL}${endpoint}`, {
		method: method,
		body: body ? JSON.stringify(body) : undefined,
		headers: { Authorization: `Bearer ${token}` }
	});
	handleError(response);
	if (response.status !== 204) {
		return await response.json();
	}
}

// Login ----------------------------------------------------------------------
export async function login(username: string, password: string) {
	const response = await fetch(`${API_URL}/login`, {
		method: 'POST',
		body: JSON.stringify({ username, password })
	});
	handleError(response);

	const { token } = await response.json();
	localStorage.setItem('token', token);
	loggedIn.set(true);
	await getProfiles();
	goto('/');
}

export async function logout() {
	localStorage.removeItem('token');
	loggedIn.set(false);
}

// Profiles -------------------------------------------------------------------
export async function getProfiles() {
	const response = await handleRequest('/profiles', 'GET');

	profiles.set(response);
	if (get(activeProfile).name === undefined) {
		activeProfile.set(get(profiles)[0]);
	}
}

export async function createProfile(profile: Profile): Promise<void> {
	await handleRequest('/profiles', 'POST', profile);
	profiles.update((profiles) => [...profiles, profile]);
}

export async function updateProfile(name: string, profile: Profile): Promise<void> {
	await handleRequest(`/profiles/${name}`, 'PUT', profile);
	profiles.update((profiles) => profiles.map((p) => (p.name === name ? profile : p)));

	if (profile.name === get(activeProfile).name) {
		activeProfile.set(profile);
	}
}

export async function deleteProfile(name: string): Promise<void> {
	await handleRequest(`/profiles/${name}`, 'DELETE');
	profiles.update((profiles) => profiles.filter((p) => p.name !== name));
}

// Routers --------------------------------------------------------------------
export async function updateRouter(
	profileName: string,
	router: Router,
	oldRouter: string
): Promise<void> {
	const response = await handleRequest(`/routers/${profileName}/${oldRouter}`, 'PUT', router);

	profiles.update((profiles) => [...profiles, response]);
	if (response.name === get(activeProfile).name) {
		activeProfile.set(response);
	}
}

export async function deleteRouter(profileName: string, routerName: string): Promise<void> {
	const response = await handleRequest(`/routers/${profileName}/${routerName}`, 'DELETE');

	profiles.update((profiles) => profiles.map((p) => (p.name === profileName ? response : p)));
	if (response.name === get(activeProfile).name) {
		activeProfile.set(response);
	}
}

// Services -------------------------------------------------------------------
export async function updateService(
	profileName: string,
	service: Service,
	oldService: string
): Promise<void> {
	const response = await handleRequest(`/services/${profileName}/${oldService}`, 'PUT', service);

	profiles.update((profiles) => [...profiles, response]);
	if (response.name === get(activeProfile).name) {
		activeProfile.set(response);
	}
}

export async function deleteService(profileName: string, serviceName: string): Promise<void> {
	const response = await handleRequest(`/services/${profileName}/${serviceName}`, 'DELETE');

	profiles.update((profiles) => profiles.map((p) => (p.name === profileName ? response : p)));
	if (response.name === get(activeProfile).name) {
		activeProfile.set(response);
	}
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

	profiles.update((profiles) => [...profiles, response]);
	if (response.name === get(activeProfile).name) {
		activeProfile.set(response);
	}
}

export async function deleteMiddleware(profileName: string, name: string): Promise<void> {
	const response = await handleRequest(`/middlewares/${profileName}/${name}`, 'DELETE');

	profiles.update((profiles) => profiles.map((p) => (p.name === name ? response : p)));
	if (response.name === get(activeProfile).name) {
		activeProfile.set(response);
	}
}
