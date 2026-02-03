import { PersistedState } from 'runed';

export const NS = 'mantrae';

export const profileID = new PersistedState<bigint>(`${NS}:profile:id`, 0n, {
	serializer: {
		serialize: (v: bigint) => v.toString(),
		deserialize: (v: string) => BigInt(v)
	}
});

// Table
export const pageIndex = new PersistedState<number>(`${NS}:page:index`, 0);
export const pageSize = new PersistedState<number>(`${NS}:page:size`, 10);
export const routerColumns = new PersistedState<string[]>(`${NS}:column:router`, []);
export const middlewareColumns = new PersistedState<string[]>(`${NS}:column:middleware`, []);

// Tabs
export const ruleTab = new PersistedState<string>(`${NS}:tab:rule`, 'simple');
