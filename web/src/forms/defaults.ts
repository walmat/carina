import * as yup from "yup";

const GenericSchema = yup.object().shape({
  store: yup
    .object()
    .shape({
      name: yup.string(),
      url: yup.string(),
    })
    .required(),
  name: yup.string().required(),
});

const ProfilesSchema = yup.object().shape({
  profileGroup: yup
    .object()
    .shape({
      id: yup.string(),
      name: yup.string(),
    })
    .required(),
  profiles: yup.array(
    yup.object().shape({
      id: yup.string(),
      name: yup.string(),
    })
  ),
  limit: yup.string(),
});

const TasksSchema = yup.object().shape({
  mode: yup.string().required(),
  proxies: yup
    .object()
    .shape({
      id: yup.string(),
      name: yup.string(),
    })
    .nullable(),
  sizes: yup.array().of(yup.string()).min(1, "Please select at least one size"),
  amount: yup.number().min(1),
  accounts: yup.array(
    yup.object().shape({
      id: yup.string(),
      name: yup.string(),
    })
  ),
  accountLimit: yup.string(),
  rate: yup.string().nullable(),
});

const validationSchema = {
  generic: GenericSchema,
  profiles: ProfilesSchema,
  tasks: TasksSchema,
};

const sections = {
  generic: {
    store: null,
    name: "",
  },
  profiles: {
    profileGroup: null,
    profiles: [],
    profileLimit: "",
  },
  tasks: {
    mode: "",
    proxies: null,
    sizes: [],
    amount: 1,
    accounts: [],
    accountLimit: "",
    rate: null,
  },
};

const initialValues = Object.values(sections).reduce(
  (prev, curr) => ({ ...prev, ...curr }),
  {}
);

export const Defaults = {
  validationSchema,
  initialValues,
  sections,
};
