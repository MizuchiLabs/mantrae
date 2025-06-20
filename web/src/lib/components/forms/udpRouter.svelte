<script lang="ts">
	import * as Select from '$lib/components/ui/select/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { type Router } from '$lib/gen/mantrae/v1/router_pb';
	import type { UDPRouter } from '$lib/gen/tygo/dynamic';
	import { Star } from '@lucide/svelte';
	import { entryPointClient } from '$lib/api';
	import { unmarshalConfig, marshalConfig } from '$lib/types';

	let { router = $bindable() }: { router: Router } = $props();

	let config = $state(unmarshalConfig(router.config) as UDPRouter);

	$effect(() => {
		if (config) router.config = marshalConfig(config);
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
				{#await entryPointClient.listEntryPoints( { profileId: router.profileId, limit: -1n, offset: 0n } ) then value}
					{#each value.entryPoints as e (e.id)}
						<Select.Item value={e.name}>
							<div class="flex items-center gap-2">
								{e.name}
								{#if e.isDefault}
									<Star size="1rem" class="text-yellow-300" />
								{/if}
							</div>
						</Select.Item>
					{/each}
				{/await}
			</Select.Content>
		</Select.Root>
	</div>
</div>
