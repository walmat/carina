import * as yup from "yup";

const validationSchema = yup.object({
  key: yup.string().max(27).required(),
  email: yup.string().max(60).email().required(),
  password: yup
    .string()
    .min(8)
    .max(60)
    .matches(
      /^.*(?=.{8,})((?=.*[!@#$%^&*()\-_=+{};:,<.>]){1})(?=.*\d)((?=.*[a-z]){1})((?=.*[A-Z]){1}).*$/
    )
    .required(),
  terms: yup.bool().oneOf([true]).required(),
});

const initialValues = {
  key: "",
  email: "",
  password: "",
  terms: false,
};

export const Register = {
  validationSchema,
  initialValues,
};
