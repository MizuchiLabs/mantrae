<script lang="ts">
	import * as Tooltip from '$lib/components/ui/tooltip';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button';
	import AgentModal from '$lib/components/modals/agent.svelte';
	import { agents, deleteAgent, getAgents, profile, upsertAgent } from '$lib/api';
	import { newAgent, type Agent } from '$lib/types/base';
	import { Bot, BotOff } from 'lucide-svelte';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';

	function checkLastSeen(agent: Agent) {
		const lastSeenDate = new Date(agent.lastSeen);
		const now = new Date();
		const diffInMinutes = (Number(now) - Number(lastSeenDate)) / (1000 * 60); // Explicitly convert to number
		return diffInMinutes < 1;
	}

	function addAgent() {
		let agent = newAgent();
		agent.profileId = $profile.id || 0;
		agent.lastSeen = new Date('2000-01-01').toISOString();
		upsertAgent(agent);

		toast.success('Agent added!', {
			description: 'Copy the agent token to setup your new agent.',
			duration: 3000
		});
	}

	let copyText = 'Copy Token';
	const copyToken = (agent: Agent) => {
		navigator.clipboard.writeText(agent.token);
		copyText = 'Copied!';
		setTimeout(() => {
			copyText = 'Copy Token';
		}, 2000);
	};

	onMount(async () => {
		await getAgents();
	});
</script>

{#if $profile}
	<div class="mt-4 flex flex-col gap-4 px-4 md:flex-row">
		<Button class="flex items-center gap-2 bg-red-400 text-black" on:click={addAgent}>
			<span>Add Agent</span>
			<iconify-icon icon="fa6-solid:plus" />
		</Button>
	</div>
{/if}

<div class="flex flex-col gap-4 px-4 md:flex-row">
	{#if $agents}
		{#each $agents as a}
			<Card.Root class="w-full md:w-[350px]">
				<Card.Header>
					<Card.Title class="flex items-center justify-between gap-2">
						<span>{a.hostname}</span>
						{#if checkLastSeen(a)}
							<Tooltip.Root>
								<Tooltip.Trigger>
									<Bot size="1.5rem" class="z-10 animate-pulse text-green-500" />
								</Tooltip.Trigger>
								<Tooltip.Content>
									<p>Agent connected</p>
								</Tooltip.Content>
							</Tooltip.Root>
						{:else}
							<Tooltip.Root>
								<Tooltip.Trigger>
									<BotOff size="1.5rem" class="text-red-500" />
								</Tooltip.Trigger>
								<Tooltip.Content>
									<p>Agent disconnected</p>
								</Tooltip.Content>
							</Tooltip.Root>
						{/if}
					</Card.Title>
				</Card.Header>
				{#if a.publicIp}
					<Card.Content class="space-y-2">
						IP: {a.publicIp}
					</Card.Content>
				{/if}
				<Card.Footer class="flex flex-col items-center gap-4">
					<div class="flex w-full flex-row gap-2">
						<AgentModal agent={a} />
						<Button
							variant="ghost"
							class="w-full bg-red-400 text-black"
							on:click={() => deleteAgent(a.id)}>Delete</Button
						>
					</div>
					<Button variant="secondary" size="sm" class="w-full" on:click={() => copyToken(a)}
						>{copyText}</Button
					>
				</Card.Footer>
			</Card.Root>
		{/each}
	{:else}
		<div class="flex h-[calc(75vh)] w-full flex-col items-center justify-center gap-2">
			<BotOff size="8rem" />
			<span class="text-xl font-semibold">No agents connected</span>
		</div>
	{/if}
</div>
