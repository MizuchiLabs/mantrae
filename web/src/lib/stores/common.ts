import { createLocalStorage } from '$lib/storage.svelte';

export const token = createLocalStorage('auth-token', null);
export const limit = createLocalStorage('limit', '10');
export const routerColumns = createLocalStorage('router-columns', []);
export const middlewareColumns = createLocalStorage('middleware-columns', []);
export const ruleTab = createLocalStorage('rule-tab', 'simple');

export const DateFormat = new Intl.DateTimeFormat('en-US', {
	year: 'numeric',
	month: 'long',
	day: 'numeric',
	hour: 'numeric',
	minute: 'numeric',
	second: 'numeric'
});
