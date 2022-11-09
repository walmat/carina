import React, { Dispatch, SetStateAction, useEffect, useState } from "react";
import moment, { Moment } from "moment";
import styled from "styled-components";
import { motion } from "framer-motion";
import { ChevronRight, ChevronLeft } from "react-feather";
import { customPickerOptions, DateOption } from './index';

export const getWeekday = (date: Moment) => {
  const weekday = (date.startOf("month").day() + 1) % 7;

  if (weekday === 0) {
    return 7;
  }

  return weekday;
};

const MONTHS = {
  ADD: (month: Moment, amount = 1) => moment(month).add(amount, "month"),
  SUB: (month: Moment, amount = 1) => moment(month).subtract(amount, "month"),
};

interface CalendarProps {
  focus: "start" | "end";
  start: Moment;
  end: Moment | null;
  setTimeframe: Dispatch<SetStateAction<DateOption>>;
  setFocus: Dispatch<SetStateAction<"start" | "end">>;
  timeframe: DateOption;
}

const Calendar = ({ focus, start, end, timeframe, setTimeframe, setFocus }: CalendarProps) => {
  const [startDate, setStartDate] = useState<Moment>(
    start
  );
  const [endDate, setEndDate] = useState<Moment | null>(
    end
  );
  const [currentDate, setCurrentDate] = useState(startDate);

  const [currentMonth, setCurrentMonth] = useState(currentDate.month() + 1);

  const weekdays = moment.weekdaysMin();
  const firstWeekday = getWeekday(moment(currentDate));
  const daysInMonth = moment(currentDate).daysInMonth();
  const daysInLastMonth = moment(currentDate)
    .subtract(1, "month")
    .daysInMonth();

  const numberOfDaysLastMonth = firstWeekday-1;

  const getDaysLastMonth = () => {
    // generate days for last month
    const daysLastMonth = Array.from(Array(numberOfDaysLastMonth).keys()).map(i => daysInLastMonth - (numberOfDaysLastMonth -i) + 1);

    return daysLastMonth.map(d => moment(currentDate)
        .set("month", currentDate.month() - 1)
        .set("date", d)
    );
  };

  const getDaysNextMonth = () => {
    const numberOfDaysNextMonth = 42 - (numberOfDaysLastMonth + daysInMonth);

    return Array.from(Array(numberOfDaysNextMonth).keys())
      .map(i => i + 1)
      .map(d => moment(currentDate)
        .set("month", currentDate.month() + 1)
        .set("date", d)
      );
  };

  const getDaysMonth = () => {
    return Array.from(Array(daysInMonth).keys())
      .map(i => i + 1)
      .map(d => moment(currentDate)
        .set("month", currentDate.month())
        .set("date", d)
      );
  };

  const dates = [
    ...getDaysLastMonth(),
    ...getDaysMonth(),
    ...getDaysNextMonth(),
  ];

  const isDaySelected = (day: Moment) => {
    if (!endDate) {
      return day.isSame(startDate, 'days');
    }

    const isSelected = day.isBetween(startDate, endDate, "days", "[]");

    return isSelected;
  }

  const onDayClick = (d: Moment) => {
    const isMonth = d.month() === currentDate.month();

    if (!isMonth) {
      return;
    }

    if (focus === 'start') {
      setTimeframe({
        ...timeframe,
        ...customPickerOptions,
        value: {
          ...timeframe.value,
          start: d,
        },
      });
      setFocus("end");
      return
    }
    //
    

    if (d.isSame(startDate, 'days')) {
      // @ts-ignore
      // setStartDate(null);
    } else if (d.isBefore(startDate)) {
      setStartDate(d);
      setEndDate(null);
    } else {
      setEndDate(d);

      setTimeframe({
        ...timeframe,
        ...customPickerOptions,
        value: {
          ...timeframe.value,
          start: startDate,
          end: d,
        },
      });
    }
  }

  const handleChangeMonth = (dir: "add" | "sub") => {
    switch (dir) {
      default:
      case "add":
        return setCurrentDate(MONTHS.ADD(currentDate));
      case "sub":
        return setCurrentDate(MONTHS.SUB(currentDate));
    }
  };

  useEffect(() => {
    const s = moment(start);
    const e = moment(end);

    setStartDate(s);
    setCurrentDate(s);

    if (focus !== 'end') {
      setEndDate(e);
    }

    setCurrentMonth(s.month() + 1);
  }, [focus, start, end]);

  return (
    <Container>
      <Flex jc="center">
        <Flex flex={1} fd="column" jc="center">
          <Flex m="8px" jc="center">
            <Arrow
              onClick={() => handleChangeMonth("sub")}
              component={<ChevronLeft />}
            />
            <MonthHeader>{currentDate.format("MMMM YYYY")}</MonthHeader>
            <Arrow
              onClick={() => handleChangeMonth("add")}
              component={<ChevronRight />}
            />
          </Flex>
        </Flex>
      </Flex>
      <Weekday>
        {weekdays.map((wk) => (
          <WeekdayText key={wk}>{wk}</WeekdayText>
        ))}
      </Weekday>
      <Days>
        {dates.map((d, i) => {
          const isMonth = d.month() === currentDate.month();
          const isSelected = isDaySelected(d);
          const day = d.date();

          return (
            <Day
              whileHover={{ scale: !isMonth ? 1 : 1.05 }}
              whileTap={{ scale: 1 }}
              key={i}
              month={isMonth}
              selected={isSelected}
              onClick={() => onDayClick(d)}
            >
              <Flex jc="flex-end">
                <DayText month={isMonth} selected={isSelected}>
                  {day}
                </DayText>
              </Flex>
            </Day>
          );
        })}
      </Days>
    </Container>
  );
};

const Container = styled.div`
  background-color: ${({ theme }) => theme.colors.sidebar};
  display: flex;
  flex-direction: column;
  border-radius: 8px;
  flex: 1;
`;

const Flex = styled.div<{
  flex?: number;
  fd?: string;
  jc?: string;
  ai?: string;
  m?: string;
}>`
  display: flex;
  ${({ flex }) => (typeof flex !== undefined ? `flex: ${flex};` : "")}
  ${({ fd }) => (fd ? `flex-direction: ${fd};` : "")}
	${({ jc }) => (jc ? `justify-content: ${jc};` : "")}
	${({ ai }) => (ai ? `align-items: ${ai};` : "")}
	${({ m }) => (m ? `margin: ${m};` : "")}
`;

const MonthHeader = styled.div`
  font-size: 12px;
  font-weight: 700;
  text-align: center;
  display: flex;
  margin: 0 auto;
  color: ${({ theme }) => theme.colors.subHeading};
  justify-content: center;
`;

const Arrow = styled(({ component, ...props }) =>
  React.cloneElement(component, props)
)`
  width: 12px;
  height: 12px;
  cursor: pointer;
  opacity: 0.6;
  color: ${({ theme }) => theme.colors.subHeading};
  margin: auto 0;
`;

const Weekday = styled.div`
  font-weight: 500;
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  grid-gap: 0;
`;

const WeekdayText = styled.div`
  color: ${({ theme }) => theme.colors.h3};
  background-color: transparent;
  padding: 2px 4px;
  text-align: center;
  font-size: 10px;
  font-weight: 700;
  opacity: 1;
`;

const Days = styled.div`
  display: grid;
  height: 100%;
  grid-template-columns: repeat(7, 1fr);
  grid-gap: 0;
  margin-top: 2px;
`;

const Day = styled(motion.div)<{ month: boolean; selected: boolean }>`
  display: flex;
  flex-direction: column;
  border-radius: 0;
  cursor: ${({ month }) => (!month ? "default" : "pointer")};
  border: 0.5px solid ${({ theme }) => theme.colors.border};
  border-bottom: none;
  border-left: none;
  background-color: ${({ theme, selected }) =>
    selected ? theme.colors.secondary : "transparent"};
  opacity: ${({ month }) => (!month ? "0.6" : "1")};

  &:nth-child(7n) {
    border-right: none;
  }

  &:hover {
    opacity: 0.6;
  }
`;

const DayText = styled.div<{ month: boolean; selected: boolean }>`
  color: ${({ theme, month, selected }) =>
    month && selected ? theme.colors.primary : theme.colors.h2};
  opacity: ${({ month }) => (!month ? 0.6 : 1)};
  padding: 2px 4px;
  cursor: ${({ month }) => (!month ? "default" : "pointer")};
  text-align: center;
  border-radius: 4px;
  font-size: 10px;
  font-weight: 700;
`;

export default Calendar;
