<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { toast } from 'svelte-sonner';
	import Separator from '../ui/separator/separator.svelte';
	import { entryPointClient } from '$lib/api';
	import { ConnectError } from '@connectrpc/connect';
	import type { EntryPoint } from '$lib/gen/mantrae/v1/entry_point_pb';
	import { profile } from '$lib/stores/profile';
	import CustomSwitch from '../ui/custom-switch/custom-switch.svelte';

	interface Props {
		item: EntryPoint;
		open?: boolean;
	}

	let { item = $bindable(), open = $bindable(false) }: Props = $props();

	const handleSubmit = async () => {
		try {
			if (item.id) {
				await entryPointClient.updateEntryPoint({
					id: item.id,
					profileId: item.profileId,
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
		} catch (err) {
			const e = ConnectError.from(err);
			toast.error('Failed to delete entry point', { description: e.message });
		}
		open = false;
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[95vh] w-[425px] overflow-y-auto">
		<Dialog.Header>
			<Dialog.Title>{item.id ? 'Edit' : 'Create'} EntryPoint</Dialog.Title>
			<Dialog.Description>Configure how external traffic reaches your services</Dialog.Description>
		</Dialog.Header>

		<form class="space-y-6" onsubmit={handleSubmit}>
			<!-- Main Configuration -->
			<div class="space-y-4">
				<div class="space-y-2">
					<Label for="name" class="flex items-center gap-2 text-sm font-medium">Name</Label>
					<Input
						id="name"
						bind:value={item.name}
						placeholder="e.g., web, api, postgres"
						class="transition-colors"
					/>
					<p class="text-xs text-muted-foreground">A descriptive name for this entry point</p>
				</div>

				<div class="space-y-2">
					<Label for="address" class="flex items-center gap-2 text-sm font-medium">Port</Label>
					<Input
						id="address"
						bind:value={item.address}
						placeholder="80, 443, 8080..."
						min="1"
						max="65535"
						class="transition-colors"
					/>
					<div class="flex items-center justify-between">
						<p class="text-xs text-muted-foreground">
							Port number (1-65535) where your service listens
						</p>
					</div>
				</div>

				<!-- Default Setting -->
				<div class="space-y-3">
					<div class="flex items-center justify-between">
						<div class="space-y-1">
							<Label class="flex items-center gap-2 text-sm font-medium">Default Entry Point</Label>
							<p class="text-xs text-muted-foreground">
								Use this as the primary entry point for new routers
							</p>
						</div>
						<CustomSwitch bind:checked={item.isDefault} size="md" />
					</div>

					{#if item.isDefault}
						<div class="rounded-lg border-l-2 border-primary bg-muted/50 p-3">
							<p class="text-xs text-muted-foreground">
								<strong>Note:</strong> Setting this as default will remove the default status from other
								entry points.
							</p>
						</div>
					{/if}
				</div>

				<Separator />

				<div class="flex w-full flex-row gap-2">
					{#if item.id}
						<Button type="button" variant="destructive" onclick={handleDelete} class="flex-1">
							Delete
						</Button>
					{/if}
					<Button type="submit" class="flex-1">
						{item.id ? 'Update' : 'Create'}
					</Button>
				</div>
			</div>
		</form>
	</Dialog.Content>
</Dialog.Root>
