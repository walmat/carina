import React, { Fragment, useState } from "react";
import styled from "styled-components";
import moment from "moment";
import { useTranslation } from "react-i18next";

import { Content, Toolbar } from "../../components";
import { Typography } from "../../elements";

import Month from "./month";
import Releases from "./releases";
import ReleaseComponent from "./release";
import { graphql, useLazyLoadQuery, useRefetchableFragment } from "react-relay";
import { OrderHistoryQuery } from "../Dashboard/OrderHistory/__generated__/OrderHistoryQuery.graphql";
import { withRelayBoundary } from "../../relay/withRelayBoundary";

const now = moment();

type Props = {
  fetchKey?: number;
};
const Calendar = (props: Props) => {
  const { t } = useTranslation();
  const [date, setDate] = useState(moment());
  const [releases, setReleases] = useState<Releases>({});
  const [release, setRelease] = useState<string>("");

  const response = useLazyLoadQuery<OrderHistoryQuery>(
    graphql`
      query CalendarQuery($filters: ReleaseFilter) {
        ...Calendar_query @arguments(filters: $filters)
      }
    `,
    {
      // TODO - filter by date range
      filters: {},
    },
    {
      fetchPolicy: "store-or-network",
      fetchKey: props.fetchKey,
    }
  );

  const [data, refetch] = useRefetchableFragment(
    graphql`
      fragment Calendar_query on Query
      @argumentDefinitions(
        first: { type: Int, defaultValue: 10 }
        after: { type: String }
        filters: { type: ReleaseFilter }
      )
      @refetchable(queryName: "CalendarRefetchQuery") {
        releases(first: $first, after: $after, filters: $filters) {
          edges {
            node {
              name
              image
              stores(first: 10) {
                edges {
                  node {
                    name
                    date
                  }
                }
              }
            }
          }
        }
      }
    `,
    response
  );

  // TODO: fetch this months timestamps from GraphQL API to display on the Calendar
  // When the select a certain date that has data, we perform the request to
  // get all releases for that date and their info (minus inner `stores` data) as well.
  const mockReleases = {
    "1620799253": {
      oonfSxoZrLm1NDbspsRKc: {
        name: "Air Jordan 1 Shadow 2.0",
        image: "https://i.ibb.co/Fsbyxnf/shadow-2-0.jpg",
        stores: {
          "1": {
            name: "Kith",
            date: "1621972000",
            datum: [],
          },
          "2": {
            name: "A-Ma-Maniere",
            date: "1620979000",
            datum: [],
          },
        },
      },
      oonfSxoZrLm1NDbsp3RKc: {
        name: "Dunk Low Grey",
        image:
          "https://images.stockx.com/images/Nike-Dunk-Low-College-Navy-Grey-W-Product.jpg",
        stores: {},
      },
    },
    "1621058453": {
      oonfSxoZrLm1NDbspsRKc: {
        name: "Dunk Low Grey",
        image:
          "https://images.stockx.com/images/Nike-Dunk-Low-College-Navy-Grey-W-Product.jpg",
        stores: {},
      },
    },
    "1621211210": {
      oonfSxoZrLm1NDbspsRKc: {
        name: "Yeezy QNTM Orange",
        image:
          "https://images.stockx.com/images/adidas-Yeezy-QNTM-Flash-Orange.jpg",
        stores: {},
      },
    },
  } as ReleasesRoot;

  return (
    <Content>
      <Fragment key="Calendar">
        <Toolbar simple title={t(`Pages.Calendar`)} />
        <DayText>Today is {now.format("MMMM Do")}</DayText>
        <Container>
          <Grid>
            <Month
              {...{
                date,
                setRelease,
                setReleases,
                setDate,
                releases: mockReleases,
              }}
            />
            <Releases {...{ date, release, setRelease, releases }} />
          </Grid>
          <ReleaseComponent
            {...{ date, release: release ? releases[release] : null }}
          />
        </Container>
      </Fragment>
    </Content>
  );
};

const Container = styled.div`
  display: flex;
  flex: 1;
  min-height: 0;
`;

const DayText = styled(Typography.Paragraph)`
  margin: 0;
  font-size: 12px;
  font-weight: 400;
  color: ${({ theme }) => theme.colors.paragraph};
`;

const Grid = styled.div`
  gap: 16px;
  margin: 16px 0 0 0;
  display: grid;
  flex: 1;
  margin-right: 16px;
  grid-template-columns: repeat(1, minmax(0, 1fr));
  grid-template-rows: repeat(2, minmax(0, 1fr));
`;

export default withRelayBoundary(Calendar);
