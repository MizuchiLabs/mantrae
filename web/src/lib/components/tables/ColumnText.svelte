<script lang="ts">
	import { Label } from '$lib/components/ui/label/index.js';
	import { type IconProps } from '@lucide/svelte';
	import { type Component } from 'svelte';

	type IconComponent = Component<IconProps, Record<string, never>, ''>;

	type Props = {
		label: string;
		icon?: IconComponent;
		iconProps?: IconProps;
	};

	let { label, icon, iconProps, ...restProps }: Props = $props();
</script>

<div class="flex flex-row items-center gap-2">
	<Label for={label} {...restProps}>{label}</Label>
	{#if icon}
		{@const Icon = icon}
		<Icon {...iconProps} />
	{/if}
</div>
