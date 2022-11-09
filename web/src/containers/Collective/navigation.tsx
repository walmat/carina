import styled from "styled-components";
import { ChevronLeft, ChevronRight } from "react-feather";

import { useLongPress } from "../../hooks";

const sideScroll = (
  element: HTMLDivElement,
  speed: number,
  distance: number,
  step: number
) => {
  let scrollAmount = 0;
  const slideTimer = setInterval(() => {
    element.scrollLeft += step;
    scrollAmount += Math.abs(step);
    if (scrollAmount >= distance) {
      clearInterval(slideTimer);
    }
  }, speed);
};

const defaultOptions = {
  shouldPreventDefault: true,
  delay: 500,
};

const LeftArrow = ({ contentWrapper }: any) => {
  const onLongPress = () => {
    sideScroll(contentWrapper.current, 50, 150, -50);
  };

  const onClick = () => {
    sideScroll(contentWrapper.current, 50, 150, -50);
  };

  const longPressEvent = useLongPress(onLongPress, onClick, defaultOptions);

  return <Left {...longPressEvent} />;
};

const RightArrow = ({ contentWrapper }: any) => {
  const onLongPress = () => {
    sideScroll(contentWrapper.current, 50, 150, 50);
  };

  const onClick = () => {
    sideScroll(contentWrapper.current, 50, 150, 50);
  };

  const longPressEvent = useLongPress(onLongPress, onClick, defaultOptions);

  return <Right {...longPressEvent} />;
};

const Navigation = ({ contentWrapper }: any) => (
  <Row m="0">
    <LeftArrow {...{ contentWrapper }} />
    <RightArrow {...{ contentWrapper }} />
  </Row>
);

const Row = styled.div<{ m?: string }>`
  display: flex;
  justify-content: center;
  align-items: center;
  ${({ m }) =>
    typeof m !== "undefined"
      ? `
		margin: ${m};
	`
      : ""}
`;

const Right = styled(ChevronRight)`
  display: flex;
  width: auto;
  height: 16px;
  margin: auto 8px auto 4px;
  color: ${({ theme }) => theme.colors.h2};
`;

const Left = styled(ChevronLeft)`
  display: flex;
  width: auto;
  height: 16px;
  stroke-width: 2.5px;
  margin: auto 4px auto 0;
  color: ${({ theme }) => theme.colors.h2};
`;

export default Navigation;
