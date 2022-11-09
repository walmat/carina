const container = {
  hidden: { opacity: 0 },
  show: {
    opacity: 1,
    transition: {
      staggerChildren: 0.5,
    },
  },
};

const item = {
  hidden: {
    opacity: 0,
    y: 5,
    transition: { ease: [0.87, 0, 0.13, 1] },
  },
  show: {
    opacity: 1,
    y: 0,
    transition: { ease: [0.87, 0, 0.13, 1] },
  },
};

export const InfiniteScroll = { container, item };
