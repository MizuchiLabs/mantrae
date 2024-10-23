<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';

	export let middleware: Middleware;
	export let disabled = false;
	middleware.plugin = { ...middleware.plugin };

	function extractInnerPluginData() {
		if (!middleware.plugin) return;
		const outerKey = Object.keys(middleware.plugin)[0];
		const innerKey = Object.keys(middleware.plugin[outerKey])[0];

		const data = middleware.plugin[outerKey][innerKey];
		return JSON.stringify(data, null, 2);
	}

	$: pluginData = extractInnerPluginData() || '{}';
	let error = '';
	function validateJSON() {
		if (!pluginData) return;
		try {
			JSON.parse(pluginData);
		} catch (e: any) {
			error = e;
		}
	}
</script>

{#if middleware.plugin}
	<div class="grid grid-cols-4 items-center gap-2">
		<Label for="config" class="text-right">Config</Label>
		<Textarea
			id="config"
			name="config"
			rows={pluginData ? pluginData.split('\n').length + 1 : 3}
			bind:value={pluginData}
			on:input={validateJSON}
			class="col-span-3"
			{disabled}
		/>
		{#if error}
			<div class="col-span-4 text-right text-sm text-red-500">{error}</div>
		{/if}
	</div>
{/if}
