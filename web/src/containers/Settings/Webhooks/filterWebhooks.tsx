import { Webhook } from '../../../stores/Main/reducers/webhooks';

export const filterWebhooks = (webhook: Webhook, search: string) => {
	const safeSearch = search.trim().toLowerCase();

	const matches = (value: string) => value.toLowerCase().includes(safeSearch);

	if (matches(webhook.url) || matches(webhook.type)) {
		return true;
	}

	if (webhook.active && matches('active')) {
		return true;
	}

	if (!webhook.active && (matches('inactive') || matches('not active'))) {
		return true;
	}

	if (webhook.declines && matches('declines')) {
		return true;
	}

	if (!webhook.declines && (matches('no declines'))) {
		return true;
	}

	if (webhook.sensitivity && matches('hide fields')) {
		return true;
	}

	return !webhook.sensitivity && (matches('show shields'));
}
