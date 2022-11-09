import React, { useState } from "react";
import styled from "styled-components";
import { Plus, Trash } from "react-feather";

import InfiniteScroll from "react-infinite-scroller";
import { Typography, Buttons, Loader } from "../../../elements";
import { useLazyLoadQuery, usePaginationFragment } from "react-relay";
import { graphql } from "react-relay";
import { withRelayBoundary } from "../../../relay/withRelayBoundary";
import { ShipmentsQuery } from "./__generated__/ShipmentsQuery.graphql";
import { ShipmentsPaginationQuery } from "./__generated__/ShipmentsPaginationQuery.graphql";
import {
  Shipments_query,
  Shipments_query$key,
} from "./__generated__/Shipments_query.graphql";
import { ExtractRelayEdgeNode } from "../../../relay/relayTypes";
import AddShipment from "./Create";
import ShipmentRow from "./ShipmentRow";

export const FullPageLoading = styled.div`
  display: flex;
  flex: 1;
  justify-content: center;
  align-items: center;
`;

// eslint-disable-next-line
type ShipmentTracker = ExtractRelayEdgeNode<
  Shipments_query["shipmentTrackers"]
>;

type Props = {
  fetchKey?: number;
};
const Shipments = (props: Props) => {
  const [create, setCreate] = useState<boolean>(false);

  const onCreate = () => {
    setCreate((create) => !create);
  };

  const response = useLazyLoadQuery<ShipmentsQuery>(
    graphql`
      query ShipmentsQuery {
        ...Shipments_query
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
    ShipmentsPaginationQuery,
    Shipments_query$key
  >(
    graphql`
      fragment Shipments_query on Query
      @argumentDefinitions(
        first: { type: Int, defaultValue: 3 }
        after: { type: String }
      )
      @refetchable(queryName: "ShipmentsPaginationQuery") {
        shipmentTrackers(first: $first, after: $after)
          @connection(key: "Shipments_shipmentTrackers", filters: []) {
          pageInfo {
            hasNextPage
            hasPreviousPage
            startCursor
            endCursor
          }
          edges {
            node {
              id
              name
              trackingID
              ...ShipmentRow_shipmentTracker
            }
          }
        }
      }
    `,
    response
  );

  // use loadMore to get more items on this connection
  const loadMore = () => {
    // Don't fetch again if we're already loading the next page
    if (isLoadingNext) {
      return;
    }
    loadNext(3);
  };

  const { shipmentTrackers } = data;
  const { pageInfo } = shipmentTrackers;

  const infiniteScrollerLoader = (
    <FullPageLoading>
      <Loader height={32} width={32} />
    </FullPageLoading>
  );

  return (
    <Container>
      <Row useBorder m="" p="0 16px" height="48px">
        <Title>Shipment Tracker</Title>
        <CreateButtonContainer>
          <Buttons.Primary
            variant="IconButton"
            width={128}
            height={28}
            onClick={onCreate}
          >
            <PlusIcon height={14} />
            <ActionText>Add Shipment</ActionText>
          </Buttons.Primary>
        </CreateButtonContainer>
      </Row>
      <AddShipment open={create} setOpen={setCreate} />
      <Col useOverflow height="calc(100% - 48px)" p="0 16px">
        {!shipmentTrackers?.edges?.length ? (
          <Center>
            <Flex fd="column" jc="center" ai="center">
              <Flex m="0 0 2px 0">
                <Title>No active shipments.</Title>
              </Flex>
              <Flex>
                {/* @ts-ignore */}
                <Date>
                  To get started,{" "}
                  <Emphasis onClick={onCreate}>
                    add a new shipment now.
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
            {shipmentTrackers.edges.map((edge) => (
              <ShipmentRow
                key={edge?.node?.id}
                shipmentTracker={edge?.node || null}
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
  overflow: hidden;
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

const CreateButtonContainer = styled.div`
  margin-left: auto;
`;

export const Row = styled.div<{
  m?: string;
  p?: string;
  height?: string;
  useBorder?: boolean;
  useHover?: boolean;
}>`
  display: flex;
  align-items: center;
  margin: ${({ m }) => m};
  padding: ${({ p }) => p};
  border-radius: 8px;
  ${({ height }) => (height ? `min-height: ${height}` : "")};
  ${({ height }) => (height ? `height: ${height}` : "")};
  ${({ theme, useBorder }) =>
    useBorder ? `border-bottom: 1px solid ${theme.colors.separator}` : ""};

  ${({ theme, useHover }) =>
    useHover
      ? `
	  &:hover {
	    background-color: ${theme.colors.secondary};
	    cursor: pointer;
	  }
	`
      : ""}
`;

export const Col = styled.div<{
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
  ${({ useOverflow }) => (useOverflow ? `overflow-y: auto;` : "")}
`;

const Title = styled(Typography.H4)`
  display: flex;
  font-size: 1.35vw;
  font-weight: 500;
  color: ${({ theme }) => theme.colors.h2};
  margin: 0 8px 0 0;
`;

const ActionText = styled(Typography.Paragraph)`
  color: #fff;
  cursor: pointer;
  font-size: 12px;
  font-weight: 500;
  margin: auto 8px;
`;

const PlusIcon = styled(Plus)`
  display: flex;
  cursor: pointer;
  justify-content: center;
  align-items: center;
`;

export const RemoveIcon = styled(Trash)`
  width: 16px;
  height: 16px;
  cursor: pointer;
  margin: 0 0 0 4px;
  fill: rgba(242, 110, 134, 0.4);
  stroke: ${({ theme }) => theme.colors.failed};
`;

export const ShipmentName = styled(Typography.H2)`
  font-size: 14px;
  text-transform: capitalize;
  color: ${({ theme }) => theme.colors.h2};
  font-weight: 500;
  margin: 0;
  line-height: 1.65em;
`;

export const TrackingNumber = styled(Typography.H2)`
  font-size: 12px;
  color: ${({ theme }) => theme.colors.paragraph};
  font-weight: 400;
  margin: 0;
  line-height: 1em;
`;

export const ShipmentStatus = styled.div<{ delivered: boolean }>`
  font-size: 14px;
  border-radius: 2px;
  color: ${({ theme, delivered }) =>
    !delivered ? theme.colors.paragraph : theme.colors.primary};
  background-color: ${({ theme, delivered }) =>
    !delivered ? theme.colors.separator : theme.colors.secondary};
  padding: 4px 8px;
  font-weight: 400;
  margin: 0;
`;

export default withRelayBoundary(Shipments);
