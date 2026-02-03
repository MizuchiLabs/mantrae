import { toast } from 'svelte-sonner';
import { useMutation, useQuery } from '$lib/query';
import { RouterService } from '$lib/gen/mantrae/v1/router_pb';
import { profileID } from '$lib/store.svelte';

export const router = {
	// Queries
	get: (id: string) =>
		useQuery(
			RouterService.method.getRouter,
			{ id },
			{ enabled: !!profileID.current, select: (res) => res.router }
		),
	list: (pid?: bigint) =>
		useQuery(
			RouterService.method.listRouters,
			{ profileId: pid ?? profileID.current },
			{ enabled: !!profileID.current, select: (res) => res.routers }
		),

	// Mutations
	create: () =>
		useMutation(RouterService.method.createRouter, {
			onSuccess: () => toast.success('Router created!'),
			transform: (variables) => ({ ...variables, profileId: profileID.current })
		}),
	update: () =>
		useMutation(RouterService.method.updateRouter, {
			onSuccess: () => toast.success('Router updated!')
		}),
	delete: () =>
		useMutation(RouterService.method.deleteRouter, {
			onSuccess: () => toast.success('Router deleted!')
		})
};
