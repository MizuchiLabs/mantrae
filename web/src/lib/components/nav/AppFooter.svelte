<script lang="ts">
	import { BookText } from 'lucide-svelte';
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
	class="sticky bottom-0 right-0 flex flex-row items-center justify-end bg-transparent px-2 py-1"
>
	<div class="flex items-center text-muted-foreground">
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
				v{version}
			{:else}
				<span class="italic">latest</span>
			{/if}
		</Button>
	</div>
</footer>
