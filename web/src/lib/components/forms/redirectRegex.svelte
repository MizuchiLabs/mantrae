<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';

	export let middleware: Middleware;
	middleware.redirectRegex = {
		regex: '',
		replacement: '',
		permanent: false,
		...middleware.redirectRegex
	};
</script>

{#if middleware.redirectRegex}
	<div class="grid grid-cols-4 items-center gap-2">
		<Label for="permanent" class="text-right">Permanent</Label>
		<Switch id="permanent" bind:checked={middleware.redirectRegex.permanent} class="col-span-3" />
	</div>
	<div class="grid grid-cols-4 items-center gap-4">
		<Label for="regex" class="text-right">Regex</Label>
		<Input
			id="regex"
			name="regex"
			type="text"
			bind:value={middleware.redirectRegex.regex}
			class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
			placeholder="^http://localhost/(.*)"
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
			placeholder="http://mydomain/$1"
		/>
	</div>
{/if}
