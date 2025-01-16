import { goto } from '$app/navigation';
import { api, user } from '$lib/api';
import { TOKEN_SK } from '$lib/store';
import type { LayoutLoad } from './$types';

export const ssr = false;
export const prerender = true;
export const trailingSlash = 'always';

const PUBLIC_ROUTES = ['/login/', '/reset/'];

export const load: LayoutLoad = async ({ url }) => {
	const isPublicRoute = PUBLIC_ROUTES.includes(url.pathname);
	const token = localStorage.getItem(TOKEN_SK);

	// Case 1: No token and accessing protected route
	if (!token && !isPublicRoute) {
		await goto('/login/');
		user.set(null);
		return;
	}

	// Case 2: Has token, verify it
	if (token) {
		try {
			const verified = await api.verify();

			// Token is valid
			if (verified) {
				user.set(verified);
			}

			// Trying to access public route
			if (isPublicRoute) {
				await goto('/');
			}
			return;
		} catch (error) {
			// Token verification failed
			api.logout();
			if (!isPublicRoute) {
				await goto('/login/');
			}
			user.set(null);
			return;
		}
	}

	// Case 3: No token and accessing public route
	user.set(null);
	return;
};
