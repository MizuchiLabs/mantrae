import type { LayoutLoad } from './$types';
import { api } from '$lib/api';
import { goto } from '$app/navigation';
import { user } from '$lib/stores/user';

export const ssr = false;
export const prerender = true;
export const trailingSlash = 'always';

const isPublicRoute = (path: string) => {
	return path.startsWith('/login') || path === '/login';
};

export const load: LayoutLoad = async ({ url, fetch }) => {
	const currentPath = url.pathname;
	const isPublic = isPublicRoute(currentPath);

	// Try to verify authentication via cookie
	try {
		const isVerified = await api.verify(fetch);

		if (isVerified) {
			// User is authenticated
			if (isPublic) {
				// Authenticated user trying to access login page - redirect to home
				await goto('/');
				return;
			}
			// Continue to protected route
			return;
		} else {
			// Verification failed but no exception thrown
			throw new Error('Authentication failed');
		}
	} catch (_) {
		// Authentication failed
		user.clear();

		if (!isPublic) {
			// User trying to access protected route without auth - redirect to login
			await goto('/login');
		}
		// If already on public route, stay there
		return;
	}
};
