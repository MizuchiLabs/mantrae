export interface Plugin {
	id: string;
	name: string;
	displayName: string;
	author: string;
	type: string;
	import: string;
	summary: string;
	iconUrl: string;
	bannerUrl: string;
	readme: string;
	latestVersion: string;
	versions: string[];
	stars: number;
	snippet: Record<string, string>;
	createdAt: string;
}
