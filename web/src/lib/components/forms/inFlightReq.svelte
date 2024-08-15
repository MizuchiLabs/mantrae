<script lang="ts">
	import type { HttpMiddleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import * as Card from '$lib/components/ui/card';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let middleware: HttpMiddleware;
	middleware.inFlightReq = middleware.inFlightReq ?? {
		amount: 0,
		sourceCriterion: {
			ipStrategy: {
				depth: 0,
				excludedIPs: []
			},
			requestHeaderName: '',
			requestHost: false
		}
	};
</script>

<div class="grid grid-cols-4 items-center gap-4">
	<Label for="amount" class="text-right">Amount</Label>
	<Input
		id="amount"
		name="amount"
		type="number"
		bind:value={middleware.inFlightReq.amount}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Amount"
	/>
</div>

<Card.Root>
	<Card.Header>
		<Card.Title>Source Criterion</Card.Title>
		<Card.Description>Add a source criterion to limit the rate of requests.</Card.Description>
	</Card.Header>
	<Card.Content class="space-y-2">
		<div class="grid grid-cols-5 items-center gap-4">
			<Label for="depth" class="col-span-2 text-right">Depth</Label>
			<Input
				id="depth"
				name="depth"
				type="number"
				bind:value={middleware.inFlightReq.sourceCriterion.ipStrategy.depth}
				class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
				placeholder="Depth"
			/>
		</div>
		<div class="grid grid-cols-5 items-center gap-2">
			<Label for="request-header-name" class="col-span-2 text-right">Request Header Name</Label>
			<Input
				id="request-header-name"
				name="request-header-name"
				type="text"
				bind:value={middleware.inFlightReq.sourceCriterion.requestHeaderName}
				class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
				placeholder="Request Header Name"
			/>
		</div>
		<div class="grid grid-cols-5 items-center gap-2">
			<Label for="request-host" class="col-span-2 text-right">Request Host</Label>
			<Switch
				id="request-host"
				bind:checked={middleware.inFlightReq.sourceCriterion.requestHost}
				class="col-span-3 justify-self-end"
			/>
		</div>
		<ArrayInput
			bind:items={middleware.inFlightReq.sourceCriterion.ipStrategy.excludedIPs}
			label="Excluded IPs"
		/>
	</Card.Content>
</Card.Root>
