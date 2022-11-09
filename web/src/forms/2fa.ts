import * as yup from "yup";

const validationSchema = yup.object({
	code: yup.string().max(6).required(),
	remember: yup.boolean(),
});

const initialValues = {
	code: "",
	remember: false
};

export const TwoFactor = {
	validationSchema,
	initialValues,
};
