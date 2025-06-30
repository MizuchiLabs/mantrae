import type { Component } from 'svelte';
import {
	Bot,
	EthernetPort,
	Gauge,
	Globe,
	Layers,
	Route,
	Settings,
	Users,
	type IconProps
} from '@lucide/svelte';

type IconComponent = Component<IconProps, Record<string, never>, ''>;

type Routes = {
	title: string;
	url: string;
	icon: IconComponent;
	adminOnly?: boolean;
	subItems?: Routes[];
};

export const mainRoutes: Routes[] = [
	{ title: 'Dashboard', url: '/', icon: Gauge },
	{ title: 'Router', url: '/router/', icon: Route },
	{ title: 'Middlewares', url: '/middlewares/', icon: Layers },
	{ title: 'EntryPoints', url: '/entrypoints/', icon: EthernetPort }
];
export const adminRoutes: Routes[] = [
	{ title: 'Users', url: '/users/', icon: Users },
	{ title: 'Agents', url: '/agents/', icon: Bot },
	{ title: 'DNS', url: '/dns/', icon: Globe },
	{ title: 'Settings', url: '/settings/', icon: Settings }
];
export const SiteRoutes = [...mainRoutes, ...adminRoutes];
