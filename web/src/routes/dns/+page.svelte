<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Button } from '$lib/components/ui/button';
	import CreateProvider from '$lib/components/modals/createProvider.svelte';
	import UpdateProvider from '$lib/components/modals/updateProvider.svelte';
	import { deleteProvider, provider } from '$lib/api';
</script>

<CreateProvider />

<div class="flex flex-row items-center gap-2">
	{#each Object.values($provider ?? []) as p}
		<Card.Root class="w-[400px]">
			<Card.Header>
				<Card.Title class="flex items-center justify-between gap-2">
					<span>{p.name}</span>
					<Badge variant="secondary" class="bg-blue-400">
						{p.type}
					</Badge>
				</Card.Title>
			</Card.Header>
			<Card.Content class="space-y-2"></Card.Content>
			<Card.Footer class="grid grid-cols-2 items-center gap-2">
				<Button variant="ghost" class="w-full bg-red-400" on:click={() => deleteProvider(p.name)}
					>Delete</Button
				>
				<UpdateProvider {p} />
			</Card.Footer>
		</Card.Root>
	{/each}
</div>
