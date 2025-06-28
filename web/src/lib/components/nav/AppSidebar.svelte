<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import logo from '$lib/images/logo-white.svg';
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
		Tag,
		Pencil,
		CircleUserRound,
		Sun,
		Moon,
		type IconProps,
		EthernetPort,
		Gauge,
		Layers2
	} from '@lucide/svelte';
	import { profile } from '$lib/stores/profile';
	import { user } from '$lib/stores/user';
	import { profileClient, userClient } from '$lib/api';
	import type { Profile } from '$lib/gen/mantrae/v1/profile_pb';
	import ProfileModal from '$lib/components/modals/profile.svelte';
	import UserModal from '$lib/components/modals/user.svelte';
	import { toggleMode, mode } from 'mode-watcher';
	import { goto } from '$app/navigation';

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
		{ title: 'Router', url: '/router/', icon: Route },
		{ title: 'Middlewares', url: '/middlewares/', icon: Layers },
		{ title: 'EntryPoints', url: '/entrypoints/', icon: EthernetPort }
	];
	const adminRoutes: Route[] = [
		{ title: 'Users', url: '/users/', icon: Users },
		{ title: 'Agents', url: '/agents/', icon: Bot },
		{ title: 'DNS', url: '/dns/', icon: Globe },
		{ title: 'Settings', url: '/settings/', icon: Settings }
	];

	let modalProfile = $state({} as Profile);
	let modalProfileOpen = $state(false);

	let modalUserOpen = $state(false);

	// let modalInfo = $state({} as TraefikInfo);
	let modalInfoOpen = $state(false);
</script>

<ProfileModal bind:item={modalProfile} bind:open={modalProfileOpen} />
<!-- <InfoModal bind:open={infoModalOpen} /> -->

{#if user.isLoggedIn() && user.value}
	<UserModal bind:open={modalUserOpen} bind:item={user.value} data={undefined} />
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
									class="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg"
								>
									<img src={logo} alt="Mantrae Logo" width="18" />
								</div>
								<div class="grid flex-1 text-left text-sm leading-tight">
									<span class="truncate font-semibold">
										{profile.name ? profile.name : 'Select Profile'}
									</span>
									<span class="truncate text-xs">{profile.value?.description ?? ''}</span>
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
						{#await profileClient.listProfiles({ limit: -1n, offset: 0n }) then value}
							{#each value.profiles || [] as p (p.id)}
								<DropdownMenu.Item
									onSelect={() => (profile.value = p)}
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
						{/await}
						<DropdownMenu.Separator />
						<DropdownMenu.Item
							class="gap-2 p-2"
							onSelect={() => {
								modalProfile = {} as Profile;
								modalProfileOpen = true;
							}}
						>
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

		<!-- Extra buttons (Traefik, etc.) -->
		<!-- <Sidebar.Group class="mt-auto"> -->
		<!-- 	<Sidebar.GroupContent> -->
		<!-- 		<Sidebar.GroupLabel>Status</Sidebar.GroupLabel> -->
		<!-- 		<Sidebar.Menu> -->
		<!-- 			{#if $profiles} -->
		<!-- 				<Sidebar.MenuItem> -->
		<!-- 					<Sidebar.MenuButton> -->
		<!-- 						{#snippet child({ props })} -->
		<!-- 							<button {...props} onclick={() => (infoModalOpen = true)}> -->
		<!-- 								<Zap /> -->
		<!-- 								<span>Traefik Status</span> -->
		<!-- 							</button> -->
		<!-- 						{/snippet} -->
		<!-- 					</Sidebar.MenuButton> -->
		<!-- 				</Sidebar.MenuItem> -->
		<!-- 			{/if} -->
		<!-- 		</Sidebar.Menu> -->
		<!-- 	</Sidebar.GroupContent> -->
		<!-- </Sidebar.Group> -->
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
						<DropdownMenu.Item
							onSelect={() => {
								userClient.logoutUser({});
								user.clear();
								goto('/login');
							}}
						>
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
