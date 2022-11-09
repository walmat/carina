import { useSelector } from "react-redux";
import styled from "styled-components";
import { Hash } from 'react-feather';

import { makeVersion } from "../../../../stores/Main/reducers/version";

const Version = () => {
  const version = useSelector(makeVersion);

  return (
    <Row>
      <VersionHash>
        <VersionIcon />
        {version}
      </VersionHash>
    </Row>
  );
};

const Row = styled.div`
  display: flex;
  margin: 6px 0 0 0;
`;

const VersionIcon = styled(Hash)`
  height: 14px;
  width: auto;
  margin: 0 4px 0 0;
  color: ${({ theme }) => theme.colors.primary};
`;

const VersionHash = styled.div`
  font-size: 14px;
  font-weight: 400;
  color: ${({ theme }) => theme.colors.primary};
  padding: 6px 8px;
  background-color: ${({ theme }) => theme.colors.secondary};
  border-radius: 4px;
  font-weight: 500;
  display: flex;
  cursor: pointer;
  justify-content: center;
  align-items: center;
`;

export default Version;
