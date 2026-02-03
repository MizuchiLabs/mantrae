import { toast } from 'svelte-sonner';
import { useMutation, useQuery } from '$lib/query';
import { AgentService } from '$lib/gen/mantrae/v1/agent_pb';
import { profileID } from '$lib/store.svelte';

export const agent = {
	// Queries
	get: (id: string) =>
		useQuery(
			AgentService.method.getAgent,
			{ id },
			{ enabled: !!profileID.current, select: (res) => res.agent }
		),
	list: (pid?: bigint) =>
		useQuery(
			AgentService.method.listAgents,
			{ profileId: pid ?? profileID.current },
			{ enabled: !!profileID.current, select: (res) => res.agents }
		),

	// Mutations
	create: () =>
		useMutation(AgentService.method.createAgent, {
			onSuccess: () => toast.success('Agent created!'),
			transform: (variables) => ({ ...variables, profileId: profileID.current })
		}),
	update: () =>
		useMutation(AgentService.method.updateAgent, {
			onSuccess: () => toast.success('Agent updated!')
		}),
	delete: () =>
		useMutation(AgentService.method.deleteAgent, {
			onSuccess: () => toast.success('Agent deleted!')
		})
};
