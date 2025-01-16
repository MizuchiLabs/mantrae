<script lang="ts">
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { CircleHelp } from 'lucide-svelte';

	interface Props {
		text: string;
	}

	let { text }: Props = $props();
</script>

<Tooltip.Root openDelay={500}>
	<Tooltip.Trigger>
		<CircleHelp size="1rem" />
	</Tooltip.Trigger>
	<Tooltip.Content side="top" align="center" class="max-w-sm">
		{text}
	</Tooltip.Content>
</Tooltip.Root>
