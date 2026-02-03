<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import Separator from '../ui/separator/separator.svelte';
	import { type Middleware } from '$lib/gen/mantrae/v1/middleware_pb';
	import HTTPMiddlewareForm from '../forms/HTTPMiddlewareForm.svelte';
	import TCPMiddlewareForm from '../forms/TCPMiddlewareForm.svelte';
	import { ProtocolType } from '$lib/gen/mantrae/v1/protocol_pb';
	import { protocolTypes } from '$lib/types';
	import { middleware } from '$lib/api/middleware.svelte';

	interface Props {
		data?: Middleware;
		open?: boolean;
	}
	let { data, open = $bindable(false) }: Props = $props();

	let mwData = $state({} as Middleware);
	$effect(() => {
		if (data) mwData = { ...data };
	});
	$effect(() => {
		if (!open) mwData = {} as Middleware;
	});

	const createMutation = middleware.create();
	const updateMutation = middleware.update();
	function onsubmit() {
		if (mwData.id) {
			updateMutation.mutate({ ...mwData });
		} else {
			createMutation.mutate({ ...mwData });
		}
		open = false;
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[80vh] max-w-2xl overflow-y-auto">
		<Dialog.Header>
			<Dialog.Title>{mwData.id ? 'Edit' : 'Create'} Middleware</Dialog.Title>
			<Dialog.Description>Configure your Traefik middleware</Dialog.Description>
		</Dialog.Header>

		<form {onsubmit} class="flex flex-col gap-4">
			<div class="grid w-full grid-cols-3 gap-2">
				<div class="col-span-2 flex flex-col gap-2">
					<Label for="name">Name</Label>
					<Input id="name" bind:value={mwData.name} required placeholder="Middleware Name" />
				</div>

				<div class="col-span-1 flex flex-col gap-2">
					<Label for="type" class="text-right">Protocol</Label>
					<Select.Root
						type="single"
						name="type"
						value={mwData.type?.toString()}
						onValueChange={(value) => (mwData.type = parseInt(value, 10))}
					>
						<Select.Trigger class="w-full">
							{protocolTypes.find((t) => t.value === mwData.type)?.label ?? 'Select type'}
						</Select.Trigger>
						<Select.Content>
							<Select.Group>
								<Select.Label>Middleware Type</Select.Label>
								{#each protocolTypes as t (t.value)}
									<!-- Skip UDP -->
									{#if t.value !== ProtocolType.UDP}
										<Select.Item value={t.value.toString()} label={t.label}>
											{t.label}
										</Select.Item>
									{/if}
								{/each}
							</Select.Group>
						</Select.Content>
					</Select.Root>
				</div>
			</div>

			{#if mwData.type === ProtocolType.HTTP}
				<HTTPMiddlewareForm bind:middleware={mwData} />
			{/if}
			{#if mwData.type === ProtocolType.TCP}
				<TCPMiddlewareForm bind:middleware={mwData} />
			{/if}

			<Separator />

			<Button type="submit" class="w-full">{mwData.id ? 'Update' : 'Create'}</Button>
		</form>
	</Dialog.Content>
</Dialog.Root>
