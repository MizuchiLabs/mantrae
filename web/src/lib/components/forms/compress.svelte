<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let middleware: Middleware;
	middleware.compress = {
		excludedContentTypes: [],
		includeContentTypes: [],
		defaultEncoding: '',
		...middleware.compress
	};
</script>

{#if middleware.compress}
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="min-response-body-bytes" class="text-right">Min Response Body Bytes</Label>
		<Input
			id="min-response-body-bytes"
			name="min-response-body-bytes"
			type="number"
			bind:value={middleware.compress.minResponseBodyBytes}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="1024"
		/>
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="default-encoding" class="text-right">Default Encoding</Label>
		<Input
			id="default-encoding"
			name="default-encoding"
			type="text"
			bind:value={middleware.compress.defaultEncoding}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="gzip"
		/>
	</div>
	<ArrayInput
		bind:items={middleware.compress.excludedContentTypes}
		label="Excluded Content Types"
		placeholder="text/event-stream"
	/>
	<ArrayInput
		bind:items={middleware.compress.includeContentTypes}
		label="Include Content Types"
		placeholder="application/json"
	/>
{/if}
