<script lang="ts">
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { Button, type ButtonProps } from '$lib/components/ui/button/index.js';
	import { type Component } from 'svelte';
	import { Ellipsis, type IconProps } from '@lucide/svelte';

	type IconComponent = Component<IconProps, Record<string, never>, ''>;

	type Action = {
		type: 'dropdown' | 'button';
		label: string;
		icon?: IconComponent;
		onClick: () => void;
		variant?: ButtonProps['variant'];
		classProps?: ButtonProps['class'];
		disabled?: boolean;
	};

	let { actions }: { actions?: Action[] } = $props();

	function showDropdown() {
		const hasActions = actions?.some((action) => action.type === 'dropdown') ?? false;
		return hasActions;
	}
</script>

<div class="flex flex-row items-center">
	{#each actions ?? [] as action (action.label)}
		{#if !action.disabled && action.type === 'button'}
			<Tooltip.Provider>
				<Tooltip.Root delayDuration={300}>
					<Tooltip.Trigger>
						<Button
							variant={action.variant ?? 'ghost'}
							onclick={action.onClick}
							class={action.classProps + ' rounded-full'}
							size="icon"
							disabled={action.disabled}
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
		{/if}
	{/each}

	{#if showDropdown()}
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				{#snippet child({ props })}
					<Button {...props} variant="ghost" size="icon">
						<span class="sr-only">Open menu</span>
						<Ellipsis size={16} />
					</Button>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content align="end">
				<DropdownMenu.Group>
					<DropdownMenu.GroupHeading>Actions</DropdownMenu.GroupHeading>
					<DropdownMenu.Separator />
					{#each actions ?? [] as action (action.label)}
						{#if !action.disabled && action.type === 'dropdown' && action.label !== 'Share'}
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
						{/if}
					{/each}
				</DropdownMenu.Group>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	{/if}
</div>
