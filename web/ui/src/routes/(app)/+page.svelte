<script lang="ts">
	import { agent } from '$lib/api/agents.svelte';
	import { dns } from '$lib/api/dns.svelte';
	import { entrypoint } from '$lib/api/entrypoints.svelte';
	import { middleware } from '$lib/api/middleware.svelte';
	import { profile } from '$lib/api/profiles.svelte';
	import { router } from '$lib/api/router.svelte';
	import { service } from '$lib/api/service.svelte';
	import { setting } from '$lib/api/settings.svelte';
	import { transport } from '$lib/api/transport.svelte';
	import { user } from '$lib/api/users.svelte';
	import { audit } from '$lib/api/util.svelte';
	import AuditLogModal from '$lib/components/modals/AuditLogModal.svelte';
	import ConfigModal from '$lib/components/modals/ConfigModal.svelte';
	import ProfileModal from '$lib/components/modals/ProfileModal.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import CopyButton from '$lib/components/ui/copy-button/copy-button.svelte';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import Progress from '$lib/components/ui/progress/progress.svelte';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import type { Profile } from '$lib/gen/mantrae/v1/profile_pb';
	import { ProtocolType } from '$lib/gen/mantrae/v1/protocol_pb';
	import { profileID } from '$lib/store.svelte';
	import { cn, toDate } from '$lib/utils';
	import type { Timestamp } from '@bufbuild/protobuf/wkt';
	import {
		Activity,
		Bot,
		Clock,
		Database,
		EthernetPort,
		GaugeIcon,
		Globe,
		Layers,
		Layers2,
		Link,
		Pen,
		Route,
		Server,
		StarIcon,
		TruckIcon,
		Users
	} from '@lucide/svelte';

	// Queries
	const profileList = profile.list();
	const userList = user.list();
	const dnsList = dns.list();
	const logs = audit.logs(8n);

	// Contextual queries (depend on profileID.current)
	const routerList = $derived(router.list());
	const middlewareList = $derived(middleware.list());
	const serviceList = $derived(service.list());
	const agentList = $derived(agent.list());
	const transportList = $derived(transport.list());
	const entrypointList = $derived(entrypoint.list());

	const currentProfile = $derived(profileList.data?.find((p) => p.id === profileID.current));

	// Derived stats
	let onlineAgents = $derived.by(() => {
		let activeAgents =
			agentList.data?.reduce((count, agent) => {
				if (!agent.updatedAt) return 0;
				const lastSeen = toDate(agent.updatedAt);
				const diffSeconds = (Date.now() - lastSeen.getTime()) / 1000;
				return diffSeconds <= 20 ? count + 1 : count;
			}, 0) || 0;
		return BigInt(activeAgents);
	});

	// Traefik Connection Details
	const serverURL = setting.get('server_url');
	let connString = $derived(
		currentProfile
			? `${serverURL.data?.value}/api/${currentProfile.name}?token=${currentProfile.token}`
			: ''
	);

	let yamlStr = $derived(`providers:
  http:
    endpoint: "${connString}"`);

	let cliStr = $derived(`--providers.http.endpoint=${connString}`);

	function timeAgo(date: Timestamp) {
		const dateTime = toDate(date);
		const seconds = Math.floor((new Date().getTime() - dateTime.getTime()) / 1000);

		if (seconds < 60) return `${seconds} second${seconds !== 1 ? 's' : ''} ago`;

		const intervals = [
			{ label: 'year', seconds: 31536000 },
			{ label: 'month', seconds: 2592000 },
			{ label: 'day', seconds: 86400 },
			{ label: 'hour', seconds: 3600 },
			{ label: 'minute', seconds: 60 }
		];

		for (const interval of intervals) {
			const count = Math.floor(seconds / interval.seconds);
			if (count >= 1) {
				return `${count} ${interval.label}${count !== 1 ? 's' : ''} ago`;
			}
		}
	}

	let modalProfile = $state({} as Profile);
	let modalProfileOpen = $state(false);
	let modalConfigOpen = $state(false);
	let modalAuditLogOpen = $state(false);
</script>

<svelte:head>
	<title>Dashboard - Mantrae</title>
	<meta
		name="description"
		content="Monitor and manage your reverse proxy configurations, agents, DNS providers, and system users"
	/>
</svelte:head>

<ProfileModal data={modalProfile} bind:open={modalProfileOpen} />
<ConfigModal bind:open={modalConfigOpen} />
<AuditLogModal bind:open={modalAuditLogOpen} />

<div class="container mx-auto space-y-8 px-6">
	<!-- Header -->
	<div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
		<div>
			<h1 class="flex items-center gap-3 text-3xl font-bold tracking-tight">
				<div class="rounded-xl bg-primary/10 p-2.5 shadow-sm">
					<GaugeIcon class="h-6 w-6 text-primary" />
				</div>
				Dashboard
			</h1>
			<p class="mt-2 text-muted-foreground">Overview of your system and active configuration.</p>
		</div>

		{#if currentProfile}
			<div class="flex items-center gap-3 rounded-lg border bg-card p-2 px-4 shadow-sm">
				<div class="flex flex-col items-end">
					<span class="text-xs font-medium text-muted-foreground">Active Profile</span>
					<span class="font-bold text-primary">{currentProfile.name}</span>
				</div>
				<div class="h-8 w-px bg-border"></div>
				<Popover.Root>
					<Popover.Trigger>
						<Button size="sm" variant="ghost" class="h-8 w-8 p-0">
							<Link class="h-4 w-4 text-muted-foreground" />
							<span class="sr-only">Connection Info</span>
						</Button>
					</Popover.Trigger>
					<Popover.Content class="w-96 p-0" align="end">
						<Tabs.Root value="yaml" class="w-full">
							<div class="flex items-center justify-between border-b px-4 py-2">
								<h4 class="font-medium">Traefik Provider</h4>
								<Tabs.List class="h-7">
									<Tabs.Trigger value="yaml" class="h-5 text-xs">YAML</Tabs.Trigger>
									<Tabs.Trigger value="cli" class="h-5 text-xs">CLI</Tabs.Trigger>
								</Tabs.List>
							</div>
							<div class="px-4 py-2">
								<Tabs.Content value="yaml" class="mt-0 space-y-3">
									<p class="text-xs text-muted-foreground">
										Add this to your dynamic configuration file:
									</p>
									<div class="relative rounded-md bg-muted p-3">
										<pre class="overflow-x-auto text-xs"><code>{yamlStr}</code></pre>
										<div class="absolute top-2 right-2 bg-muted">
											<CopyButton text={yamlStr} />
										</div>
									</div>
								</Tabs.Content>
								<Tabs.Content value="cli" class="mt-0 space-y-3">
									<p class="text-xs text-muted-foreground">Use this flag when starting Traefik:</p>
									<div class="relative rounded-md bg-muted p-3">
										<code class="text-xs break-all">{cliStr}</code>
										<div class="absolute top-2 right-2 bg-muted">
											<CopyButton text={cliStr} />
										</div>
									</div>
								</Tabs.Content>
							</div>
						</Tabs.Root>
					</Popover.Content>
				</Popover.Root>
			</div>
		{/if}
	</div>

	<!-- System Overview Grid -->
	<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
		<!-- Total Profiles -->
		<Card.Root
			class="overflow-hidden border-b-4 border-b-primary shadow-sm transition-all hover:shadow-md"
		>
			<Card.Header class="flex flex-row items-center justify-between pb-2">
				<Card.Title class="text-sm font-medium text-muted-foreground">Profiles</Card.Title>
				<Layers2 class="h-4 w-4 text-muted-foreground" />
			</Card.Header>
			<Card.Content>
				<div class="text-3xl font-bold">{profileList.data?.length ?? 0}</div>
				<p class="mt-1 text-xs text-muted-foreground">Total configurations</p>
			</Card.Content>
		</Card.Root>

		<!-- Users -->
		<Card.Root
			class="overflow-hidden border-b-4 border-b-purple-500 shadow-sm transition-all hover:shadow-md"
		>
			<Card.Header class="flex flex-row items-center justify-between pb-2">
				<Card.Title class="text-sm font-medium text-muted-foreground">Users</Card.Title>
				<Users class="h-4 w-4 text-muted-foreground" />
			</Card.Header>
			<Card.Content>
				<div class="text-3xl font-bold">{userList.data?.length ?? 0}</div>
				<p class="mt-1 text-xs text-muted-foreground">Active accounts</p>
			</Card.Content>
		</Card.Root>

		<!-- DNS -->
		<Card.Root
			class="overflow-hidden border-b-4 border-b-green-500 shadow-sm transition-all hover:shadow-md"
		>
			<Card.Header class="flex flex-row items-center justify-between pb-2">
				<Card.Title class="text-sm font-medium text-muted-foreground">DNS Providers</Card.Title>
				<Globe class="h-4 w-4 text-muted-foreground" />
			</Card.Header>
			<Card.Content class="flex items-end justify-between">
				<div class="text-3xl font-bold">{dnsList.data?.length ?? 0}</div>
				<Badge variant="secondary" class="bg-muted">
					{#if dnsList.data?.find((p) => p.isDefault)}
						<StarIcon class="h-3 w-3 text-yellow-500" />
						<span title="Default provider">{dnsList.data?.find((p) => p.isDefault)?.name}</span>
					{:else}
						No default
					{/if}
				</Badge>
			</Card.Content>
		</Card.Root>

		<!-- Online Agents -->
		<Card.Root
			class="overflow-hidden border-b-4 border-b-orange-500 shadow-sm transition-all hover:shadow-md"
		>
			<Card.Header class="flex flex-row items-center justify-between pb-2">
				<Card.Title class="text-sm font-medium text-muted-foreground">Active Agents</Card.Title>
				<Bot class="h-4 w-4 text-muted-foreground" />
			</Card.Header>
			<Card.Content>
				<div class="flex items-baseline gap-2">
					<span class="text-3xl font-bold">{onlineAgents}</span>
					<span class="text-sm text-muted-foreground">/ {agentList.data?.length ?? 0}</span>
				</div>
				<Progress
					value={agentList.data?.length ? (Number(onlineAgents) / agentList.data.length) * 100 : 0}
					class="mt-2 h-1.5"
				/>
			</Card.Content>
		</Card.Root>
	</div>

	<div class="grid items-start gap-6 lg:grid-cols-3">
		<!-- Left Column: Active Profile Detail -->
		<div class="space-y-6 lg:col-span-2">
			{#if currentProfile}
				<Card.Root class="shadow-sm">
					<Card.Header>
						<div class="flex items-center justify-between">
							<div class="space-y-1">
								<Card.Title class="flex items-center gap-2 text-xl">
									<Database class="h-5 w-5 text-primary" />
									{currentProfile.name}
								</Card.Title>
								<Card.Description
									>{currentProfile.description || 'No description provided'}</Card.Description
								>
							</div>
							<div class="flex gap-2">
								<Button
									size="sm"
									variant="outline"
									class="h-8 gap-2"
									onclick={() => {
										modalProfile = currentProfile!;
										modalProfileOpen = true;
									}}
								>
									<Pen class="h-3.5 w-3.5" /> Edit
								</Button>
								<Button
									size="sm"
									variant="outline"
									class="h-8 gap-2"
									onclick={() => {
										modalProfile = currentProfile!;
										modalConfigOpen = true;
									}}
								>
									<Activity class="h-3.5 w-3.5" /> Config
								</Button>
							</div>
						</div>
					</Card.Header>
					<Card.Content>
						<div class="grid grid-cols-2 sm:grid-cols-5">
							<div class="rounded-lg bg-secondary/30 p-4 text-center">
								<div class="mb-1 flex justify-center text-blue-500">
									<Route class="h-6 w-6" />
								</div>
								<div class="text-2xl font-bold">{routerList.data?.length ?? 0}</div>
								<div class="text-xs text-muted-foreground">Routers</div>
							</div>
							<div class="rounded-lg bg-secondary/30 p-4 text-center">
								<div class="mb-1 flex justify-center text-orange-500">
									<Server class="h-6 w-6" />
								</div>
								<div class="text-2xl font-bold">{serviceList.data?.length ?? 0}</div>
								<div class="text-xs text-muted-foreground">Services</div>
							</div>
							<div class="rounded-lg bg-secondary/30 p-4 text-center">
								<div class="mb-1 flex justify-center text-purple-500">
									<Layers class="h-6 w-6" />
								</div>
								<div class="text-2xl font-bold">{middlewareList.data?.length ?? 0}</div>
								<div class="text-xs text-muted-foreground">Middlewares</div>
							</div>
							<div class="rounded-lg bg-secondary/30 p-4 text-center">
								<div class="mb-1 flex justify-center text-green-500">
									<EthernetPort class="h-6 w-6" />
								</div>
								<div class="text-2xl font-bold">{entrypointList.data?.length ?? 0}</div>
								<div class="text-xs text-muted-foreground">Entrypoints</div>
							</div>
							<div class="rounded-lg bg-secondary/30 p-4 text-center">
								<div class="mb-1 flex justify-center text-red-500">
									<TruckIcon class="h-6 w-6" />
								</div>
								<div class="text-2xl font-bold">{transportList.data?.length ?? 0}</div>
								<div class="text-xs text-muted-foreground">Server Transports</div>
							</div>
						</div>

						<Separator class="my-4" />

						<div class="grid gap-6 md:grid-cols-2">
							<!-- Router Breakdown -->
							<div class="space-y-3">
								<h4 class="text-sm font-semibold text-muted-foreground">Protocols</h4>
								<div class="space-y-2">
									<div class="flex items-center justify-between text-sm">
										<span class="flex items-center gap-2"
											><div class="h-2 w-2 rounded-full bg-blue-500"></div>
											HTTP</span
										>
										<span class="font-mono"
											>{routerList.data?.filter((r) => r.type === ProtocolType.HTTP).length ??
												0}</span
										>
									</div>
									<div class="flex items-center justify-between text-sm">
										<span class="flex items-center gap-2"
											><div class="h-2 w-2 rounded-full bg-green-500"></div>
											TCP</span
										>
										<span class="font-mono"
											>{routerList.data?.filter((r) => r.type === ProtocolType.TCP).length ??
												0}</span
										>
									</div>
									<div class="flex items-center justify-between text-sm">
										<span class="flex items-center gap-2"
											><div class="h-2 w-2 rounded-full bg-purple-500"></div>
											UDP</span
										>
										<span class="font-mono"
											>{routerList.data?.filter((r) => r.type === ProtocolType.UDP).length ??
												0}</span
										>
									</div>
								</div>
							</div>

							<!-- Middleware Breakdown -->
							<div class="space-y-3">
								<h4 class="text-sm font-semibold text-muted-foreground">Middlewares</h4>
								<div class="space-y-2">
									<div class="flex items-center justify-between text-sm">
										<span class="flex items-center gap-2"
											><div class="h-2 w-2 rounded-full bg-red-500"></div>
											HTTP</span
										>
										<span class="font-mono"
											>{middlewareList.data?.filter((m) => m.type === ProtocolType.HTTP).length ??
												0}</span
										>
									</div>
									<div class="flex items-center justify-between text-sm">
										<span class="flex items-center gap-2"
											><div class="h-2 w-2 rounded-full bg-yellow-500"></div>
											TCP</span
										>
										<span class="font-mono"
											>{middlewareList.data?.filter((m) => m.type === ProtocolType.TCP).length ??
												0}</span
										>
									</div>
								</div>
							</div>
						</div>
					</Card.Content>
				</Card.Root>
			{:else}
				<div
					class="flex h-full flex-col items-center justify-center rounded-lg border border-dashed p-8 text-center text-muted-foreground"
				>
					<Database class="mb-4 h-12 w-12 opacity-20" />
					<h3 class="text-lg font-medium">No Profile Selected</h3>
					<p class="text-sm">Select or create a profile to view its configuration.</p>
				</div>
			{/if}

			<!-- Recent Activity -->
			<Card.Root>
				<Card.Header class="flex flex-row items-center justify-between">
					<Card.Title class="flex items-center gap-2 text-base">
						<Clock class="h-4 w-4 text-muted-foreground" />
						Recent Activity
					</Card.Title>
					<Button
						variant="ghost"
						size="sm"
						class="h-8 text-xs"
						onclick={() => (modalAuditLogOpen = true)}
					>
						View All
					</Button>
				</Card.Header>
				<Card.Content>
					<div class="space-y-4">
						{#each logs.data || [] as log (log.id)}
							<div class="flex items-start gap-3 text-sm">
								<div
									class={cn(
										'mt-1.5 h-2 w-2 shrink-0 rounded-full',
										log.agentId ? 'bg-blue-500' : log.userId ? 'bg-green-500' : 'bg-orange-500'
									)}
								></div>
								<div class="grid gap-1">
									<p class="text-foreground">{log.details}</p>
									<div class="flex items-center gap-2 text-xs text-muted-foreground">
										<span>{log.createdAt ? timeAgo(log.createdAt) : 'Unknown'}</span>
										<span>•</span>
										<span>{log.agentId ? 'Agent' : log.userId ? 'User' : 'System'}</span>
									</div>
								</div>
							</div>
						{/each}
						{#if (logs.data?.length ?? 0) === 0}
							<p class="text-center text-xs text-muted-foreground">No recent logs found.</p>
						{/if}
					</div>
				</Card.Content>
			</Card.Root>
		</div>

		<!-- Right Column: Profile Switcher -->
		<Card.Root class="max-h-[80vh] overflow-hidden">
			<Card.Header class="border-b">
				<Card.Title>Profile Switcher</Card.Title>
			</Card.Header>
			<Card.Content class="space-y-4 p-0">
				<div class="space-y-2 px-4">
					{#each profileList.data || [] as p (p.id)}
						<button
							class={cn(
								'flex w-full items-center justify-between rounded-lg border p-3 text-left transition-colors hover:bg-accent hover:text-accent-foreground',
								p.id === profileID.current
									? 'border-primary bg-primary/5 shadow-sm'
									: 'border-transparent bg-muted/30'
							)}
							onclick={() => (profileID.current = p.id)}
						>
							<div class="grid gap-0.5">
								<span class="text-sm font-medium">{p.name}</span>
								<span class="line-clamp-1 text-xs text-muted-foreground"
									>{p.description || 'No description'}</span
								>
							</div>
							{#if p.id === profileID.current}
								<div
									class="h-2 w-2 rounded-full bg-primary shadow-[0_0_8px] shadow-primary/50"
								></div>
							{/if}
						</button>
					{/each}
					{#if (profileList.data?.length ?? 0) === 0}
						<div class="py-8 text-center text-sm text-muted-foreground">No profiles found.</div>
					{/if}
				</div>
				<div class="border-t bg-muted/30 px-3 pt-4">
					<Button
						class="w-full gap-2"
						variant="outline"
						onclick={() => {
							modalProfile = {} as Profile;
							modalProfileOpen = true;
						}}
					>
						<Layers2 class="h-4 w-4" />
						Create New Profile
					</Button>
				</div>
			</Card.Content>
		</Card.Root>
	</div>
</div>
