<script lang="ts">
	import { BookText } from '@lucide/svelte';
	import { Button } from '../ui/button';
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import Separator from '../ui/separator/separator.svelte';

	let version = $state('');
	onMount(async () => {
		const data = await api.getVersion();
		version = data.version;
	});
</script>

<footer
	class="sticky right-0 bottom-0 flex flex-row items-center justify-end bg-transparent px-2 py-1 shadow-[0_-1px_2px_rgba(0,0,0,0.05)]"
>
	<div class="text-muted-foreground flex items-center">
		<Button
			variant="ghost"
			href="https://mizuchi.dev/mantrae/"
			target="_blank"
			rel="noreferrer"
			size="sm"
			class="flex items-center gap-1 text-xs "
		>
			<BookText size={16} />
			Docs
		</Button>
		<Separator orientation="vertical" class="h-5" />
		<Button
			variant="ghost"
			href="https://github.com/mizuchilabs/mantrae"
			target="_blank"
			rel="noreferrer"
			size="sm"
			class="flex items-center gap-1 text-xs "
		>
			Mantrae
			{#if version && version !== 'unknown'}
				{version}
			{:else}
				<span class="italic">latest</span>
			{/if}
		</Button>
	</div>
</footer>
