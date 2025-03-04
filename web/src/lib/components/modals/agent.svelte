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
	import CopyButton from '../ui/copy-button/copy-button.svelte';
	import { DateFormat } from '$lib/stores/common';

	interface Props {
		agent: Agent | undefined;
		open?: boolean;
	}

	let { agent = $bindable({} as Agent), open = $bindable(false) }: Props = $props();

	let newIP = $state('');
	const handleSubmit = async (ip: string | undefined) => {
		if (agent.id) {
			if (!ip) return;
			const params: UpdateAgentIPParams = {
				id: agent.id,
				activeIp: ip
			};
			await api.updateAgentIP(params);
			toast.success(`Agent ${agent.hostname} updated successfully`);
		} else {
			await api.createAgent();
			toast.success(`Agent ${agent.hostname} created successfully`);
		}
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header>
			<Dialog.Title>{agent.hostname ? 'Update' : 'Connect your'} Agent</Dialog.Title>
			<Dialog.Description>
				{agent.hostname
					? 'Update the active IP address for your agent'
					: 'Copy the token for your agent below'}
			</Dialog.Description>
		</Dialog.Header>

		<Separator />

		<form onsubmit={() => handleSubmit(newIP)} class="flex flex-col gap-4">
			{#if agent.hostname}
				<div class="grid grid-cols-4 items-center gap-2">
					<Label for="hostname">Hostname</Label>
					<div class="col-span-3 space-x-2">
						<Badge variant="secondary">{agent.hostname}</Badge>
					</div>
				</div>
			{/if}

			{#if agent.publicIp}
				<div class="grid grid-cols-4 items-center gap-2">
					<Label for="publicip">Public IP</Label>
					<div class="col-span-3 space-x-2">
						{#if agent.activeIp === agent.publicIp || !agent.activeIp}
							<Badge variant="default">{agent.publicIp ?? 'None'}</Badge>
						{:else}
							<button onclick={() => handleSubmit(agent.publicIp)}>
								<Badge variant="secondary">{agent.publicIp}</Badge>
							</button>
						{/if}
					</div>
				</div>
			{/if}

			{#if agent.privateIps?.privateIps?.length > 0}
				<div class="grid grid-cols-4 items-center gap-2">
					<Label for="privateip">Private IPs</Label>
					<div class="col-span-3 flex flex-wrap gap-2">
						{#each agent.privateIps?.privateIps ?? [] as ip}
							{#if agent.activeIp === ip}
								<Badge variant="default">{ip ?? 'None'}</Badge>
							{:else}
								<button onclick={() => handleSubmit(ip)}>
									<Badge variant="secondary">{ip}</Badge>
								</button>
							{/if}
						{/each}
					</div>
				</div>
			{/if}

			{#if agent.containers?.length > 0}
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
			{/if}

			{#if agent.hostname}
				<div class="grid grid-cols-4 items-center gap-2">
					<Label for="lastseen">Last Seen</Label>
					<div class="col-span-3 flex flex-wrap gap-2">
						<Badge variant="secondary">{DateFormat.format(new Date(agent.updatedAt))}</Badge>
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
				<div class="relative flex">
					<Input id="token" name="token" type="text" value={agent.token} class="pr-10" readonly />
					<CopyButton text={agent.token} class="absolute right-0" />
				</div>
			</div>

			{#if agent.hostname}
				<Separator />
				<Button type="submit" class="w-full" disabled={$loading}>
					{agent.id ? 'Update' : 'Save'}
				</Button>
			{/if}
		</form>
	</Dialog.Content>
</Dialog.Root>
