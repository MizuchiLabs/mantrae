<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Collapsible from '$lib/components/ui/collapsible/index.js';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import logo from '$lib/images/logo-white.svg';
	import type { Component, ComponentProps } from 'svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import {
		Blocks,
		ChevronRight,
		ChevronsUpDown,
		Globe,
		House,
		Layers,
		LogOut,
		Plus,
		Route,
		Settings,
		Users,
		Bot,
		Tag,
		Pencil,
		Wrench,
		CircleUserRound,
		Sun,
		Moon,
		Zap,
		type IconProps
	} from '@lucide/svelte';
	import InfoModal from '../modals/info.svelte';
	import ProfileModal from '../modals/profile.svelte';
	import UserModal from '../modals/user.svelte';
	import { profiles, api } from '$lib/api';
	import { slide } from 'svelte/transition';
	import type { Profile } from '$lib/types';
	import { theme } from '$lib/stores/theme';
	import { profile } from '$lib/stores/profile';
	import { user } from '$lib/stores/user';

	let {
		ref = $bindable(null),
		collapsible = 'icon',
		...restProps
	}: ComponentProps<typeof Sidebar.Root> = $props();

	const sidebar = Sidebar.useSidebar();

	type IconComponent = Component<IconProps, Record<string, never>, ''>;

	type Route = {
		title: string;
		url: string;
		icon: IconComponent;
		adminOnly?: boolean;
		subItems?: Route[];
	};
	const routes: Route[] = [
		{ title: 'Overview', url: '/', icon: House },
		{ title: 'Router', url: '/router/', icon: Route },
		{ title: 'Middlewares', url: '/middlewares/', icon: Layers },
		{ title: 'Plugins', url: '/plugins/', icon: Blocks },
		{ title: 'Agents', url: '/agents/', icon: Bot },
		{
			title: 'Settings',
			url: '/settings/',
			icon: Settings,
			subItems: [
				{ title: 'General', url: '/settings/', icon: Wrench, adminOnly: true },
				{ title: 'Users', url: '/users/', icon: Users, adminOnly: true },
				{ title: 'DNS', url: '/dns/', icon: Globe, adminOnly: true }
			]
		}
	];

	// Filter out any routes that the user doesn't have access to
	const usableRoutes = $derived(
		routes
			.filter((route) => !user || !route.adminOnly || user.isAdmin)
			.map((route) => {
				if (route.subItems) {
					const filteredSubs = route.subItems.filter(
						(sub) => !user || !sub.adminOnly || user.isAdmin
					);
					if (filteredSubs.length === 0) return null; // skip this route if no visible subs
					return { ...route, subItems: filteredSubs };
				}
				return route;
			})
			.filter((r): r is Route => r !== null)
	);

	interface ModalState {
		isOpen: boolean;
		profile?: Profile;
	}

	const initialModalState: ModalState = { isOpen: false };
	let modalState = $state(initialModalState);
	let infoModalOpen = $state(false);
	let userModalOpen = $state(false);
</script>

<ProfileModal profile={modalState.profile} bind:open={modalState.isOpen} />
<InfoModal bind:open={infoModalOpen} />

{#if user.isLoggedIn() && user.value}
	<UserModal bind:open={userModalOpen} user={user.value} />
{/if}

<Sidebar.Root bind:ref {collapsible} {...restProps}>
	<!-- Profile Selection -->
	<Sidebar.Header>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<DropdownMenu.Root>
					<DropdownMenu.Trigger>
						{#snippet child({ props })}
							<Sidebar.MenuButton
								{...props}
								size="lg"
								class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
							>
								<div
									class="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg"
								>
									<img src={logo} alt="Mantrae Logo" width="18" />
								</div>
								<div class="grid flex-1 text-left text-sm leading-tight">
									<span class="truncate font-semibold">
										{profile.name ? profile.name : 'Select Profile'}
									</span>
									<span class="truncate text-xs">{profile.value?.url ?? ''}</span>
								</div>
								<ChevronsUpDown class="ml-auto" />
							</Sidebar.MenuButton>
						{/snippet}
					</DropdownMenu.Trigger>
					<DropdownMenu.Content
						class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg"
						align="start"
						side={sidebar.isMobile ? 'bottom' : 'right'}
						sideOffset={4}
					>
						<DropdownMenu.Label class="text-muted-foreground text-xs">Profiles</DropdownMenu.Label>
						{#each $profiles as p (p.name)}
							<DropdownMenu.Item
								onSelect={() => (profile.value = p)}
								class="flex justify-between gap-2"
							>
								<div class="flex items-center gap-2">
									<div class="flex size-6 items-center justify-center rounded-sm border">
										<Tag class="size-4 shrink-0" />
									</div>
									{p.name}
								</div>
								<Button
									variant="secondary"
									class="h-8 w-4 rounded-full"
									onclick={() => (modalState = { isOpen: true, profile: p })}
								>
									<Pencil />
								</Button>
							</DropdownMenu.Item>
						{/each}
						<DropdownMenu.Separator />
						<DropdownMenu.Item class="gap-2 p-2" onSelect={() => (modalState = { isOpen: true })}>
							<div class="bg-background flex size-6 items-center justify-center rounded-md border">
								<Plus class="size-4" />
							</div>
							<div class="text-muted-foreground font-medium">Add Profile</div>
						</DropdownMenu.Item>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.Header>

	<!-- Navigation -->
	<Sidebar.Content>
		<Sidebar.Group>
			<Sidebar.GroupLabel>Dashboard</Sidebar.GroupLabel>
			<Sidebar.Menu>
				{#each usableRoutes as route (route.title)}
					{#if route.subItems}
						<Collapsible.Root class="group/collapsible">
							<Sidebar.MenuItem>
								<Collapsible.Trigger>
									{#snippet child({ props })}
										<Sidebar.MenuButton {...props}>
											{#snippet tooltipContent()}
												{route.title}
											{/snippet}
											{#if route.icon}
												<route.icon />
											{/if}
											<span>{route.title}</span>
											<ChevronRight
												class="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90"
											/>
										</Sidebar.MenuButton>
									{/snippet}
								</Collapsible.Trigger>
								<Collapsible.Content forceMount>
									{#snippet child({ props, open })}
										{#if open}
											<div {...props} transition:slide={{ duration: 200 }}>
												{#each route.subItems || [] as subItem (subItem.title)}
													<Sidebar.MenuSub>
														<Sidebar.MenuSubButton>
															{#snippet child({ props })}
																<a href={subItem.url} {...props}>
																	<subItem.icon />
																	<span>{subItem.title}</span>
																</a>
															{/snippet}
														</Sidebar.MenuSubButton>
													</Sidebar.MenuSub>
												{/each}
											</div>
										{/if}
									{/snippet}
								</Collapsible.Content>
							</Sidebar.MenuItem>
						</Collapsible.Root>
					{:else}
						<Sidebar.MenuItem>
							<Sidebar.MenuButton>
								{#snippet child({ props })}
									<a href={route.url} {...props}>
										<route.icon />
										<span>{route.title}</span>
									</a>
								{/snippet}
							</Sidebar.MenuButton>
						</Sidebar.MenuItem>
					{/if}
				{/each}
			</Sidebar.Menu>
		</Sidebar.Group>

		<!-- Extra buttons (Traefik, etc.) -->
		<Sidebar.Group class="mt-auto">
			<Sidebar.GroupContent>
				<Sidebar.GroupLabel>Status</Sidebar.GroupLabel>
				<Sidebar.Menu>
					{#if $profiles}
						<Sidebar.MenuItem>
							<Sidebar.MenuButton>
								{#snippet child({ props })}
									<button {...props} onclick={() => (infoModalOpen = true)}>
										<Zap />
										<span>Traefik Status</span>
									</button>
								{/snippet}
							</Sidebar.MenuButton>
						</Sidebar.MenuItem>
					{/if}
				</Sidebar.Menu>
			</Sidebar.GroupContent>
		</Sidebar.Group>
	</Sidebar.Content>

	<!-- User Profile -->
	<Sidebar.Footer>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<DropdownMenu.Root>
					<DropdownMenu.Trigger>
						{#snippet child({ props })}
							<Sidebar.MenuButton
								size="lg"
								class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
								{...props}
							>
								<Avatar.Root class="h-8 w-8 rounded-lg">
									<!-- <Avatar.Image src={user.avatar} alt={'@' + $user?.username} /> -->
									<Avatar.Fallback class="rounded-lg">
										{user.username?.slice(0, 2).toUpperCase()}
									</Avatar.Fallback>
								</Avatar.Root>
								<div class="grid flex-1 text-left text-sm leading-tight">
									<span class="truncate font-semibold">{user?.username}</span>
									<span class="truncate text-xs">{user?.email}</span>
								</div>
								<ChevronsUpDown class="ml-auto size-4" />
							</Sidebar.MenuButton>
						{/snippet}
					</DropdownMenu.Trigger>
					<DropdownMenu.Content
						class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg"
						side={sidebar.isMobile ? 'bottom' : 'right'}
						align="end"
						sideOffset={4}
					>
						<DropdownMenu.Label class="p-0 font-normal">
							<div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
								<Avatar.Root class="h-8 w-8 rounded-lg">
									<!-- <Avatar.Image src={user?.avatar} alt={'@' + $user?.username} /> -->
									<Avatar.Fallback class="rounded-lg">
										{user.username?.slice(0, 2).toUpperCase()}
									</Avatar.Fallback>
								</Avatar.Root>
								<div class="grid flex-1 text-left text-sm leading-tight">
									<span class="truncate font-semibold">{user?.username}</span>
									<span class="truncate text-xs">{user?.email}</span>
								</div>
							</div>
						</DropdownMenu.Label>
						<DropdownMenu.Separator />
						<DropdownMenu.Group>
							<DropdownMenu.Item onSelect={() => (userModalOpen = true)}>
								<CircleUserRound />
								Account
							</DropdownMenu.Item>
							<DropdownMenu.Item onSelect={() => theme.toggle()}>
								{#if theme.value === 'dark'}
									<Sun class="size-4" />
									<span>Light Mode</span>
								{:else}
									<Moon class="size-4" />
									<span>Dark Mode</span>
								{/if}
							</DropdownMenu.Item>
						</DropdownMenu.Group>
						<DropdownMenu.Separator />
						<DropdownMenu.Item onSelect={() => api.logout()}>
							<LogOut />
							Log out
						</DropdownMenu.Item>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.Footer>
	<Sidebar.Rail />
</Sidebar.Root>
