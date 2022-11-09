import * as React from "react";
import styled from "styled-components";
import { motion, AnimatePresence } from "framer-motion";
import { useHover, useLayer } from "react-laag";

interface TooltipProps {
  children: any;
  text: string;
  triggerOffset?: number;
  autoAdjust?: boolean;
  anchor?: any;
}

function Tooltip({
  children,
  text,
  anchor = "top-center",
  autoAdjust = true,
  triggerOffset = 4,
}: TooltipProps) {
  const [isOpen, hoverProps] = useHover({
    delayEnter: 250,
    delayLeave: 125,
  });
  const { layerProps, triggerProps, renderLayer, arrowProps } = useLayer({
    isOpen,
    placement: anchor,
    auto: autoAdjust,
    triggerOffset,
  });

  const { layerSide } = arrowProps;

  return (
    <>
      <span {...triggerProps} {...hoverProps}>
        {children}
      </span>
      {renderLayer(
        <AnimatePresence>
          {isOpen && (
            <TooltipBox
              {...layerProps}
              // provide config for animated styles
              initial={{
                opacity: 0,
                scale: 0.8,
                y: layerSide === "top" ? -8 : 8,
              }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{
                opacity: 0,
                scale: 0.8,
                y: layerSide === "top" ? -8 : 8,
              }}
              transition={{
                type: "spring",
                damping: 30,
                stiffness: 800,
              }}
            >
              {text}
            </TooltipBox>
          )}
        </AnimatePresence>
      )}
    </>
  );
}

const TooltipBox = styled(motion.div)`
  color: ${({ theme }) => theme.colors.background};
  background-color: ${({ theme }) => theme.colors.tooltip};
  border: none;
  padding: 4px 12px;
  border-radius: 4px;
  transform-origin: center center;
  z-index: 9998;
  font-size: 12px;
`;

export default Tooltip;
