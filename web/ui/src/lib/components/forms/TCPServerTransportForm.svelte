<script lang="ts">
	import type { ServersTransport } from '$lib/gen/mantrae/v1/servers_transport_pb';
	import { type TCPServersTransport } from '$lib/gen/zen/traefik-schemas';
	import { marshalConfig, unmarshalConfig } from '$lib/types';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Trash2, Plus } from '@lucide/svelte';
	import CustomSwitch from '../ui/custom-switch/custom-switch.svelte';
	import Separator from '../ui/separator/separator.svelte';
	import { parseGoDuration } from '$lib/utils';

	let { transport = $bindable() }: { transport: ServersTransport } = $props();
	let config = $state<TCPServersTransport>(
		unmarshalConfig(transport.config) as TCPServersTransport
	);

	$effect(() => {
		if (config) transport.config = marshalConfig(config);
	});

	function addRootCA() {
		if (!config.tls) config.tls = {};
		config.tls.rootCAs = [...(config.tls.rootCAs || []), ''];
	}

	function removeRootCA(index: number) {
		if (config.tls?.rootCAs) {
			config.tls.rootCAs = config.tls.rootCAs.filter((_, i) => i !== index);
		}
	}

	function addCertificate() {
		if (!config.tls) config.tls = {};
		config.tls.certificates = [...(config.tls.certificates || []), { certFile: '', keyFile: '' }];
	}

	function removeCertificate(index: number) {
		if (config.tls?.certificates) {
			config.tls.certificates = config.tls.certificates.filter((_, i) => i !== index);
		}
	}

	function addSpiffeId() {
		if (!config.tls) config.tls = {};
		if (!config.tls.spiffe) config.tls.spiffe = {};
		config.tls.spiffe.ids = [...(config.tls.spiffe.ids || []), ''];
	}

	function removeSpiffeId(index: number) {
		if (config.tls?.spiffe?.ids) {
			config.tls.spiffe.ids = config.tls.spiffe.ids.filter((_, i) => i !== index);
		}
	}
</script>

<div class="flex flex-col gap-4">
	<Separator />

	<div class="space-y-2">
		<div class="flex flex-col gap-1 pb-2">
			<Label for="connectionTimeouts" class="text-sm font-medium">Connection Timeouts</Label>
			<p class="text-xs text-muted-foreground">Configure connection timeouts</p>
		</div>

		<div class="space-y-2">
			<Label for="dialKeepAlive">Dial Keep Alive</Label>
			<Input
				id="dialKeepAlive"
				value={config.dialKeepAlive}
				onblur={(e) => {
					let input = e.target as HTMLInputElement;
					const parsed = parseGoDuration(input.value);
					if (parsed) input.value = parsed;
				}}
				placeholder="15s"
			/>
		</div>
		<div class="space-y-2">
			<Label for="dialTimeout">Dial Timeout</Label>
			<Input
				id="dialTimeout"
				value={config.dialTimeout}
				onblur={(e) => {
					let input = e.target as HTMLInputElement;
					const parsed = parseGoDuration(input.value);
					if (parsed) input.value = parsed;
				}}
				placeholder="30s"
			/>
		</div>
		<div class="space-y-2">
			<Label for="terminationDelay">Termination Delay</Label>
			<Input
				id="terminationDelay"
				value={config.terminationDelay}
				onblur={(e) => {
					let input = e.target as HTMLInputElement;
					const parsed = parseGoDuration(input.value);
					if (parsed) input.value = parsed;
				}}
				placeholder="100ms"
			/>
		</div>
	</div>

	<Separator />

	<div class="space-y-2">
		<div class="flex flex-col gap-1 pb-2">
			<Label for="tls" class="text-sm font-medium">TLS Configuration</Label>
			<p class="text-xs text-muted-foreground">Configure TLS settings</p>
		</div>

		<div class="grid grid-cols-2 gap-2">
			<div class="space-y-2">
				<Label for="serverName">Server Name</Label>
				<Input
					id="serverName"
					value={config.tls?.serverName}
					oninput={(e) => {
						let input = e.target as HTMLInputElement;
						if (!config.tls) config.tls = {};
						config.tls.serverName = input.value;
					}}
					placeholder="server.example.com"
				/>
			</div>
			<div class="space-y-2">
				<Label for="peerCertURI">Peer Certificate URI</Label>
				<Input
					id="peerCertURI"
					value={config.tls?.peerCertURI}
					oninput={(e) => {
						let input = e.target as HTMLInputElement;
						if (!config.tls) config.tls = {};
						config.tls.peerCertURI = input.value;
					}}
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
				checked={config.tls?.insecureSkipVerify}
				onCheckedChange={(checked) => {
					if (!config.tls) config.tls = {};
					config.tls.insecureSkipVerify = checked;
				}}
			/>
		</div>
	</div>

	<Separator />

	<div class="space-y-2">
		<div class="flex flex-col gap-1 pb-2">
			<Label for="rootCAs" class="text-sm font-medium">Certificates</Label>
			<p class="text-xs text-muted-foreground">Configure TLS certificates</p>
		</div>

		<div class="space-y-2">
			<Label for="rootCAs">Root CAs</Label>

			{#each config.tls?.rootCAs || [] as rootCA, index (index)}
				<div class="flex gap-2">
					<Input
						value={rootCA}
						oninput={(e) => {
							let input = e.target as HTMLInputElement;
							if (!config.tls) config.tls = {};
							if (!config.tls.rootCAs) config.tls.rootCAs = [];
							config.tls.rootCAs[index] = input.value;
						}}
						placeholder="/path/to/ca.pem"
					/>
					<Button variant="outline" size="icon" onclick={() => removeRootCA(index)}>
						<Trash2 />
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

			{#each config.tls?.certificates || [] as cert, index (index)}
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
									if (!config.tls) config.tls = {};
									if (!config.tls.certificates) config.tls.certificates = [];
									config.tls.certificates[index].certFile = input.value;
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
									if (!config.tls) config.tls = {};
									if (!config.tls.certificates) config.tls.certificates = [];
									config.tls.certificates[index].keyFile = input.value;
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
			<Label for="spiffe" class="text-sm font-medium">SPIFFE Configuration</Label>
			<p class="text-xs text-muted-foreground">Configure SPIFFE integration</p>
		</div>

		<div class="space-y-2">
			<Label for="trustDomain">Trust Domain</Label>
			<Input
				id="trustDomain"
				value={config.tls?.spiffe?.trustDomain}
				oninput={(e) => {
					let input = e.target as HTMLInputElement;
					if (!config.tls) config.tls = {};
					if (!config.tls.spiffe) config.tls.spiffe = {};
					config.tls.spiffe.trustDomain = input.value;
				}}
				placeholder="example.org"
			/>
		</div>

		<div class="space-y-2">
			<Label>SPIFFE IDs</Label>
			{#each config.tls?.spiffe?.ids || [] as spiffeId, index (index)}
				<div class="flex gap-2">
					<Input
						value={spiffeId}
						oninput={(e) => {
							let input = e.target as HTMLInputElement;
							if (!config.tls) config.tls = {};
							if (!config.tls.spiffe) config.tls.spiffe = {};
							if (!config.tls.spiffe.ids) config.tls.spiffe.ids = [];
							config.tls.spiffe.ids[index] = input.value;
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
