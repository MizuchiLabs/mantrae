<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import type { Agent } from '$lib/types';
	import { toast } from 'svelte-sonner';
	import { api, loading } from '$lib/api';

	let newIP = $state('');

	interface Props {
		agent: Agent | undefined;
		open?: boolean;
	}

	let { agent = $bindable({} as Agent), open = $bindable(false) }: Props = $props();

	const handleSubmit = async () => {
		if (!newIP) return;
		if (agent.id) {
			agent.activeIp = newIP;
			await api.updateAgent(agent);
			toast.success(`Agent ${agent.hostname} updated successfully`);
		} else {
			await api.createAgent(agent);
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

{#if agent.id}
	<Dialog.Root bind:open>
		<Dialog.Content class="sm:max-w-[425px]">
			<Dialog.Header>
				<Dialog.Title>{agent.id ? 'Update' : 'Add'} Agent</Dialog.Title>
			</Dialog.Header>

			<form onsubmit={handleSubmit} class="space-y-4">
				<div class="grid grid-cols-4 items-center gap-2">
					<span class="col-span-1 font-mono">Hostname</span>
					<div class="col-span-3 space-x-2">
						<Badge variant="secondary">{agent.hostname}</Badge>
					</div>
				</div>
				<div class="grid grid-cols-4 items-center gap-2">
					<span class="col-span-1 font-mono">Public IP</span>
					<div class="col-span-3 space-x-2">
						{#if agent.activeIp === agent.publicIp || !agent.activeIp}
							<Badge variant="default">{agent.publicIp}</Badge>
						{:else}
							<button onclick={handleSubmit}>
								<Badge variant="secondary">{agent.publicIp}</Badge>
							</button>
						{/if}
					</div>
				</div>
				<div class="grid grid-cols-4 items-center gap-2">
					<span class="col-span-1 font-mono">Private IPs</span>
					<div class="col-span-3 flex flex-wrap gap-2">
						<!-- {#each agent.privateIps ?? [] as ip} -->
						<!-- 	{#if agent.activeIp === ip} -->
						<!-- 		<Badge variant="default">{ip}</Badge> -->
						<!-- 	{:else} -->
						<!-- 		<button onclick={() => (agent.activeIp = ip)}> -->
						<!-- 			<Badge variant="secondary">{ip}</Badge> -->
						<!-- 		</button> -->
						<!-- 	{/if} -->
						<!-- {/each} -->
					</div>
				</div>
				<div class="grid grid-cols-4 items-center gap-2">
					<span class="col-span-1 font-mono">Containers</span>
					<div class="col-span-3 flex flex-wrap gap-2">
						{#each agent.containers ?? [] as container}
							<Badge variant="secondary">{container.name.slice(1)}</Badge>
						{/each}
					</div>
				</div>
				<div class="grid grid-cols-4 items-center gap-2">
					<span class="col-span-1 font-mono">Last Seen</span>
					<div class="col-span-3 flex flex-wrap gap-2">
						<Badge variant="secondary">{new Date(agent.updatedAt).toLocaleString()}</Badge>
					</div>
				</div>
				<div class="mt-4 grid grid-cols-4 items-center gap-2">
					<Label for="ip" class="mr-2">Custom IP</Label>
					<div class="col-span-3 flex flex-wrap gap-2">
						<Input
							id="ip"
							name="ip"
							type="text"
							bind:value={newIP}
							placeholder="Enter a custom ip"
						/>
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
{/if}
