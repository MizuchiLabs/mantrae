<script lang="ts">
	import type { ComponentProps } from 'svelte';
	import { Badge } from '$lib/components/ui/badge';

	type Props = ComponentProps<typeof Badge> & {
		label: string;
	};
	let { variant = 'default', label, ...restProps }: Props = $props();
</script>

<Badge {variant} {...restProps}>{label}</Badge>
