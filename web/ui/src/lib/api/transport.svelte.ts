import { toast } from 'svelte-sonner';
import { useMutation, useQuery } from '$lib/query';
import { profileID } from '$lib/store.svelte';
import { ServersTransportService } from '$lib/gen/mantrae/v1/servers_transport_pb';

export const transport = {
	// Queries
	get: (id: string) =>
		useQuery(
			ServersTransportService.method.getServersTransport,
			{ id },
			{ enabled: !!profileID.current, select: (res) => res.serversTransport }
		),
	list: (pid?: bigint) =>
		useQuery(
			ServersTransportService.method.listServersTransports,
			{ profileId: pid ?? profileID.current },
			{ enabled: !!profileID.current, select: (res) => res.serversTransports }
		),

	// Mutations
	create: () =>
		useMutation(ServersTransportService.method.createServersTransport, {
			onSuccess: () => toast.success('Server transport created!'),
			transform: (variables) => ({ ...variables, profileId: profileID.current })
		}),
	update: () =>
		useMutation(ServersTransportService.method.updateServersTransport, {
			onSuccess: () => toast.success('Server transport updated!')
		}),
	delete: () =>
		useMutation(ServersTransportService.method.deleteServersTransport, {
			onSuccess: () => toast.success('Server transport deleted!')
		})
};
