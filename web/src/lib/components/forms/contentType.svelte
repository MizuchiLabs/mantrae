<script lang="ts">
	import type { HttpMiddleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';

	export let middleware: HttpMiddleware;
	middleware.contentType = middleware.contentType ?? {
		autoDetect: false
	};
</script>

<div class="grid grid-cols-4 items-center gap-4">
	<Label for="auto-detect" class="text-right">Auto Detect</Label>
	<Switch
		id="auto-detect"
		bind:checked={middleware.contentType.autoDetect}
		class="col-span-3 justify-self-end"
	/>
</div>
