import React, { useState } from "react";
import { useDispatch } from "react-redux";
import { motion, AnimatePresence } from "framer-motion";
import { Check, X as Close } from "react-feather";
import styled from "styled-components";

import { Buttons } from "../../../../elements";

import { logout } from '../../../../stores/Main/reducers/user';
import { close } from '../../../../utils';

type SignOffProps = {
  loading: boolean;
  onClick: any;
  children: any;
};

function SignOffButton({ loading, children, ...props }: SignOffProps) {
  return (
    <Buttons.Secondary width={95} height={32} text={children} {...props}>
      {children}
    </Buttons.Secondary>
  );
}

type ConfirmsProps = {
  confirmProps: object;
  cancelProps: object;
  indeterminate: boolean;
};

function Confirms({
  confirmProps,
  cancelProps,
  indeterminate,
  ...props
}: ConfirmsProps) {
  return (
    <ConfirmContainer
      key="button-confirms"
      initial={{ opacity: 0 }}
      animate={{
        opacity: indeterminate ? 1 : 0,
      }}
      transition={{
        opacity: { duration: 0.2, type: "spring" },
      }}
      exit={{ opacity: 0 }}
      {...props}
    >
      <Confirm {...confirmProps} />
      <Cancel {...cancelProps} />
    </ConfirmContainer>
  );
}

const SignOff = () => {
  const dispatch = useDispatch();

  const [message, setMessage] = useState("Sign Out");
  const [indeterminate, setIndeterminate] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleIndeterminate = (e: any) => {
    e.stopPropagation();

    setIndeterminate((prev) => !prev);
  };

  const handleLogout = async (e: any) => {
    e.stopPropagation();

    if (loading) {
      return;
    }

    setLoading(true);
    setIndeterminate(false);
    setMessage("");

    try {
      // await deactivate();
    } catch (e) {
      setLoading(false);
      return setTimeout(() => setMessage("Sign Out"), 1000);
    }

    setMessage("Sign Out");
    setLoading(false);
    dispatch(logout());
    return close();
  };

  return (
    <Container>
      <SignOffButton
        onClick={!indeterminate ? handleIndeterminate : () => {}}
        loading={loading}
      >
        <AnimatePresence exitBeforeEnter>
          {indeterminate && (
            <Confirms
              indeterminate={indeterminate}
              cancelProps={{
                onClick: handleIndeterminate,
              }}
              confirmProps={{
                onClick: handleLogout,
              }}
            />
          )}
          {!indeterminate && (
            <Row
              animate={{
                opacity: !indeterminate ? 1 : 0,
              }}
              transition={{
                opacity: { duration: 0.2, type: "spring" },
              }}
              exit={{ opacity: 0 }}
              key="button-message"
            >
              {message}
            </Row>
          )}
        </AnimatePresence>
      </SignOffButton>
    </Container>
  );
};

const Container = styled.div`
  display: flex;
  flex-direction: column;
  margin-left: auto;
  margin-top: 10px;
`;

const Row = styled(motion.div)`
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const ConfirmContainer = styled(motion.div)`
  display: flex;
  align-items: center;
  justify-content: space-around;
`;

const Cancel = styled(Close)`
  height: 16px;
  width: 16px;
  cursor: pointer;
  opacity: 0.6;
  transition-duration: 240ms;
  transition-property: scale, opacity;

  &:hover {
    opacity: 1;
    transform: scale(1.25);
  }
`;

const Confirm = styled(Check)`
  height: 16px;
  width: 16px;
  cursor: pointer;
  opacity: 0.6;
  transition-duration: 240ms;
  transition-property: scale, opacity;

  &:hover {
    opacity: 1;
    transform: scale(1.25);
  }
`;

export default SignOff;