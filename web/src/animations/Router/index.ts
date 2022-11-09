const transition = {
  duration: 0.15,
  ease: [0.43, 0.13, 0.23, 0.96],
};

const variants = {
  enter: {
    opacity: 1,
    transition,
  },
  exit: {
    opacity: 0,
    transition,
  },
};

export const Router = { transition, variants };
