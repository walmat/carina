import * as yup from "yup";

const validationSchema = yup.object({
  name: yup.string().required(),
  tracking: yup.string().required(),
});

const initialValues = {
  name: "",
  tracking: "",
};

export const Shipment = {
  validationSchema,
  initialValues,
};
