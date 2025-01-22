// Some global constants
export const DateFormat = new Intl.DateTimeFormat('en-US', {
	year: 'numeric',
	month: 'long',
	day: 'numeric',
	hour: 'numeric',
	minute: 'numeric',
	second: 'numeric'
});

// Localstorage keys
export const PROFILE_SK = 'profile';
export const TOKEN_SK = 'token';
export const LIMIT_SK = 'limit';
export const LOCAL_PROVIDER_SK = 'local-provider';
export const ROUTER_COLUMN_SK = 'router-columns';
export const MIDDLEWARE_COLUMN_SK = 'middleware-columns';
export const RULE_EDITOR_TAB_SK = 'rule-editor-tab';
export const SOURCE_TAB_SK = 'traefik-source-tab';
