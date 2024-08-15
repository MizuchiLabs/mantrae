<script lang="ts">
	import type { HttpMiddleware } from '$lib/types/middlewares';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import ArrayInput from '../ui/array-input/array-input.svelte';

	export let middleware: HttpMiddleware;
	middleware.headers = middleware.headers ?? {
		customRequestHeaders: {},
		customResponseHeaders: {},
		accessControlAllowCredentials: false,
		accessControlAllowHeaders: [''],
		accessControlAllowMethods: [''],
		accessControlAllowOriginList: [''],
		accessControlAllowOriginListRegex: [''],
		accessControlExposeHeaders: [''],
		accessControlMaxAge: 0,
		addVaryHeader: false,
		allowedHosts: [''],
		hostsProxyHeaders: [''],
		sslRedirect: false,
		sslTemporaryRedirect: false,
		sslHost: '',
		sslProxyHeaders: {},
		sslForceHost: false,
		stsSeconds: 0,
		stsIncludeSubdomains: false,
		stsPreload: false,
		forceSTSHeader: false,
		frameDeny: false,
		customFrameOptionsValue: '',
		contentTypeNosniff: false,
		browserXssFilter: false,
		customBrowserXSSValue: '',
		contentSecurityPolicy: '',
		publicKey: '',
		referrerPolicy: '',
		featurePolicy: '',
		permissionsPolicy: ''
	};
</script>

<div class="grid grid-cols-4 items-center gap-4">
	<Label for="custom-request-headers" class="text-right">Custom Request Headers</Label>
	<Input
		id="custom-request-headers"
		name="custom-request-headers"
		type="text"
		bind:value={middleware.headers.customRequestHeaders}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Custom Request Headers"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="custom-response-headers" class="text-right">Custom Response Headers</Label>
	<Input
		id="custom-response-headers"
		name="custom-response-headers"
		type="text"
		bind:value={middleware.headers.customResponseHeaders}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Custom Response Headers"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="access-control-allow-credentials" class="text-right"
		>Access Control Allow Credentials</Label
	>
	<Switch
		id="access-control-allow-credentials"
		bind:checked={middleware.headers.accessControlAllowCredentials}
		class="col-span-3 justify-self-end"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="access-control-allow-headers" class="text-right">Access Control Allow Headers</Label>
	<Input
		id="access-control-allow-headers"
		name="access-control-allow-headers"
		type="text"
		bind:value={middleware.headers.accessControlAllowHeaders}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Access Control Allow Headers"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="access-control-allow-methods" class="text-right">Access Control Allow Methods</Label>
	<Input
		id="access-control-allow-methods"
		name="access-control-allow-methods"
		type="text"
		bind:value={middleware.headers.accessControlAllowMethods}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Access Control Allow Methods"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="access-control-allow-origin-list" class="text-right"
		>Access Control Allow Origin List</Label
	>
	<Input
		id="access-control-allow-origin-list"
		name="access-control-allow-origin-list"
		type="text"
		bind:value={middleware.headers.accessControlAllowOriginList}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Access Control Allow Origin List"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="access-control-allow-origin-list-regex" class="text-right"
		>Access Control Allow Origin List Regex</Label
	>
	<Input
		id="access-control-allow-origin-list-regex"
		name="access-control-allow-origin-list-regex"
		type="text"
		bind:value={middleware.headers.accessControlAllowOriginListRegex}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Access Control Allow Origin List Regex"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="access-control-expose-headers" class="text-right">Access Control Expose Headers</Label
	>
	<Input
		id="access-control-expose-headers"
		name="access-control-expose-headers"
		type="text"
		bind:value={middleware.headers.accessControlExposeHeaders}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Access Control Expose Headers"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="access-control-max-age" class="text-right">Access Control Max Age</Label>
	<Input
		id="access-control-max-age"
		name="access-control-max-age"
		type="number"
		bind:value={middleware.headers.accessControlMaxAge}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Access Control Max Age"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="add-vary-header" class="text-right">Add Vary Header</Label>
	<Switch
		id="add-vary-header"
		bind:checked={middleware.headers.addVaryHeader}
		class="col-span-3 justify-self-end"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="allowed-hosts" class="text-right">Allowed Hosts</Label>
	<Input
		id="allowed-hosts"
		name="allowed-hosts"
		type="text"
		bind:value={middleware.headers.allowedHosts}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Allowed Hosts"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="hosts-proxy-headers" class="text-right">Hosts Proxy Headers</Label>
	<Input
		id="hosts-proxy-headers"
		name="hosts-proxy-headers"
		type="text"
		bind:value={middleware.headers.hostsProxyHeaders}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Hosts Proxy Headers"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="ssl-redirect" class="text-right">SSL Redirect</Label>
	<Switch
		id="ssl-redirect"
		bind:checked={middleware.headers.sslRedirect}
		class="col-span-3 justify-self-end"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="ssl-temporary-redirect" class="text-right">SSL Temporary Redirect</Label>
	<Switch
		id="ssl-temporary-redirect"
		bind:checked={middleware.headers.sslTemporaryRedirect}
		class="col-span-3 justify-self-end"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="ssl-host" class="text-right">SSL Host</Label>
	<Input
		id="ssl-host"
		name="ssl-host"
		type="text"
		bind:value={middleware.headers.sslHost}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="SSL Host"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="ssl-proxy-headers" class="text-right">SSL Proxy Headers</Label>
	<Input
		id="ssl-proxy-headers"
		name="ssl-proxy-headers"
		type="text"
		bind:value={middleware.headers.sslProxyHeaders}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="SSL Proxy Headers"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="ssl-force-host" class="text-right">SSL Force Host</Label>
	<Switch
		id="ssl-force-host"
		bind:checked={middleware.headers.sslForceHost}
		class="col-span-3 justify-self-end"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="sts-seconds" class="text-right">STS Seconds</Label>
	<Input
		id="sts-seconds"
		name="sts-seconds"
		type="number"
		bind:value={middleware.headers.stsSeconds}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="STS Seconds"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="sts-include-sub-domains" class="text-right">STS Include Sub Domains</Label>
	<Switch
		id="sts-include-sub-domains"
		bind:checked={middleware.headers.stsIncludeSubdomains}
		class="col-span-3 justify-self-end"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="sts-pre-load" class="text-right">STS Pre Load</Label>
	<Switch
		id="sts-pre-load"
		bind:checked={middleware.headers.stsPreload}
		class="col-span-3 justify-self-end"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="force-sts-header" class="text-right">Force STS Header</Label>
	<Switch
		id="force-sts-header"
		bind:checked={middleware.headers.forceSTSHeader}
		class="col-span-3 justify-self-end"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="frame-deny" class="text-right">Frame Deny</Label>
	<Switch
		id="frame-deny"
		bind:checked={middleware.headers.frameDeny}
		class="col-span-3 justify-self-end"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="custom-frame-options-value" class="text-right">Custom Frame Options Value</Label>
	<Input
		id="custom-frame-options-value"
		name="custom-frame-options-value"
		type="text"
		bind:value={middleware.headers.customFrameOptionsValue}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Custom Frame Options Value"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="content-type-no-sniff" class="text-right">Content Type No Sniff</Label>
	<Switch
		id="content-type-no-snuff"
		bind:checked={middleware.headers.contentTypeNosniff}
		class="col-span-3 justify-self-end"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="browser-xss-filter" class="text-right">Browser XSS Filter</Label>
	<Switch
		id="browser-xss-filter"
		bind:checked={middleware.headers.browserXssFilter}
		class="col-span-3 justify-self-end"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="custom-browser-xss-value" class="text-right">Custom Browser XSS Value</Label>
	<Input
		id="custom-browser-xss-value"
		name="custom-browser-xss-value"
		type="text"
		bind:value={middleware.headers.customBrowserXSSValue}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Custom Browser XSS Value"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="content-security-policy" class="text-right">Content Security Policy</Label>
	<Input
		id="content-security-policy"
		name="content-security-policy"
		type="text"
		bind:value={middleware.headers.contentSecurityPolicy}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Content Security Policy"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="public-key" class="text-right">Public Key</Label>
	<Input
		id="public-key"
		name="public-key"
		type="text"
		bind:value={middleware.headers.publicKey}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Public Key"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="referrer-policy" class="text-right">Referrer Policy</Label>
	<Input
		id="referrer-policy"
		name="referrer-policy"
		type="text"
		bind:value={middleware.headers.referrerPolicy}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Referrer Policy"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="feature-policy" class="text-right">Feature Policy</Label>
	<Input
		id="feature-policy"
		name="feature-policy"
		type="text"
		bind:value={middleware.headers.featurePolicy}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Feature Policy"
	/>
</div>
<div class="grid grid-cols-4 items-center gap-4">
	<Label for="permissions-policy" class="text-right">Permissions Policy</Label>
	<Input
		id="permissions-policy"
		name="permissions-policy"
		type="text"
		bind:value={middleware.headers.permissionsPolicy}
		class="col-span-3 focus-visible:ring-0 focus-visible:ring-offset-0"
		placeholder="Permissions Policy"
	/>
</div>
