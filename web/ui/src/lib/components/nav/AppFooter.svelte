<script lang="ts">
	import { util } from '$lib/api/util.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Separator } from '$lib/components/ui/separator';
	import { DocsURL } from '$lib/config';
	import { BookText } from '@lucide/svelte';

	const version = util.version();
</script>

<footer
	class="sticky right-0 bottom-0 flex flex-row items-center justify-end rounded-b-xl bg-background px-2 py-1 shadow-[0_-1px_2px_rgba(0,0,0,0.05)]"
>
	<div class="flex items-center text-muted-foreground">
		<Button
			variant="ghost"
			href={DocsURL}
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
			{#if version.isSuccess && version.data !== 'debug'}
				{version.data}
			{:else}
				<span class="italic">latest</span>
			{/if}
		</Button>
	</div>
</footer>
