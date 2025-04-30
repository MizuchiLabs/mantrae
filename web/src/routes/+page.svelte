<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import {
		Globe,
		Shield,
		Bot,
		LayoutDashboard,
		Origami,
		Users,
		TriangleAlert,
		Trash2
	} from '@lucide/svelte';
	import { onMount } from 'svelte';
	import { api, errors, profiles, stats } from '$lib/api';
	import { TraefikSource } from '$lib/types';
	import Button from '$lib/components/ui/button/button.svelte';
	import { toast } from 'svelte-sonner';

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

	const clearErrors = async () => {
		await api.deleteErrorsByProfile();
		toast.success('Errors cleared successfully');
	};

	onMount(async () => {
		await api.loadStats();
		await api.listErrors();

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
				<Origami class="text-muted-foreground h-4 w-4" />
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">{$stats.profiles}</div>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between pb-2">
				<Card.Title class="text-sm font-medium">Connected Agents</Card.Title>
				<Bot class="text-muted-foreground h-4 w-4" />
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">{$stats.agents}</div>
				<p class="text-muted-foreground text-xs">Across all profiles</p>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between pb-2">
				<Card.Title class="text-sm font-medium">Active DNS Provider</Card.Title>
				<Globe class="text-muted-foreground h-4 w-4" />
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">
					{$stats.activeDNS || 'None'}
				</div>
				<p class="text-muted-foreground text-xs">
					{$stats.dnsProviders} providers configured
				</p>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between pb-2">
				<Card.Title class="text-sm font-medium">Total Users</Card.Title>
				<Users class="text-muted-foreground h-4 w-4" />
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">{$stats.users}</div>
				<p class="text-muted-foreground text-xs"></p>
			</Card.Content>
		</Card.Root>
	</div>

	<div class="mt-6 flex flex-row items-start gap-6">
		<!-- Profile Status -->
		<Card.Root class="flex-1">
			<Card.Header>
				<Card.Title>Profile Status</Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="space-y-4">
					{#each profileStats as profile (profile.id)}
						<div class="flex items-center justify-between">
							<div class="flex items-center space-x-4">
								<Shield class="h-4 w-4" />
								<div class="space-y-1">
									<p class="text-sm leading-none font-medium">
										{profile.name}
									</p>
									<p class="text-muted-foreground text-xs">
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

		<!-- Errors -->
		<Card.Root class="flex-1">
			<Card.Header>
				<Card.Title class="flex items-center justify-between gap-2">
					System Errors
					<Button
						variant="ghost"
						size="icon"
						class="rounded-full hover:bg-red-300"
						onclick={clearErrors}
					>
						<Trash2 />
					</Button>
				</Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="space-y-4">
					{#each $errors as error (error.id)}
						<div class="flex items-center">
							<div class="relative mr-4">
								<TriangleAlert class="h-4 w-4 text-red-500" />
							</div>
							<div class="space-y-1">
								<p class="text-sm">
									{error.message}
								</p>
								<p class="text-muted-foreground text-sm">
									{error.details}
								</p>
							</div>
						</div>
					{/each}
				</div>
			</Card.Content>
		</Card.Root>
	</div>
</div>
