import * as React from "react";
import { motion } from "framer-motion";

type Props = {
  isOpen: boolean;
};

const Path = (props: any) => (
  <motion.path
    fill="transparent"
    strokeWidth="2"
    stroke="#786EF2"
    strokeLinecap="round"
    {...props}
  />
);

export const Toggle = ({ isOpen }: Props) => (
  <svg width="14" height="14" viewBox="0 0 23 23">
    <Path
      initial="closed"
      animate={isOpen ? "open" : "closed"}
      variants={{
        closed: { d: "M 2 2.5 L 20 2.5" },
        open: { d: "M 3 16.5 L 17 2.5" },
      }}
    />
    <Path
      initial="closed"
      animate={isOpen ? "open" : "closed"}
      d="M 2 9.423 L 20 9.423"
      variants={{
        closed: { opacity: 1 },
        open: { opacity: 0 },
      }}
      transition={{ duration: 0.1 }}
    />
    <Path
      initial="closed"
      animate={isOpen ? "open" : "closed"}
      variants={{
        closed: { d: "M 2 16.346 L 20 16.346" },
        open: { d: "M 3 2.5 L 17 16.346" },
      }}
    />
  </svg>
);
