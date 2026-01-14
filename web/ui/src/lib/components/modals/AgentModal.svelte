<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { toast } from 'svelte-sonner';
	import Separator from '../ui/separator/separator.svelte';
	import { Check, Container, Copy, RotateCcw, Trash2, X } from '@lucide/svelte';
	import { agentClient, settingClient } from '$lib/api';
	import { profile } from '$lib/stores/profile';
	import { ConnectError } from '@connectrpc/connect';
	import type { Agent } from '$lib/gen/mantrae/v1/agent_pb';
	import { UseClipboard } from '$lib/hooks/use-clipboard.svelte';
	import { scale } from 'svelte/transition';
	import ConfirmButton from '../ui/confirm-button/confirm-button.svelte';
	import { CopyInput } from '../ui/input-group';
	import { formatTs } from '$lib/utils';

	interface Props {
		item: Agent;
		open?: boolean;
	}

	let { item = $bindable(), open = $bindable(false) }: Props = $props();

	let newIP = $state('');

	const handleSubmit = async (ip: string | undefined) => {
		try {
			if (item.id) {
				const response = await agentClient.updateAgent({ id: item.id, ip });
				if (!response.agent) throw new Error('Failed to create agent');
				item = response.agent;
				toast.success(`Agent ${item.hostname} updated successfully`);
			} else {
				await agentClient.createAgent({ profileId: profile.id });
				toast.success(`Agent ${item.hostname} created successfully`);
				open = false;
			}
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to save agent', { description: e.message });
		}
	};
	const handleDelete = async () => {
		if (!item.id) return;

		try {
			await agentClient.deleteAgent({ id: item.id });
			toast.success('Agent deleted successfully');
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete agent', { description: e.message });
		}
		open = false;
	};

	const handleRotate = async () => {
		try {
			const response = await agentClient.updateAgent({ id: item.id, rotateToken: true });
			if (!response.agent) throw new Error('Failed to rotate token');
			item.token = response.agent.token;
			toast.success('Token rotated successfully');
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to rotate token', { description: e.message });
		}
	};

	const dockerComposeText = $derived.by(async () => {
		const response = await settingClient.getSetting({ key: 'server_url' });
		const serverURL = response.value ?? 'http://127.0.0.1:3000';
		return `services:
  mantrae-agent:
    image: ghcr.io/mizuchilabs/mantrae-agent:latest
    container_name: mantrae-agent
    network_mode: host
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
    environment:
      - TOKEN=${item.token}
      - HOST=${serverURL}
    restart: unless-stopped`;
	});

	const dockerRunText = $derived.by(async () => {
		const response = await settingClient.getSetting({ key: 'server_url' });
		const serverURL = response.value ?? 'http://127.0.0.1:3000';
		return `docker run -d --name mantrae-agent --network host -v /var/run/docker.sock:/var/run/docker.sock:ro -e TOKEN=${item.token} -e HOST=${serverURL} ghcr.io/mizuchilabs/mantrae-agent:latest`;
	});

	const dockerComposeClipboard = new UseClipboard();
	const dockerRunClipboard = new UseClipboard();

	const handleCopyCompose = async () => {
		await dockerComposeClipboard.copy(await dockerComposeText);
	};

	const handleCopyRun = async () => {
		await dockerRunClipboard.copy(await dockerRunText);
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
								<p class="text-xs text-muted-foreground">Hostname</p>
								<Badge variant="secondary" class="w-full justify-center">
									{item.hostname}
								</Badge>
							</div>
							{#if item.containers?.length > 0}
								<div class="space-y-1">
									<p class="text-xs text-muted-foreground">Containers</p>
									<Badge variant="secondary" class="w-full justify-center">
										{item.containers.length}
									</Badge>
								</div>
							{/if}
							{#if item.updatedAt}
								<div class="space-y-1">
									<p class="text-xs text-muted-foreground">Last Seen</p>
									<Badge variant="secondary" class="justify-center text-xs">
										{formatTs(item.updatedAt, 'relative')}
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
						<p class="text-xs text-muted-foreground">
							Choose which IP address to use for connecting to this agent
						</p>
					</div>

					<div class="space-y-3">
						{#if item.publicIp}
							<div class="flex items-center justify-between rounded-lg border p-3">
								<div class="space-y-1">
									<Label class="text-sm">Public IP</Label>
									<p class="text-xs text-muted-foreground">External network address</p>
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
									<p class="text-xs text-muted-foreground">Internal network address</p>
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
							<p class="text-xs text-muted-foreground">
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
						<p class="text-xs text-muted-foreground">
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
			<div class="flex items-center justify-between gap-2">
				<DropdownMenu.Root>
					<DropdownMenu.Trigger>
						{#snippet child({ props })}
							<Button {...props} variant="outline" class="gap-2">
								<Container class="h-4 w-4" />
								Docker Setup
							</Button>
						{/snippet}
					</DropdownMenu.Trigger>
					<DropdownMenu.Content class="w-80" align="center" side="top" sideOffset={4}>
						<DropdownMenu.Label class="flex items-center gap-2 px-2 py-3">
							<Container class="h-4 w-4" />
							Docker Configuration
						</DropdownMenu.Label>
						<DropdownMenu.Separator />

						<div class="space-y-1 p-1">
							<!-- Docker Compose Item -->
							<DropdownMenu.Item
								closeOnSelect={false}
								class="relative h-auto cursor-pointer flex-col items-start overflow-hidden p-3 transition-all duration-200 hover:bg-accent/50 focus:bg-accent/50"
								onclick={handleCopyCompose}
							>
								{#if dockerComposeClipboard.status === 'success'}
									<div
										in:scale={{ duration: 500, start: 0.85 }}
										class="absolute inset-0 animate-pulse bg-green-500/20"
									></div>
								{:else if dockerComposeClipboard.status === 'failure'}
									<div
										in:scale={{ duration: 500, start: 0.85 }}
										class="absolute inset-0 animate-pulse bg-red-500/20"
									></div>
								{/if}

								<div class="relative z-10 flex w-full items-center justify-between">
									<div class="flex items-center gap-2">
										<Container
											class="h-4 w-4 text-blue-500 transition-transform duration-200 group-hover:scale-110"
										/>
										<span class="font-medium">Docker Compose</span>
									</div>
									<div class="flex items-center gap-1 text-xs text-muted-foreground">
										{#if dockerComposeClipboard.status === 'success'}
											<Check class="h-3 w-3 text-green-500" />
											<span class="text-green-600">Copied!</span>
										{:else if dockerComposeClipboard.status === 'failure'}
											<X class="h-3 w-3 text-red-500" />
											<span class="text-red-600">Failed</span>
										{:else}
											<Copy class="h-3 w-3" />
											Click to copy
										{/if}
									</div>
								</div>
								<p class="relative z-10 mt-1 text-xs leading-relaxed text-muted-foreground">
									Complete docker-compose.yml configuration with volumes and environment setup
								</p>
							</DropdownMenu.Item>

							<!-- Docker Run Item -->
							<DropdownMenu.Item
								closeOnSelect={false}
								class="relative h-auto cursor-pointer flex-col items-start overflow-hidden p-3 transition-all duration-200 hover:bg-accent/50 focus:bg-accent/50"
								onclick={handleCopyRun}
							>
								{#if dockerRunClipboard.status === 'success'}
									<div
										in:scale={{ duration: 500, start: 0.85 }}
										class="absolute inset-0 animate-pulse bg-green-500/20"
									></div>
								{:else if dockerRunClipboard.status === 'failure'}
									<div
										in:scale={{ duration: 500, start: 0.85 }}
										class="absolute inset-0 animate-pulse bg-red-500/20"
									></div>
								{/if}

								<div class="relative z-10 flex w-full items-center justify-between">
									<div class="flex items-center gap-2">
										<Container
											class="h-4 w-4 text-green-500 transition-transform duration-200 group-hover:scale-110"
										/>
										<span class="font-medium">Docker Run</span>
									</div>
									<div class="flex items-center gap-1 text-xs text-muted-foreground">
										{#if dockerRunClipboard.status === 'success'}
											<Check class="h-3 w-3 text-green-500" />
											<span class="text-green-600">Copied!</span>
										{:else if dockerRunClipboard.status === 'failure'}
											<X class="h-3 w-3 text-red-500" />
											<span class="text-red-600">Failed</span>
										{:else}
											<Copy class="h-3 w-3" />
											Click to copy
										{/if}
									</div>
								</div>
								<p class="relative z-10 mt-1 text-xs leading-relaxed text-muted-foreground">
									Single command to run the agent container with all required parameters
								</p>
							</DropdownMenu.Item>
						</div>

						<DropdownMenu.Separator />

						<div class="p-2">
							<div class="rounded-lg border border-blue-200 bg-blue-50 p-3">
								<p class="text-xs text-blue-600 dark:text-blue-800">
									<strong>💡 Quick Start:</strong> Click any option above to copy. Make sure Docker is
									installed and running on your target machine.
								</p>
							</div>
						</div>
					</DropdownMenu.Content>
				</DropdownMenu.Root>

				<ConfirmButton
					title="Delete Agent"
					description="This agent and all associated data will be permanently deleted."
					confirmLabel="Delete"
					cancelLabel="Cancel"
					icon={Trash2}
					class="text-destructive"
					onclick={handleDelete}
				/>
			</div>
		</div>
	</Dialog.Content>
</Dialog.Root>
