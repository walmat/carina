import React, {
  Dispatch,
  SetStateAction,
  useEffect,
  useState,
  useMemo,
  useCallback,
} from "react";
import styled from "styled-components";
import { Calendar } from "react-feather";
import { Popover } from "react-tiny-popover";
import moment, { Moment } from "moment";

import { Typography, DateSelect } from "../../elements";

import Picker from "./picker";

export const customPickerOptions = {
  id: "custom",
  label: "Custom",
};

// @ts-ignore
const DATE_PICKER_OPTIONS: DateOption = [
  {
    id: "day",
    label: "Today",
    value: {
      start: moment().startOf("day"),
      end: moment().endOf("day"),
    },
  },
  {
    id: "week",
    label: "Last 7 days",
    value: {
      start: moment().subtract("6", "days"),
      end: moment().endOf("day"),
    },
  },
  {
    id: "month",
    label: "Last 4 weeks",
    value: {
      start: moment().subtract("29", "days"),
      end: moment().endOf("day"),
    },
  },
  {
    id: "3mo",
    label: "Last 3 months",
    value: {
      start: moment().subtract("3", "months").add("1", "day"),
      end: moment().endOf("day"),
    },
  },
  {
    id: "year",
    label: "Last 12 months",
    value: {
      start: moment().subtract("11", "months").add("1", "day"),
      end: moment().endOf("day"),
    },
  },
  {
    id: "all",
    label: "All time",
    value: {
      start: "∞",
      end: "∞",
    },
  },
];

// @ts-ignore
const COMPARISON_OPTIONS: DateOption = [
  {
    label: "No comparison",
    value: {
      start: null,
      end: null,
    },
  },
  {
    label: "Previous period",
    value: {
      start: moment().startOf("day"),
      end: moment().endOf("day"),
    },
  },
  {
    label: "Previous month",
    value: {
      start: moment().subtract("7", "days"),
      end: moment().endOf("day"),
    },
  },
  {
    label: "Previous year",
    value: {
      start: moment().subtract("4", "weeks"),
      end: moment().endOf("day"),
    },
  },
];

export interface DateValue {
  start: Moment;
  end: Moment | null;
}

export interface DateOption {
  label: string;
  value: DateValue;
}

interface DatePickerProps {
  positions?: string[];
  timeframe: DateOption;
  setTimeframe: Dispatch<SetStateAction<DateOption>>;
}

const StartTime = ({ timeframe, diff }: any) => {
  const { start } = timeframe.value;

  if (typeof start === "string") {
    return <>{start}</>;
  }

  return (
    <>
      {moment(timeframe.value.start).format("MMM")}{" "}
      {moment(timeframe.value.start).format("DD")}
      {diff ? `, ${moment(timeframe.value.start).format("YYYY")}` : ""}
    </>
  );
};

const EndTime = ({ timeframe }: any) => {
  const { end } = timeframe.value;

  if (typeof end === "string") {
    return <>{end}</>;
  }

  return (
    <>
      {moment(timeframe.value.end).format("MMM")}{" "}
      {moment(timeframe.value.end).format("DD")}
    </>
  );
};

const DatePicker = ({ timeframe, setTimeframe }: DatePickerProps) => {
  const [diff, setDiff] = useState(false);
  const [open, setOpen] = useState(false);

  const togglePicker = useCallback(() => {
    setOpen((prev) => !prev);
  }, []);

  useEffect(() => {
    // @ts-ignore
    const { start, end } = timeframe.value;

    const startYear = moment(start).year();
    const endYear = moment(end).year();

    if (diff) {
      setDiff(false);
    }

    if (startYear !== endYear) {
      return setDiff(true);
    }
  }, [diff, timeframe]);

  return useMemo(
    () => (
      <Row>
        <DateSelect
          required
          type="first"
          isClearable={false}
          isMulti={false}
          placeholder="Today"
          name="datetime"
          // @ts-ignore
          options={DATE_PICKER_OPTIONS}
          onChange={setTimeframe}
          value={timeframe}
        />
        <Popover
          isOpen={open}
          padding={4}
          reposition={false}
          positions={["bottom"]}
          align="end"
          onClickOutside={togglePicker}
          content={<Picker setTimeframe={setTimeframe} timeframe={timeframe} />}
        >
          <CustomDatePicker onClick={togglePicker}>
            <Cal />
            <DateTime>
              <StartTime timeframe={timeframe} diff={diff} />
              <span> – </span>
              <EndTime timeframe={timeframe} />
            </DateTime>
          </CustomDatePicker>
        </Popover>
      </Row>
    ),
    [timeframe, open, diff]
  );
};

DatePicker.DATE_PICKER_OPTIONS = DATE_PICKER_OPTIONS;
DatePicker.COMPARISON_OPTIONS = COMPARISON_OPTIONS;

const Row = styled.div`
  display: flex;
  box-shadow: 0px 1px 2px rgba(0, 0, 0, 0.1);
`;

const CustomDatePicker = styled.div`
  height: 28px;
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 0 8px;
  box-sizing: border-box;
  border-radius: 0 4px 4px 0;
  border: 1px solid ${({ theme }) => theme.colors.border};
  transition: background-color, color, border 0.15 ease-in-out;
  margin-left: -1px;

  & > * {
    cursor: pointer;
  }

  &:hover {
    // hack to show full border
    z-index: 900;
    border-color: hsl(0, 0%, 70%);
  }

  &:active {
    border-color: ${({ theme }) => theme.colors.h2};
  }
`;

const Cal = styled(Calendar)`
  margin: 0 8px 0 0;
  width: 14px;
  color: ${({ theme }) => theme.colors.paragraph};
`;

const DateTime = styled(Typography.Paragraph)`
  font-size: 12px;
  font-weight: 500;
  margin: 0;
  cursor: pointer;
  color: ${({ theme }) => theme.colors.paragraph};

  & > span {
    cursor: pointer;
  }
`;

export default DatePicker;
