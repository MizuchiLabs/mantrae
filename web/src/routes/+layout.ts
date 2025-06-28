import type { LayoutLoad } from "./$types";
import { goto } from "$app/navigation";
import { useClient } from "$lib/api";
import { profile } from "$lib/stores/profile";
import { user } from "$lib/stores/user";
import { UserService } from "$lib/gen/mantrae/v1/user_pb";
import { ProfileService } from "$lib/gen/mantrae/v1/profile_pb";

export const ssr = false;
export const prerender = false;
export const trailingSlash = "always";

const isPublicRoute = (path: string) => {
	return path.startsWith("/login") || path === "/login";
};

export const load: LayoutLoad = async ({ url, fetch }) => {
	const currentPath = url.pathname;
	const isPublic = isPublicRoute(currentPath);
	// Check if cookie is set

	try {
		const userClient = useClient(UserService, fetch);
		const resUser = await userClient.verifyJWT({});

		if (resUser.user) {
			user.value = resUser.user;

			// Update profile if not set
			if (!profile.id) {
				const profileClient = useClient(ProfileService, fetch);
				const resProfile = await profileClient.listProfiles({});
				profile.value = resProfile.profiles[0];
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
