import * as yup from "yup";

const validationSchema = yup.object({
	email: yup.string().max(60).email().required()
});

const initialValues = {
	email: ""
};

export const ForgotPassword = {
	validationSchema,
	initialValues,
};
