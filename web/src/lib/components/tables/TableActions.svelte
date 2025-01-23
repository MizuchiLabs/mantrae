<script lang="ts">
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { Button, type ButtonProps } from '$lib/components/ui/button/index.js';
	import type { SvelteComponent } from 'svelte';
	import type { IconProps } from 'lucide-svelte';

	type Action = {
		label: string;
		icon?: typeof SvelteComponent<IconProps>;
		onClick: () => void;
		variant?: ButtonProps['variant'];
		classProps?: ButtonProps['class'];
	};

	let { actions }: { actions: Action[] } = $props();
</script>

<div class="flex flex-row items-center gap-2">
	{#each actions as action}
		<Tooltip.Provider>
			<Tooltip.Root delayDuration={300}>
				<Tooltip.Trigger>
					<Button
						variant={action.variant ?? 'ghost'}
						onclick={action.onClick}
						class={action.classProps}
						size="icon"
					>
						{#if action.icon}
							{@const Icon = action.icon}
							<Icon size={16} />
						{:else}
							{action.label}
						{/if}
					</Button>
				</Tooltip.Trigger>
				<Tooltip.Content side="top" align="center" class="max-w-sm">
					{action.label}
				</Tooltip.Content>
			</Tooltip.Root>
		</Tooltip.Provider>
	{/each}
</div>
