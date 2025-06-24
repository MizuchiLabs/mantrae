import { goto } from "$app/navigation";
import { profileClient, userClient } from "$lib/api";
import { profile } from "$lib/stores/profile";
import { user } from "$lib/stores/user";
import type { LayoutLoad } from "./$types";

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

// export const load: LayoutLoad = async ({ url }) => {
// 	// Case 1: No token and accessing protected route
// 	if (!token.value && !isPublicRoute(url.pathname)) {
// 		await goto("/login/");
// 		user.clear();
// 		return;
// 	}
//
// 	// If we have a token, verify it
// 	if (token.value) {
// 		try {
// 			const verified = await userClient.verifyJWT({});
// 			if (!verified.user) {
// 				throw new Error("Invalid token");
// 			}
// 			user.value = verified.user;
// 			if (!profile.id) {
// 				const response = await profileClient.listProfiles({});
// 				profile.value = response.profiles[0];
// 			}
//
// 			// Redirect to home if trying to access login page while authenticated
// 			if (isPublicRoute(url.pathname) && user.isLoggedIn()) {
// 				await goto("/");
// 			}
// 			return;
// 		} catch (error) {
// 			// Token verification failed, clean up
// 			logout();
// 			user.clear();
// 			throw new Error("Token verification failed: " + error);
// 		}
// 	}
//
// 	// No token and trying to access protected route
// 	if (!isPublicRoute) {
// 		await goto("/login");
// 	}
//
// 	return;
// };
