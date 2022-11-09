import React from "react";
import { useSelector } from "react-redux";
import styled from "styled-components";
import { useFragment, graphql } from "react-relay";
import moment from "moment";

import { Row, Col } from "../Shipments/Shipments";
import { OrderHistoryRow_checkout$key } from "./__generated__/OrderHistoryRow_checkout.graphql";
import { Typography } from "../../../elements";
import { CurrencyConverter } from '../../../components';
import { makeSettings } from '../../../stores/Main/reducers/settings';

type Props = {
  checkout: OrderHistoryRow_checkout$key | null;
};
const OrderHistoryRow = (props: Props) => {
  const { currency } = useSelector(makeSettings);

  const checkout = useFragment<OrderHistoryRow_checkout$key>(
    graphql`
      fragment OrderHistoryRow_checkout on Checkout {
        id
        date
        amount
        product {
          image
        }
      }
    `,
    props.checkout
  );

  if (!checkout) {
    return null;
  }

  return (
    <Row useHover p="8px 8px" m="8px 0" useBorder={false}>
      <Col>
        <Row useBorder={false}>
          {/* TODO: Change this to `checkout.product.image` after updating schema */}
          {/* @ts-ignore */}
          <Image
            src={
              checkout?.product?.image ||
              "https://nebulabots.s3.amazonaws.com/default.jpeg"
            }
          />
        </Row>
      </Col>
      <Col m="0 auto 0 0">
        <Row useBorder={false}>
          {/* TODO: Change this to `checkout.product.name` after updating schema */}
          <ProductName>{checkout.date?.slice(0, 15)}</ProductName>
        </Row>
        <Row useBorder={false}>
          <ProductPrice>
            {/* TODO: Cache the amount somehow in case the api has rate limiting or w/e */}
            <CurrencyConverter
              from="USD"
              to={currency.value}
              value={Number(checkout.amount)}
            />
          </ProductPrice>
        </Row>
      </Col>
      <Col m="4px 0 auto 0">
        <Row useBorder={false}>
          <OrderDate>{moment(checkout.date).format("MMMM D, YYYY")}</OrderDate>
        </Row>
      </Col>
    </Row>
  );
};

const Image = styled.img`
  display: flex;
  height: 40px;
  width: 40px;
  border-radius: 4px;
  margin: 0 16px 0 0;
`;

export const ProductName = styled(Typography.H2)`
  font-size: 14px;
  text-transform: capitalize;
  cursor: pointer;
  color: ${({ theme }) => theme.colors.h2};
  font-weight: 500;
  margin: 0;
  line-height: 1.65em;
`;

export const ProductPrice = styled(Typography.H2)`
  font-size: 12px;
  color: ${({ theme }) => theme.colors.paragraph};
  font-weight: 400;
  cursor: pointer;
  margin: 0 0 4px 0;
  line-height: 1.35em;
`;

export const OrderDate = styled(Typography.H2)`
  font-size: 12px;
  color: ${({ theme }) => theme.colors.paragraph};
  font-weight: 400;
  cursor: pointer;
  margin: 0;
  line-height: 1em;
`;

export default OrderHistoryRow;
