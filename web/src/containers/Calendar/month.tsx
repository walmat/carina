import React, { Dispatch, SetStateAction, useState } from "react";
import moment, { Moment } from "moment";
import styled from "styled-components";
import { motion } from "framer-motion";
import { ChevronRight, ChevronLeft } from "react-feather";

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

interface MonthProps {
  date: Moment;
  setDate: Dispatch<SetStateAction<Moment>>;
  setRelease: Dispatch<SetStateAction<string>>;
  setReleases: Dispatch<SetStateAction<Releases>>;
  releases: ReleasesRoot;
}

const Month = ({
  date,
  setDate,
  setRelease,
  setReleases,
  releases,
}: MonthProps) => {
  const [month, setMonth] = useState(moment(date));
  const weekdays = moment.weekdaysMin();
  const firstWeekday = getWeekday(moment(month));
  const daysInMonth = moment(month).daysInMonth();
  const startOfMonth = moment(month).startOf("month");
  const endOfMonth = moment(month).endOf("month");

  let days = Array.from(Array(daysInMonth).keys()).map((i) => {
    return {
      date: moment(startOfMonth).add(i, "days"),
      month: 0,
    };
  });

  for (let i = 1; i < firstWeekday; i += 1) {
    days.unshift({ date: moment(startOfMonth).subtract(i, "days"), month: -1 });
  }

  const added = 42 - days.length;
  for (let i = 1; i < added + 1; i += 1) {
    days = [...days, { date: moment(endOfMonth).add(i, "days"), month: 1 }];
  }

  const handleChangeMonth = (dir: "add" | "sub") => {
    switch (dir) {
      default:
      case "add":
        return setMonth(MONTHS.ADD(month));
      case "sub":
        return setMonth(MONTHS.SUB(month));
    }
  };

  const handleChangeDate = (month: number, _date: Moment) => {
    if (month !== 0) {
      handleChangeMonth(month > 0 ? "add" : "sub");
    }

    if (moment(date).isSame(_date, "day")) {
      return;
    }

    const selected = _date.format("L");
    const found = Object.keys(releases).find((d) => {
      const entry = moment.unix(Number(d)).format("L");
      if (entry === selected) {
        return d;
      }
    });

    if (found) {
      setReleases(releases[found]);
    } else {
      setReleases({});
    }

    setRelease("");
    setDate(_date);
  };

  return (
    <Container>
      <Flex jc="flex-start">
        <Flex fd="column" jc="center" ai="center" mt={16}>
          <Flex jc="flex-start" ai="center">
            <MonthHeader>{month.format("MMM YYYY")}</MonthHeader>
            <Arrow
              onClick={() => handleChangeMonth("sub")}
              component={<ChevronLeft />}
            />
            <Arrow
              onClick={() => handleChangeMonth("add")}
              component={<ChevronRight />}
            />
          </Flex>
        </Flex>
      </Flex>
      <Weekday>
        {weekdays.map((wk) => (
          <WeekdayText key={wk}>{wk.slice(0, 1)}</WeekdayText>
        ))}
      </Weekday>
      <Days>
        {days.map(({ date: _date, month }) => {
          const selected = moment(_date).isSame(date, "day");
          const hasData = Object.keys(releases).some((d) => {
            const first = _date.startOf("day").unix();
            const second = moment.unix(Number(d)).startOf("day").unix();
            return first === second;
          });

          return (
            <Day
              hasData={hasData}
              selected={selected}
              onClick={() => handleChangeDate(month, _date)}
              key={`${moment(_date).date()}-${month}`}
              month={month}
            >
              <Flex jc="flex-end">
                <DayText hasData={hasData} selected={selected} month={month}>
                  {moment(_date).format("D")}
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
  border-radius: 16px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  flex: 1;
`;

const Flex = styled.div<{
  fd?: string;
  jc?: string;
  ai?: string;
  mt?: number;
}>`
  display: flex;
  ${({ fd }) => (fd ? `flex-direction: ${fd};` : "")}
  ${({ jc }) => (jc ? `justify-content: ${jc};` : "")}
	${({ ai }) => (ai ? `align-items: ${ai};` : "")}
	${({ mt }) => (mt ? `margin-top: ${mt};` : "")}
`;

const MonthHeader = styled.div`
  font-size: 18px;
  font-weight: 700;
  text-align: center;
  display: flex;
  color: ${({ theme }) => theme.colors.subHeading};
  min-width: 90px;
  justify-content: center;
`;

const Arrow = styled(({ component, ...props }) =>
  React.cloneElement(component, props)
)`
  width: 16px;
  height: 16px;
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
  margin-top: 16px;
`;

const WeekdayText = styled.div`
  color: ${({ theme }) => theme.colors.h3};
  background-color: transparent;
  text-align: center;
  font-size: 10px;
  font-weight: 700;
  opacity: 1;
`;

const Days = styled.div`
  display: grid;
  height: 100%;
  grid-template-columns: repeat(7, 1fr);
  grid-gap: 4px;
  margin-top: 8px;
`;

const Day = styled(motion.div)<{
  month: number;
  hasData: boolean;
  selected: boolean;
}>`
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  border-radius: 4px;
  box-sizing: border-box;
  opacity: ${({ month }) => (!month ? 1 : 0.4)};
  background-color: ${({ theme, hasData, selected }) =>
    selected
      ? theme.colors.dayActive
      : hasData
      ? theme.colors.primary
      : theme.colors.day};
  cursor: ${({ month }) => (month ? "default" : "pointer")};
  border: 1px solid
    ${({ theme, selected }) =>
      selected ? theme.colors.primary : theme.colors.day};
`;

const DayText = styled.div<{
  month: number;
  hasData: boolean;
  selected: boolean;
}>`
  color: ${({ theme, selected, hasData }) =>
    selected ? theme.colors.dayText : hasData ? "#fff" : theme.colors.h2};
  cursor: ${({ month }) => (month ? "default" : "pointer")};
  text-align: center;
  border-radius: 4px;
  font-size: 10px;
  font-weight: 700;
`;

export default Month;
