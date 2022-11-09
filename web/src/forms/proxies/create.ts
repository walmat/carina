import * as yup from "yup";

const validationSchema = yup.object({
	group: yup.object({ id: yup.string(), name: yup.string() }).required(),
	proxies: yup.string().required(),
});

const initialValues = {
	group: { id: "default", name: "default" },
	proxies: "",
};

export const Create = {
	validationSchema,
	initialValues,
};
