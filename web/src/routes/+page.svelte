<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import {
		Globe,
		Shield,
		Bot,
		Users,
		Gauge,
		Layers2,
		Activity,
		Server,
		Network,
		Route,
		CheckCircle,
		AlertCircle,
		Clock,
		TrendingUp,
		Wifi,
		Database,
		Pen,
		Eye
	} from '@lucide/svelte';
	import { profile } from '$lib/stores/profile';
	import {
		agentClient,
		auditLogClient,
		dnsClient,
		middlewareClient,
		profileClient,
		routerClient,
		serviceClient,
		userClient
	} from '$lib/api';
	import { RouterType } from '$lib/gen/mantrae/v1/router_pb';
	import { MiddlewareType } from '$lib/gen/mantrae/v1/middleware_pb';
	import { DateFormat } from '$lib/stores/common';
	import { timestampDate, type Timestamp } from '@bufbuild/protobuf/wkt';
	import ProfileModal from '$lib/components/modals/ProfileModal.svelte';
	import ConfigModal from '$lib/components/modals/ConfigModal.svelte';
	import type { Profile } from '$lib/gen/mantrae/v1/profile_pb';
	import AuditLogModal from '$lib/components/modals/AuditLogModal.svelte';

	let totalAgents = $derived.by(async () => {
		const response = await agentClient.listAgents({
			profileId: profile.id,
			limit: -1n,
			offset: 0n
		});
		return response.totalCount;
	});

	let onlineAgents = $derived.by(async () => {
		const response = await agentClient.listAgents({
			profileId: profile.id,
			limit: -1n,
			offset: 0n
		});
		const now = Date.now();

		let activeAgents = response.agents.reduce((count, agent) => {
			if (!agent.updatedAt) return 0;
			const lastSeen = new Date(timestampDate(agent.updatedAt));
			const diffSeconds = (now - lastSeen.getTime()) / 1000;

			return diffSeconds <= 30 ? count + 1 : count;
		}, 0);
		return BigInt(activeAgents);
	});

	// Helper function to calculate percentage
	function getPercentage(value: number, total: number) {
		return total > 0 ? Math.round((value / total) * 100) : 0;
	}

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

<ProfileModal bind:item={modalProfile} bind:open={modalProfileOpen} />
<ConfigModal bind:open={modalConfigOpen} />
<AuditLogModal bind:open={modalAuditLogOpen} />

<div class="container mx-auto space-y-6 p-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h1 class="flex items-center gap-3 text-3xl font-bold tracking-tight">
				<div class="bg-primary/10 rounded-lg p-2">
					<Gauge class="text-primary h-6 w-6" />
				</div>
				Dashboard
			</h1>
			<p class="text-muted-foreground mt-1">Monitor and manage your reverse proxy configurations</p>
		</div>
	</div>

	{#if profile.id}
		<!-- Main Stats Grid -->
		<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
			<!-- Total Profiles -->
			<Card.Root class="relative overflow-hidden">
				<Card.Header class="flex flex-row items-center justify-between pb-2">
					<Card.Title class="text-muted-foreground text-sm font-medium">Total Profiles</Card.Title>
					<Layers2 class="text-muted-foreground h-4 w-4" />
				</Card.Header>
				<Card.Content>
					{#await profileClient.listProfiles({ limit: -1n, offset: 0n }) then result}
						<div class="text-3xl font-bold">{result.totalCount}</div>
						<div class="mt-2 flex items-center text-sm">
							<TrendingUp class="mr-1 h-3 w-3 text-green-500" />
							<span class="text-green-500">Active configurations</span>
						</div>
					{/await}
				</Card.Content>
				<div
					class="from-primary/5 absolute top-0 right-0 h-16 w-16 rounded-bl-full bg-gradient-to-bl to-transparent"
				></div>
			</Card.Root>

			<!-- Connected Agents -->
			<Card.Root class="relative overflow-hidden">
				<Card.Header class="flex flex-row items-center justify-between pb-2">
					<Card.Title class="text-muted-foreground text-sm font-medium">Connected Agents</Card.Title
					>
					<Bot class="text-muted-foreground h-4 w-4" />
				</Card.Header>
				<Card.Content>
					{#await onlineAgents}
						<div class="text-3xl font-bold">0</div>
					{:then result}
						<div class="text-3xl font-bold">{result}</div>
						<div class="mt-2 flex items-center text-sm">
							<Wifi class="mr-1 h-3 w-3 text-blue-500" />
							<span class="text-blue-500">Online now</span>
						</div>
					{/await}
				</Card.Content>
				<div
					class="absolute top-0 right-0 h-16 w-16 rounded-bl-full bg-gradient-to-bl from-blue-500/5 to-transparent"
				></div>
			</Card.Root>

			<!-- DNS Provider -->
			<Card.Root class="relative overflow-hidden">
				<Card.Header class="flex flex-row items-center justify-between pb-2">
					<Card.Title class="text-muted-foreground text-sm font-medium">DNS Provider</Card.Title>
					<Globe class="text-muted-foreground h-4 w-4" />
				</Card.Header>
				<Card.Content>
					{#await dnsClient.listDnsProviders({ limit: -1n, offset: 0n }) then result}
						<div class="text-3xl font-bold">
							{result.totalCount}
						</div>
						<div class="mt-2 flex items-center text-sm">
							{#if result.dnsProviders.find((p) => p.isDefault)}
								<CheckCircle class="mr-1 h-3 w-3 text-green-500" />
							{:else}
								<AlertCircle class="mr-1 h-3 w-3 text-yellow-500" />
							{/if}
							<span class="text-muted-foreground">
								{result.dnsProviders.find((p) => p.isDefault)?.name || 'None'}
								set as default
							</span>
						</div>
					{/await}
				</Card.Content>
				<div
					class="absolute top-0 right-0 h-16 w-16 rounded-bl-full bg-gradient-to-bl from-green-500/5 to-transparent"
				></div>
			</Card.Root>

			<!-- Total Users -->
			<Card.Root class="relative overflow-hidden">
				<Card.Header class="flex flex-row items-center justify-between pb-2">
					<Card.Title class="text-muted-foreground text-sm font-medium">System Users</Card.Title>
					<Users class="text-muted-foreground h-4 w-4" />
				</Card.Header>
				<Card.Content>
					{#await userClient.listUsers({}) then result}
						<div class="text-3xl font-bold">{result.totalCount}</div>
						<div class="mt-2 flex items-center text-sm">
							<Shield class="mr-1 h-3 w-3 text-purple-500" />
							<span class="text-muted-foreground">Access managed</span>
						</div>
					{/await}
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
					{#await routerClient.listRouters( { profileId: profile.id, limit: -1n, offset: 0n } ) then result}
						<div class="space-y-3">
							<div class="flex items-center justify-between">
								<div class="flex items-center gap-2">
									<div class="h-3 w-3 rounded-full bg-blue-500"></div>
									<span class="text-sm">HTTP Routers</span>
								</div>
								<Badge variant="secondary">
									{result.routers.filter((r) => r.type === RouterType.HTTP).length}
								</Badge>
							</div>
							<div class="flex items-center justify-between">
								<div class="flex items-center gap-2">
									<div class="h-3 w-3 rounded-full bg-green-500"></div>
									<span class="text-sm">TCP Routers</span>
								</div>
								<Badge variant="secondary">
									{result.routers.filter((r) => r.type === RouterType.TCP).length}
								</Badge>
							</div>
							<div class="flex items-center justify-between">
								<div class="flex items-center gap-2">
									<div class="h-3 w-3 rounded-full bg-purple-500"></div>
									<span class="text-sm">UDP Routers</span>
								</div>
								<Badge variant="secondary">
									{result.routers.filter((r) => r.type === RouterType.UDP).length}
								</Badge>
							</div>
						</div>
					{/await}
				</Card.Content>
			</Card.Root>

			<!-- Services Overview -->
			<Card.Root>
				<Card.Header>
					<Card.Title class="flex items-center gap-2">
						<Server class="h-5 w-5" />
						Services Overview
					</Card.Title>
				</Card.Header>
				<Card.Content class="space-y-4">
					{#await serviceClient.listServices({ profileId: profile.id }) then result}
						<div class="space-y-3">
							<div class="flex items-center justify-between">
								<div class="flex items-center gap-2">
									<div class="h-3 w-3 rounded-full bg-orange-500"></div>
									<span class="text-sm">HTTP Services</span>
								</div>
								<Badge variant="secondary">{result.totalCount}</Badge>
							</div>
							<div class="w-full">
								<div class="mb-1 flex items-center justify-between text-sm">
									<span class="text-muted-foreground">Health Status</span>
									<span>100%</span>
								</div>
								<Progress value={100} class="h-2" />
							</div>
						</div>
					{/await}
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
					{#await middlewareClient.listMiddlewares({ profileId: profile.id }) then result}
						<div class="space-y-3">
							<div class="flex items-center justify-between">
								<div class="flex items-center gap-2">
									<div class="h-3 w-3 rounded-full bg-red-500"></div>
									<span class="text-sm">HTTP Middlewares</span>
								</div>
								<Badge variant="secondary">
									{result.middlewares.filter((m) => m.type === MiddlewareType.HTTP).length}
								</Badge>
							</div>
							<div class="flex items-center justify-between">
								<div class="flex items-center gap-2">
									<div class="h-3 w-3 rounded-full bg-yellow-500"></div>
									<span class="text-sm">TCP Middlewares</span>
								</div>
								<Badge variant="secondary">
									{result.middlewares.filter((m) => m.type === MiddlewareType.TCP).length}
								</Badge>
							</div>
						</div>
					{/await}
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
					<div class="space-y-6">
						{#await profileClient.listProfiles({ limit: -1n, offset: 0n }) then result}
							{#each result.profiles || [] as profile (profile.id)}
								<div class="space-y-4 rounded-lg border p-4">
									<div class="flex items-start justify-between">
										<div class="flex items-start gap-3">
											<div class="bg-primary/10 rounded-lg p-2">
												<Database class="text-primary h-4 w-4" />
											</div>
											<div class="space-y-1">
												<h3 class="text-lg font-semibold">{profile.name}</h3>
												<p class="text-muted-foreground text-sm">{profile.description}</p>
												<p class="text-muted-foreground text-xs">
													{#if profile.createdAt}
														Created {DateFormat.format(timestampDate(profile.createdAt))}
													{/if}
												</p>
											</div>
										</div>
									</div>

									<Separator />

									<div class="grid grid-cols-2 gap-4 md:grid-cols-4">
										{#await agentClient.listAgents( { profileId: profile.id, limit: -1n, offset: 0n } ) then agentResult}
											<div class="text-center">
												<div class="text-2xl font-bold text-blue-600">{agentResult.totalCount}</div>
												<div class="text-muted-foreground text-xs">Agents</div>
											</div>
										{/await}

										{#await routerClient.listRouters( { profileId: profile.id, limit: -1n, offset: 0n } ) then routerResult}
											<div class="text-center">
												<div class="text-2xl font-bold text-green-600">
													{routerResult.totalCount}
												</div>
												<div class="text-muted-foreground text-xs">Routers</div>
											</div>
										{/await}

										{#await serviceClient.listServices( { profileId: profile.id, limit: -1n, offset: 0n } ) then serviceResult}
											<div class="text-center">
												<div class="text-2xl font-bold text-orange-600">
													{serviceResult.totalCount}
												</div>
												<div class="text-muted-foreground text-xs">Services</div>
											</div>
										{/await}

										{#await middlewareClient.listMiddlewares( { profileId: profile.id, limit: -1n, offset: 0n } ) then middlewareResult}
											<div class="text-center">
												<div class="text-2xl font-bold text-purple-600">
													{middlewareResult.totalCount}
												</div>
												<div class="text-muted-foreground text-xs">Middlewares</div>
											</div>
										{/await}
									</div>

									<div class="grid grid-cols-2 gap-2 md:grid-cols-6">
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
									</div>
								</div>
							{/each}
						{/await}
					</div>
				</Card.Content>
			</Card.Root>

			<!-- System Health & Quick Actions -->
			<div class="space-y-6">
				<!-- System Health -->
				<Card.Root>
					<Card.Header>
						<Card.Title class="flex items-center gap-2">
							<CheckCircle class="h-5 w-5 text-green-500" />
							System Health
						</Card.Title>
					</Card.Header>
					<Card.Content class="space-y-4">
						<div class="space-y-3">
							<div class="flex items-center justify-between">
								<span class="text-sm">Configuration Valid</span>
								<CheckCircle class="h-4 w-4 text-green-500" />
							</div>
							<div class="flex items-center justify-between">
								<span class="text-sm">DNS Resolution</span>
								<CheckCircle class="h-4 w-4 text-green-500" />
							</div>
							<div class="flex items-center justify-between">
								<span class="text-sm">Agent Connectivity</span>
								{#await Promise.all([totalAgents, onlineAgents]) then [total, online]}
									{#if total === online}
										<CheckCircle class="h-4 w-4 text-green-500" />
									{:else}
										<AlertCircle class="h-4 w-4 text-yellow-500" />
									{/if}
								{/await}
							</div>
							<div class="flex items-center justify-between">
								<span class="text-sm">SSL Certificates</span>
								<AlertCircle class="h-4 w-4 text-yellow-500" />
							</div>
						</div>
					</Card.Content>
				</Card.Root>

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
							{#await auditLogClient.listAuditLogs({ limit: 8n, offset: 0n }) then result}
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
											<p class="text-sm">{log.details}</p>

											<div class="text-muted-foreground flex items-center gap-2 text-xs">
												{#if log.createdAt}
													<span class="text-muted-foreground text-xs">
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
