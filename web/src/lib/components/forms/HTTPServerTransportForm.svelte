<script lang="ts">
	import type { ServersTransport } from '$lib/gen/mantrae/v1/servers_transport_pb';
	import { type ServersTransport as HTTPServersTransport } from '$lib/gen/zen/traefik-schemas';
	import { marshalConfig, unmarshalConfig } from '$lib/types';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Trash2, Plus } from '@lucide/svelte';
	import CustomSwitch from '../ui/custom-switch/custom-switch.svelte';
	import Separator from '../ui/separator/separator.svelte';
	import { parseGoDuration } from '$lib/utils';

	let { transport = $bindable() }: { transport: ServersTransport } = $props();
	let config = $state<HTTPServersTransport>(
		unmarshalConfig(transport.config) as HTTPServersTransport
	);

	$effect(() => {
		if (config) transport.config = marshalConfig(config);
	});

	function addRootCA() {
		config.rootCAs = [...(config.rootCAs || []), ''];
	}

	function removeRootCA(index: number) {
		config.rootCAs = config.rootCAs?.filter((_, i) => i !== index);
	}

	function addCertificate() {
		config.certificates = [...(config.certificates || []), { certFile: '', keyFile: '' }];
	}

	function removeCertificate(index: number) {
		config.certificates = config.certificates?.filter((_, i) => i !== index);
	}

	function addSpiffeId() {
		if (!config.spiffe) config.spiffe = {};
		config.spiffe.ids = [...(config.spiffe.ids || []), ''];
	}

	function removeSpiffeId(index: number) {
		if (config.spiffe?.ids) {
			config.spiffe.ids = config.spiffe.ids.filter((_, i) => i !== index);
		}
	}
</script>

<div class="flex flex-col gap-4">
	<Separator />

	<div class="space-y-4">
		<div class="flex flex-col gap-1">
			<Label for="serverName" class="text-sm font-medium">General</Label>
			<p class="text-xs text-muted-foreground">Configure general settings</p>
		</div>

		<div class="grid grid-cols-2 gap-2">
			<div class="space-y-2">
				<Label for="serverName">Server Name</Label>
				<Input id="serverName" bind:value={config.serverName} placeholder="server.example.com" />
			</div>
			<div class="space-y-2">
				<Label for="peerCertURI">Peer Certificate URI</Label>
				<Input
					id="peerCertURI"
					bind:value={config.peerCertURI}
					placeholder="spiffe://example.org/service"
				/>
			</div>
		</div>

		<div class="flex items-center justify-between rounded-lg border p-3">
			<div class="space-y-1">
				<Label class="text-sm">Skip TLS verification (insecure)</Label>
				<!-- <p class="text-muted-foreground text-xs">Skip TLS verification</p> -->
			</div>

			<CustomSwitch
				checked={config.insecureSkipVerify}
				onCheckedChange={(checked) => (config.insecureSkipVerify = checked)}
			/>
		</div>

		<div class="flex items-center justify-between rounded-lg border p-3">
			<div class="space-y-1">
				<Label class="text-sm">Disable HTTP/2</Label>
				<!-- <p class="text-muted-foreground text-xs">Disable HTTP/2</p> -->
			</div>
			<CustomSwitch
				checked={config.disableHTTP2}
				onCheckedChange={(checked) => (config.disableHTTP2 = checked)}
			/>
		</div>

		<div class="space-y-2">
			<div class="space-y-2">
				<Label for="maxIdleConnsPerHost">Max Idle Connections Per Host</Label>
				<Input
					id="maxIdleConnsPerHost"
					type="number"
					bind:value={config.maxIdleConnsPerHost}
					placeholder="10"
				/>
			</div>
		</div>
	</div>

	<Separator />

	<div class="space-y-2">
		<div class="flex flex-col gap-1 pb-2">
			<Label class="text-sm font-medium">Certificates</Label>
			<p class="text-xs text-muted-foreground">Configure TLS certificates</p>
		</div>

		<div class="space-y-2">
			<Label for="rootCAs">Root CAs</Label>
			{#each config.rootCAs || [] as rootCA, index (index)}
				<div class="flex gap-2">
					<Input
						value={rootCA}
						oninput={(e) => {
							let input = e.target as HTMLInputElement;
							if (!config.rootCAs) config.rootCAs = [];
							config.rootCAs[index] = input.value;
						}}
						placeholder="/path/to/ca.pem"
					/>
					<Button variant="outline" size="icon" onclick={() => removeRootCA(index)}>
						<Trash2 class="h-4 w-4" />
					</Button>
				</div>
			{/each}
			<Button variant="outline" onclick={addRootCA} class="w-full">
				<Plus />
				Add Root CA
			</Button>
		</div>

		<div class="space-y-2">
			<Label for="certificates">Client Certificates</Label>

			{#each config.certificates || [] as cert, index (index)}
				<div class="py-2">
					<Separator />
				</div>
				<div class="flex justify-between gap-2">
					<div class="grid grid-cols-2 gap-2">
						<div class="space-y-2">
							<Label>Certificate File</Label>
							<Input
								value={cert.certFile}
								oninput={(e) => {
									let input = e.target as HTMLInputElement;
									if (!config.certificates) config.certificates = [];
									config.certificates[index].certFile = input.value;
								}}
								placeholder="/path/to/cert.pem"
							/>
						</div>
						<div class="space-y-2">
							<Label>Key File</Label>
							<Input
								value={cert.keyFile}
								oninput={(e) => {
									let input = e.target as HTMLInputElement;
									if (!input.value) {
										config.certificates = config.certificates?.filter((_, i) => i !== index);
										return;
									}
									if (!config.certificates) config.certificates = [];
									config.certificates[index].keyFile = input.value;
								}}
								placeholder="/path/to/key.pem"
							/>
						</div>
					</div>
					<Button
						variant="outline"
						size="icon"
						onclick={() => removeCertificate(index)}
						class="self-end"
					>
						<Trash2 />
					</Button>
				</div>
			{/each}
			<Button variant="outline" onclick={addCertificate} class="w-full">
				<Plus />
				Add Certificate
			</Button>
		</div>
	</div>

	<Separator />

	<div class="space-y-2">
		<div class="flex flex-col gap-1 pb-2">
			<Label for="forwardingTimeouts" class="text-sm font-medium">Forwarding Timeouts</Label>
			<p class="text-xs text-muted-foreground">Configure timeouts for forwarding</p>
		</div>

		<div class="flex flex-col gap-4">
			<div class="space-y-2">
				<Label for="dialTimeout">Dial Timeout</Label>
				<Input
					id="dialTimeout"
					value={config.forwardingTimeouts?.dialTimeout}
					oninput={(e) => {
						let input = e.target as HTMLInputElement;
						if (!config.forwardingTimeouts) config.forwardingTimeouts = {};
						config.forwardingTimeouts.dialTimeout = input.value;
					}}
					onblur={(e) => {
						let input = e.target as HTMLInputElement;
						const parsed = parseGoDuration(input.value);
						if (parsed) input.value = parsed;
					}}
					placeholder="30s"
				/>
			</div>
			<div class="space-y-2">
				<Label for="responseHeaderTimeout">Response Header Timeout</Label>
				<Input
					id="responseHeaderTimeout"
					value={config.forwardingTimeouts?.responseHeaderTimeout}
					oninput={(e) => {
						let input = e.target as HTMLInputElement;
						if (!config.forwardingTimeouts) config.forwardingTimeouts = {};
						config.forwardingTimeouts.responseHeaderTimeout = input.value;
					}}
					onblur={(e) => {
						let input = e.target as HTMLInputElement;
						const parsed = parseGoDuration(input.value);
						if (parsed) input.value = parsed;
					}}
					placeholder="10s"
				/>
			</div>
		</div>
		<div class="grid grid-cols-2 gap-2">
			<div class="space-y-2">
				<Label for="idleConnTimeout">Idle Connection Timeout</Label>
				<Input
					id="idleConnTimeout"
					value={config.forwardingTimeouts?.idleConnTimeout}
					oninput={(e) => {
						let input = e.target as HTMLInputElement;
						if (!config.forwardingTimeouts) config.forwardingTimeouts = {};
						config.forwardingTimeouts.idleConnTimeout = input.value;
					}}
					onblur={(e) => {
						let input = e.target as HTMLInputElement;
						const parsed = parseGoDuration(input.value);
						if (parsed) input.value = parsed;
					}}
					placeholder="1m30s"
				/>
			</div>
			<div class="space-y-2">
				<Label for="readIdleTimeout">Read Idle Timeout</Label>
				<Input
					id="readIdleTimeout"
					value={config.forwardingTimeouts?.readIdleTimeout}
					oninput={(e) => {
						let input = e.target as HTMLInputElement;
						if (!config.forwardingTimeouts) config.forwardingTimeouts = {};
						config.forwardingTimeouts.readIdleTimeout = input.value;
					}}
					onblur={(e) => {
						let input = e.target as HTMLInputElement;
						const parsed = parseGoDuration(input.value);
						if (parsed) input.value = parsed;
					}}
					placeholder="10s"
				/>
			</div>
		</div>
		<div class="space-y-2">
			<Label for="pingTimeout">Ping Timeout</Label>
			<Input
				id="pingTimeout"
				value={config.forwardingTimeouts?.pingTimeout}
				oninput={(e) => {
					let input = e.target as HTMLInputElement;
					if (!config.forwardingTimeouts) config.forwardingTimeouts = {};
					config.forwardingTimeouts.pingTimeout = input.value;
				}}
				onblur={(e) => {
					let input = e.target as HTMLInputElement;
					const parsed = parseGoDuration(input.value);
					if (parsed) input.value = parsed;
				}}
				placeholder="15s"
			/>
		</div>
	</div>

	<Separator />

	<div class="space-y-2">
		<div class="flex flex-col gap-1 pb-2">
			<Label for="spiffe" class="text-sm font-medium">SPIFFE Configuration</Label>
			<p class="text-xs text-muted-foreground">Configure SPIFFE integration</p>
		</div>

		<div class="space-y-2">
			<Label for="trustDomain">Trust Domain</Label>
			<Input
				id="trustDomain"
				value={config.spiffe?.trustDomain}
				oninput={(e) => {
					let input = e.target as HTMLInputElement;
					if (!config.spiffe) config.spiffe = {};
					config.spiffe.trustDomain = input.value;
				}}
				placeholder="example.org"
			/>
		</div>

		<div class="space-y-2">
			<Label>SPIFFE IDs</Label>
			{#each config.spiffe?.ids || [] as spiffeId, index (index)}
				<div class="flex gap-2">
					<Input
						value={spiffeId}
						oninput={(e) => {
							let input = e.target as HTMLInputElement;
							if (!config.spiffe) config.spiffe = {};
							if (!config.spiffe.ids) config.spiffe.ids = [];
							config.spiffe.ids[index] = input.value;
						}}
						placeholder="spiffe://example.org/service"
					/>
					<Button variant="outline" size="icon" onclick={() => removeSpiffeId(index)}>
						<Trash2 class="h-4 w-4" />
					</Button>
				</div>
			{/each}
			<Button variant="outline" onclick={addSpiffeId} class="w-full">
				<Plus />
				Add SPIFFE ID
			</Button>
		</div>
	</div>
</div>
