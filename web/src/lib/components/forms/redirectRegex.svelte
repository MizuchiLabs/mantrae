<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';

	export let middleware: Middleware;
	middleware.redirectRegex = middleware.redirectRegex ?? {
		regex: '',
		replacement: ''
	};
</script>

<div class="grid grid-cols-4 items-center gap-4">
	<Label for="regex" class="text-right">Regex</Label>
	<Input
		id="regex"
		name="regex"
		type="text"
		bind:value={middleware.redirectRegex.regex}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Regex"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="replacement" class="text-right">Replacement</Label>
	<Input
		id="replacement"
		name="replacement"
		type="text"
		bind:value={middleware.redirectRegex.replacement}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Replacement"
	/>
</div>
