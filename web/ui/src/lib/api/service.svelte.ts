import { toast } from 'svelte-sonner';
import { useMutation, useQuery } from '$lib/query';
import { profileID } from '$lib/store.svelte';
import { ServiceService } from '$lib/gen/mantrae/v1/service_pb';

export const service = {
	// Queries
	get: (name: string) =>
		useQuery(
			ServiceService.method.getService,
			{ identifier: { case: 'name', value: name }, profileId: profileID.current },
			{ enabled: !!profileID.current, select: (res) => res.service }
		),
	list: (pid?: bigint) =>
		useQuery(
			ServiceService.method.listServices,
			{ profileId: pid ?? profileID.current },
			{ enabled: !!profileID.current, select: (res) => res.services }
		),

	// Mutations
	create: () =>
		useMutation(ServiceService.method.createService, {
			onSuccess: () => toast.success('Service created!'),
			transform: (variables) => ({ ...variables, profileId: profileID.current })
		}),
	update: () =>
		useMutation(ServiceService.method.updateService, {
			onSuccess: () => toast.success('Service updated!')
		}),
	delete: () =>
		useMutation(ServiceService.method.deleteService, {
			onSuccess: () => toast.success('Service deleted!')
		})
};
