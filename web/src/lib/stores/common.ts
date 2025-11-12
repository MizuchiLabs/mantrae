import { env } from '$env/dynamic/public';
import { createLocalStorage } from '$lib/storage.svelte';

const fallbackURL = import.meta.env.PROD ? window?.location?.origin : 'http://localhost:3000';

export const BackendURL = env.PUBLIC_BACKEND_URL || fallbackURL;
export const pageIndex = createLocalStorage('page_index', 0);
export const pageSize = createLocalStorage('page_size', 10);
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
