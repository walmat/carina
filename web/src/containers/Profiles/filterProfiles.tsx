import { Profile } from '../../stores/Main/reducers/profiles';

export const filterProfiles = (profile: Profile, search: string) => {
	const safeSearch = search.trim().toLowerCase();

	const matches = (value: string) => value.toLowerCase().includes(safeSearch);

	if (matches(profile.name)) {
		return true;
	}

	if (matches(profile.shipping.name)) {
		return true;
	}

	if (matches(profile.shipping.city)) {
		return true;
	}

	if (matches(profile.shipping.country.name) || matches(profile.shipping.country.code)) {
		return true;
	}

	if (matches(profile.shipping.state.name) || matches(profile.shipping.state.code)) {
		return true;
	}

	if (matches(profile.shipping.line1) || matches(profile.shipping.line2) || matches(profile.shipping.line3)) {
		return true;
	}

	if (matches(profile.shipping.postCode)) {
		return true;
	}

	if (matches(profile.billing.name)) {
		return true;
	}

	if (matches(profile.billing.city)) {
		return true;
	}

	if (matches(profile.billing.country.name) || matches(profile.billing.country.code)) {
		return true;
	}

	if (matches(profile.billing.state.name) || matches(profile.billing.state.code)) {
		return true;
	}

	if (matches(profile.billing.line1) || matches(profile.billing.line2) || matches(profile.billing.line3)) {
		return true;
	}

	if (matches(profile.billing.postCode)) {
		return true;
	}

	if (matches(profile.payment.name)) {
		return true;
	}

	if (matches(profile.payment.email)) {
		return true;
	}

	if (matches(profile.payment.phone)) {
		return true;
	}

	if (matches(profile.payment.type)) {
		return true;
	}

	if (matches(profile.payment.number)) {
		return true;
	}

	if (matches(profile.payment.cvv)) {
		return true;
	}

	if (matches(profile.payment.expMonth)) {
		return true;
	}

	return matches(profile.payment.expYear);
}
