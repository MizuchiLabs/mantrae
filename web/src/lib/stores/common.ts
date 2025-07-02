import { createLocalStorage } from "$lib/storage.svelte";
import { PUBLIC_BACKEND_URL } from "$env/static/public";

export const DEFAULT_URL = import.meta.env.PROD ? "/" : `http://127.0.0.1:3000`;

export const baseURL = createLocalStorage("base_url", DEFAULT_URL);
export const pageIndex = createLocalStorage("page_index", 0);
export const pageSize = createLocalStorage("page_size", 10);
export const routerColumns = createLocalStorage("router_columns", []);
export const middlewareColumns = createLocalStorage("middleware_columns", []);
export const ruleTab = createLocalStorage("rule_tab", "simple");

export const DateFormat = new Intl.DateTimeFormat("en-US", {
	year: "numeric",
	month: "long",
	day: "numeric",
	hour: "numeric",
	minute: "numeric",
	second: "numeric",
});
