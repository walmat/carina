import * as yup from "yup";

const validationSchema = yup.object({
  token: yup.string().required(),
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
  token: "",
  password: "",
};

export const ChangePassword = {
  validationSchema,
  initialValues,
};
