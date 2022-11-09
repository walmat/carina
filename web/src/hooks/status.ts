import { useEffect } from "react";
import { isEmpty } from "lodash";

import { useRecoilCallback } from "recoil";
import { taskStatusAtomFamily } from '../stores/Main/reducers/tasksAtom';

interface Events {
  type: string;
  data: any;
}

interface EventHandlerProps {
  type: string;
  data: any;
}

export const useTaskStatus = () => {
  const setTaskStatus = useRecoilCallback<any, any>(
    ({ set }: any) =>
      (data: any) => {
        Object.keys(data).forEach((group) => {
          Object.keys(data[group]).forEach((id) => {
            const { message, color } = data[group][id];

            set(taskStatusAtomFamily(id), {
              status: message,
              color: `rgba(${color.R}, ${color.G}, ${color.B}, ${color.A / 255
                })`,
            });
          });
        });
      },
    []
  );

  const stopTaskRecoil = useRecoilCallback<any, any>(
    ({ set }: any) =>
      (data: any) => {
        const { taskId } = data;

        set(taskStatusAtomFamily(taskId), {
          status: "Idle",
          color: null,
        });
      },
    []
  );

  const eventHandler = ({ type, data }: EventHandlerProps) => {
    if (type != "taskEvents") {
      return;
    }

    data.forEach(({ type, data }: Events) => {
      switch (type) {
        case "status": {
          if (!isEmpty(data)) {
            // update recoil atom family
            setTaskStatus(data);
          }
          break;
        }
        case "taskStart": {
          break;
        }
        case "harvesterData": {
          break;
        }
        case "taskCancelled": {
          // set task to idle state
          stopTaskRecoil(data);
          break;
        }
        case "taskCompleted": {
          break;
        }
        case "taskProductUpdate": {
          break;
        }
      }
    });
  }

  useEffect(() => {
    astilectron.onMessage(eventHandler)

    return () => {
      // TODO: Research if we can remove astilectron.onMessage
    };
  }, []);
};
