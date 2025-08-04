import type { LayoutLoad } from './$types';
import { goto } from '$app/navigation';
import { checkHealth, useClient } from '$lib/api';
import { profile } from '$lib/stores/profile';
import { user } from '$lib/stores/user';
import { UserService } from '$lib/gen/mantrae/v1/user_pb';
import { ProfileService } from '$lib/gen/mantrae/v1/profile_pb';

export const ssr = false;
export const prerender = true;
export const trailingSlash = 'always';

export const load: LayoutLoad = async ({ url, fetch }) => {
	const currentPath = url.pathname;
	const isPublic = currentPath.startsWith('/login') || currentPath.startsWith('/welcome');

	const healthy = await checkHealth(fetch);
	if (!healthy) {
		// No backend, force redirect to welcome screen to enter backend URL
		if (currentPath !== '/welcome') {
			await goto('/welcome');
			user.clear();
			return {};
		}
	} else {
		// Backend reachable
		if (currentPath === '/welcome') {
			// Backend is back, redirect from welcome to login
			await goto('/login');
			return {};
		}

		try {
			const userClient = useClient(UserService, fetch);
			const resUser = await userClient.getUser({});
			if (!resUser.user) throw new Error('Authentication failed');
			user.value = resUser.user;

			const profileClient = useClient(ProfileService, fetch);
			if (!profile.id) {
				const response = await profileClient.listProfiles({});
				profile.value = response.profiles[0];
			} else {
				const response = await profileClient.getProfile({ id: profile.id });
				if (!response.profile) throw new Error('Profile not found');
				profile.value = response.profile;
			}

			if (isPublic) await goto('/');
		} catch (_) {
			user.clear();
			profile.clear();
			if (!isPublic) await goto('/login');
		}
	}
};
