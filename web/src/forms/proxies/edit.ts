import * as yup from "yup";

const validationSchema = yup.object({
	group: yup.object({ id: yup.string(), name: yup.string() }).required(),
	proxy: yup.string().required(),
});

const initialValues = {
	group: { id: "default", name: "default" },
	proxy: "",
};

export const Edit = {
	validationSchema,
	initialValues,
};
