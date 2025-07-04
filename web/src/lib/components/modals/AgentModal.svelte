<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { toast } from 'svelte-sonner';
	import Separator from '../ui/separator/separator.svelte';
	import CopyButton from '../ui/copy-button/copy-button.svelte';
	import { DateFormat, pageIndex, pageSize } from '$lib/stores/common';
	import { RotateCcw } from '@lucide/svelte';
	import { agentClient } from '$lib/api';
	import { profile } from '$lib/stores/profile';
	import { ConnectError } from '@connectrpc/connect';
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import type { Agent } from '$lib/gen/mantrae/v1/agent_pb';

	interface Props {
		data: Agent[];
		item: Agent;
		open?: boolean;
	}

	let { data = $bindable(), item = $bindable(), open = $bindable(false) }: Props = $props();

	let newIP = $state('');

	const handleSubmit = async (ip: string | undefined) => {
		try {
			if (item.id) {
				const response = await agentClient.updateAgentIP({ id: item.id, ip });
				if (!response.agent) throw new Error('Failed to create agent');
				item = response.agent;
				toast.success(`Agent ${item.hostname} updated successfully`);
			} else {
				await agentClient.createAgent({ profileId: profile.id });
				toast.success(`Agent ${item.hostname} created successfully`);
				open = false;
			}

			// Refresh data
			let response = await agentClient.listAgents({
				profileId: profile.id,
				limit: BigInt(pageSize.value ?? 10),
				offset: BigInt(pageIndex.value ?? 0)
			});
			data = response.agents;
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to save agent', { description: e.message });
		}
		// open = false;
	};
	const handleDelete = async () => {
		if (!item.id) return;

		try {
			await agentClient.deleteAgent({ id: item.id });
			toast.success('Agent deleted successfully');

			// Refresh data
			let response = await agentClient.listAgents({
				profileId: profile.id,
				limit: BigInt(pageSize.value ?? 10),
				offset: BigInt(pageIndex.value ?? 0)
			});
			data = response.agents;
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete agent', { description: e.message });
		}
		open = false;
	};
	const handleRotate = async () => {
		const response = await agentClient.rotateAgentToken({ id: item.id });
		if (!response.agent) throw new Error('Failed to rotate token');
		item.token = response.agent.token;
		toast.success('Token rotated successfully');

		// Refresh data
		let response2 = await agentClient.listAgents({
			profileId: profile.id,
			limit: BigInt(pageSize.value ?? 10),
			offset: BigInt(pageIndex.value ?? 0)
		});
		data = response2.agents;
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header>
			<Dialog.Title>{item.hostname ? 'Update' : 'Connect your'} Agent</Dialog.Title>
			<Dialog.Description>
				{item.hostname
					? 'Update the active IP address for your agent'
					: 'Copy the token for your agent below'}
			</Dialog.Description>
		</Dialog.Header>

		<Separator />

		<div class="flex flex-col gap-4">
			{#if item.hostname}
				<div class="grid grid-cols-4 items-center gap-2">
					<Label for="hostname">Hostname</Label>
					<div class="col-span-3 space-x-2">
						<Badge variant="secondary">{item.hostname}</Badge>
					</div>
				</div>
			{/if}

			{#if item.publicIp}
				<div class="grid grid-cols-4 items-center gap-2">
					<Label for="publicip">Public IP</Label>
					<div class="col-span-3 space-x-2">
						{#if item.activeIp === item.publicIp || !item.activeIp}
							<Badge variant="default">{item.publicIp ?? 'None'}</Badge>
						{:else}
							<button onclick={() => handleSubmit(item.publicIp)}>
								<Badge variant="secondary">{item.publicIp}</Badge>
							</button>
						{/if}
					</div>
				</div>
			{/if}

			{#if item.privateIp !== ''}
				<div class="grid grid-cols-4 items-center gap-2">
					<Label for="privateip">Private IPs</Label>
					<div class="col-span-3 flex flex-wrap gap-2">
						{#if item.activeIp === item.privateIp}
							<Badge variant="default">{item.privateIp ?? 'None'}</Badge>
						{:else}
							<button onclick={() => handleSubmit(item.privateIp)}>
								<Badge variant="secondary">{item.privateIp}</Badge>
							</button>
						{/if}
					</div>
				</div>
			{/if}

			{#if item.containers?.length > 0}
				<div class="grid grid-cols-4 items-center gap-2">
					<Label for="containers">Containers</Label>
					<div class="col-span-3 flex flex-wrap gap-2">
						{#each item.containers ?? [] as container (container.id)}
							{#if container.name}
								<Badge variant="secondary">
									{typeof container.name === 'string' ? container.name.slice(1) : ''}
								</Badge>
							{/if}
						{/each}
					</div>
				</div>
			{/if}

			{#if item.hostname}
				<div class="grid grid-cols-4 items-center gap-2">
					<Label for="lastseen">Last Seen</Label>
					<div class="col-span-3 flex flex-wrap gap-2">
						<Badge variant="secondary">
							{#if item.updatedAt}
								{DateFormat.format(timestampDate(item.updatedAt))}
							{/if}
						</Badge>
					</div>
				</div>

				<Separator />
				<div class="space-y-1">
					<Label for="ip">Custom IP</Label>
					<Input
						id="ip"
						name="ip"
						type="text"
						bind:value={newIP}
						placeholder="Use a custom IP address"
					/>
				</div>
			{/if}

			<div class="space-y-1">
				<Label for="token">Token</Label>
				<div class="flex w-full items-center gap-1">
					<div class="relative flex w-full">
						<Input id="token" name="token" type="text" value={item.token} class="pr-10" readonly />
						<CopyButton text={item.token} class="absolute right-0" />
					</div>
					<Button
						variant="ghost"
						class="h-10 w-10 cursor-pointer hover:bg-red-300"
						onclick={handleRotate}
					>
						<RotateCcw />
					</Button>
				</div>
			</div>

			<Separator />

			<Button type="button" variant="destructive" class="w-full" onclick={handleDelete}>
				Delete
			</Button>
			{#if item.id && newIP}
				<Button type="submit" class="w-full cursor-pointer" onclick={() => handleSubmit(newIP)}>
					{item.id ? 'Update' : 'Save'}
				</Button>
			{/if}
		</div>
	</Dialog.Content>
</Dialog.Root>
