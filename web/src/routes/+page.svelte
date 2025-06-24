<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Globe, Shield, Bot, Origami, Users, Gauge } from '@lucide/svelte';
	import { profile } from '$lib/stores/profile';
	import {
		agentClient,
		dnsClient,
		middlewareClient,
		profileClient,
		routerClient,
		serviceClient,
		userClient
	} from '$lib/api';
</script>

<div class="container mx-auto p-6">
	<h2 class="mb-6 text-3xl font-bold tracking-tight">
		<div class="flex flex-row items-center gap-2">
			<Gauge />
			Dashboard
		</div>
	</h2>

	<!-- Stats Grid -->
	{#if profile.id}
		<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
			<Card.Root>
				<Card.Header class="flex flex-row items-center justify-between pb-2">
					<Card.Title class="text-sm font-medium">Total Profiles</Card.Title>
					<Origami class="text-muted-foreground h-4 w-4" />
				</Card.Header>
				<Card.Content>
					{#await profileClient.listProfiles({}) then result}
						<div class="text-2xl font-bold">{result.totalCount}</div>
					{/await}
				</Card.Content>
			</Card.Root>

			<Card.Root>
				<Card.Header class="flex flex-row items-center justify-between pb-2">
					<Card.Title class="text-sm font-medium">Connected Agents</Card.Title>
					<Bot class="text-muted-foreground h-4 w-4" />
				</Card.Header>
				<Card.Content>
					{#await agentClient.listAgents({ profileId: profile.id }) then result}
						<div class="text-2xl font-bold">{result.totalCount}</div>
					{/await}
				</Card.Content>
			</Card.Root>

			<Card.Root>
				<Card.Header class="flex flex-row items-center justify-between pb-2">
					<Card.Title class="text-sm font-medium">Active DNS Provider</Card.Title>
					<Globe class="text-muted-foreground h-4 w-4" />
				</Card.Header>
				<Card.Content>
					{#await dnsClient.listDnsProviders({}) then result}
						<div class="text-2xl font-bold">
							{result.dnsProviders.find((p) => p.isActive)?.name || 'None'}
						</div>
						<p class="text-muted-foreground text-xs">
							{result.totalCount} providers configured
						</p>
					{/await}
				</Card.Content>
			</Card.Root>

			<Card.Root>
				<Card.Header class="flex flex-row items-center justify-between pb-2">
					<Card.Title class="text-sm font-medium">Total Users</Card.Title>
					<Users class="text-muted-foreground h-4 w-4" />
				</Card.Header>
				<Card.Content>
					{#await userClient.listUsers({}) then result}
						<div class="text-2xl font-bold">{result.totalCount}</div>
					{/await}
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
						{#await profileClient.listProfiles({}) then result}
							{#each result.profiles || [] as profile (profile.id)}
								<div class="flex items-center justify-between">
									<div class="flex items-center space-x-4">
										<Shield class="h-4 w-4" />
										<div class="space-y-1">
											<p class="text-sm leading-none font-medium">
												{profile.name}
											</p>
											<p class="text-muted-foreground text-xs">
												{profile.description}
											</p>
										</div>
									</div>
									<div class="flex items-center gap-2">
										{#await agentClient.listAgents({ profileId: profile.id }) then result}
											<Badge variant={result.totalCount > 0 ? 'default' : 'secondary'}>
												{result.totalCount}
												{result.totalCount === 1n ? 'Agent' : 'Agents'}
											</Badge>
										{/await}
										{#await routerClient.listRouters({ profileId: profile.id }) then result}
											<Badge variant={result.totalCount > 0 ? 'default' : 'secondary'}>
												{result.totalCount}
												{result.totalCount === 1n ? 'Router' : 'Routers'}
											</Badge>
										{/await}
										{#await serviceClient.listServices({ profileId: profile.id }) then result}
											<Badge variant={result.totalCount > 0 ? 'default' : 'secondary'}>
												{result.totalCount}
												{result.totalCount === 1n ? 'Service' : 'Services'}
											</Badge>
										{/await}
										{#await middlewareClient.listMiddlewares({ profileId: profile.id }) then result}
											<Badge variant={result.totalCount > 0 ? 'default' : 'secondary'}>
												{result.totalCount}
												{result.totalCount === 1n ? 'Middleware' : 'Middlewares'}
											</Badge>
										{/await}
									</div>
								</div>
							{/each}
						{/await}
					</div>
				</Card.Content>
			</Card.Root>

			<!-- Errors -->
			<Card.Root class="flex-1">
				<Card.Header>
					<Card.Title class="flex items-center justify-between gap-2">
						System Errors
						<!-- <Button -->
						<!-- 	variant="ghost" -->
						<!-- 	size="icon" -->
						<!-- 	class="rounded-full hover:bg-red-300" -->
						<!-- 	onclick={() => errorClient.deleteErrorsByProfile({})} -->
						<!-- > -->
						<!-- 	<Trash2 /> -->
						<!-- </Button> -->
					</Card.Title>
				</Card.Header>
				<Card.Content>
					<div class="space-y-4">
						<!-- {#await errorClient.listErrors({}) then result} -->
						<!-- 	{#each result.errors || [] as error (error.id)} -->
						<!-- 		<div class="flex items-center"> -->
						<!-- 			<div class="relative mr-4"> -->
						<!-- 				<TriangleAlert class="h-4 w-4 text-red-500" /> -->
						<!-- 			</div> -->
						<!-- 			<div class="space-y-1"> -->
						<!-- 				<p class="text-sm"> -->
						<!-- 					{error.message} -->
						<!-- 				</p> -->
						<!-- 				<p class="text-muted-foreground text-sm"> -->
						<!-- 					{error.details} -->
						<!-- 				</p> -->
						<!-- 			</div> -->
						<!-- 		</div> -->
						<!-- 	{/each} -->
						<!-- {/await} -->
					</div>
				</Card.Content>
			</Card.Root>
		</div>
	{/if}
</div>
