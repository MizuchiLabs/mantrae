import { goto } from '$app/navigation';
import { API_URL, loggedIn, logout } from '$lib/api';
import { TOKEN_SK } from '$lib/store';
import type { LayoutLoad } from './$types';

export const ssr = false;
export const prerender = true;
export const trailingSlash = 'always';

export const load: LayoutLoad = async ({ fetch, url }) => {
	const token = localStorage.getItem(TOKEN_SK);

	if (token === null) {
		logout();
		if (url.pathname !== '/login' && url.pathname !== '/login/reset/') {
			console.log(url.pathname);
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
		if (url.pathname !== '/login' && url.pathname !== '/login/reset/') {
			console.log(url.pathname);
			goto('/login');
		}
		return {};
	}

	loggedIn.set(true);
};
