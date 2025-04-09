<script lang="ts">
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { Button, type ButtonProps } from '$lib/components/ui/button/index.js';
	import { type SvelteComponent } from 'svelte';
	import { Ellipsis, Tag, Tags, type IconProps } from 'lucide-svelte';
	import { api, profiles, type RouterWithService } from '$lib/api';
	import type { UpsertRouterParams } from '$lib/types/router';
	import { toast } from 'svelte-sonner';
	import type { Profile } from '$lib/types';
	import type { Middleware, UpsertMiddlewareParams } from '$lib/types/middlewares';
	import { profile } from '$lib/stores/profile';

	interface Props {
		actions?: Action[];
		shareObject?: RouterWithService | Middleware;
	}

	type Action = {
		type: 'dropdown' | 'button';
		label: string;
		icon?: typeof SvelteComponent<IconProps>;
		onClick: () => void;
		variant?: ButtonProps['variant'];
		classProps?: ButtonProps['class'];
		disabled?: boolean;
	};

	let { actions, shareObject }: Props = $props();

	// TODO: Maybe simplify this in the future
	function isRouterWithService(item: RouterWithService | Middleware): item is RouterWithService {
		return 'router' in item && 'service' in item;
	}
	function isMiddleware(item: RouterWithService | Middleware): item is Middleware {
		return 'type' in item;
	}
	async function handleProfileShare(item: RouterWithService | Middleware, profile: Profile) {
		try {
			if (isRouterWithService(item)) {
				let params: UpsertRouterParams = {
					name: item.router.name,
					protocol: item.router.protocol
				};
				switch (item.router.protocol) {
					case 'http':
						params.router = item.router;
						params.service = item.service;
						break;
					case 'tcp':
						params.tcpRouter = item.router;
						params.tcpService = item.service;
						break;
					case 'udp':
						params.udpRouter = item.router;
						params.udpService = item.service;
						break;
				}
				await api.shareRouter(params, profile.id);
				toast.success(`Sent ${item.router.name} routers to profile ${profile.name}`);
			}
			if (isMiddleware(item)) {
				let params: UpsertMiddlewareParams = {
					name: item.name,
					protocol: item.protocol,
					type: item.type
				};
				if (item.protocol === 'http' && item.type) {
					params.middleware = {
						[item.type]: item
					};
				} else if (item.protocol === 'tcp' && item.type) {
					params.tcpMiddleware = {
						[item.type]: item
					};
				}
				await api.shareMiddleware(params, profile.id);
				toast.success(`Sent ${item.name} middleware to profile ${profile.name}`);
				return;
			}
		} catch (err: unknown) {
			const e = err as Error;
			toast.error(`Failed to share ${isMiddleware(item) ? 'middleware' : 'routers'}: ${e.message}`);
		}
	}

	function showDropdown() {
		const hasActions = actions?.some((action) => action.type === 'dropdown') ?? false;
		return hasActions || !!shareObject;
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

					{#if shareObject}
						<DropdownMenu.Sub>
							<DropdownMenu.SubTrigger>
								<Tags class="mr-2 size-4" />
								<span>Send to...</span>
							</DropdownMenu.SubTrigger>
							<DropdownMenu.SubContent>
								{#each $profiles as p (p.id)}
									{#if p.id !== profile.id}
										<DropdownMenu.Item
											onclick={() => {
												if (shareObject) handleProfileShare(shareObject, p);
											}}
										>
											<Tag class="mr-2 size-4" />
											<span>{p.name}</span>
										</DropdownMenu.Item>
									{/if}
								{/each}
							</DropdownMenu.SubContent>
						</DropdownMenu.Sub>
					{/if}
				</DropdownMenu.Group>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	{/if}
</div>
