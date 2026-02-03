import { toast } from 'svelte-sonner';
import { useMutation, useQuery } from '$lib/query';
import { ProfileService } from '$lib/gen/mantrae/v1/profile_pb';
import { profileID } from '$lib/store.svelte';

export const profile = {
	// Queries
	get: (id?: bigint) =>
		useQuery(
			ProfileService.method.getProfile,
			{ id: id ?? profileID.current },
			{ enabled: !!profileID.current, select: (res) => res.profile }
		),
	list: () =>
		useQuery(
			ProfileService.method.listProfiles,
			{},
			{
				select: (res) => {
					if (!profileID.current && res.profiles.length > 0) {
						profileID.current = res.profiles[0].id;
					}
					return res.profiles;
				}
			}
		),

	// Mutations
	create: () =>
		useMutation(ProfileService.method.createProfile, {
			onSuccess: (data) => {
				if (data.profile) profileID.current = data.profile.id;
				toast.success('Profile created!');
			}
		}),
	update: () =>
		useMutation(ProfileService.method.updateProfile, {
			onSuccess: () => toast.success('Profile updated!'),
			transform: (variables) => ({ ...variables, id: profileID.current })
		}),
	delete: () =>
		useMutation(ProfileService.method.deleteProfile, {
			onSuccess: (_, variables) => {
				if (variables.id === profileID.current) {
					profileID.current = 0n;
				}
				toast.success('Profile deleted');
			}
		})
};
