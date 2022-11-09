import * as yup from "yup";

const validationSchema = yup.object({
	group: yup.object({ id: yup.string(), name: yup.string() }).required(),
	store: yup
		.object({
			label: yup.string(),
			supported: yup.boolean(),
			value: yup.string(),
		})
		.required(),
	accounts: yup.string().required(),
});

const initialValues = {
	group: { id: "default", name: "Default" },
	store: null,
	accounts: "",
};

export const Create = {
	validationSchema,
	initialValues,
};
