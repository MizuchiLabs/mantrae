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
	import CopyInput from '../ui/copy-input/copy-input.svelte';

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
	};
	const handleDelete = async () => {
		if (!item.id) return;

		try {
			await agentClient.deleteAgent({ id: item.id });
			data = data.filter((e) => e.id !== item.id);
			toast.success('Agent deleted successfully');
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete agent', { description: e.message });
		}
		open = false;
	};

	const handleRotate = async () => {
		try {
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
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to rotate token', { description: e.message });
		}
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[95vh] w-[500px] overflow-y-auto">
		<Dialog.Header>
			<Dialog.Title>{item.hostname ? 'Manage' : 'Connect'} Agent</Dialog.Title>
			<Dialog.Description>
				{item.hostname
					? 'Configure agent settings and manage network connections'
					: 'Copy the token below to connect your agent'}
			</Dialog.Description>
		</Dialog.Header>

		<div class="space-y-4">
			{#if item.hostname}
				<!-- Agent Information -->
				<div class="space-y-4">
					<div class="space-y-2">
						<Label class="text-sm font-medium">Agent Information</Label>
						<div class="flex gap-2">
							<div class="space-y-1">
								<p class="text-muted-foreground text-xs">Hostname</p>
								<Badge variant="secondary" class="w-full justify-center">
									{item.hostname}
								</Badge>
							</div>
							{#if item.containers?.length > 0}
								<div class="space-y-1">
									<p class="text-muted-foreground text-xs">Containers</p>
									<Badge variant="secondary" class="w-full justify-center">
										{item.containers.length}
									</Badge>
								</div>
							{/if}
							{#if item.updatedAt}
								<div class="space-y-1">
									<p class="text-muted-foreground text-xs">Last Seen</p>
									<Badge variant="secondary" class="justify-center text-xs">
										{DateFormat.format(timestampDate(item.updatedAt))}
									</Badge>
								</div>
							{/if}
						</div>
					</div>

					{#if item.containers?.length > 0}
						<div class="space-y-2">
							<Label class="text-sm font-medium">Running Containers</Label>
							<div class="flex flex-wrap gap-2">
								{#each item.containers as container (container.id)}
									{#if container.name}
										<Badge variant="outline" class="text-xs">
											{typeof container.name === 'string' ? container.name.slice(1) : ''}
										</Badge>
									{/if}
								{/each}
							</div>
						</div>
					{/if}
				</div>

				<Separator />

				<!-- Network Configuration -->
				<div class="space-y-4">
					<div class="space-y-2">
						<Label class="text-sm font-medium">Network Configuration</Label>
						<p class="text-muted-foreground text-xs">
							Choose which IP address to use for connecting to this agent
						</p>
					</div>

					<div class="space-y-3">
						{#if item.publicIp}
							<div class="flex items-center justify-between rounded-lg border p-3">
								<div class="space-y-1">
									<Label class="text-sm">Public IP</Label>
									<p class="text-muted-foreground text-xs">External network address</p>
								</div>
								<div class="flex items-center gap-2">
									{#if item.activeIp === item.publicIp || !item.activeIp}
										<Badge variant="default">Active</Badge>
										<Badge variant="secondary">{item.publicIp}</Badge>
									{:else}
										<Button variant="outline" size="sm" onclick={() => handleSubmit(item.publicIp)}>
											Use {item.publicIp}
										</Button>
									{/if}
								</div>
							</div>
						{/if}

						{#if item.privateIp}
							<div class="flex items-center justify-between rounded-lg border p-3">
								<div class="space-y-1">
									<Label class="text-sm">Private IP</Label>
									<p class="text-muted-foreground text-xs">Internal network address</p>
								</div>
								<div class="flex items-center gap-2">
									{#if item.activeIp === item.privateIp}
										<Badge variant="default">Active</Badge>
										<Badge variant="secondary">{item.privateIp}</Badge>
									{:else}
										<Button
											variant="outline"
											size="sm"
											onclick={() => handleSubmit(item.privateIp)}
										>
											Use {item.privateIp}
										</Button>
									{/if}
								</div>
							</div>
						{/if}

						<!-- Custom IP -->
						<div class="space-y-2">
							<Label for="customip" class="text-sm font-medium">Custom IP Address</Label>
							<div class="flex gap-2">
								<Input
									id="customip"
									bind:value={newIP}
									placeholder="Enter custom IP address"
									class="flex-1"
								/>
								{#if newIP}
									<Button onclick={() => handleSubmit(newIP)} size="sm">Use</Button>
								{/if}
							</div>
							<p class="text-muted-foreground text-xs">
								Specify a custom IP address for this agent
							</p>
						</div>
					</div>
				</div>

				<Separator />
			{/if}

			<!-- Token Management -->
			<div class="space-y-2">
				{#if item.hostname}
					<div class="space-y-2">
						<Label class="text-sm font-medium">Agent Token</Label>
						<p class="text-muted-foreground text-xs">
							{item.hostname
								? 'Secure token for agent authentication'
								: 'Copy this token to connect your agent'}
						</p>
					</div>
				{/if}

				<div class="flex gap-2">
					<CopyInput value={item.token} />
					<Button variant="outline" size="icon" onclick={handleRotate} title="Rotate token">
						<RotateCcw class="h-4 w-4" />
					</Button>
				</div>

				{#if item.hostname}
					<div class="rounded-lg border border-amber-200 bg-amber-50 p-3">
						<p class="text-xs text-amber-800">
							<strong>Warning:</strong> Rotating the token will invalidate the current token and require
							updating your agent configuration.
						</p>
					</div>
				{/if}
			</div>

			<Separator />

			<!-- Actions -->
			<div class="flex gap-2">
				<Button type="button" variant="destructive" onclick={handleDelete} class="flex-1">
					Delete Agent
				</Button>
			</div>
		</div>
	</Dialog.Content>
</Dialog.Root>
