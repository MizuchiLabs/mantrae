import { toast } from 'svelte-sonner';
import { useMutation, useQuery } from '$lib/query';
import { SettingService } from '$lib/gen/mantrae/v1/setting_pb';

export const setting = {
	// Queries
	get: (key: string) => useQuery(SettingService.method.getSetting, { key }, {}),
	list: () => useQuery(SettingService.method.listSettings, {}, {}),

	// Mutations
	update: () =>
		useMutation(SettingService.method.updateSetting, {
			onSuccess: () => toast.success('Settings updated!')
		})
};
