<script lang="ts">
	import { Button, type ButtonProps } from '$lib/components/ui/button/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import { TriangleAlert, type IconProps } from '@lucide/svelte';
	import type { Component } from 'svelte';

	export type IconComponent = Component<IconProps, Record<string, never>, ''>;
	interface Props extends Omit<ButtonProps, 'href'> {
		title: string;
		description: string;
		confirmLabel: string;
		cancelLabel: string;
		icon?: IconComponent;
		align?: 'start' | 'end' | 'center';
		onclick?: () => void;
	}
	let {
		title,
		description,
		confirmLabel = 'Confirm',
		cancelLabel = 'Cancel',
		variant = 'ghost',
		size = 'sm',
		icon,
		align = 'center',
		onclick,
		...restProps
	}: Props = $props();

	let open = $state(false);
</script>

<Popover.Root bind:open>
	<Popover.Trigger>
		<Button {variant} {size} {...restProps}>
			{#if icon}
				{@const Icon = icon}
				<Icon />
			{:else}
				{title}
			{/if}
		</Button>
	</Popover.Trigger>
	<Popover.Content class="w-80" {align}>
		<div class="space-y-4">
			<div class="flex items-start gap-3">
				<div
					class="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-destructive/10"
				>
					<TriangleAlert class="h-4 w-4 text-destructive" />
				</div>
				<div class="flex-1 space-y-2">
					<h4 class="text-sm leading-none font-semibold">
						{title}
					</h4>
					<p class="text-sm leading-relaxed text-muted-foreground">
						{description}
					</p>
				</div>
			</div>
			<div class="flex justify-end gap-2">
				<Button variant="outline" size="sm" onclick={() => (open = false)}>
					{cancelLabel}
				</Button>
				<Button variant="destructive" size="sm" {onclick}>
					{confirmLabel}
				</Button>
			</div>
		</div>
	</Popover.Content>
</Popover.Root>
