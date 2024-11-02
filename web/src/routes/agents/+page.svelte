<script lang="ts">
	import * as Tooltip from '$lib/components/ui/tooltip';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Button } from '$lib/components/ui/button';
	import AgentModal from '$lib/components/modals/agent.svelte';
	import {
		getAgents,
		agents,
		settings,
		profile,
		getAgentToken,
		getSettings,
		deleteAgent
	} from '$lib/api';
	import { type Agent } from '$lib/types/base';
	import { Bot, BotOff, Copy, CopyCheck } from 'lucide-svelte';
	import HoverInfo from '$lib/components/utils/hoverInfo.svelte';

	let token = '';

	let copyText = 'Copy';
	const copyToken = () => {
		navigator.clipboard.writeText(token);
		copyText = 'Copied!';
		setTimeout(() => {
			copyText = 'Copy';
		}, 2000);
	};

	function checkLastSeen(agent: Agent) {
		const lastSeenDate = new Date(agent.lastSeen);
		const now = new Date();
		const diffInMinutes = (Number(now) - Number(lastSeenDate)) / (1000 * 60); // Explicitly convert to number
		return diffInMinutes < 1;
	}

	// Get routers when the profile changes
	profile.subscribe(async (value) => {
		if (!value?.id) return;
		await getAgents(value.id);

		if (!token) {
			token = await getAgentToken();
		}
		if (!$settings) {
			await getSettings();
		}
	});
</script>

<div class="my-6 flex flex-col gap-2 px-4">
	<h1 class="text-xl font-semibold">
		Agent Token <HoverInfo
			text="The token used to authenticate with the server. Initially it will be valid for 14 days. The agents will renew their token automatically."
		/>
	</h1>
	<div class="relative flex max-w-[380px]">
		<Input
			id="agent-token"
			name="agent-token"
			type="text"
			bind:value={token}
			class="overflow-hidden text-ellipsis whitespace-nowrap pr-10"
			readonly
		/>
		<Button
			variant="ghost"
			class="absolute right-0 hover:bg-transparent hover:text-red-400"
			on:click={copyToken}
		>
			{#if copyText === 'Copied!'}
				<CopyCheck size="1rem" />
			{:else}
				<Copy size="1rem" />
			{/if}
		</Button>
	</div>
</div>

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
				<Card.Content class="space-y-2">
					IP: {a.publicIp}
				</Card.Content>
				<Card.Footer class="flex items-center gap-2">
					<AgentModal agent={a} />
					<Button
						variant="ghost"
						class="w-full bg-red-400 text-black"
						on:click={() => deleteAgent(a.id)}>Delete</Button
					>
				</Card.Footer>
			</Card.Root>
		{/each}
	{/if}
</div>
