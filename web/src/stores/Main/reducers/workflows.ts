import {
  createAsyncThunk,
  createSelector,
  createSlice,
} from "@reduxjs/toolkit";

import { RootState } from ".";

export interface WorkflowNode {
  name: string;
  description: string;
  data: any;
}

export interface WorkFlow {
  id: string;
  action: WorkflowNode;
  reaction?: WorkflowNode;
  prev?: string; // previous action id
  next?: string; // next action id
}

export type WorkFlowGroup = {
  id: string;
  name: string;
  nodes: {
    [id: string]: WorkFlow;
  };
};

export type WorkFlows = {
  [key: string]: WorkFlowGroup;
};

export const initialWorkflowsState: WorkFlows = {
  default: {
    id: "default",
    name: "Default",
    nodes: {},
  },
};

export const createGroup = createAsyncThunk(
  "workflows/createGroup",
  async (groupName: string) => {
    if (!groupName.trim()) {
      return null;
    }

    // @ts-ignore
    const id = await groupsHelper.addGroup(groupName);
    return { id, groupName };
  }
);

export const removeWorkflow = createAsyncThunk(
  "workflows/removeGroup",
  async (groups: string[]) => {
    if (!groups?.length) {
      return null;
    }

    try {
      // @ts-ignore
      await workflowsHelper.removeGroup(groups);
    } catch (e) {
      console.log(e);
      // noop for now
    }

    return groups;
  }
);

export const workflows = createSlice({
  name: "workflows",
  initialState: initialWorkflowsState,
  reducers: {
    importWorkflows: (state, action) => {},
  },
  extraReducers: {
    // @ts-ignore
    [createGroup.fulfilled]: (state, action) => {
      const { payload } = action;

      if (!payload) {
        return;
      }

      const { id, groupName } = payload;

      // eslint-disable-next-line no-param-reassign
      state[id] = {
        id: id,
        name: groupName,
        nodes: {},
      };
    },
    // @ts-ignore
    [removeWorkflow.fulfilled]: (state, action) => {
      const { payload } = action;

      if (!payload.length) {
        return;
      }

      payload.forEach((id: string) => {
        delete state[id];
      });
    },
  },
});

export const selectWorkflows = ({ Workflows }: RootState) => Workflows;

export const makeWorkflows = createSelector(
  selectWorkflows,
  (state: RootState["Workflows"]) => state
);

const { actions, reducer } = workflows;

export const { importWorkflows } = actions;

export default reducer;
