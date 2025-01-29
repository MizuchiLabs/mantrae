import type { LayoutLoad } from './$types';
import { TOKEN_SK } from '$lib/store';
import { api, user } from '$lib/api';
import { goto } from '$app/navigation';
import { toast } from 'svelte-sonner';

export const ssr = false;
export const prerender = true;
export const trailingSlash = 'always';

const isPublicRoute = (path: string) => path.startsWith('/login/');

export const load: LayoutLoad = async ({ url, fetch }) => {
	const token = localStorage.getItem(TOKEN_SK);

	// Case 1: No token and accessing protected route
	if (!token && !isPublicRoute(url.pathname)) {
		await goto('/login/');
		user.set(null);
		return;
	}

	// Case 2: Has token, verify it
	if (token) {
		try {
			const verified = await api.verify(fetch);

			// Token is valid
			if (verified) {
				user.set(verified);
			}

			// Trying to access public route
			if (isPublicRoute(url.pathname)) {
				await goto('/');
			}
			return;
		} catch (err: unknown) {
			const error = err instanceof Error ? err : new Error(String(err));
			// Token verification failed
			api.logout();
			if (!isPublicRoute(url.pathname)) {
				await goto('/login');
			}
			user.set(null);
			toast.error('Session expired', { description: error.message });
			return;
		}
	}

	// Case 3: No token and accessing public route
	user.set(null);
	return;
};
