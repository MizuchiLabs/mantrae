<script lang="ts">
	import { BookText } from 'lucide-svelte';
	import { Button } from '../ui/button';
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { writable } from 'svelte/store';

	let version = writable('');
	onMount(async () => {
		const data = await api.getVersion();
		version.set(data);
	});
</script>

<footer
	class="bottom-0 left-16 right-0 flex flex-row items-center justify-end bg-background px-2 py-1"
>
	<div class="flex flex-row items-center divide-x text-xs text-muted-foreground">
		<Button
			variant="ghost"
			href="https://mizuchi.dev/mantrae/"
			target="_blank"
			rel="noreferrer"
			size="sm"
			class="flex items-center gap-1 text-xs"
		>
			<BookText size={16} />
			Docs
		</Button>
		<Button
			variant="ghost"
			href="https://github.com/mizuchilabs/mantrae"
			target="_blank"
			rel="noreferrer"
			size="sm"
			class="flex items-center gap-1 text-xs"
		>
			Mantrae
			{#if $version && $version !== 'unknown'}
				v{$version}
			{:else}
				<span class="italic">latest</span>
			{/if}
		</Button>
	</div>
</footer>
