import { get, writable, type Writable } from 'svelte/store';
import { toast } from 'svelte-sonner';
import type { Profile } from './types/dynamic';
import type { Router, Service } from './types/config';
import type { Middleware } from './types/middlewares';

export const profiles: Writable<Profile[]> = writable([]);
export const activeProfile: Writable<Profile> = writable({} as Profile);
export const API_URL = import.meta.env.PROD ? '/api' : 'http://localhost:3000/api';

async function handleError(response: Response) {
	if (!response.ok) {
		toast.error('Request failed', {
			description: await response.text(),
			duration: 3000
		});
		throw new Error(`Failed to fetch: ${response}`);
	}
}

// Profiles -------------------------------------------------------------------
export async function getProfiles() {
	const response = await fetch(`${API_URL}/profiles`);
	handleError(response);

	profiles.set(await response.json());
}

export async function createProfile(profile: Profile): Promise<void> {
	const response = await fetch(`${API_URL}/profiles`, {
		method: 'POST',
		body: JSON.stringify(profile)
	});
	handleError(response);

	profiles.update((profiles) => [...profiles, profile]);
}

export async function updateProfile(name: string, profile: Profile): Promise<void> {
	const response = await fetch(`${API_URL}/profiles/${name}`, {
		method: 'PUT',
		body: JSON.stringify(profile)
	});
	handleError(response);

	profiles.update((profiles) => profiles.map((p) => (p.name === name ? profile : p)));

	if (profile.name === get(activeProfile).name) {
		activeProfile.set(profile);
	}
}

export async function deleteProfile(name: string): Promise<void> {
	const response = await fetch(`${API_URL}/profiles/${name}`, {
		method: 'DELETE'
	});
	handleError(response);

	profiles.update((profiles) => profiles.filter((p) => p.name !== name));
}

// Routers --------------------------------------------------------------------
export async function updateRouter(
	profileName: string,
	router: Router,
	oldRouter: string
): Promise<void> {
	const response = await fetch(`${API_URL}/routers/${profileName}/${oldRouter}`, {
		method: 'PUT',
		body: JSON.stringify(router)
	});
	handleError(response);

	const updatedProfile = await response.json();
	profiles.update((profiles) => [...profiles, updatedProfile]);
	if (updatedProfile.name === get(activeProfile).name) {
		activeProfile.set(updatedProfile);
	}
}

export async function deleteRouter(profileName: string, routerName: string): Promise<void> {
	const response = await fetch(`${API_URL}/routers/${profileName}/${routerName}`, {
		method: 'DELETE'
	});
	handleError(response);

	const updatedProfile = await response.json();
	profiles.update((profiles) => profiles.map((p) => (p.name === profileName ? updatedProfile : p)));

	if (updatedProfile.name === get(activeProfile).name) {
		activeProfile.set(updatedProfile);
	}
}

// Services -------------------------------------------------------------------
export async function updateService(
	profileName: string,
	service: Service,
	oldService: string
): Promise<void> {
	const response = await fetch(`${API_URL}/services/${profileName}/${oldService}`, {
		method: 'PUT',
		body: JSON.stringify(service)
	});
	handleError(response);

	const updatedProfile = await response.json();
	profiles.update((profiles) => [...profiles, updatedProfile]);
	if (updatedProfile.name === get(activeProfile).name) {
		activeProfile.set(updatedProfile);
	}
}

export async function deleteService(profileName: string, serviceName: string): Promise<void> {
	const response = await fetch(`${API_URL}/services/${profileName}/${serviceName}`, {
		method: 'DELETE'
	});
	handleError(response);

	const updatedProfile = await response.json();
	profiles.update((profiles) => profiles.map((p) => (p.name === profileName ? updatedProfile : p)));

	if (updatedProfile.name === get(activeProfile).name) {
		activeProfile.set(updatedProfile);
	}
}

// Middlewares ----------------------------------------------------------------
export async function updateMiddleware(profileName: string, middleware: Middleware): Promise<void> {
	const response = await fetch(`${API_URL}/middlewares/${profileName}/${middleware.name}`, {
		method: 'PUT',
		body: JSON.stringify(middleware)
	});
	handleError(response);

	const updatedProfile = await response.json();
	profiles.update((profiles) => [...profiles, updatedProfile]);
	if (updatedProfile.name === get(activeProfile).name) {
		activeProfile.set(updatedProfile);
	}
}

export async function deleteMiddleware(name: string): Promise<void> {
	const response = await fetch(`${API_URL}/middlewares/${name}`, {
		method: 'DELETE'
	});
	handleError(response);

	const updatedProfile = await response.json();
	profiles.update((profiles) => profiles.map((p) => (p.name === name ? updatedProfile : p)));

	if (updatedProfile.name === get(activeProfile).name) {
		activeProfile.set(updatedProfile);
	}
}

// Helpers
// export function deleteRouter(name: string) {
// 	activeProfile.update((p) => {
// 		if (p.instance.dynamic?.routers === undefined) return p;
// 		p.instance.dynamic.routers = p.instance.dynamic.routers.filter((r) => r.service !== name);
// 		p.instance.dynamic.services = p.instance.dynamic.services?.filter((s) => s.name !== name);
// 		return p;
// 	});
// 	updateProfile(get(activeProfile).name, get(activeProfile));
// }
//
// export function deleteMiddleware(name: string) {
// 	activeProfile.update((p) => {
// 		if (p.instance.dynamic?.httpmiddlewares === undefined) return p;
// 		p.instance.dynamic.httpmiddlewares = p.instance.dynamic.httpmiddlewares.filter(
// 			(m) => m.name !== name
// 		);
// 		if (p.instance.dynamic?.tcpmiddlewares === undefined) return p;
// 		p.instance.dynamic.tcpmiddlewares = p.instance.dynamic.tcpmiddlewares.filter(
// 			(m) => m.name !== name
// 		);
// 		return p;
// 	});
// 	updateProfile(get(activeProfile).name, get(activeProfile));
// }
