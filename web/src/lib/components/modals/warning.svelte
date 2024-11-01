<script lang="ts">
	import * as Popover from '$lib/components/ui/popover';
	import { Button } from '$lib/components/ui/button/index.js';
	import { configError } from '$lib/api';
	import { TriangleAlert } from 'lucide-svelte';
</script>

{#if $configError}
	<Popover.Root>
		<Popover.Trigger>
			<Button variant="ghost" class="text-yellow-500" size="icon">
				<TriangleAlert size="24" />
			</Button>
		</Popover.Trigger>
		<Popover.Content class="text-sm">
			{$configError}
		</Popover.Content>
	</Popover.Root>
{/if}
