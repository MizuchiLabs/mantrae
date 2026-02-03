<script lang="ts">
	import * as Select from '$lib/components/ui/select/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { type Router } from '$lib/gen/mantrae/v1/router_pb';
	import type { UDPRouter } from '$lib/gen/zen/traefik-schemas';
	import { Star } from '@lucide/svelte';
	import { unmarshalConfig, marshalConfig } from '$lib/types';
	import { entrypoint } from '$lib/api/entrypoints.svelte';

	interface Props {
		data: Router;
	}
	let { data = $bindable() }: Props = $props();

	let config = $state(unmarshalConfig(data.config) as UDPRouter);
	const epList = entrypoint.list();

	$effect(() => {
		if (config) data.config = marshalConfig(config);
	});
	$effect(() => {
		if (epList.isSuccess && epList.data && !data.id) {
			config.entryPoints = [epList.data.find((ep) => ep.isDefault)?.name ?? ''];
		}
	});
</script>

<div class="flex flex-col gap-3">
	<!-- Entrypoints -->
	<div class="flex flex-col gap-2">
		<Label class="mr-2">Entrypoints</Label>
		<Select.Root type="multiple" bind:value={config.entryPoints}>
			<Select.Trigger class="w-full">
				{config.entryPoints?.join(', ') || 'Select entrypoints'}
			</Select.Trigger>
			<Select.Content>
				{#each epList.data || [] as e (e.id)}
					<Select.Item value={e.name}>
						<div class="flex items-center gap-2">
							{e.name}
							{#if e.isDefault}
								<Star size="1rem" class="text-yellow-300" />
							{/if}
						</div>
					</Select.Item>
				{/each}
			</Select.Content>
		</Select.Root>
	</div>
</div>
