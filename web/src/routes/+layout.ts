import type { LayoutLoad } from './$types';
import { api } from '$lib/api';
import { goto } from '$app/navigation';
import { toast } from 'svelte-sonner';
import { user } from '$lib/stores/user';
import { token } from '$lib/stores/common';

export const ssr = false;
export const prerender = true;
export const trailingSlash = 'always';

const isPublicRoute = (path: string) => path.startsWith('/login/');

export const load: LayoutLoad = async ({ url, fetch }) => {
	// Case 1: No token and accessing protected route
	if (!token.value && !isPublicRoute(url.pathname)) {
		await goto('/login/');
		user.clear();
		return;
	}

	// Case 2: Has token, verify it
	if (token.value) {
		try {
			await api.verify(fetch);

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
			user.clear();
			toast.error('Session expired', { description: error.message });
			return;
		}
	}

	// Case 3: No token and accessing public route
	user.clear();
	return;
};
