<script lang="ts">
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import { Button, type ButtonProps } from '$lib/components/ui/button/index.js';
	import { type Component } from 'svelte';
	import { AlertTriangle, Ellipsis, type IconProps } from '@lucide/svelte';

	type IconComponent = Component<IconProps, Record<string, never>, ''>;

	type Action = {
		type: 'button' | 'dropdown' | 'popover';
		label: string;
		icon?: IconComponent;
		onClick: () => void;
		variant?: ButtonProps['variant'];
		classProps?: string;
		iconProps?: IconProps;
		disabled?: boolean;
		popover?: PopoverAction;
	};

	interface PopoverAction {
		title?: string;
		description?: string;
		confirmLabel?: string;
		cancelLabel?: string;
	}

	let { actions = [] }: { actions?: Action[] } = $props();

	let popoverOpen = $state(false);

	const inlineActions = actions.filter(
		(a) => ['button', 'popover'].includes(a.type) && !a.disabled
	);
	const dropdownActions = actions.filter((a) => a.type === 'dropdown' && !a.disabled);
</script>

<div class="flex flex-row items-center">
	{#each inlineActions as action (action.label)}
		{#if action.type === 'popover' && action.popover}
			<Popover.Root bind:open={popoverOpen}>
				<Popover.Trigger>
					<Button
						variant={action.variant ?? 'ghost'}
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
				</Popover.Trigger>
				<Popover.Content class="w-80" align="end">
					<div class="space-y-4">
						<div class="flex items-start gap-3">
							<div
								class="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-destructive/10"
							>
								<AlertTriangle class="h-4 w-4 text-destructive" />
							</div>
							<div class="flex-1 space-y-2">
								<h4 class="text-sm leading-none font-semibold">
									{action.popover.title}
								</h4>
								<p class="text-sm leading-relaxed text-muted-foreground">
									{action.popover.description}
								</p>
							</div>
						</div>
						<div class="flex justify-end gap-2">
							<Button variant="outline" size="sm" onclick={() => (popoverOpen = false)}>
								{action.popover.cancelLabel}
							</Button>
							<Button variant="destructive" size="sm" onclick={action.onClick}>
								{action.popover.confirmLabel}
							</Button>
						</div>
					</div>
				</Popover.Content>
			</Popover.Root>
		{:else if action.type === 'button'}
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
		{/if}
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
