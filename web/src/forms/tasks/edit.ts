import * as yup from "yup";

import { TaskGroup } from "../../stores/Main/reducers/tasks";
import { Store } from "../../stores/Main/reducers/stores";
import { Profile } from "../../stores/Main/reducers/profiles";

const validationSchema = yup.object({
  group: yup.object().shape({
    id: yup.string(),
    name: yup.string(),
  }),
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
  sizes: yup.array(yup.string()),
  product: yup.string(),
  profiles: yup
    .array(
      yup.object().shape({
        id: yup.string(),
        name: yup.string(),
      })
    )
    .nullable(),
  proxies: yup
    .object()
    .shape({
      id: yup.string(),
      name: yup.string(),
    })
    .nullable(),
});

const initialValues = (group: TaskGroup) => ({
  group: group,
  store: null,
  mode: null,
  sizes: [],
  product: "",
  profiles: [],
  proxies: null,
});

interface Proxies {
  id: string;
  name: string;
}

interface StrippedTaskGroup {
  id: string;
  name: string;
}

export interface EditTaskForm {
  group: null | StrippedTaskGroup;
  store: null | Store;
  mode: null | string;
  sizes: string[];
  product: string;
  profiles: Profile[];
  proxies: null | Proxies;
}

export const Edit = {
  validationSchema,
  initialValues,
};
