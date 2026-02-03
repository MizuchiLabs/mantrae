<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Select from '$lib/components/ui/select';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Separator } from '$lib/components/ui/separator';
	import { type ServersTransport } from '$lib/gen/mantrae/v1/servers_transport_pb';
	import {
		type TCPServersTransport,
		type ServersTransport as HTTPServersTransport
	} from '$lib/gen/zen/traefik-schemas';
	import HTTPServerTransportForm from '$lib/components/forms/HTTPServerTransportForm.svelte';
	import TCPServerTransportForm from '$lib/components/forms/TCPServerTransportForm.svelte';
	import { ProtocolType } from '$lib/gen/mantrae/v1/protocol_pb';
	import { protocolTypes } from '$lib/types';
	import { transport } from '$lib/api/transport.svelte';

	interface Props {
		data?: ServersTransport;
		open?: boolean;
	}
	let { data, open = $bindable(false) }: Props = $props();

	let transportData = $state({} as ServersTransport);
	$effect(() => {
		if (data) transportData = { ...data };
	});
	$effect(() => {
		if (!open) transportData = {} as ServersTransport;
	});

	const createMutation = transport.create();
	const updateMutation = transport.update();
	function onsubmit() {
		if (transportData.id) {
			updateMutation.mutate({ ...transportData });
		} else {
			createMutation.mutate({ ...transportData });
		}
		open = false;
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[95vh] w-[500px] overflow-y-auto">
		<Dialog.Header>
			<Dialog.Title>{transportData.id ? 'Edit' : 'Create'} ServersTransport</Dialog.Title>
			<Dialog.Description>Configure how external traffic reaches your services</Dialog.Description>
		</Dialog.Header>

		<form {onsubmit} class="space-y-4">
			<div class="grid w-full grid-cols-1 gap-4 sm:grid-cols-3 sm:gap-2">
				<div class="flex flex-col gap-2 {transportData.id ? 'sm:col-span-3' : 'sm:col-span-2'}">
					<Label for="name">Name</Label>
					<Input
						id="name"
						bind:value={transportData.name}
						placeholder="e.g., web, api, postgres"
						class="transition-colors"
					/>
				</div>

				{#if !transportData.id}
					<div class="flex flex-col gap-2 sm:col-span-1">
						<Label for="type">Protocol</Label>
						<Select.Root
							type="single"
							name="type"
							value={transportData.type?.toString()}
							onValueChange={(value) => {
								// Reset config
								transportData.type = parseInt(value, 10);
								switch (transportData.type) {
									case ProtocolType.HTTP:
										transportData.config = {} as HTTPServersTransport;
										break;
									case ProtocolType.TCP:
										transportData.config = {} as TCPServersTransport;
										break;
								}
							}}
						>
							<Select.Trigger class="w-full">
								<span class="truncate">
									{protocolTypes.find((t) => t.value === transportData.type)?.label ?? 'Select'}
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

			{#if transportData.type === ProtocolType.HTTP}
				<HTTPServerTransportForm bind:transport={transportData} />
			{/if}
			{#if transportData.type === ProtocolType.TCP}
				<TCPServerTransportForm bind:transport={transportData} />
			{/if}

			<Separator />

			<Button type="submit" class="w-full">{transportData.id ? 'Update' : 'Create'}</Button>
		</form>
	</Dialog.Content>
</Dialog.Root>
