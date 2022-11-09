import * as yup from "yup";

const validationSchema = yup.object({
  group: yup.array(
    yup.object().shape({
      id: yup.string(),
      name: yup.string(),
    })
  ),
});

const initialValues = {
  groups: [],
};

export const Remove = {
  validationSchema,
  initialValues,
};
