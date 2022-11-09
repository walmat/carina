import React, {
  Dispatch,
  SetStateAction,
  ReactElement,
  useCallback,
} from "react";
import styled from "styled-components";
import { AnimatePresence, motion } from "framer-motion";
import { X } from "react-feather";

import { Modal as ModalAnimations } from "../../animations";
import { Typography } from "../../elements";

import Crumbs from "./crumbs";

import { useEscape } from "../../hooks";

interface ModalProps {
  width: number | string;
  height: number | string;
  amount?: string | number;
  crumbs?: any;
  show: boolean;
  title: string;
  setShow: Dispatch<SetStateAction<boolean>>;
  children: ReactElement;
}

const Modal = (props: ModalProps) => {
  const {
    width,
    height,
    amount,
    crumbs,
    show,
    title,
    setShow,
    children
  } = props;

  const handleClose = useCallback(() => setShow(false), []);
  useEscape(handleClose);

  return (
    <AnimatePresence exitBeforeEnter>
      {show && (
        <Backdrop
          onClick={handleClose}
          variants={ModalAnimations.backdrop}
          initial="hide"
          animate="show"
          exit="hide"
        >
          <Container
            onClick={(e) => e.stopPropagation()}
            {...props}
            width={width}
            height={height}
            variants={ModalAnimations.modal}
            initial="hide"
            animate="show"
            exit="hide"
          >
            <Row bottom={16}>
              <Title>{title}</Title>
              {typeof crumbs !== "undefined" ? <Crumbs {...crumbs} /> : null}
              {typeof amount !== "undefined" ? <Amount>{amount}</Amount> : null}
              <Close onClick={() => setShow(false)} />
            </Row>
            <Col bottom={0}>{children}</Col>
          </Container>
        </Backdrop>
      )}
    </AnimatePresence>
  );
};

const Backdrop = styled(motion.div)`
  position: fixed;
  display: flex;
  justify-content: center;
  align-items: center;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  opacity: 0.6;
  background-color: ${({ theme }) => theme.colors.backdrop};
  z-index: 900;
`;

const Container = styled(motion.div)<{
  width: number | string;
  height: number | string;
}>`
  width: ${({ width }) => (typeof width === "string" ? width : `${width}px`)};
  height: ${({ height }) =>
    typeof height === "string" ? height : `${height}px`};
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  padding: 16px;
  box-shadow: rgba(0, 0, 0, 0.15) 0px 5px 15px 0px;
  background-color: ${({ theme }) => theme.colors.sidebar};
  z-index: 1000;
`;

const Title = styled(Typography.H2)`
  margin: 0 16px 0 0;
  font-size: 22px;
  color: ${({ theme }) => theme.colors.h2};
  font-weight: 700;
  display: flex;
`;

const Amount = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 6.5px;
  font-size: 14px;
  font-weight: 700;
  min-width: 40px;
  max-height: 32px;
  border-radius: 4px;
  color: ${({ theme }) => theme.colors.primary};
  background-color: ${({ theme }) => theme.colors.secondary};
`;

const Row = styled.div<{ bottom: number }>`
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: ${({ bottom }) => `${bottom}px`};
`;

const Col = styled.div<{ bottom: number }>`
  display: flex;
  flex-direction: column;
  margin-bottom: ${({ bottom }) => `${bottom}px`};
  flex: 1;
`;

const Close = styled(X)`
  margin-left: auto;
  color: ${({ theme }) => theme.colors.paragraph};
`;

export default Modal;
