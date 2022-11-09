const transition = {
  duration: 0.15,
  ease: [0.43, 0.13, 0.23, 0.96],
};

const backdrop = {
  show: {
    opacity: 1,
    transition,
  },
  hide: {
    opacity: 0,
    transition,
  },
};

const modal = {
  show: {
    opacity: 1,
    scale: 1,
    transition,
  },
  hide: {
    opacity: 0,
    scale: 0.75,
    transition,
  },
};

export const Modal = { transition, modal, backdrop };
