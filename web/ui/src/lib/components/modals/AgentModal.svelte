<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import Separator from '../ui/separator/separator.svelte';
	import { Check, Container, Copy, RotateCcw, X } from '@lucide/svelte';
	import type { Agent } from '$lib/gen/mantrae/v1/agent_pb';
	import { UseClipboard } from '$lib/hooks/use-clipboard.svelte';
	import { scale } from 'svelte/transition';
	import { CopyInput } from '../ui/input-group';
	import { formatTs } from '$lib/utils';
	import { agent } from '$lib/api/agents.svelte';
	import { setting } from '$lib/api/settings.svelte';

	interface Props {
		data?: Agent;
		open?: boolean;
	}
	let { data, open = $bindable(false) }: Props = $props();

	let agentData = $state({} as Agent);
	$effect(() => {
		if (data) agentData = { ...data };
	});
	$effect(() => {
		if (!open) agentData = {} as Agent;
	});

	const serverURL = setting.get('server_url');
	const createMutation = agent.create();
	const updateMutation = agent.update();
	async function onsubmit() {
		if (agentData.id) {
			updateMutation.mutate({ ...agentData });
		} else {
			createMutation.mutate({});
		}
	}

	const dockerComposeText = $derived.by(async () => {
		return `services:
  mantrae-agent:
    image: ghcr.io/mizuchilabs/mantrae-agent:latest
    container_name: mantrae-agent
    network_mode: host
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
    environment:
      - TOKEN=${agentData.token}
      - HOST=${serverURL.data?.value}
    restart: unless-stopped`;
	});

	const dockerRunText = $derived.by(async () => {
		return `docker run -d --name mantrae-agent --network host -v /var/run/docker.sock:/var/run/docker.sock:ro -e TOKEN=${agentData.token} -e HOST=${serverURL.data?.value} ghcr.io/mizuchilabs/mantrae-agent:latest`;
	});

	const dockerComposeClipboard = new UseClipboard();
	const dockerRunClipboard = new UseClipboard();

	const handleCopyCompose = async () => {
		await dockerComposeClipboard.copy(await dockerComposeText);
	};
	const handleCopyRun = async () => {
		await dockerRunClipboard.copy(await dockerRunText);
	};

	let newAgent = $derived(agent.get(agentData.id));
	function rotate() {
		updateMutation.mutate({ ...agentData, rotateToken: true });
		if (newAgent.isSuccess && newAgent.data?.token) {
			agentData.token = newAgent.data.token;
		}
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[95vh] w-125 overflow-y-auto">
		<Dialog.Header>
			<Dialog.Title>{agentData.hostname ? 'Manage' : 'Connect'} Agent</Dialog.Title>
			<Dialog.Description>
				{agentData.hostname
					? 'Configure agent settings and manage network connections'
					: 'Copy the token below to connect your agent'}
			</Dialog.Description>
		</Dialog.Header>

		<form {onsubmit} class="space-y-6">
			{#if agentData.hostname}
				<!-- Agent Information -->
				<div class="space-y-4">
					<div class="space-y-2">
						<Label class="text-sm font-medium">Agent Information</Label>
						<div class="flex gap-2">
							<div class="space-y-1">
								<p class="text-xs text-muted-foreground">Hostname</p>
								<Badge variant="secondary" class="w-full justify-center">
									{agentData.hostname}
								</Badge>
							</div>
							{#if agentData.containers?.length > 0}
								<div class="space-y-1">
									<p class="text-xs text-muted-foreground">Containers</p>
									<Badge variant="secondary" class="w-full justify-center">
										{agentData.containers.length}
									</Badge>
								</div>
							{/if}
							{#if agentData.updatedAt}
								<div class="space-y-1">
									<p class="text-xs text-muted-foreground">Last Seen</p>
									<Badge variant="secondary" class="justify-center text-xs">
										{formatTs(agentData.updatedAt, 'relative')}
									</Badge>
								</div>
							{/if}
						</div>
					</div>

					{#if agentData.containers?.length > 0}
						<div class="space-y-2">
							<Label class="text-sm font-medium">Running Containers</Label>
							<div class="flex flex-wrap gap-2">
								{#each agentData.containers as container (container.id)}
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
						{#if agentData.publicIp}
							<div class="flex items-center justify-between rounded-lg border p-3">
								<div class="space-y-1">
									<Label class="text-sm">Public IP</Label>
									<p class="text-xs text-muted-foreground">External network address</p>
								</div>
								<div class="flex items-center gap-2">
									{#if agentData.activeIp === agentData.publicIp || !agentData.activeIp}
										<Badge variant="default">Active</Badge>
										<Badge variant="secondary">{agentData.publicIp}</Badge>
									{:else}
										<Button
											variant="outline"
											size="sm"
											onclick={() => (agentData.activeIp = agentData.publicIp)}
										>
											Use {agentData.publicIp}
										</Button>
									{/if}
								</div>
							</div>
						{/if}

						{#if agentData.privateIp}
							<div class="flex items-center justify-between rounded-lg border p-3">
								<div class="space-y-1">
									<Label class="text-sm">Private IP</Label>
									<p class="text-xs text-muted-foreground">Internal network address</p>
								</div>
								<div class="flex items-center gap-2">
									{#if agentData.activeIp === agentData.privateIp}
										<Badge variant="default">Active</Badge>
										<Badge variant="secondary">{agentData.privateIp}</Badge>
									{:else}
										<Button
											variant="outline"
											size="sm"
											onclick={() => (agentData.activeIp = agentData.privateIp)}
										>
											Use {agentData.privateIp}
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
									bind:value={agentData.activeIp}
									placeholder="Enter custom IP address"
									class="flex-1"
								/>
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
				{#if agentData.hostname}
					<div class="space-y-2">
						<Label class="text-sm font-medium">Agent Token</Label>
						<p class="text-xs text-muted-foreground">
							{agentData.hostname
								? 'Secure token for agent authentication'
								: 'Copy this token to connect your agent'}
						</p>
					</div>
				{/if}

				<div class="flex gap-2">
					<CopyInput value={agentData.token} />
					<Button variant="outline" size="icon" onclick={rotate} title="Rotate token">
						<RotateCcw class="h-4 w-4" />
					</Button>
				</div>

				{#if agentData.hostname}
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
			</div>
		</form>
	</Dialog.Content>
</Dialog.Root>
