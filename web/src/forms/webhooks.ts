import * as yup from "yup";

import { fields } from '../constants';

const validationSchema = yup.object({
  profiles: yup.array().of(
    yup.object().shape({
      id: yup.string(),
      name: yup.string(),
    })
  ).required(),
  name: yup.string().required(),
  url: yup.string().required(),
  fields: yup.array().of(
    yup.object().shape({
      name: yup.string(),
      enabled: yup.boolean()
    })
  ).required(),
  declines: yup.boolean(),
  sensitivity: yup.boolean()
});

const initialValues = {
  profiles: [{ id: 'default', name: 'Default' }],
  name: "",
  url: "",
  fields,
  declines: true,
  sensitivity: true
};

export const Webhooks = {
  validationSchema,
  initialValues,
};
