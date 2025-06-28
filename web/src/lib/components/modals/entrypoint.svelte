<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { toast } from 'svelte-sonner';
	import Separator from '../ui/separator/separator.svelte';
	import { entryPointClient } from '$lib/api';
	import { ConnectError } from '@connectrpc/connect';
	import type { EntryPoint } from '$lib/gen/mantrae/v1/entry_point_pb';
	import { profile } from '$lib/stores/profile';
	import { pageIndex, pageSize } from '$lib/stores/common';

	interface Props {
		data: EntryPoint[];
		item: EntryPoint;
		open?: boolean;
	}

	let { data = $bindable(), item = $bindable(), open = $bindable(false) }: Props = $props();

	const handleSubmit = async () => {
		try {
			if (item.id) {
				await entryPointClient.updateEntryPoint({
					id: item.id,
					name: item.name,
					address: item.address,
					isDefault: item.isDefault
				});
				toast.success('EntryPoint updated successfully');
			} else {
				await entryPointClient.createEntryPoint({
					profileId: profile.id,
					name: item.name,
					address: item.address,
					isDefault: item.isDefault
				});
				toast.success('EntryPoint created successfully');
			}

			// Refresh data
			let response = await entryPointClient.listEntryPoints({
				profileId: profile.id,
				limit: BigInt(pageSize.value ?? 10),
				offset: BigInt(pageIndex.value ?? 0)
			});
			data = response.entryPoints;
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to save entry point', { description: e.message });
		}
		open = false;
	};

	const handleDelete = async () => {
		if (!item.id) return;

		try {
			await entryPointClient.deleteEntryPoint({ id: item.id });
			toast.success('EntryPoint deleted successfully');

			// Refresh data
			let response = await entryPointClient.listEntryPoints({
				profileId: profile.id,
				limit: BigInt(pageSize.value ?? 10),
				offset: BigInt(pageIndex.value ?? 0)
			});
			data = response.entryPoints;
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete entry point', { description: e.message });
		}
		open = false;
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[95vh] w-[425px] overflow-y-auto">
		<Dialog.Header class="flex flex-row items-center justify-between">
			<div>
				<Dialog.Title>{item.id ? 'Edit' : 'Create'} EntryPoint</Dialog.Title>
				<Dialog.Description>Configure your entry point settings</Dialog.Description>
			</div>
			<div class="mr-4 flex items-center gap-2">
				<Label for="default">Default</Label>
				<Switch
					id="default"
					checked={item.isDefault}
					onCheckedChange={(value) => (item.isDefault = value)}
				/>
			</div>
		</Dialog.Header>

		<form class="flex flex-col gap-4">
			<div class="flex flex-col gap-2">
				<Label for="name">Name</Label>
				<Input id="name" bind:value={item.name} required placeholder="web" />
			</div>

			<div class="flex flex-col gap-2">
				<Label for="address">Port</Label>
				<Input id="address" bind:value={item.address} required placeholder="80" />
			</div>

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
