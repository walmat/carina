import * as yup from "yup";

const validationSchema = yup.object({
	store: yup.object({
		name: yup.string(),
		url: yup.string(),
	}),
	account: yup.string().required(),
});

const initialValues = {
	group: { id: "default", name: "Default" },
	store: null,
	account: "",
};

export const Edit = {
	validationSchema,
	initialValues,
};
