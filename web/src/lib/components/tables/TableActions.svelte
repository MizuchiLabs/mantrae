<script lang="ts">
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { Button, type ButtonProps } from '$lib/components/ui/button/index.js';
	import { type Component } from 'svelte';
	import { Ellipsis, type IconProps } from '@lucide/svelte';

	type IconComponent = Component<IconProps, Record<string, never>, ''>;

	type Action = {
		type: 'button' | 'dropdown';
		label: string;
		icon?: IconComponent;
		onClick: () => void;
		variant?: ButtonProps['variant'];
		classProps?: string;
		iconProps?: IconProps;
		disabled?: boolean;
	};

	let { actions = [] }: { actions?: Action[] } = $props();

	const inlineActions = actions.filter((a) => a.type === 'button' && !a.disabled);
	const dropdownActions = actions.filter((a) => a.type === 'dropdown' && !a.disabled);
</script>

<div class="flex flex-row items-center">
	{#each inlineActions as action (action.label)}
		<Tooltip.Provider>
			<Tooltip.Root delayDuration={300}>
				<Tooltip.Trigger>
					<Button
						variant={action.variant ?? 'ghost'}
						onclick={action.onClick}
						class={`rounded-full ${action.classProps ?? ''}`}
						size={action.icon ? 'icon' : 'sm'}
						disabled={action.disabled}
					>
						{#if action.icon}
							{@const Icon = action.icon}
							<Icon {...action.iconProps} class={`${action.iconProps?.class ?? ''}`} />
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

	{#if dropdownActions.length > 0}
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				<Button variant="ghost" size="icon" class="rounded-full">
					<span class="sr-only">Open menu</span>
					<Ellipsis size={16} />
				</Button>
			</DropdownMenu.Trigger>
			<DropdownMenu.Content align="end">
				<DropdownMenu.Group>
					<DropdownMenu.GroupHeading>Actions</DropdownMenu.GroupHeading>
					<DropdownMenu.Separator />
					{#each dropdownActions as action (action.label)}
						<DropdownMenu.Item
							onclick={action.onClick}
							class={`flex items-center gap-2 ${action.variant === 'destructive' ? 'text-destructive' : ''}`}
						>
							{#if action.icon}
								{@const Icon = action.icon}
								<Icon {...action.iconProps} />
							{/if}
							<span>{action.label}</span>
						</DropdownMenu.Item>
					{/each}
				</DropdownMenu.Group>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	{/if}
</div>
