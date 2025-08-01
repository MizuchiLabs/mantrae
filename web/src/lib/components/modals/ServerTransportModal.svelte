<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { toast } from 'svelte-sonner';
	import Separator from '../ui/separator/separator.svelte';
	import { serversTransportClient } from '$lib/api';
	import { ConnectError } from '@connectrpc/connect';
	import { profile } from '$lib/stores/profile';
	import { pageIndex, pageSize } from '$lib/stores/common';
	import { type ServersTransport } from '$lib/gen/mantrae/v1/servers_transport_pb';
	import {
		type TCPServersTransport,
		type ServersTransport as HTTPServersTransport
	} from '$lib/gen/zen/traefik-schemas';
	import HTTPServerTransportForm from '$lib/components/forms/HTTPServerTransportForm.svelte';
	import TCPServerTransportForm from '$lib/components/forms/TCPServerTransportForm.svelte';
	import { ProtocolType } from '$lib/gen/mantrae/v1/protocol_pb';
	import { protocolTypes } from '$lib/types';

	interface Props {
		item: ServersTransport;
		open?: boolean;
	}

	let { item = $bindable(), open = $bindable(false) }: Props = $props();

	const handleSubmit = async () => {
		try {
			if (item.id) {
				await serversTransportClient.updateServersTransport({
					id: item.id,
					name: item.name,
					config: item.config,
					type: item.type,
					enabled: item.enabled
				});
				toast.success('Transport updated successfully');
			} else {
				await serversTransportClient.createServersTransport({
					profileId: profile.id,
					name: item.name,
					config: item.config,
					type: item.type
				});
				toast.success('Transport created successfully');
			}
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to save transport', { description: e.message });
		}
		open = false;
	};

	const handleDelete = async () => {
		if (!item.id) return;

		try {
			await serversTransportClient.deleteServersTransport({ id: item.id });
			toast.success('Transport deleted successfully');
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete transport', { description: e.message });
		}
		open = false;
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[95vh] w-[500px] overflow-y-auto">
		<Dialog.Header>
			<Dialog.Title>{item.id ? 'Edit' : 'Create'} ServersTransport</Dialog.Title>
			<Dialog.Description>Configure how external traffic reaches your services</Dialog.Description>
		</Dialog.Header>

		<form onsubmit={handleSubmit} class="space-y-4">
			<div class="grid w-full grid-cols-1 gap-4 sm:grid-cols-3 sm:gap-2">
				<div class="flex flex-col gap-2 {item.id ? 'sm:col-span-3' : 'sm:col-span-2'}">
					<Label for="name">Name</Label>
					<Input
						id="name"
						bind:value={item.name}
						placeholder="e.g., web, api, postgres"
						class="transition-colors"
					/>
				</div>

				{#if !item.id}
					<div class="flex flex-col gap-2 sm:col-span-1">
						<Label for="type">Protocol</Label>
						<Select.Root
							type="single"
							name="type"
							value={item.type?.toString()}
							onValueChange={(value) => {
								// Reset config
								item.type = parseInt(value, 10);
								switch (item.type) {
									case ProtocolType.HTTP:
										item.config = {} as HTTPServersTransport;
										break;
									case ProtocolType.TCP:
										item.config = {} as TCPServersTransport;
										break;
								}
							}}
						>
							<Select.Trigger class="w-full">
								<span class="truncate">
									{protocolTypes.find((t) => t.value === item.type)?.label ?? 'Select'}
								</span>
							</Select.Trigger>
							<Select.Content>
								{#each protocolTypes as t (t.value)}
									<!-- Skip UDP -->
									{#if t.value !== ProtocolType.UDP}
										<Select.Item value={t.value.toString()} label={t.label}>
											{t.label}
										</Select.Item>
									{/if}
								{/each}
							</Select.Content>
						</Select.Root>
					</div>
				{/if}
			</div>

			{#if item.type === ProtocolType.HTTP}
				<HTTPServerTransportForm bind:transport={item} />
			{/if}
			{#if item.type === ProtocolType.TCP}
				<TCPServerTransportForm bind:transport={item} />
			{/if}

			<Separator />

			<div class="flex w-full flex-col gap-2 sm:flex-row">
				{#if item.id}
					<Button type="button" variant="destructive" onclick={handleDelete} class="flex-1">
						Delete
					</Button>
				{/if}
				<Button type="submit" class="flex-1 text-sm">
					{item.id ? 'Update' : 'Create'}
				</Button>
			</div>
		</form>
	</Dialog.Content>
</Dialog.Root>
