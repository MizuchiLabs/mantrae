<script lang="ts">
	import Ellipsis from 'lucide-svelte/icons/ellipsis';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import type { SvelteComponent } from 'svelte';
	import type { IconProps } from 'lucide-svelte';

	type Action = {
		label: string;
		icon?: typeof SvelteComponent<IconProps>;
		onClick: () => void;
		variant?: 'default' | 'destructive';
	};

	let { actions }: { actions: Action[] } = $props();
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger>
		{#snippet child({ props })}
			<Button {...props} variant="ghost" size="icon">
				<span class="sr-only">Open menu</span>
				<Ellipsis class="size-4" />
			</Button>
		{/snippet}
	</DropdownMenu.Trigger>
	<DropdownMenu.Content align="end">
		<DropdownMenu.Group>
			{#each actions as action}
				<DropdownMenu.Item
					onclick={action.onClick}
					class={action.variant === 'destructive' ? 'text-destructive' : ''}
				>
					<div class="flex flex-row items-center justify-between gap-4">
						{#if action.icon}
							{@const Icon = action.icon}
							<Icon size={16} />
						{/if}
						<span>{action.label}</span>
					</div>
				</DropdownMenu.Item>
			{/each}
		</DropdownMenu.Group>
	</DropdownMenu.Content>
</DropdownMenu.Root>
