<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { toast } from 'svelte-sonner';
	import Separator from '../ui/separator/separator.svelte';
	import { MiddlewareType, type Middleware } from '$lib/gen/mantrae/v1/middleware_pb';
	import { middlewareClient } from '$lib/api';
	import { middlewareTypes } from '$lib/types';
	import { ConnectError } from '@connectrpc/connect';
	import { profile } from '$lib/stores/profile';
	import { pageIndex, pageSize } from '$lib/stores/common';
	import HTTPMiddlewareForm from '../forms/HTTPMiddlewareForm.svelte';
	import TCPMiddlewareForm from '../forms/TCPMiddlewareForm.svelte';

	interface Props {
		data: Middleware[];
		item: Middleware;
		open?: boolean;
	}
	let { data = $bindable(), item = $bindable(), open = $bindable(false) }: Props = $props();

	const handleSubmit = async () => {
		try {
			if (item.id) {
				await middlewareClient.updateMiddleware({
					id: item.id,
					name: item.name,
					config: item.config,
					type: item.type,
					enabled: item.enabled
				});
				toast.success('Middleware updated successfully');
			} else {
				await middlewareClient.createMiddleware({
					profileId: profile.id,
					name: item.name,
					config: item.config,
					type: item.type
				});
				toast.success('Middleware created successfully');
			}

			// Refresh data
			const response = await middlewareClient.listMiddlewares({
				profileId: profile.id,
				limit: BigInt(pageSize.value ?? 10),
				offset: BigInt(pageIndex.value ?? 0)
			});
			data = response.middlewares;
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to save item', { description: e.message });
		}
		open = false;
	};

	const handleDelete = async () => {
		if (!item.id || !item.type) return;

		try {
			await middlewareClient.deleteMiddleware({ id: item.id, type: item.type });
			toast.success('Middleware deleted successfully');

			// Refresh data
			const response = await middlewareClient.listMiddlewares({
				profileId: profile.id,
				limit: BigInt(pageSize.value ?? 10),
				offset: BigInt(pageIndex.value ?? 0)
			});
			data = response.middlewares;
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete item', { description: e.message });
		}
		open = false;
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[80vh] max-w-2xl overflow-y-auto">
		<Dialog.Header>
			<Dialog.Title>{item.id ? 'Edit' : 'Create'} Middleware</Dialog.Title>
			<Dialog.Description>Configure your Traefik middleware</Dialog.Description>
		</Dialog.Header>

		<form class="flex flex-col gap-4">
			<div class="grid w-full grid-cols-3 gap-2">
				<div class="col-span-2 flex flex-col gap-2">
					<Label for="name">Name</Label>
					<Input id="name" bind:value={item.name} required placeholder="Middleware Name" />
				</div>

				<div class="col-span-1 flex flex-col gap-2">
					<Label for="type" class="text-right">Protocol</Label>
					<Select.Root
						type="single"
						name="type"
						value={item.type?.toString()}
						onValueChange={(value) => (item.type = parseInt(value, 10))}
					>
						<Select.Trigger class="w-full">
							{middlewareTypes.find((t) => t.value === item.type)?.label ?? 'Select type'}
						</Select.Trigger>
						<Select.Content>
							<Select.Group>
								<Select.Label>Middleware Type</Select.Label>
								{#each middlewareTypes as t (t.value)}
									<Select.Item value={t.value.toString()} label={t.label}>
										{t.label}
									</Select.Item>
								{/each}
							</Select.Group>
						</Select.Content>
					</Select.Root>
				</div>
			</div>

			{#if item.type === MiddlewareType.HTTP}
				<HTTPMiddlewareForm bind:middleware={item} />
			{/if}
			{#if item.type === MiddlewareType.TCP}
				<TCPMiddlewareForm bind:middleware={item} />
			{/if}

			<Separator />

			<div class="flex w-full flex-row gap-2">
				{#if item.id}
					<Button type="button" variant="destructive" onclick={handleDelete} class="flex-1">
						Delete
					</Button>
				{/if}
				<Button type="submit" class="flex-1" onclick={handleSubmit}>
					{item.id ? 'Update' : 'Create'}
				</Button>
			</div>
		</form>
	</Dialog.Content>
</Dialog.Root>
