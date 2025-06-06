import { createLocalStorage } from '$lib/storage.svelte';

export const limit = createLocalStorage('limit', '10');
export const routerColumns = createLocalStorage('router_columns', []);
export const middlewareColumns = createLocalStorage('middleware_columns', []);
export const ruleTab = createLocalStorage('rule_tab', 'simple');

export const DateFormat = new Intl.DateTimeFormat('en-US', {
	year: 'numeric',
	month: 'long',
	day: 'numeric',
	hour: 'numeric',
	minute: 'numeric',
	second: 'numeric'
});
