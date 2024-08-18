import { goto } from '$app/navigation';
import { loggedIn, logout } from '$lib/api';
import type { LayoutLoad } from './$types';

export const ssr = false;
export const prerender = true;

export const load: LayoutLoad = async ({ url }) => {
	const token = localStorage.getItem('token');
	const expiry = localStorage.getItem('expiry') as string;

	if (token === null || expiry === null) {
		logout();
		if (url.pathname !== '/login') {
			goto('/login');
		}
		return {};
	}

	const expiryDate = new Date(expiry);
	if (Date.now() > expiryDate.getTime()) {
		logout();
		if (url.pathname !== '/login') {
			goto('/login');
		}
		return {};
	}

	try {
		if (url.pathname === '/login') {
			goto('/');
		}
		loggedIn.set(true);
	} catch (e) {
		logout();
		if (url.pathname !== '/login') {
			goto('/login');
		}
	}
};
