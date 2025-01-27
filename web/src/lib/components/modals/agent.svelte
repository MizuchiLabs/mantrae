<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import type { Agent, UpdateAgentIPParams } from '$lib/types';
	import { toast } from 'svelte-sonner';
	import { api, loading } from '$lib/api';
	import Separator from '../ui/separator/separator.svelte';

	interface Props {
		agent: Agent | undefined;
		open?: boolean;
	}

	let { agent = $bindable({} as Agent), open = $bindable(false) }: Props = $props();

	let newIP: string | undefined = $state();
	const handleSubmit = async () => {
		if (agent.id) {
			if (!newIP) return;
			const params: UpdateAgentIPParams = {
				id: agent.id,
				activeIp: newIP
			};
			await api.updateAgentIP(params);
			toast.success(`Agent ${agent.hostname} updated successfully`);
		} else {
			await api.createAgent();
			toast.success(`Agent ${agent.hostname} created successfully`);
		}
		open = false;
	};

	const handleDelete = async () => {
		try {
			await api.deleteAgent(agent.id);
			toast.success(`Agent ${agent.hostname} deleted successfully`);
			open = false;
		} catch (err: unknown) {
			const e = err as Error;
			toast.error('Failed to delete agent', {
				description: e.message
			});
		}
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header>
			<Dialog.Title>{agent.id ? 'Update' : 'Add'} Agent</Dialog.Title>
			<Dialog.Description>Update the active IP address for your agent</Dialog.Description>
		</Dialog.Header>

		<Separator />

		<form onsubmit={handleSubmit} class="flex flex-col gap-4">
			<div class="grid grid-cols-4 items-center gap-2">
				<Label for="hostname">Hostname</Label>
				<div class="col-span-3 space-x-2">
					<Badge variant="secondary">{agent.hostname}</Badge>
				</div>
			</div>
			<div class="grid grid-cols-4 items-center gap-2">
				<Label for="publicip">Public IP</Label>
				<div class="col-span-3 space-x-2">
					{#if agent.activeIp === agent.publicIp || !agent.activeIp}
						<Badge variant="default">{agent.publicIp}</Badge>
					{:else}
						<button onclick={() => (newIP = agent.publicIp)}>
							<Badge variant="secondary">{agent.publicIp}</Badge>
						</button>
					{/if}
				</div>
			</div>
			<div class="grid grid-cols-4 items-center gap-2">
				<Label for="privateip">Private IPs</Label>
				<div class="col-span-3 flex flex-wrap gap-2">
					{#each agent.privateIps.privateIps ?? [] as ip}
						{#if agent.activeIp === ip}
							<Badge variant="default">{ip}</Badge>
						{:else}
							<button onclick={() => (newIP = ip)}>
								<Badge variant="secondary">{ip}</Badge>
							</button>
						{/if}
					{/each}
				</div>
			</div>
			<div class="grid grid-cols-4 items-center gap-2">
				<Label for="containers">Containers</Label>
				<div class="col-span-3 flex flex-wrap gap-2">
					{#each agent.containers ?? [] as container}
						{#if container.name}
							<Badge variant="secondary">{container.name.slice(1)}</Badge>
						{/if}
					{/each}
				</div>
			</div>
			<div class="grid grid-cols-4 items-center gap-2">
				<Label for="lastseen">Last Seen</Label>
				<div class="col-span-3 flex flex-wrap gap-2">
					<Badge variant="secondary">{new Date(agent.updatedAt).toLocaleString()}</Badge>
				</div>
			</div>

			<Separator />

			<div class="grid grid-cols-4 items-center gap-2">
				<Label for="ip" class="mr-2">Custom IP</Label>
				<div class="col-span-3 flex flex-wrap gap-2">
					<Input id="ip" name="ip" type="text" bind:value={newIP} placeholder="Enter a custom ip" />
				</div>
			</div>

			<Dialog.Footer>
				<Button type="button" variant="destructive" onclick={handleDelete} disabled={$loading}>
					Delete
				</Button>
				<Button type="submit" disabled={$loading}>Update</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>
