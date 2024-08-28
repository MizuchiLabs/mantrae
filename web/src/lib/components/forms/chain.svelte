<script lang="ts">
	import { Label } from '$lib/components/ui/label/index.js';
	import type { Middleware } from '$lib/types/middlewares';
	import * as Select from '$lib/components/ui/select';
	import type { Selected } from 'bits-ui';
	import { middlewares } from '$lib/api';

	export let middleware: Middleware;
	export let disabled = false;
	middleware.chain = { middlewares: [], ...middleware.chain };

	let selectedMiddlewares: Selected<string>[] | undefined = middleware.chain.middlewares?.map(
		(m) => ({ value: m, label: m })
	);
	const changeMiddlewares = (middlewares: Selected<string>[] | undefined) => {
		if (middlewares === undefined || !middleware.chain) return;
		middleware.chain.middlewares = middlewares.map((m) => m.value);
	};
</script>

{#if middleware.chain}
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="middlewares" class="text-right">Middlewares</Label>
		<div class="col-span-3 space-y-2">
			<Select.Root multiple selected={selectedMiddlewares} onSelectedChange={changeMiddlewares}>
				<Select.Trigger>
					<Select.Value placeholder="Select middlewares to chain" />
				</Select.Trigger>
				<Select.Content class="no-scrollbar max-h-[300px] overflow-y-auto">
					{#each $middlewares as m}
						<Select.Item value={m.name} label={m.name}>
							{m.name}
						</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>
		</div>
	</div>
{/if}
