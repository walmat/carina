import * as yup from "yup";

const validationSchema = yup.object({
	email: yup.string().max(60).email().required(),
	password: yup
		.string()
		.min(8)
		.max(60)
		.matches(
			/^.*(?=.{8,})((?=.*[!@#$%^&*()\-_=+{};:,<.>]){1})(?=.*\d)((?=.*[a-z]){1})((?=.*[A-Z]){1}).*$/
		)
		.required(),
});

const initialValues = {
	email: "",
	password: ""
};

export const Login = {
	validationSchema,
	initialValues,
};
