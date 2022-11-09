import {
  createAsyncThunk,
  createSelector,
  createSlice,
} from "@reduxjs/toolkit";

import { RootState } from ".";

export interface Store {
  url: string;
  name: string;
}

export interface Product {
  raw: string;
  variant: null | string;
  pos: null | string[];
  neg: null | string[];
  url: null | string;
}

export interface Task {
  id: string;
  type: string;
  mode: string;
  state: number;
  store: {
    name: string;
    meta: string;
    details: string;
  };
  profile: {
    name: string;
    meta: string;
    details: string;
  };
  proxies: {
    name: string;
    meta: string;
    details: string;
  };
  product: {
    name: string;
    meta: string;
    details: string;
  };
  color: null | string;
  status: string;
}

export type TaskGroup = {
  id: string;
  name: string;
  byId: {
    [id: string]: Task;
  };
};

export type Tasks = {
  [key: string]: TaskGroup;
};

const extractTableData = (data: any, isEditing = false) => {
  const base: any = {};

  console.log(data);

  if (data?.store) {
    if (data.store.type) {
      base.type = data.store.type;
    }
    base.store = {
      name: data.store.name,
      meta: data.store.url,
      details: data.mode.name,
    };
  }

  if (data?.profiles?.length) {
    base.profile = {
      name: data.profiles[0].name,
      meta: data.profiles[0].id,
      details: data.profiles[0].group,
    };
  }

  if ((isEditing && data?.proxies) || !isEditing) {
    base.proxies = {
      name: data.proxies ? data.proxies.name : "No Group",
      meta: data.proxies ? data.proxies.id : "",
      details: "None", // NOTE: No immediately assigned proxy yet
    };
  }

  if (data?.product) {
    base.product = {
      name: data.product,
      meta: data.product, // NOTE: Store this in case the product name changes
      details: data.sizes.join(","),
    };
  }

  return base;
};

export const initialTasksState: Tasks = {
  default: {
    id: "default",
    name: "Default",
    byId: {},
  },
};

export const addTaskGroup = createAsyncThunk(
  "tasks/addGroup",
  async (groupName: string) => {
    if (!groupName.trim()) {
      return null;
    }

    const id = await window.RPCAction("tasks:addGroup", [groupName.trim()]);

    return { id, groupName };
  }
);

export const removeTaskGroup = createAsyncThunk(
  "tasks/removeGroup",
  async (groups: string[]) => {
    if (!groups?.length) {
      return null;
    }

    groups.forEach((id) => {
      try {
        window.RPCAction("tasks:deleteGroup", [id]);
      } catch (e) {
        console.log(e);
        // noop...
      }
    });

    return groups;
  }
);

interface GroupMap {
  [groupId: string]: {
    [taskId: string]: string;
  };
}

export const addTasks = createAsyncThunk("tasks/add", async (values: any) => {
  const { groups, amount, ...data } = values;

  const tableData = extractTableData(data);

  const ctx = {
    default: {
      store: {
        name: values.store.name,
        url: values.store.url,
      },
      sku: values.product,
      sizes: values.sizes,
    },
    // TODO: custom inputs
  };

  const groupMap: GroupMap = await window.RPCAction("tasks:add", [
    values.store.type,
    values.mode.name,
    ctx,
    groups,
    values.profiles[0].id,
    values.proxies ? values.proxies.id : "",
    amount || 1,
  ]);

  return { data: tableData, groupMap };
});

export const copyTask = createAsyncThunk(
  "tasks/copy",
  async ({ ids, group }: any, { getState }) => {
    const { Tasks } = getState() as RootState;
    const state = Tasks[group];

    let tasks: any = {};
    for (const id of ids) {
      const original = state.byId[id];
      if (!original) {
        return;
      }

      const ctx = {
        default: {
          store: {
            name: original.store.name,
            url: original.store.meta,
          },
          sku: original.product.meta,
          sizes: original.product.details.split(","),
        },
        // TODO: custom inputs
      };

      const groupMap: GroupMap = await window.RPCAction("tasks:add", [
        original.type,
        original.mode,
        ctx,
        [{ id: group }],
        original.profile.meta,
        original.proxies.meta,
        1,
      ]);

      Object.keys(groupMap).map((groupId) => {
        Object.values(groupMap[groupId]).map((taskId) => {
          tasks[taskId] = {
            ...original,
            id: taskId,
            color: null,
            status: "Idle",
            state: 0,
          };
        });
      });
    }

    return { group, tasks };
  }
);

export const editTasks = createAsyncThunk("tasks/edit", async (values: any) => {
  // TODO: RPC Action
  const { ids, group, ...newData } = values;

  return { data: extractTableData(newData, true), group, ids };
});

export const copyToGroup = createAsyncThunk(
  "tasks/copyToGroup",
  async ({ oldGroup, newGroup, ids }: any, { getState }) => {
    const { Tasks } = getState() as RootState;
    const state = Tasks[oldGroup];

    const tasks: any[] = [];
    const newIds: string[] = [];
    for (const id of ids) {
      const original = state.byId[id];
      if (!original) {
        return;
      }

      const ctx = {
        default: {
          store: {
            name: original.store.name,
            url: original.store.meta,
          },
          sku: original.product.meta,
          sizes: original.product.details.split(","),
        },
        // TODO: custom inputs
      };

      const newId = await window.RPCAction("tasks:add", [
        original.type,
        original.mode,
        ctx,
        [{ id: newGroup }],
        original.profile.meta,
        original.proxies.meta,
        1,
      ]);

      newIds.push(...newId[newGroup]);
      tasks.push({
        ...original,
        color: null,
        status: "Idle",
        state: 0,
      });
    }

    return { group: newGroup, tasks, ids: newIds };
  }
);

export const moveToGroup = createAsyncThunk(
  "tasks/moveToGroup",
  async ({ oldGroup, newGroup, ids }: any) => {
    window.RPCAction("tasks:move", [ids, newGroup]);

    return { oldGroup, newGroup, ids };
  }
);

export const launchGroup = createAsyncThunk(
  "tasks/launch",
  async ({ group }: any) => {
    // TODO: Launch the webview, do the dance, etc.

    return { group };
  }
);

interface TaskProps {
  ids: string[];
  group: string;
}

export const removeTask = createAsyncThunk(
  "task/remove",
  async ({ ids, group }: TaskProps) => {
    window.RPCAction("tasks:delete", [ids]);

    return { ids, group };
  }
);

export const startTask = createAsyncThunk(
  "task/start",
  async ({ ids, group }: TaskProps) => {
    window.RPCAction("tasks:start", [ids]);

    return { ids, group };
  }
);

export const stopTask = createAsyncThunk(
  "task/stop",
  async ({ ids, group }: TaskProps) => {
    window.RPCAction("tasks:stop", [ids]);

    return { ids, group };
  }
);

export const tasks = createSlice({
  name: "tasks",
  initialState: initialTasksState,
  reducers: {
    importTasks: (state, action) => {},
    status: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      Object.keys(payload).forEach((group) => {
        Object.keys(payload[group]).forEach((id) => {
          const { message, color } = payload[group][id];

          state[group].byId[id] = {
            ...state[group].byId[id],
            status: message,
            color: `rgba(${color.R}, ${color.G}, ${color.B}, ${color.A / 255})`,
          };
        });
      });
    },
  },
  extraReducers: {
    // @ts-ignore
    [addTaskGroup.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { id, groupName } = payload;

      // eslint-disable-next-line no-param-reassign
      state[id] = {
        id: id,
        name: groupName,
        byId: {},
      };
    },
    // @ts-ignore
    [removeTaskGroup.fulfilled]: (state, action) => {
      const { payload } = action;

      if (!payload.length) {
        return;
      }

      payload.forEach((id: string) => {
        delete state[id];
      });
    },
    // @ts-ignore
    [addTasks.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { groupMap, data } = payload;

      Object.keys(groupMap).forEach((group: string) => {
        // @ts-ignore
        for (const taskId of groupMap[group]) {
          state[group].byId[taskId] = {
            ...data,
            id: taskId,
            status: "Idle",
            state: 0,
          };
        }
      });
    },
    // @ts-ignore
    [editTasks.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { ids, group, data } = payload;

      ids.forEach((id: string) => {
        state[group.id].byId[id] = {
          ...state[group.id].byId[id],
          ...data,
        };
      });
    },
    // @ts-ignore
    [copyTask.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { group, tasks } = payload;

      // @ts-ignore
      Object.values(tasks).forEach((task: Task) => {
        // @ts-ignore
        state[group].byId[task.id] = task;
      });
    },
    // @ts-ignore
    [removeTask.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      // TODO: Make sure the tasks are stopped first and cleaned up properly
      const { ids, group } = payload;
      ids.forEach((id: string) => {
        delete state[group].byId[id];
      });
    },
    // @ts-ignore
    [startTask.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { ids, group } = payload;
      ids.forEach((id: string) => {
        state[group].byId[id] = {
          ...state[group].byId[id],
          // @ts-ignore
          color: null,
          state: 1,
        };
      });
    },
    // @ts-ignore
    [stopTask.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { ids, group } = payload;
      ids.forEach((id: string) => {
        state[group].byId[id] = {
          ...state[group].byId[id],
          status: "Idle",
          color: null,
          state: 0,
        };
      });
    },
    // @ts-ignore
    [copyToGroup.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { group, tasks, ids } = payload;
      tasks.forEach((task: Task, i: number) => {
        state[group].byId[ids[i]] = {
          ...task,
          id: ids[i],
        };
      });
    },
    // @ts-ignore
    [moveToGroup.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      const { oldGroup, ids, newGroup } = payload;

      ids.forEach((id: string) => {
        const old = state[oldGroup].byId[id];
        delete state[oldGroup].byId[id];
        state[newGroup].byId[id] = old;
      });
    },
    // @ts-ignore
    [launchGroup.fulfilled]: (state, action) => {
      const { payload } = action;
      if (!payload) {
        return;
      }

      // TODO: whatever we need to do here?
    },
  },
});

export const selectTasks = ({ Tasks }: RootState) => Tasks;

export const makeTasks = createSelector(
  selectTasks,
  (state: RootState["Tasks"]) => state
);

const { actions, reducer } = tasks;

export const { importTasks, status } = actions;

export default reducer;
