import * as yup from "yup";

const validationSchema = yup.object({
  store: yup
    .object({
      label: yup.string(),
      supported: yup.boolean(),
      value: yup.string(),
    })
    .required(),
  profile: yup.object({ id: yup.string(), name: yup.string() }).required(),
  value: yup.string().required(),
});

const initialValues = {
  store: null,
  profile: null,
  value: "",
};

export const Rates = {
  validationSchema,
  initialValues,
};
