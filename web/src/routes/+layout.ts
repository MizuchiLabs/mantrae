import type { LayoutLoad } from "./$types";
import { goto } from "$app/navigation";
import { profileClient, userClient } from "$lib/api";
import { profile } from "$lib/stores/profile";
import { user } from "$lib/stores/user";

export const ssr = false;
export const prerender = true;
export const trailingSlash = "always";

const isPublicRoute = (path: string) => {
	return path.startsWith("/login") || path === "/login";
};

export const load: LayoutLoad = async ({ url }) => {
	const currentPath = url.pathname;
	const isPublic = isPublicRoute(currentPath);

	try {
		const verified = await userClient.verifyJWT({});

		if (verified.user) {
			user.value = verified.user;

			// Update profile if not set
			if (!profile.id) {
				const response = await profileClient.listProfiles({});
				profile.value = response.profiles[0];
			}

			if (isPublic) {
				// Authenticated user trying to access login page - redirect to home
				await goto("/");
				return;
			}
			return;
		} else {
			throw new Error("Authentication failed");
		}
	} catch (_) {
		user.clear();
		if (!isPublic) {
			await goto("/login");
		}
		return;
	}
};
