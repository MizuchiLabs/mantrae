import type { LayoutLoad } from './$types';
import { logout, useClient } from '$lib/api';
import { goto } from '$app/navigation';
import { user } from '$lib/stores/user';
import { UserService } from '$lib/gen/mantrae/v1/user_pb';
import { token } from '$lib/stores/common';

export const ssr = false;
export const prerender = true;
export const trailingSlash = 'always';

const isPublicRoute = (path: string) => {
	return path.startsWith('/login') || path === '/login';
};

export const load: LayoutLoad = async ({ url, fetch }) => {
	// Case 1: No token and accessing protected route
	if (!token.value && !isPublicRoute(url.pathname)) {
		await goto('/login/');
		user.clear();
		return;
	}

	// If we have a token, verify it
	if (token.value) {
		try {
			const client = useClient(UserService, fetch);
			const userId = (await client.verifyJWT({ token: token.value })).userId;
			if (!userId) {
				throw new Error('Invalid token');
			}
			const data = await client.getUser({ identifier: { value: userId, case: 'id' } });
			if (!data.user || !data.user.id) {
				throw new Error('User not found');
			}

			// Redirect to home if trying to access login page while authenticated
			if (isPublicRoute(url.pathname) && user.isLoggedIn()) {
				await goto('/');
			}
			return;
		} catch (error) {
			// Token verification failed, clean up
			logout();
			user.clear();
			throw new Error('Token verification failed: ' + error);
		}
	}

	// No token and trying to access protected route
	if (!isPublicRoute) {
		await goto('/login');
	}

	return;
};
