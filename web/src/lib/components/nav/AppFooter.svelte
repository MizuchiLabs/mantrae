<script lang="ts">
	import { BookText } from '@lucide/svelte';
	import { Button } from '../ui/button';
	import { utilClient } from '$lib/api';
	import Separator from '../ui/separator/separator.svelte';
</script>

<footer
	class="sticky right-0 bottom-0 flex flex-row items-center justify-end rounded-b-xl bg-background px-2 py-1 shadow-[0_-1px_2px_rgba(0,0,0,0.05)]"
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
			{#await utilClient.getVersion({}) then result}
				{#if result.version && result.version !== 'unknown'}
					{result.version}
				{:else}
					<span class="italic">latest</span>
				{/if}
			{/await}
		</Button>
	</div>
</footer>
