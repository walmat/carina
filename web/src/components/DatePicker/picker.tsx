import React, { Dispatch, SetStateAction, useCallback, useState } from "react";
import styled from "styled-components";
import { ArrowRight } from "react-feather";
import moment from "moment";

import Input from "./input";
import Calendar from "./calendar";

import { customPickerOptions, DateOption } from ".";

interface PickerProps {
  setTimeframe: Dispatch<SetStateAction<DateOption>>;
  timeframe: DateOption;
}

const Picker = ({ timeframe, setTimeframe }: PickerProps) => {
  const [focus, setFocus] = useState<"start" | "end">("start");
  const { start, end } = timeframe.value;

  const format = useCallback(() => {}, []);

  return (
    <Container>
      <Row useBorder p="16px">
        <Col>
          <Input
            autoFocus
            focused={focus === "start"}
            id="start"
            name="start"
            onFocus={() => setFocus("start")}
            value={start ? moment(start).format("MM/DD/YYYY") : ""}
            onChange={(e) => {
              const value = e.target.value;
              const m = moment(value, "MM/DD/YYYY", true);

              if (m.isValid()) {
                setTimeframe({
                  ...timeframe,
                  ...customPickerOptions,
                  value: {
                    ...timeframe.value,
                    start: m,
                  },
                })
              }
            }}
          />
        </Col>
        <Col>
          <Through />
        </Col>
        <Col>
          <Input
            focused={focus === "end"}
            id="end"
            name="end"
            onFocus={() => setFocus("end")}
            value={end ? moment(end).format("MM/DD/YYYY") : ""}
            onChange={(e) => {
              const value = e.target.value;
              const m = moment(value, "MM/DD/YYYY", true);

              if (m.isValid()) {
                setTimeframe({
                  ...timeframe,
                  ...customPickerOptions,
                  value: {
                    ...timeframe.value,
                    end: m,
                  },
                })
              }
            }}
          />
        </Col>
      </Row>
      <Row flex={1} p="0">
        <Calendar
          focus={focus}
          setFocus={setFocus}
          start={timeframe.value.start}
          end={timeframe.value.end}
          timeframe={timeframe}
          setTimeframe={setTimeframe}
        />
      </Row>
    </Container>
  );
};

const Container = styled.div`
  display: flex;
  height: 288px;
  width: 240px;
  flex-direction: column;
  border-radius: 8px;
  background-color: ${({ theme }) => theme.colors.sidebar};
  box-shadow: rgba(0, 0, 0, 0.15) 0px 5px 15px 0px;
`;

const Row = styled.div<{ p?: string; flex?: number; useBorder?: boolean }>`
  display: flex;
  ${({ flex }) => (flex ? `flex: ${flex};` : "")}
  ${({ p }) => (p ? `padding: ${p};` : "")}
	${({ theme, useBorder }) =>
    useBorder ? `border-bottom: 1px solid ${theme.colors.separator}` : ""};
`;

const Col = styled.div`
  display: flex;
  flex-direction: column;
`;

const Through = styled(ArrowRight)`
  margin: 0 8px;
  width: 14px;
  cursor: default;
  color: ${({ theme }) => theme.colors.paragraph};

  & > * {
    cursor: default;
  }
`;

export default Picker;
