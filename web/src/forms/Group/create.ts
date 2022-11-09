import * as yup from "yup";

const validationSchema = yup.object({
  name: yup.string().required(),
});

const initialValues = {
  name: "",
};

export const Create = {
  validationSchema,
  initialValues,
};
