<script lang="ts">
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { CircleHelp } from 'lucide-svelte';

	interface Props {
		text: string;
	}

	let { text }: Props = $props();
</script>

<Tooltip.Root delayDuration={300}>
	<Tooltip.Trigger>
		<CircleHelp size="1rem" />
	</Tooltip.Trigger>
	<Tooltip.Content side="top" align="center" class="max-w-sm">
		{text}
	</Tooltip.Content>
</Tooltip.Root>
