import styled from "styled-components";

const Dragbar = () => {
  return (
    <Container />
  );
};

const Container = styled.div`
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 28px;
  -webkit-app-region: drag;
`;

export default Dragbar;
