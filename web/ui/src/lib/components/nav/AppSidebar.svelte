<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { LogoLight } from '$lib/assets';
	import type { Component, ComponentProps } from 'svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import {
		ChevronsUpDown,
		Globe,
		Layers,
		LogOut,
		Plus,
		Route,
		Settings,
		Users,
		Bot,
		Pencil,
		CircleUserRound,
		Sun,
		Moon,
		type IconProps,
		EthernetPort,
		Gauge,
		Layers2,
		Truck,
		Cog
	} from '@lucide/svelte';
	import type { Profile } from '$lib/gen/mantrae/v1/profile_pb';
	import ProfileModal from '$lib/components/modals/ProfileModal.svelte';
	import UserModal from '$lib/components/modals/UserModal.svelte';
	import { toggleMode, mode } from 'mode-watcher';
	import RandomAvatar from '../utils/RandomAvatar.svelte';
	import { BackendURL } from '$lib/config';
	import { profile } from '$lib/api/profiles.svelte';
	import { user } from '$lib/api/users.svelte';
	import { profileID } from '$lib/store.svelte';

	let { ...restProps }: ComponentProps<typeof Sidebar.Root> = $props();

	const sidebar = Sidebar.useSidebar();

	type IconComponent = Component<IconProps, Record<string, never>, ''>;

	type Route = {
		title: string;
		url: string;
		icon: IconComponent;
		adminOnly?: boolean;
		subItems?: Route[];
	};
	const mainRoutes: Route[] = [
		{ title: 'Dashboard', url: '/', icon: Gauge },
		{ title: 'Routers', url: '/router/', icon: Route },
		{ title: 'Middlewares', url: '/middlewares/', icon: Layers },
		{ title: 'Entry Points', url: '/entrypoints/', icon: EthernetPort },
		{ title: 'Server Transports', url: '/transport/', icon: Truck }
	];
	const adminRoutes: Route[] = [
		{ title: 'Users', url: '/users/', icon: Users },
		{ title: 'Agents', url: '/agents/', icon: Bot },
		{ title: 'DNS', url: '/dns/', icon: Globe },
		{ title: 'Settings', url: '/settings/', icon: Settings }
	];
	const supportRoutes: Route[] = [
		{ title: 'API Reference', url: `${BackendURL}/openapi`, icon: Cog }
	];

	const profileList = profile.list();
	const currentProfile = $derived(profile.get());
	const currentUser = $derived(user.self());
	const logout = user.logout();

	let modalProfile = $state({} as Profile);
	let modalProfileOpen = $state(false);
	let modalUserOpen = $state(false);
</script>

<ProfileModal data={modalProfile} bind:open={modalProfileOpen} />

{#if currentUser.data}
	<UserModal bind:open={modalUserOpen} data={currentUser.data} />
{/if}

<Sidebar.Root collapsible="offcanvas" {...restProps}>
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
									class="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground"
								>
									<LogoLight class="size-5" />
								</div>
								<div class="grid flex-1 text-left text-sm leading-tight">
									<span class="truncate font-semibold">
										{currentProfile.data?.name ? currentProfile.data.name : 'Select Profile'}
									</span>
									<span class="truncate text-xs">{currentProfile.data?.description ?? ''}</span>
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
						<DropdownMenu.Label class="text-xs text-muted-foreground">Profiles</DropdownMenu.Label>
						{#each profileList.data || [] as p (p.id)}
							<DropdownMenu.Item
								onSelect={() => (profileID.current = p.id)}
								class="flex justify-between gap-2"
							>
								<div class="flex items-center gap-2">
									<Layers2 class="size-4 shrink-0" />
									{p.name}
								</div>
								<Button
									variant="outline"
									class="rounded-full"
									size="sm"
									onclick={() => {
										modalProfile = p;
										modalProfileOpen = true;
									}}
								>
									<Pencil />
									Edit
								</Button>
							</DropdownMenu.Item>
						{/each}
						<DropdownMenu.Separator />
						<DropdownMenu.Item
							class="gap-2 p-2"
							onSelect={() => {
								modalProfile = {} as Profile;
								modalProfileOpen = true;
							}}
						>
							<div class="flex size-6 items-center justify-center rounded-md border bg-background">
								<Plus class="size-4" />
							</div>
							<div class="font-medium text-muted-foreground">Add Profile</div>
						</DropdownMenu.Item>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.Header>

	<Sidebar.Content>
		<Sidebar.Group>
			<Sidebar.GroupLabel>Overview</Sidebar.GroupLabel>
			<Sidebar.GroupContent class="flex flex-col gap-2">
				<Sidebar.Menu>
					{#each mainRoutes as r (r.title)}
						<Sidebar.MenuItem>
							<Sidebar.MenuButton tooltipContent={r.title}>
								{#snippet child({ props })}
									<a href={r.url} {...props}>
										<r.icon />
										<span>{r.title}</span>
									</a>
								{/snippet}
							</Sidebar.MenuButton>
						</Sidebar.MenuItem>
					{/each}
				</Sidebar.Menu>
			</Sidebar.GroupContent>
		</Sidebar.Group>

		<Sidebar.Group class="group-data-[collapsible=icon]:hidden">
			<Sidebar.GroupLabel>Management</Sidebar.GroupLabel>
			<Sidebar.Menu>
				{#each adminRoutes as r (r.title)}
					<Sidebar.MenuItem>
						<Sidebar.MenuButton tooltipContent={r.title}>
							{#snippet child({ props })}
								<a href={r.url} {...props}>
									<r.icon />
									<span>{r.title}</span>
								</a>
							{/snippet}
						</Sidebar.MenuButton>
					</Sidebar.MenuItem>
				{/each}
			</Sidebar.Menu>
		</Sidebar.Group>
	</Sidebar.Content>

	<!-- Bottom buttons -->
	<Sidebar.Group class="mt-auto">
		<Sidebar.GroupContent>
			<Sidebar.GroupLabel>Support</Sidebar.GroupLabel>
			<Sidebar.Menu>
				{#each supportRoutes as r (r.title)}
					<Sidebar.MenuItem>
						<Sidebar.MenuButton>
							{#snippet child({ props })}
								<a href={r.url} {...props}>
									<r.icon />
									<span>{r.title}</span>
								</a>
							{/snippet}
						</Sidebar.MenuButton>
					</Sidebar.MenuItem>
				{/each}
			</Sidebar.Menu>
		</Sidebar.GroupContent>
	</Sidebar.Group>

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
								<div class="relative">
									<RandomAvatar name={currentUser.data?.username} />
								</div>
								<div class="grid flex-1 text-left text-sm leading-tight">
									<span class="truncate font-semibold">{currentUser.data?.username}</span>
									<span class="truncate text-xs">{currentUser.data?.email}</span>
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
								<RandomAvatar name={currentUser.data?.username} />
								<div class="grid flex-1 text-left text-sm leading-tight">
									<span class="truncate font-semibold">{currentUser.data?.username}</span>
									<span class="truncate text-xs">{currentUser.data?.email}</span>
								</div>
							</div>
						</DropdownMenu.Label>
						<DropdownMenu.Separator />
						<DropdownMenu.Group>
							<DropdownMenu.Item onSelect={() => (modalUserOpen = true)}>
								<CircleUserRound />
								Account
							</DropdownMenu.Item>
							<DropdownMenu.Item onSelect={toggleMode}>
								{#if mode.current === 'dark'}
									<Sun class="size-4" />
									<span>Light Mode</span>
								{:else}
									<Moon class="size-4" />
									<span>Dark Mode</span>
								{/if}
							</DropdownMenu.Item>
						</DropdownMenu.Group>
						<DropdownMenu.Separator />
						<DropdownMenu.Item onSelect={() => logout.mutate({})}>
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
