import { toast } from 'svelte-sonner';
import { useMutation, useQuery } from '$lib/query';
import { profileID } from '$lib/store.svelte';
import { EntryPointService } from '$lib/gen/mantrae/v1/entry_point_pb';

export const entrypoint = {
	// Queries
	get: (id: string) =>
		useQuery(
			EntryPointService.method.getEntryPoint,
			{ id },
			{ enabled: !!profileID.current, select: (res) => res.entryPoint }
		),
	list: (pid?: bigint) =>
		useQuery(
			EntryPointService.method.listEntryPoints,
			{ profileId: pid ?? profileID.current },
			{ enabled: !!profileID.current, select: (res) => res.entryPoints }
		),

	// Mutations
	create: () =>
		useMutation(EntryPointService.method.createEntryPoint, {
			onSuccess: () => toast.success('Entrypoint created!'),
			transform: (variables) => ({ ...variables, profileId: profileID.current })
		}),
	update: () =>
		useMutation(EntryPointService.method.updateEntryPoint, {
			onSuccess: () => toast.success('Entrypoint updated!')
		}),
	delete: () =>
		useMutation(EntryPointService.method.deleteEntryPoint, {
			onSuccess: () => toast.success('Entrypoint deleted!')
		})
};
