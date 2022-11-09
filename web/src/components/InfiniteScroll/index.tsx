import React, { useEffect, useState } from "react";
import InfiniteScroll from "react-infinite-scroll-component";
import { motion, useAnimation } from "framer-motion";
import { useInView } from "react-intersection-observer";
import styled from "styled-components";

import { Loader } from "../../elements";
import { InfiniteScroll as InfiniteScroller } from "../../animations";

const { container, item } = InfiniteScroller;

interface ScrollProps {
  fetch: (page: number) => Promise<any>;
  children: React.ReactElement;
}

export default function InfiniteScrollComponent({
  fetch,
  children,
}: ScrollProps) {
  const [page, setPage] = useState(1);
  const [data, setData] = useState([]);

  const doFetch = (pageNumber: number = 1) => {
    return fetch(pageNumber).then((res) => {
      setData((d) => d.concat(res.data));
      setPage((p) => p + 1);
    });
  };

  const fetchMoreData = (page: number) => {
    console.log("fetching more data");
    return doFetch(page);
  };

  useEffect(() => {
    doFetch(page);
  }, []);

  if (!data.length) {
    return (
      <FullPageLoading>
        <Loader height={32} width={32} />
      </FullPageLoading>
    );
  }

  return (
    <InfiniteScroll
      dataLength={data.length}
      next={() => fetchMoreData(page)}
      hasMore={true}
      loader={
        <Loading>
          <Loader height={16} width={16} />
        </Loading>
      }
    >
      <Items initial="hidden" animate="show" variants={container}>
        {data.map((datum: any) => {
          return <Item key={datum.id} children={children} {...datum} />;
        })}
      </Items>
    </InfiniteScroll>
  );
}

const Item = ({ children, ...data }: any): React.ReactElement => {
  const controls = useAnimation();
  const [ref, inView] = useInView();

  useEffect(() => {
    if (inView) {
      controls.start("show");
    }
  }, [controls, inView]);

  return (
    <ListItem
      key={data.id}
      variants={item}
      initial="hidden"
      animate={controls}
      ref={ref}
    >
      {children(data)}
    </ListItem>
  );
};

const Items = styled(motion.ul)`
  margin: 0 8px 0 20px;
  height: 100%;
  list-style: none;
  padding: 0;
`;

const ListItem = styled(motion.li)`
  margin: 16px 0;

  &:nth-child(1) {
    margin-top: 0;
  }
`;

const FullPageLoading = styled.div`
  display: flex;
  flex: 1;
  justify-content: center;
  align-items: center;
`;

const Loading = styled.div`
  display: flex;
  padding: 12px 8px;
  justify-content: center;
`;
