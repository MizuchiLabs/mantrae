<script lang="ts">
	import type { HttpMiddleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let middleware: HttpMiddleware;
	middleware.compress = middleware.compress ?? {
		excludedContentTypes: [''],
		minResponseBodyBytes: 0
	};
</script>

<div class="grid grid-cols-4 items-center gap-4">
	<Label for="min-response-body-bytes" class="text-right">Min Response Body Bytes</Label>
	<Input
		id="min-response-body-bytes"
		name="min-response-body-bytes"
		type="number"
		bind:value={middleware.compress.minResponseBodyBytes}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Min Response Body Bytes"
	/>
</div>
<ArrayInput bind:items={middleware.compress.excludedContentTypes} label="Excluded Content Types" />
