<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Globe, Shield, Bot, LayoutDashboard, Origami, Users } from 'lucide-svelte';
	import { onMount } from 'svelte';
	import { api, profiles, stats } from '$lib/api';
	import { TraefikSource } from '$lib/types';

	interface ProfileStats {
		id: number;
		name: string;
		url: string;
		routers: number;
		services: number;
		middlewares: number;
		agents: number;
	}
	let profileStats: ProfileStats[] = $state([]);

	onMount(async () => {
		await api.loadStats();

		if (!$profiles) return;
		await api.getTraefikConfig(TraefikSource.LOCAL);

		// Get profile stats
		const t = await api.getTraefikStats();
		profileStats =
			t?.map((cfg) => ({
				id: cfg.id,
				name: cfg.name,
				url: cfg.url,
				routers: cfg.routers,
				services: cfg.services,
				middlewares: cfg.middlewares,
				agents: cfg.agents
			})) || [];
	});
</script>

<div class="container mx-auto p-6">
	<h2 class="mb-6 text-3xl font-bold tracking-tight">
		<div class="flex flex-row items-center gap-2">
			<LayoutDashboard />
			Dashboard
		</div>
	</h2>

	<!-- Stats Grid -->
	<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between pb-2">
				<Card.Title class="text-sm font-medium">Total Profiles</Card.Title>
				<Origami class="h-4 w-4 text-muted-foreground" />
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">{$stats.profiles}</div>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between pb-2">
				<Card.Title class="text-sm font-medium">Connected Agents</Card.Title>
				<Bot class="h-4 w-4 text-muted-foreground" />
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">{$stats.agents}</div>
				<p class="text-xs text-muted-foreground">Across all profiles</p>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between pb-2">
				<Card.Title class="text-sm font-medium">Active DNS Provider</Card.Title>
				<Globe class="h-4 w-4 text-muted-foreground" />
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">
					{$stats.activeDNS || 'None'}
				</div>
				<p class="text-xs text-muted-foreground">
					{$stats.dnsProviders} providers configured
				</p>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between pb-2">
				<Card.Title class="text-sm font-medium">Total Users</Card.Title>
				<Users class="h-4 w-4 text-muted-foreground" />
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">{$stats.users}</div>
				<p class="text-xs text-muted-foreground"></p>
			</Card.Content>
		</Card.Root>
	</div>

	<!-- Recent Activity -->
	<div class="mt-6 grid gap-6 md:grid-cols-2">
		<!-- <Card.Root class="col-span-1"> -->
		<!-- 	<Card.Header> -->
		<!-- 		<Card.Title>Recent Activity</Card.Title> -->
		<!-- 	</Card.Header> -->
		<!-- 	<Card.Content> -->
		<!-- 		<div class="space-y-4"> -->
		<!-- 			{#each stats.recentActivity as activity} -->
		<!-- 				<div class="flex items-center"> -->
		<!-- 					<div class="relative mr-4"> -->
		<!-- 						<Activity class="h-4 w-4" /> -->
		<!-- 						<span class="absolute right-0 top-0 h-2 w-2 rounded-full bg-green-500"></span> -->
		<!-- 					</div> -->
		<!-- 					<div class="space-y-1"> -->
		<!-- 						<p class="text-sm font-medium leading-none"> -->
		<!-- 							{activity.description} -->
		<!-- 						</p> -->
		<!-- 						<p class="text-sm text-muted-foreground"> -->
		<!-- 							{activity.timestamp} -->
		<!-- 						</p> -->
		<!-- 					</div> -->
		<!-- 				</div> -->
		<!-- 			{/each} -->
		<!-- 		</div> -->
		<!-- 	</Card.Content> -->
		<!-- </Card.Root> -->

		<!-- Profile Status -->
		<Card.Root class="col-span-2">
			<Card.Header>
				<Card.Title>Profile Status</Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="space-y-4">
					{#each profileStats as profile}
						<div class="flex items-center justify-between">
							<div class="flex items-center space-x-4">
								<Shield class="h-4 w-4" />
								<div class="space-y-1">
									<p class="text-sm font-medium leading-none">
										{profile.name}
									</p>
									<p class="text-xs text-muted-foreground">
										{profile.url}
									</p>
								</div>
							</div>
							<div class="flex items-center gap-2">
								<Badge variant={profile.agents > 0 ? 'default' : 'secondary'}>
									{profile.agents}
									{profile.agents === 1 ? 'Agent' : 'Agents'}
								</Badge>
								<Badge variant={profile.routers > 0 ? 'default' : 'secondary'}>
									{profile.routers}
									{profile.routers === 1 ? 'Router' : 'Routers'}
								</Badge>
								<Badge variant={profile.middlewares > 0 ? 'default' : 'secondary'}>
									{profile.middlewares}
									{profile.middlewares === 1 ? 'Middleware' : 'Middlewares'}
								</Badge>
							</div>
						</div>
					{/each}
				</div>
			</Card.Content>
		</Card.Root>
	</div>
</div>
