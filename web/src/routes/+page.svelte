<script lang="ts">
	import { auditLogClient } from '$lib/api';
	import AuditLogModal from '$lib/components/modals/AuditLogModal.svelte';
	import ConfigModal from '$lib/components/modals/ConfigModal.svelte';
	import ProfileModal from '$lib/components/modals/ProfileModal.svelte';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import TraefikConnection from '$lib/components/utils/TraefikConnection.svelte';
	import type { Profile } from '$lib/gen/mantrae/v1/profile_pb';
	import { ProtocolType } from '$lib/gen/mantrae/v1/protocol_pb';
	import { DateFormat } from '$lib/stores/common';
	import { profile } from '$lib/stores/profile';
	import {
		agents,
		dnsProviders,
		middlewares,
		profiles,
		routers,
		services,
		traefikInstances,
		users
	} from '$lib/stores/realtime';
	import { timestampDate, type Timestamp } from '@bufbuild/protobuf/wkt';
	import {
		Activity,
		Bot,
		CircleAlert,
		CircleCheck,
		Clock,
		Database,
		Eye,
		Gauge,
		Globe,
		Layers2,
		Network,
		Pen,
		Route,
		Server,
		Shield,
		TrendingUp,
		Users,
		Wifi
	} from '@lucide/svelte';

	let onlineAgents = $derived.by(() => {
		let activeAgents = $agents.reduce((count, agent) => {
			if (!agent.updatedAt) return 0;
			const lastSeen = new Date(timestampDate(agent.updatedAt));
			const diffSeconds = (Date.now() - lastSeen.getTime()) / 1000;

			return diffSeconds <= 20 ? count + 1 : count;
		}, 0);
		return BigInt(activeAgents);
	});

	function timeAgo(date: Timestamp) {
		const dateTime = new Date(timestampDate(date));
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
	<meta name="description" content="Monitor and manage your reverse proxy configurations, agents, DNS providers, and system users" />
</svelte:head>

<ProfileModal bind:item={modalProfile} bind:open={modalProfileOpen} />
<ConfigModal bind:open={modalConfigOpen} />
<AuditLogModal bind:open={modalAuditLogOpen} />

<div class="container mx-auto space-y-6 p-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h1 class="flex items-center gap-3 text-3xl font-bold tracking-tight">
				<div class="rounded-lg bg-primary/10 p-2">
					<Gauge class="h-6 w-6 text-primary" />
				</div>
				Dashboard
			</h1>
			<p class="mt-1 text-muted-foreground">Monitor and manage your reverse proxy configurations</p>
		</div>
	</div>

	{#if profile.id}
		<!-- Main Stats Grid -->
		<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
			<!-- Total Profiles -->
			<Card.Root class="relative overflow-hidden">
				<Card.Header class="flex flex-row items-center justify-between pb-2">
					<Card.Title class="text-sm font-medium text-muted-foreground">Total Profiles</Card.Title>
					<Layers2 class="h-4 w-4 text-muted-foreground" />
				</Card.Header>
				<Card.Content>
					<div class="text-3xl font-bold">{$profiles?.length}</div>
					<div class="mt-2 flex items-center text-sm">
						<TrendingUp class="mr-1 h-3 w-3 text-green-500" />
						<span class="text-green-500">Active configurations</span>
					</div>
				</Card.Content>
				<div
					class="absolute top-0 right-0 h-16 w-16 rounded-bl-full bg-gradient-to-bl from-primary/5 to-transparent"
				></div>
			</Card.Root>

			<!-- Connected Agents -->
			<Card.Root class="relative overflow-hidden">
				<Card.Header class="flex flex-row items-center justify-between pb-2">
					<Card.Title class="text-sm font-medium text-muted-foreground">Connected Agents</Card.Title
					>
					<Bot class="h-4 w-4 text-muted-foreground" />
				</Card.Header>
				<Card.Content>
					<div class="text-3xl font-bold">{onlineAgents}</div>
					<div class="mt-2 flex items-center text-sm">
						<Wifi class="mr-1 h-3 w-3 text-blue-500" />
						<span class="text-blue-500">Online now</span>
					</div>
				</Card.Content>
				<div
					class="absolute top-0 right-0 h-16 w-16 rounded-bl-full bg-gradient-to-bl from-blue-500/5 to-transparent"
				></div>
			</Card.Root>

			<!-- DNS Provider -->
			<Card.Root class="relative overflow-hidden">
				<Card.Header class="flex flex-row items-center justify-between pb-2">
					<Card.Title class="text-sm font-medium text-muted-foreground">DNS Provider</Card.Title>
					<Globe class="h-4 w-4 text-muted-foreground" />
				</Card.Header>
				<Card.Content>
					<div class="text-3xl font-bold">
						{$dnsProviders?.length}
					</div>
					<div class="mt-2 flex items-center text-sm">
						{#if $dnsProviders.find((p) => p.isDefault)}
							<CircleCheck class="mr-1 h-3 w-3 text-green-500" />
						{:else}
							<CircleAlert class="mr-1 h-3 w-3 text-yellow-500" />
						{/if}
						<span class="text-muted-foreground">
							{$dnsProviders.find((p) => p.isDefault)?.name || 'None'}
							set as default
						</span>
					</div>
				</Card.Content>
				<div
					class="absolute top-0 right-0 h-16 w-16 rounded-bl-full bg-gradient-to-bl from-green-500/5 to-transparent"
				></div>
			</Card.Root>

			<!-- Total Users -->
			<Card.Root class="relative overflow-hidden">
				<Card.Header class="flex flex-row items-center justify-between pb-2">
					<Card.Title class="text-sm font-medium text-muted-foreground">System Users</Card.Title>
					<Users class="h-4 w-4 text-muted-foreground" />
				</Card.Header>
				<Card.Content>
					<div class="text-3xl font-bold">{$users?.length}</div>
					<div class="mt-2 flex items-center text-sm">
						<Shield class="mr-1 h-3 w-3 text-purple-500" />
						<span class="text-muted-foreground">Access managed</span>
					</div>
				</Card.Content>
				<div
					class="absolute top-0 right-0 h-16 w-16 rounded-bl-full bg-gradient-to-bl from-purple-500/5 to-transparent"
				></div>
			</Card.Root>
		</div>

		<!-- Configuration Overview -->
		<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
			<!-- Routers Overview -->
			<Card.Root>
				<Card.Header>
					<Card.Title class="flex items-center gap-2">
						<Route class="h-5 w-5" />
						Routers Overview
					</Card.Title>
				</Card.Header>
				<Card.Content class="space-y-4">
					<div class="space-y-3">
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-2">
								<div class="h-3 w-3 rounded-full bg-blue-500"></div>
								<span class="text-sm">HTTP Routers</span>
							</div>
							<Badge variant="secondary">
								{$routers.filter((r) => r.type === ProtocolType.HTTP).length}
							</Badge>
						</div>
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-2">
								<div class="h-3 w-3 rounded-full bg-green-500"></div>
								<span class="text-sm">TCP Routers</span>
							</div>
							<Badge variant="secondary">
								{$routers.filter((r) => r.type === ProtocolType.TCP).length}
							</Badge>
						</div>
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-2">
								<div class="h-3 w-3 rounded-full bg-purple-500"></div>
								<span class="text-sm">UDP Routers</span>
							</div>
							<Badge variant="secondary">
								{$routers.filter((r) => r.type === ProtocolType.UDP).length}
							</Badge>
						</div>
					</div>
				</Card.Content>
			</Card.Root>

			<!-- Middlewares Overview -->
			<Card.Root>
				<Card.Header>
					<Card.Title class="flex items-center gap-2">
						<Network class="h-5 w-5" />
						Middlewares
					</Card.Title>
				</Card.Header>
				<Card.Content class="space-y-4">
					<div class="space-y-3">
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-2">
								<div class="h-3 w-3 rounded-full bg-red-500"></div>
								<span class="text-sm">HTTP Middlewares</span>
							</div>
							<Badge variant="secondary">
								{$middlewares.filter((m) => m.type === ProtocolType.HTTP).length}
							</Badge>
						</div>
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-2">
								<div class="h-3 w-3 rounded-full bg-yellow-500"></div>
								<span class="text-sm">TCP Middlewares</span>
							</div>
							<Badge variant="secondary">
								{$middlewares.filter((m) => m.type === ProtocolType.TCP).length}
							</Badge>
						</div>
					</div>
				</Card.Content>
			</Card.Root>

			<!-- Instances Overview -->
			<Card.Root>
				<Card.Header class="flex items-center justify-between">
					<Card.Title class="flex items-center gap-2">
						<Server class="h-5 w-5" />
						Traefik Instances
					</Card.Title>
					<Badge variant="secondary">{$traefikInstances?.length}</Badge>
				</Card.Header>
				<Card.Content class="space-y-4">
					<div class="max-h-64 space-y-2 overflow-y-auto pr-2">
						{#each $traefikInstances || [] as instance (instance.id)}
							<div class="flex items-center justify-between">
								<div class="flex items-center gap-2">
									<div class="h-3 w-3 rounded-full bg-green-500"></div>
									<span class="text-sm">{instance.name}</span>
								</div>
								<Badge variant="secondary">
									{instance.url}
								</Badge>
							</div>
						{/each}
					</div>
				</Card.Content>
			</Card.Root>
		</div>

		<!-- Profile Details and Activity -->
		<div class="grid gap-6 lg:grid-cols-3">
			<!-- Profile Status - Enhanced -->
			<Card.Root class="lg:col-span-2">
				<Card.Header>
					<div class="flex items-center justify-between">
						<Card.Title class="flex items-center gap-2">
							<Shield class="h-5 w-5" />
							Profile Configuration Status
						</Card.Title>
						<Badge variant="outline" class="gap-1">
							<Activity class="h-3 w-3" />
							Live
						</Badge>
					</div>
				</Card.Header>
				<Card.Content>
					<div class="space-y-4">
						{#each $profiles || [] as profile (profile.id)}
							<div class="space-y-4 rounded-lg border p-4">
								<div class="flex items-start justify-between">
									<div class="flex items-start gap-3">
										<div class="rounded-lg bg-primary/10 p-2">
											<Database class="h-4 w-4 text-primary" />
										</div>
										<div class="space-y-1">
											<h3 class="text-lg font-semibold">{profile.name}</h3>
											<p class="text-sm text-muted-foreground">{profile.description}</p>
											<p class="text-xs text-muted-foreground">
												{#if profile.createdAt}
													Created {DateFormat.format(timestampDate(profile.createdAt))}
												{/if}
											</p>
										</div>
									</div>
								</div>

								<Separator />

								<div class="grid grid-cols-2 gap-4 md:grid-cols-4">
									<div class="text-center">
										<div class="text-2xl font-bold text-blue-600">{$agents?.length}</div>
										<div class="text-xs text-muted-foreground">Agents</div>
									</div>

									<div class="text-center">
										<div class="text-2xl font-bold text-green-600">
											{$routers?.length}
										</div>
										<div class="text-xs text-muted-foreground">Routers</div>
									</div>

									<div class="text-center">
										<div class="text-2xl font-bold text-orange-600">
											{$services?.length}
										</div>
										<div class="text-xs text-muted-foreground">Services</div>
									</div>

									<div class="text-center">
										<div class="text-2xl font-bold text-purple-600">
											{$middlewares?.length}
										</div>
										<div class="text-xs text-muted-foreground">Middlewares</div>
									</div>
								</div>

								<div class="grid grid-cols-1 gap-2 lg:grid-cols-3">
									<Button
										size="sm"
										variant="outline"
										class="gap-2"
										onclick={() => {
											modalProfile = profile;
											modalProfileOpen = true;
										}}
									>
										<Pen />
										Edit
									</Button>
									<Button
										size="sm"
										variant="outline"
										class="gap-2"
										onclick={() => (modalConfigOpen = true)}
									>
										<Activity />
										Config
									</Button>
									<TraefikConnection {profile} variant="compact" />
								</div>
							</div>
						{/each}
					</div>
				</Card.Content>
			</Card.Root>

			<div class="space-y-6">
				<!-- Traefik Connection -->
				{#if profile.value}
					<TraefikConnection profile={profile.value} variant="full" />
				{/if}

				<!-- Recent Activity -->
				<Card.Root>
					<Card.Header class="flex items-center justify-between">
						<Card.Title class="flex items-center gap-2">
							<Clock class="h-5 w-5" />
							Recent Activity
						</Card.Title>
						<Button variant="ghost" size="sm" onclick={() => (modalAuditLogOpen = true)}>
							<Eye />
							Show more
						</Button>
					</Card.Header>
					<Card.Content class="space-y-3">
						<div class="space-y-3 text-sm">
							{#await auditLogClient.listAuditLogs({ limit: 8n }) then result}
								{#each result.auditLogs || [] as log (log.id)}
									<div class="flex items-start gap-3">
										{#if log.agentId}
											<div class="mt-2 h-2 w-2 rounded-full bg-blue-500"></div>
										{:else if log.userId}
											<div class="mt-2 h-2 w-2 rounded-full bg-green-500"></div>
										{:else}
											<div class="mt-2 h-2 w-2 rounded-full bg-orange-500"></div>
										{/if}
										<div>
											<p class="line-clamp-1 text-sm" title={log.details}>{log.details}</p>

											<div class="flex items-center gap-2 text-xs text-muted-foreground">
												{#if log.createdAt}
													<span class="text-xs text-muted-foreground">
														{timeAgo(log.createdAt)}
													</span>
												{/if}
												{#if log.agentId}
													<span
														class="rounded bg-blue-100 px-1.5 py-0.5 text-blue-700"
														title={log.agentId}
													>
														Agent: {log.agentName || `...${log.agentId.slice(-8)}`}
													</span>
												{:else if log.userId}
													<span
														class="rounded bg-green-100 px-1.5 py-0.5 text-green-700"
														title={log.userId}
													>
														User: {log.userName || `...${log.userId.slice(-8)}`}
													</span>
												{/if}
											</div>
										</div>
									</div>
								{/each}
							{/await}
						</div>
					</Card.Content>
				</Card.Root>
			</div>
		</div>
	{/if}
</div>
