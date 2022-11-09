import React, { useEffect, useState } from "react";
import styled from "styled-components";
import moment from "moment";

interface PillProps {
  fetch: any; // todo: type
}

const Pill = ({ fetch }: PillProps) => {
  const [lastFetch, setLastFetch] = useState(moment());
  const [diff, setDiff] = useState("just now");

  useEffect(() => {
    let interval: any = null;

    const poll = () => {
      interval = setInterval(() => {
        const diff = moment().diff(lastFetch, "minutes");

        if (diff >= 30) {
          // TODO: Update stale data every ~30m
          return setLastFetch(moment());
        }

        if (diff === 0) {
          return setDiff("just now");
        }

        if (diff === 1) {
          return setDiff(`${diff} minute ago`);
        }

        return setDiff(`${diff} minutes ago`);
      }, 1000);
    };

    poll();

    return () => {
      clearInterval(interval);
      interval = null;
    };
  });

  return <Accent>Last updated {diff}</Accent>;
};

const Accent = styled.div`
  font-size: 12px;
  font-weight: 400;
  padding: 4px 8px;
  border-radius: 4px;
  color: ${({ theme }) => theme.colors.primary};
  background-color: ${({ theme }) => theme.colors.secondary};
`;

export default Pill;
