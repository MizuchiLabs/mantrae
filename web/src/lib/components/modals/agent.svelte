<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { upsertAgent } from '$lib/api';
	import type { Agent } from '$lib/types/base';
	import { toast } from 'svelte-sonner';

	export let agent: Agent;
	let newIP = '';

	const setActiveIP = async (ip: string) => {
		agent.activeIp = ip;
		await upsertAgent(agent);

		toast.success('IP address updated!');
	};
</script>

<Dialog.Root>
	<Dialog.Trigger class="w-full">
		<Button variant="ghost" class="w-full bg-orange-400 text-black">Details</Button>
	</Dialog.Trigger>
	<Dialog.Content class="no-scrollbar max-h-[80vh] max-w-2xl overflow-y-auto">
		<Card.Root>
			<Card.Header>
				<Card.Title class="flex items-center justify-between gap-2">Agent details</Card.Title>
				<Card.Description>
					You can update which ip address should be used for the routers reported by the agent, or
					set a custom ip.
				</Card.Description>
			</Card.Header>
			<Card.Content class="flex flex-col gap-2 text-sm">
				<div class="grid grid-cols-4 items-center gap-2">
					<span class="col-span-1 font-mono">Hostname</span>
					<div class="col-span-3 space-x-2">
						<Badge variant="secondary">{agent.hostname}</Badge>
					</div>
				</div>
				<div class="grid grid-cols-4 items-center gap-2">
					<span class="col-span-1 font-mono">Public IP</span>
					<div class="col-span-3 space-x-2">
						{#if agent.activeIp === agent.publicIp || !agent.activeIp}
							<Badge variant="default">{agent.publicIp}</Badge>
						{:else}
							<button on:click={() => setActiveIP(agent.publicIp)}>
								<Badge variant="secondary">{agent.publicIp}</Badge>
							</button>
						{/if}
					</div>
				</div>
				<div class="grid grid-cols-4 items-center gap-2">
					<span class="col-span-1 font-mono">Private IPs</span>
					<div class="col-span-3 flex flex-wrap gap-2">
						{#each agent.privateIps ?? [] as ip}
							{#if agent.activeIp === ip}
								<Badge variant="default">{ip}</Badge>
							{:else}
								<button on:click={() => setActiveIP(ip)}>
									<Badge variant="secondary">{ip}</Badge>
								</button>
							{/if}
						{/each}
					</div>
				</div>
				<div class="grid grid-cols-4 items-center gap-2">
					<span class="col-span-1 font-mono">Containers</span>
					<div class="col-span-3 flex flex-wrap gap-2">
						{#each agent.containers ?? [] as container}
							<Badge variant="secondary">{container.name.slice(1)}</Badge>
						{/each}
					</div>
				</div>
				<div class="grid grid-cols-4 items-center gap-2">
					<span class="col-span-1 font-mono">Last Seen</span>
					<div class="col-span-3 flex flex-wrap gap-2">
						<Badge variant="secondary">{new Date(agent.lastSeen).toLocaleString()}</Badge>
					</div>
				</div>
				<div class="mt-4 grid grid-cols-4 items-center gap-2">
					<Label for="ip" class="mr-2">Custom IP</Label>
					<div class="col-span-3 flex flex-wrap gap-2">
						<Input
							id="ip"
							name="ip"
							type="text"
							bind:value={newIP}
							placeholder="Enter a custom ip"
						/>
					</div>
				</div>
			</Card.Content>
		</Card.Root>
		<Dialog.Close class="w-full">
			{#if newIP}
				<Button class="w-full" on:click={() => setActiveIP(newIP)}>Save</Button>
			{:else}
				<Button class="w-full">Close</Button>
			{/if}
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
