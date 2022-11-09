import { useState } from "react";
import { useSelector } from "react-redux";
import { motion } from "framer-motion";
import styled from "styled-components";

import Circle from './circle';
import { Buttons } from "../../../../elements";

import { makeUser } from '../../../../stores/Main/reducers/user';

const License = () => {
  const { email } = useSelector(makeUser);

  const [show, setShow] = useState(false);

  const groups = (email || 'johndoe@gmail.com').split('');

  return (
    <Row>
      <Container transition={{ staggerChildren: 0.5 }}>
        {groups.map((group: string, index: number) => {
          return [...group].map((c: string, i: number) => (
            <Circle
              key={`${index}-${c}-${i}`}
              index={index * i + i}
              show={show}
              character={c}
            />
          ));
        })}
      </Container>
      <Actions>
        <Buttons.Show show={show} setShow={setShow} />
        <Buttons.Copy text={email} />
      </Actions>
    </Row>
  );
};

const Row = styled.div`
  display: flex;
`;

const Container = styled(motion.div)`
  display: flex;
  flex: 1;
  align-items: center;
  margin-right: auto;
  max-width: 400px;
  justify-content: flex-start;  
`;

const Actions = styled.div`
  display: flex;
  margin-left: 16px;
`;

export default License;
