import React, { useState } from "react";
import styled from "styled-components";
import { Search as Magnifier } from "react-feather";

import { Input } from "../../elements";

interface SearchProps {
  placeholder: string;
  list: Stores | Releases;
  setList: any;
  property: string;
}

const Search = ({ placeholder, list, setList, property }: SearchProps) => {
  const [value, setValue] = useState("");

  const handleChange = async (event: any) => {
    const { value } = event.target;

    const filtered = Object.values(list).filter((obj) =>
      obj[property].toLowerCase().includes(value.toLowerCase())
    );

    setValue(value);
    setList(filtered);
  };

  return (
    <Container>
      <Input
        onChange={handleChange}
        value={value}
        maxLength={48}
        radius={8}
        name="Search"
        id="search"
        Icon={Magnifier}
        placeholder={placeholder}
      />
    </Container>
  );
};

const Container = styled.div`
  display: flex;
  flex: 0;
  margin-bottom: 16px;
`;

export default Search;
