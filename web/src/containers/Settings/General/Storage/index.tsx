import React, { useState } from 'react';
import { useDispatch, useStore } from 'react-redux';
import styled from 'styled-components';
import { motion } from 'framer-motion';

import { Select, IndicatorSeparator, Typography } from '../../../../elements';

import { importTasks, makeTasks } from '../../../../stores/Main/reducers/tasks';
import { importProfiles, makeProfiles } from '../../../../stores/Main/reducers/profiles';
import { makeProxies } from '../../../../stores/Main/reducers/proxies';
import { makeAccounts } from '../../../../stores/Main/reducers/accounts';
import { importRates, makeRates } from '../../../../stores/Main/reducers/rates';
import { importWebhooks, makeWebhooks } from '../../../../stores/Main/reducers/webhooks';

const options = [
  "All",
  "Tasks",
  "Workflows",
  "Profiles",
  "Proxies",
  "Accounts",
  "Rates",
  "Webhooks",
];

const Storage = () => {
  const dispatch = useDispatch();

  const store = useStore();

  const [choice, setChoice] = useState<{ value: string; label: string }>({
    label: options[0],
    value: options[0],
  });

  const handleExport = () => {
    let data;

    const state = store.getState();

    switch (choice?.value) {
      case "All": {
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        const { User, News, Checkouts, Stores, ...rest } = state;
        data = rest;
        break;
      }
      case "Tasks": {
        data = makeTasks(state);
        break;
      }
      case "Profiles": {
        data = makeProfiles(state);
        break;
      }
      case "Proxies": {
        data = makeProxies(state);
        break;
      }
      case "Accounts": {
        data = makeAccounts(state);
        break;
      }
      case "Rates": {
        data = makeRates(state);
        break;
      }
      case "Webhooks": {
        data = makeWebhooks(state);
        break;
      }
      default:
        break;
    }

    if (!data) {
      return;
    }

    // TODO: Save to file
    // return saveFile(data);
  };

  const handleImport = async () => {
    // const { success, data } = await loadJsonFile();
    // if (success && choice) {
    //   const { value } = choice;
    //   switch (value) {
    //     case 'All':
    //       return dispatch(importAll(data));
    //     case 'Create':
    //       return dispatch(importTasks(data));
    //     case 'Profiles':
    //       return dispatch(importProfiles(data));
    //     case 'Proxies':
    //       return dispatch(importProxies(data));
    //     case 'Accounts':
    //       return dispatch(importAccounts(data));
    //     case 'Rates':
    //       return dispatch(importRates(data));
    //     case 'Webhooks':
    //       return dispatch(importWebhooks(data));
    //     case 'QuickTasks':
    //       return dispatch(importQuickTasks(data));
    //     default:
    //       break;
    //   }
    // }
    // return null;
  };

  return (
    <Container>
      <Heading>Software Storage</Heading>
      <RowBottom>
        <ColFill>
          <Title>Choose Data</Title>
          <SelectContainer>
            <Select
              name="software-storage"
              isClearable={false}
              placeholder="All"
              components={{ IndicatorSeparator }}
              value={choice}
              options={options.map((option) => ({
                label: option,
                value: option,
              }))}
              onChange={setChoice}
            />
          </SelectContainer>
        </ColFill>
        <ColEnd>
          <RowEnd>
            <Import
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              type="button"
              onClick={handleImport}
            >
              Import
            </Import>
            <Export
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              type="button"
              onClick={handleExport}
            >
              Export
            </Export>
          </RowEnd>
        </ColEnd>
      </RowBottom>
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
  flex-basis: 55%;
`;

const Heading = styled.div`
  font-size: 18px;
  margin-bottom: 8px;
  font-weight: 700;
  color: ${({ theme }) => theme.colors.subHeading};
`;

const RowBottom = styled.div`
  display: flex;
  flex: 1;
  align-items: flex-end;
`;

const RowEnd = styled.div`
  display: flex;
  flex: 1;
`;

const ColEnd = styled.div`
  display: flex;
  flex-direction: column;
  margin-top: auto;
`;

const ColFill = styled.div`
  display: flex;
  flex-direction: column;
  flex: 1;
`;

const Title = styled(Typography.H2)`
  font-size: 14px;
  margin: 0 0 8px 0;
  font-weight: 500;
  color: ${({ theme }) => theme.colors.paragraph};
`;

const Import = styled(motion.button)`
  height: 36px;
  font-size: 14px;
  font-weight: 700;
  background: transparent;
  color: ${({ theme }) => theme.colors.primary};
  width: 7.5vw;
  padding: 8px 16px;
  cursor: pointer;
  margin-right: 8px;
  border: 1px solid ${({ theme }) => theme.colors.primary};
  border-radius: 4px;
`;

const Export = styled(motion.button)`
  height: 36px;
  font-size: 14px;
  font-weight: 700;
  background: ${({ theme }) => theme.colors.primary};
  color: #fff;
  width: 7.5vw;
  padding: 8px 16px;
  cursor: pointer;
  margin-left: 8px;
  border: 1px solid ${({ theme }) => theme.colors.primary};
  border-radius: 4px;
`;

const SelectContainer = styled.div`
  margin-right: auto;
  width: 15vw;
  max-width: 15vw;
`;

export default Storage;
