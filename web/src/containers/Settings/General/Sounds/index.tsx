import React, { useState } from 'react';
import styled from 'styled-components';
import SystemSoundRow from './SystemSoundRow';
import { useSelector } from 'react-redux';
import { makeTheme } from '../../../../stores/Main/reducers/theme';
import { makeSounds } from "../../../../stores/Main/reducers/settings";

const Sounds = () => {
  const { checkout, harvester } = useSelector(makeSounds);

  const theme = useSelector(makeTheme);
  const [isPlaying, setIsPlaying] = useState<string>("");

  return (
    <Container>
      <Heading>System Sounds</Heading>
      <SystemSoundRow
        {...{ isPlaying, theme, setIsPlaying }}
        soundType="Checkout"
        selector={checkout}
      />
      <Border />
      <SystemSoundRow
        {...{ isPlaying, theme, setIsPlaying }}
        soundType="Harvester"
        selector={harvester}
      />
    </Container>
  );
};

const Container = styled.div`
  padding: 16px;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  background-color: ${({ theme }) => theme.colors.sidebar};
  margin: 0 0 0 8px;
  height: 190px;
  max-height: 190px;
  flex-basis: 55%;
`;

const Heading = styled.div`
  font-size: 18px;
  margin-bottom: 8px;
  font-weight: 700;
  color: ${({ theme }) => theme.colors.subHeading};
`;

const Border = styled.div`
  display: flex;
  margin: 16px 0;
  padding: 0.5px 0;
  background-color: ${({ theme }) => theme.colors.separator};
`;

export default Sounds;
