import { atomFamily } from "recoil";

const color = {
  R: 255,
  G: 255,
  B: 255,
  A: 1,
};

export const taskStatusAtomFamily = atomFamily({
  key: "TaskStatus",
  default: {
    status: "Idle",
    color: `rgba(${color.R}, ${color.G}, ${color.B}, ${color.A})`,
  },
});
