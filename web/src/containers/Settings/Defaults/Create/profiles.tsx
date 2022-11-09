import React, { useEffect, useState } from "react";
import { useSelector } from "react-redux";
import { motion } from "framer-motion";
import styled from "styled-components";

import { makeProfiles, ProfileGroup } from "../../../../stores/Main/reducers/profiles";
import { IndicatorSeparator, Control } from '../../../../elements';
import SelectFormik from '../../../../elements/Select/SelectFormik';
import InputFormik from '../../../../elements/Input/InputFormik';

interface ProfilesProps {
  values: any;
  errors: any;
  touched: any;
  handleChange: any;
  setFieldTouched: any;
  setFieldValue: any;
}

const Profiles = ({
  values,
}: ProfilesProps) => {
  const [profiles, setProfiles] = useState<ProfileGroup[]>([]);
  const profileGroups = useSelector(makeProfiles);

  useEffect(() => {
    if (!values.profileGroup) {
      return;
    }

    setProfiles(Object.values(values.profileGroup.byId));
  }, [values.profileGroup]);

  return (
    <motion.div>
      <Row m="16px 0">
        <SelectFormik
          key="profiles-group"
          name="profileGroup"
          required
          isClearable={false}
          label="Profile Group"
          placeholder="Profile Group"
          components={{
            IndicatorSeparator,
            Control,
          }}
          getOptionLabel={(option: any) => option.name}
          getOptionValue={(option: any) => option.id}
          options={Object.values(profileGroups)}
        />
      </Row>

      <Row m="0 0 16px 0">
        <SelectFormik
          isMulti
          isClearable
          closeMenuOnSelect={false}
          key="profiles-profiles"
          name="profiles"
          label="Specific Profiles"
          placeholder="Specific Profiles (Optional)"
          components={{
            IndicatorSeparator,
            Control,
          }}
          getOptionLabel={(option: any) => option.name}
          getOptionValue={(option: any) => option.id}
          options={profiles}
        />
      </Row>

      <Row>
        <InputFormik
          key="profiles-limit"
          id="profileLimit"
          name="profileLimit"
          useLabel
          placeholder="Checkout Limit (Profile)"
          textTransform="capitalize"
        />
      </Row>
    </motion.div>
  );
};

const Row = styled.div<{ m?: string }>`
  display: flex;
  justify-content: space-between;
  margin: ${({ m }) => (m ? m : "0")};
`;

export default Profiles;
