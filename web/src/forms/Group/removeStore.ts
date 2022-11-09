import * as yup from "yup";

const validationSchema = yup.object({
  stores: yup.array(
    yup.object().shape({
      name: yup.string(),
      url: yup.string(),
    })
  ),
});

const initialValues = {
  stores: [],
};

export const RemoveStore = {
  validationSchema,
  initialValues,
};
