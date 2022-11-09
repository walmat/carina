import * as yup from "yup";

import { Store, Mode } from "../../stores/Main/reducers/stores";
import { Profile } from "../../stores/Main/reducers/profiles";
import { TaskGroup } from "../../stores/Main/reducers/tasks";

const validationSchema = yup.object({
  group: yup.array(
    yup.object().shape({
      id: yup.string(),
      name: yup.string(),
    })
  ),
  store: yup
    .object()
    .shape({
      name: yup.string(),
      url: yup.string(),
    })
    .nullable(),
  mode: yup
    .object()
    .shape({
      name: yup.string(),
      label: yup.string(),
    })
    .nullable(),
  sizes: yup.array(yup.string()).required(),
  product: yup.string().required(),
  profiles: yup.array(
    yup.object().shape({
      id: yup.string(),
      name: yup.string(),
    })
  ),
  proxies: yup
    .object()
    .shape({
      id: yup.string(),
      name: yup.string(),
    })
    .nullable(),
});

const initialValues = (group: TaskGroup) => ({
  groups: [group],
  store: null,
  mode: null,
  sizes: [],
  product: "",
  profiles: [],
  proxies: null,
});

type Proxies = {
  id: string;
  name: string;
};

type StrippedTaskGroup = {
  id: string;
  name: string;
};

export interface CreateTaskForm {
  groups: StrippedTaskGroup[];
  store: null | Store;
  mode: null | Mode;
  sizes: string[];
  product: string;
  profiles: Profile[];
  proxies: null | Proxies;
}

export const Create = {
  validationSchema,
  initialValues,
};
