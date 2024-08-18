import { goto } from '$app/navigation';
import { API_URL, loggedIn, logout } from '$lib/api';
import type { LayoutLoad } from './$types';

export const ssr = false;
export const prerender = true;

export const load: LayoutLoad = async ({ fetch, url }) => {
	const token = localStorage.getItem('token');

	if (token === null) {
		logout();
		if (url.pathname !== '/login') {
			goto('/login');
		}
		return {};
	}

	const response = await fetch(`${API_URL}/verify`, {
		method: 'POST',
		headers: { Authorization: `Bearer ${token}` }
	});
	if (!response.ok) {
		logout();
		if (url.pathname !== '/login') {
			goto('/login');
		}
		return {};
	}

	loggedIn.set(true);
};
