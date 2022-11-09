import React, { Fragment, useEffect } from 'react';
import { useSelector } from 'react-redux';
import styled from "styled-components";
import { useTranslation } from 'react-i18next';
import { useIntercom } from 'react-use-intercom';

import { Content, Toolbar } from '../../../components';
import { makeUser } from '../../../stores/Main/reducers/user';

import Information from "./Information";
import Preferences from "./Preferences";
import Behaviors from "./Behaviors";
import Theme from "./Theme";
import Storage from "./Storage";
import Sounds from "./Sounds";
import Footer from "./Footer";

const General = () => {
  const { t } = useTranslation();
  const { boot, shutdown } = useIntercom();

  const { hash, type, email } = useSelector(makeUser);

  useEffect(() => {
    boot({
      name: email,
      email,
      horizontalPadding: 32,
      verticalPadding: 16,
      userHash: hash,
      customAttributes: {
        type,
        email,
      }
    });

    return () => {
      shutdown();
    }
  }, []);

  return (
    <Content>
      <Fragment key="General">
        <Toolbar simple title={t("Pages.General")} />
        <Container>
          <Information />
          <Preferences />
        </Container>
        <Container>
          <Theme />
          <Storage />
        </Container>
        <Container>
          <Behaviors />
          <Sounds />
        </Container>
        <Footer />
      </Fragment>
    </Content>
  );
};

const Container = styled.div`
  display: flex;
  margin-top: 16px;
`;

export default General;
