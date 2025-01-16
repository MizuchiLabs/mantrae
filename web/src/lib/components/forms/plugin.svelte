<script lang="ts">
	import { run } from 'svelte/legacy';

	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';

	interface Props {
		middleware: Middleware;
		disabled?: boolean;
	}

	let { middleware = $bindable(), disabled = false }: Props = $props();
	let pluginData = $state('{}');

	function extractInnerPluginData() {
		if (!middleware.content) return;
		pluginData = JSON.stringify(middleware.content, null, 2) || '{}';
	}

	run(() => {
		middleware.content, extractInnerPluginData();
	});
	let error = $state('');
	function validateJSON() {
		if (!pluginData || !middleware.content) return;
		try {
			JSON.parse(pluginData);
			middleware.content = JSON.parse(pluginData);
			error = '';
		} catch (e: any) {
			error = e;
		}
	}
</script>

<div class="grid grid-cols-8 items-center gap-2">
	<Label for="config" class="text-right">Config</Label>
	<Textarea
		id="config"
		name="config"
		rows={pluginData ? pluginData.split('\n').length + 1 : 3}
		bind:value={pluginData}
		on:input={validateJSON}
		class="col-span-7 max-h-[500px] overflow-y-auto"
		{disabled}
	/>
	{#if error}
		<div class="col-span-4 text-right text-sm text-red-500">{error}</div>
	{/if}
</div>
