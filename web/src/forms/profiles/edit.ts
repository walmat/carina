import * as yup from "yup";
import valid from "card-validator";

const AddressSchema = (step: "shipping" | "billing") =>
	yup.object({
		[step]: yup.object().shape({
			name: yup.string(),
			line1: yup.string(),
			line2: yup.string(),
			line3: yup.string(),
			postCode: yup.string(),
			city: yup.string(),
			country: yup
				.object()
				.shape({
					name: yup.string(),
					code: yup.string(),
				})
				.nullable(),
			state: yup
				.object()
				.shape({
					name: yup.string(),
					code: yup.string(),
				})
				.nullable(),
		}),
	});

const PaymentSchema = yup.object().shape({
	group: yup
		.object()
		.shape({
			id: yup.string(),
			name: yup.string(),
		})
		.nullable(),
	name: yup.string(),
	payment: yup.object().shape({
		email: yup.string().email(),
		phone: yup.string(),
		name: yup.string(),
		number: yup
			.string()
			.test(
				"test-credit-card-number",
				"Invalid card number",
				(value) => valid.number(value).isValid
			),
		exp: yup
			.string()
			.max(5)
			.matches(/([0-9]{2})\/([0-9]{2})/)
			.test(
				"test-credit-card-expiration-date",
				"Invalid expiration date",
				(expirationDate) => {
					if (!expirationDate) {
						return false;
					}

					const today = new Date();
					const monthToday = today.getMonth() + 1;
					const yearToday = today.getFullYear().toString().substr(-2);

					const [expMonth, expYear] = expirationDate.split("/");

					if (Number(expYear) < Number(yearToday)) {
						return false;
					} else if (
						Number(expMonth) < monthToday &&
						Number(expYear) <= Number(yearToday)
					) {
						return false;
					}

					return true;
				}
			)
			.test(
				"test-credit-card-expiration-date",
				"Invalid expiration month",
				(expirationDate) => {
					if (!expirationDate) {
						return false;
					}

					const [expMonth] = expirationDate.split("/");
					return Number(expMonth) <= 12;
				}
			),
		cvv: yup
			.string()
			.test("test-credit-card-cvv", "Invalid security code", (value) => {
				if (!value?.length) {
					return false;
				}

				return value.length >= 3;
			}),
	}),
});

const validationSchema = {
	shipping: AddressSchema("shipping"),
	billing: AddressSchema("billing"),
	payment: PaymentSchema,
};

const sections = {
	shipping: {
		shipping: {
			name: "",
			line1: "",
			line2: "",
			line3: "",
			postCode: "",
			city: "",
			country: null,
			state: null,
		},
	},
	billing: {
		billing: {
			name: "",
			line1: "",
			line2: "",
			line3: "",
			postCode: "",
			city: "",
			country: null,
			state: null,
		},
	},
	payment: {
		name: "",
		group: null,
		payment: {
			name: "",
			email: "",
			phone: "",
			type: "",
			number: "",
			exp: "",
			cvv: "",
		},
	},
};

const initialValues = Object.values(sections).reduce(
	(prev, curr) => ({ ...prev, ...curr }),
	{}
);

export const Edit = {
	validationSchema,
	initialValues,
	sections,
};
