import React, { useState } from "react";
import { useDispatch } from "react-redux";
import styled from "styled-components";
import { graphql, useLazyLoadQuery, usePaginationFragment } from "react-relay";
import { Settings } from "react-feather";
import { useHistory } from "react-router-dom";
import InfiniteScroll from "react-infinite-scroller";

import { Loader, Typography } from "../../../elements";
import { DatePicker } from "../../../components";
import { withRelayBoundary } from "../../../relay/withRelayBoundary";
import { OrderHistoryQuery } from "./__generated__/OrderHistoryQuery.graphql";
import { OrderHistoryPaginationQuery } from "./__generated__/OrderHistoryPaginationQuery.graphql";
import { OrderHistory_query$key } from "./__generated__/OrderHistory_query.graphql";
import { FullPageLoading } from "../Shipments/Shipments";
import OrderHistoryRow from "./OrderHistoryRow";
import { toggle, SETTINGS } from "../../../stores/Main/reducers/settings";

type Props = {
  fetchKey?: number;
};

const OrderHistory = (props: Props) => {
  const history = useHistory();
  const dispatch = useDispatch();

  // @ts-ignore
  const [timeframe, setTimeframe] = useState<any>(
    // @ts-ignore
    DatePicker.DATE_PICKER_OPTIONS[0]
  );

  const response = useLazyLoadQuery<OrderHistoryQuery>(
    graphql`
      query OrderHistoryQuery {
        ...OrderHistory_query
      }
    `,
    {},
    {
      fetchPolicy: "store-or-network",
      fetchKey: props.fetchKey,
    }
  );

  // render data.edges
  const { data, loadNext, isLoadingNext } = usePaginationFragment<
    OrderHistoryPaginationQuery,
    OrderHistory_query$key
  >(
    graphql`
      fragment OrderHistory_query on Query
      @argumentDefinitions(
        first: { type: Int, defaultValue: 3 }
        after: { type: String }
      )
      @refetchable(queryName: "OrderHistoryPaginationQuery") {
        checkouts(first: $first, after: $after)
          @connection(key: "OrderHistory_checkouts", filters: []) {
          pageInfo {
            hasNextPage
            hasPreviousPage
            startCursor
            endCursor
          }
          edges {
            node {
              id
              date
              amount
              ...OrderHistoryRow_checkout
            }
          }
        }
      }
    `,
    response
  );

  // NOTE: Navigate to tasks page, open dialog
  const handleAddTasks = () => {
    dispatch(toggle({ field: SETTINGS.IS_OPEN }));
    history.push("/tasks");
  };

  // use loadMore to get more items on this connection
  const loadMore = () => {
    // Don't fetch again if we're already loading the next page
    if (isLoadingNext) {
      return;
    }
    loadNext(3);
  };

  const { checkouts } = data;
  const { pageInfo } = checkouts;

  const infiniteScrollerLoader = (
    <FullPageLoading>
      <Loader height={32} width={32} />
    </FullPageLoading>
  );

  return (
    <Container>
      <Row useBorder m="" p="0 16px" height="48px">
        <Title>Order History</Title>

        <Selects>
          <DatePicker timeframe={timeframe} setTimeframe={setTimeframe} />
          <Expand />
        </Selects>
      </Row>
      <Col useOverflow height="calc(100% - 48px)" m="0 0 0 8px">
        {!checkouts?.edges?.length ? (
          <Center>
            <Flex fd="column" jc="center" ai="center">
              <Flex m="0 0 2px 0">
                <Title>No order history found.</Title>
              </Flex>
              <Flex>
                <Date>
                  Please choose another time period or{" "}
                  <Emphasis onClick={handleAddTasks}>
                    add tasks now.
                  </Emphasis>
                </Date>
              </Flex>
            </Flex>
          </Center>
        ) : (
          <InfiniteScroll
            pageStart={0}
            loadMore={loadMore}
            hasMore={pageInfo.hasNextPage}
            loader={infiniteScrollerLoader}
            useWindow
          >
            {checkouts.edges.map((edge) => (
              <OrderHistoryRow
                key={edge?.node?.id}
                checkout={edge?.node || null}
              />
            ))}
          </InfiniteScroll>
        )}
      </Col>
    </Container>
  );
};

const Container = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  border-radius: ${({ theme }) => theme.space.M}px;
  padding: ${({ theme }) => theme.space.NONE}px;
  background-color: ${({ theme }) => theme.colors.sidebar};
`;

const Center = styled.div`
  background-color: ${({ theme }) => theme.colors.sidebar};
  border-radius: 16px;
  padding: 16px;
  border: none;
  display: flex;
  margin-top: 16px;
  flex-direction: column;
  flex: 1;
  align-items: center;
  justify-content: center;
`;

const Flex = styled.div<{
  fd?: string;
  flex?: number | string;
  jc?: string;
  ai?: string;
  m?: string;
  of?: string;
}>`
  display: flex;
  min-width: 0;
  ${({ of }) => (typeof of !== "undefined" ? `overflow: ${of};` : "")}
  ${({ flex }) => (typeof flex !== "undefined" ? `flex: ${flex};` : "")}
	${({ fd }) => (fd ? `flex-direction: ${fd};` : "")}
	${({ jc }) => (jc ? `justify-content: ${jc};` : "")}
	${({ ai }) => (ai ? `align-items: ${ai};` : "")}
	${({ m }) => (m ? `margin: ${m};` : "")}
`;

const Date = styled(Typography.Paragraph)`
  display: flex;
  font-size: 12px;
  font-weight: 400;
  opacity: 0.8;
  white-space: pre;
  color: ${({ theme }) => theme.colors.h1};
  margin: 0 16px;
`;

const Emphasis = styled.span`
  font-weight: 500;
  cursor: pointer;
  color: ${({ theme }) => theme.colors.paragraph};
`;

const Row = styled.div<{
  m?: string;
  p?: string;
  height?: string;
  useBorder?: boolean;
}>`
  display: flex;
  align-items: center;
  margin: ${({ m }) => m};
  padding: ${({ p }) => p};
  ${({ height }) =>
    height
      ? `
    height: ${height};
    min-height: ${height};
  `
      : ""};
  ${({ theme, useBorder }) =>
    useBorder ? `border-bottom: 1px solid ${theme.colors.separator}` : ""};
`;

const Title = styled(Typography.H4)`
  display: flex;
  font-size: 1.35vw;
  font-weight: 500;
  color: ${({ theme }) => theme.colors.h2};
  margin: 0 8px 0 0;
`;

const Selects = styled.div`
  margin-left: auto;
  display: flex;
`;

const Expand = styled(Settings)`
  margin: auto 0 auto 16px;
  width: 16px;
  color: ${({ theme }) => theme.colors.paragraph};
`;

const Col = styled.div<{
  m?: string;
  p?: string;
  height?: string;
  useOverflow?: boolean;
}>`
  display: flex;
  flex-direction: column;
  margin: ${({ m }) => m};
  padding: ${({ p }) => p};
  position: relative;
  ${({ height }) => (height ? `height: ${height};` : "")}
  ${({ useOverflow }) => (useOverflow ? `overflow-y: scroll;` : "")}
`;

export default withRelayBoundary(OrderHistory);
