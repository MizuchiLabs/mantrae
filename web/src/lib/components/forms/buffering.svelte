<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';

	export let middleware: Middleware;
	middleware.buffering = middleware.buffering || {
		maxRequestBodyBytes: 0,
		memRequestBodyBytes: 0,
		maxResponseBodyBytes: 0,
		memResponseBodyBytes: 0,
		retryExpression: ''
	};
</script>

<div class="grid grid-cols-4 items-center gap-4">
	<Label for="max-request-body-bytes" class="text-right">Max Request Body Bytes</Label>
	<Input
		id="max-request-body-bytes"
		name="max-request-body-bytes"
		type="number"
		bind:value={middleware.buffering.maxRequestBodyBytes}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Max Request Body Bytes"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="mem-request-body-bytes" class="text-right">Mem Request Body Bytes</Label>
	<Input
		id="mem-request-body-bytes"
		name="mem-request-body-bytes"
		type="number"
		bind:value={middleware.buffering.memRequestBodyBytes}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Mem Request Body Bytes"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="max-response-body-bytes" class="text-right">Max Response Body Bytes</Label>
	<Input
		id="max-response-body-bytes"
		name="max-response-body-bytes"
		type="number"
		bind:value={middleware.buffering.maxResponseBodyBytes}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Max Response Body Bytes"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="mem-response-body-bytes" class="text-right">Mem Response Body Bytes</Label>
	<Input
		id="mem-response-body-bytes"
		name="mem-response-body-bytes"
		type="number"
		bind:value={middleware.buffering.memResponseBodyBytes}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Mem Response Body Bytes"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="retry-expression" class="text-right">Retry Expression</Label>
	<Input
		id="retry-expression"
		name="retry-expression"
		type="text"
		bind:value={middleware.buffering.retryExpression}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Retry Expression"
	/>
</div>
