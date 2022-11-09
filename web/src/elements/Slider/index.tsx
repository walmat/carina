import React from "react";
import styled from "styled-components";
import {
  motion,
  useMotionValue,
  useTransform,
  useDragControls,
} from "framer-motion";

interface SliderProps {
  length?: number;
  height: number;
}

const Slider = ({ length = 100, height }: SliderProps) => {
  const handleX = useMotionValue(0);
  const progressScaleX = useTransform(handleX, [0, length], [0, 1]);
  const dragControls = useDragControls();

  return (
    <Container
      {...{ height }}
      onMouseDown={(e) => dragControls.start(e, { snapToCursor: true })}
    >
      <Bar>
        <Progress style={{ scaleX: progressScaleX }} />
      </Bar>
      <Handle
        drag="x"
        dragConstraints={{ left: 0, right: length }}
        dragControls={dragControls}
        dragElastic={0.1}
        dragMomentum={false}
        dragTransition={{
          bounceDamping: 100,
          bounceStiffness: 800,
        }}
        style={{ x: handleX }}
      />
    </Container>
  );
};

const Container = styled.div<{ height: number }>`
  height: ${({ height }) => height}px;
  flex: 1;
  width: 100%;
  display: flex;
  justify-content: flex-start;
  align-items: center;
  position: relative;
`;

const Bar = styled.div`
  width: 100%;
  height: 8px;
  background: rgba(255, 255, 255, 0.5); // TODO;
  border-radius: 8px;
  overflow: hidden;
`;

const Progress = styled(motion.div)`
  width: 100%;
  height: 100%;
  background: white;
  transform-origin: 0% 50%;
`;

const Handle = styled(motion.div)`
  width: 16px;
  height: 16px;
  position: absolute;
  border-radius: 50%;
  background: white;
  box-shadow: 2px 2px 4px rgba(0, 0, 0, 0.1);
`;

export default Slider;
