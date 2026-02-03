<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import Separator from '../ui/separator/separator.svelte';
	import type { EntryPoint } from '$lib/gen/mantrae/v1/entry_point_pb';
	import CustomSwitch from '../ui/custom-switch/custom-switch.svelte';
	import { entrypoint } from '$lib/api/entrypoints.svelte';

	interface Props {
		data?: EntryPoint;
		open?: boolean;
	}
	let { data, open = $bindable(false) }: Props = $props();

	let epData = $state({} as EntryPoint);
	$effect(() => {
		if (data) epData = { ...data };
	});
	$effect(() => {
		if (!open) epData = {} as EntryPoint;
	});

	const createMutation = entrypoint.create();
	const updateMutation = entrypoint.update();
	function onsubmit() {
		if (epData.id) {
			updateMutation.mutate({ ...epData });
		} else {
			createMutation.mutate({ ...epData });
		}
		open = false;
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="no-scrollbar max-h-[95vh] w-[425px] overflow-y-auto">
		<Dialog.Header>
			<Dialog.Title>{epData.id ? 'Edit' : 'Create'} EntryPoint</Dialog.Title>
			<Dialog.Description>Configure how external traffic reaches your services</Dialog.Description>
		</Dialog.Header>

		<form class="space-y-6" {onsubmit}>
			<!-- Main Configuration -->
			<div class="space-y-4">
				<div class="space-y-2">
					<Label for="name" class="flex items-center gap-2 text-sm font-medium">Name</Label>
					<Input
						id="name"
						bind:value={epData.name}
						placeholder="e.g., web, api, postgres"
						class="transition-colors"
					/>
					<p class="text-xs text-muted-foreground">A descriptive name for this entry point</p>
				</div>

				<div class="space-y-2">
					<Label for="address" class="flex items-center gap-2 text-sm font-medium">Port</Label>
					<Input
						id="address"
						bind:value={epData.address}
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
						<CustomSwitch bind:checked={epData.isDefault} size="md" />
					</div>

					{#if epData.isDefault}
						<div class="rounded-lg border-l-2 border-primary bg-muted/50 p-3">
							<p class="text-xs text-muted-foreground">
								<strong>Note:</strong> Setting this as default will remove the default status from other
								entry points.
							</p>
						</div>
					{/if}
				</div>

				<Separator />

				<Button type="submit" class="w-full">{epData.id ? 'Update' : 'Create'}</Button>
			</div>
		</form>
	</Dialog.Content>
</Dialog.Root>
