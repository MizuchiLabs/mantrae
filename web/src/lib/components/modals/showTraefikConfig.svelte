<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import { getTraefikConfig } from '$lib/api';
	import { onMount } from 'svelte';

	let config = '';
	let rows = 10;

	onMount(async () => {
		config = await getTraefikConfig();
		rows = config.split('\n').length;
	});
</script>

<Dialog.Root>
	<Dialog.Trigger>
		<Button variant="ghost" on:click={getTraefikConfig}>
			<iconify-icon icon="devicon:traefikproxy" width="24" />
		</Button>
	</Dialog.Trigger>
	<Dialog.Content class="sm:max-w-[600px]">
		<Card.Root class="mt-4">
			<Card.Header>
				<Card.Title class="flex items-center justify-between gap-2">
					Traefik Dynamic Config
				</Card.Title>
				<Card.Description>
					This is the current dynamic configuration your traefik instance is using.
				</Card.Description>
			</Card.Header>
			<Card.Content>
				<Textarea
					value={config}
					{rows}
					class="focus-visible:ring-0 focus-visible:ring-offset-0"
					on:click={(e) => e.target?.select()}
					readonly
				/>
			</Card.Content>
		</Card.Root>
		<Dialog.Close class="w-full">
			<Button class="w-full">Close</Button>
		</Dialog.Close>
	</Dialog.Content>
</Dialog.Root>
