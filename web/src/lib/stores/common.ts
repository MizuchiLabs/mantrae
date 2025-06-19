import { createLocalStorage, LocalStorage } from '$lib/storage.svelte';

export const token: LocalStorage<string | null> = createLocalStorage('auth_token', null);
export const userId: LocalStorage<string | null> = createLocalStorage('user_id', null);
export const pageIndex: LocalStorage<number> = createLocalStorage('page_index', 0);
export const pageSize: LocalStorage<number> = createLocalStorage('page_size', 10);
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
