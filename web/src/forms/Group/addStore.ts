import * as yup from "yup";

const validationSchema = yup.object().shape({
  store: yup.object().shape({
    url: yup.string().required(),
    name: yup.string().required(),
  }),
});

const initialValues = {
  store: null,
};

export const AddStore = {
  validationSchema,
  initialValues,
};
