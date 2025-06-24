import { goto } from "$app/navigation";
import { logout, profileClient, useClient } from "$lib/api";
import { UserService } from "$lib/gen/mantrae/v1/user_pb";
import { token } from "$lib/stores/common";
import { profile } from "$lib/stores/profile";
import { user } from "$lib/stores/user";
import type { LayoutLoad } from "./$types";

export const ssr = false;
export const prerender = true;
export const trailingSlash = "always";

const isPublicRoute = (path: string) => {
	return path.startsWith("/login") || path === "/login";
};

export const load: LayoutLoad = async ({ url, fetch }) => {
	// Case 1: No token and accessing protected route
	if (!token.value && !isPublicRoute(url.pathname)) {
		await goto("/login/");
		user.clear();
		return;
	}

	// If we have a token, verify it
	if (token.value) {
		try {
			const client = useClient(UserService, fetch);
			const verified = await client.verifyJWT({ token: token.value });
			if (!verified.user) {
				throw new Error("Invalid token");
			}
			user.value = verified.user;
			if (!profile.id) {
				const response = await profileClient.listProfiles({});
				profile.value = response.profiles[0];
			}

			// Redirect to home if trying to access login page while authenticated
			if (isPublicRoute(url.pathname) && user.isLoggedIn()) {
				await goto("/");
			}
			return;
		} catch (error) {
			// Token verification failed, clean up
			logout();
			user.clear();
			throw new Error("Token verification failed: " + error);
		}
	}

	// No token and trying to access protected route
	if (!isPublicRoute) {
		await goto("/login");
	}

	return;
};
