<script lang="ts">
	import type { Middleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let middleware: Middleware;
	middleware.forwardAuth = middleware.forwardAuth ?? {
		address: '',
		tls: { insecureSkipVerify: false, ca: '', cert: '', key: '' },
		trustForwardHeader: false,
		authResponseHeaders: [],
		authResponseHeadersRegex: '',
		authRequestHeaders: []
	};
</script>

<div class="grid grid-cols-4 items-center gap-4">
	<Label for="address" class="text-right">Address</Label>
	<Input
		id="address"
		name="address"
		type="text"
		bind:value={middleware.forwardAuth.address}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Address"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="trust-forward-header" class="text-right">Trust Forward Header</Label>
	<Switch
		id="trust-forward-header"
		bind:checked={middleware.forwardAuth.trustForwardHeader}
		class="col-span-3"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="auth-response-headers-regex" class="text-right">Auth Response Headers Regex</Label>
	<Input
		id="auth-response-headers-regex"
		name="auth-response-headers-regex"
		type="text"
		bind:value={middleware.forwardAuth.authResponseHeadersRegex}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Auth Response Headers Regex"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="auth-request-headers" class="text-right">Auth Request Headers</Label>
	<Input
		id="auth-request-headers"
		name="auth-request-headers"
		type="text"
		bind:value={middleware.forwardAuth.authRequestHeaders}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Auth Request Headers"
	/>
</div>
<ArrayInput bind:items={middleware.forwardAuth.authResponseHeaders} label="Auth Response Headers" />
<ArrayInput bind:items={middleware.forwardAuth.authRequestHeaders} label="Auth Request Headers" />
