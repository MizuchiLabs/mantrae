<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';

	interface Props {
		middleware: Middleware;
		disabled?: boolean;
	}

	let { middleware = $bindable(), disabled = false }: Props = $props();
</script>

<div class="grid grid-cols-4 items-center gap-4">
	<Label for="path" class="text-right">Path</Label>
	<Input
		id="path"
		name="path"
		type="text"
		bind:value={middleware.content.path}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="/foo"
		{disabled}
	/>
</div>
