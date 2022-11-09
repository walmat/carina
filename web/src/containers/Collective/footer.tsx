import { useCallback, ChangeEvent } from "react";
import { useDispatch } from "react-redux";
import styled from "styled-components";

import {
  HARVESTER_FIELDS,
  edit,
} from "../../stores/Collective/reducers/harvester";
import { Store } from "../../stores/Main/reducers/stores";

import { Select, Input, Control, IndicatorSeparator } from "../../elements";
import { Harvester } from "./types";

const FILTER_OPTIONS = {
  platforms: [
    {
      id: "yeezysupply",
      name: "Yeezy Supply",
      sites: [
        {
          id: "yeezysupply",
          name: "Yeezy Supply",
          platform: "yeezysupply",
          url: "https://yeezysupply.com/",
        },
      ],
    },
    {
      id: "shopify",
      name: "Shopify",
      sites: [
        {
          id: "kith",
          name: "Kith",
          platform: "shopify",
          url: "https://kith.com/",
        },
      ],
    },
  ],
};

let siteOptions = FILTER_OPTIONS.platforms
  .map((platform) => platform.sites)
  .flat();

const Footer = (harvester: Harvester) => {
  const dispatch = useDispatch();

  const selectHandler = useCallback((field: string, value: Store | null) => {
    console.log("store  brh", harvester.store?.platform);
    dispatch(edit({ harvester, fieldsToUpdate: { [field]: value } }));
  }, []);

  const editHandler = useCallback(
    (e: ChangeEvent<HTMLInputElement>, field: string) => {
      dispatch(
        edit({ harvester, fieldsToUpdate: { [field]: e.target.value } })
      );
    },
    []
  );

  return (
    <Container>
      {/*<Row m="0 0 8px 0">
        <Col m="0 8px 0 16px">
          <Select
            isClearable
            name="platform"
            placeholder="Harvester Platform"
            label="Harvester Platform"
            disabled={true}
            getOptionLabel={(option: any) => option.name}
            getOptionValue={(option: any) => option.id}
            components={{ Control, IndicatorSeparator }}
            // value={harvester.store?.platform}
            value={"hello world"}
            options={FILTER_OPTIONS.platforms}
            onChange={(e: any | null) =>
              selectHandler(HARVESTER_FIELDS.STORE, {
                id: e.id,
                name: e.name,
                platform: e.id,
                url: "",
                modes: [],
              })
            }
          />
        </Col>
        <Col m="0 16px 0 8px">
          <Select
            isClearable
            name="store"
            placeholder="Harvester Site"
            label="Harvester Site"
            disabled={true}
            getOptionLabel={(option: Store) => option.name}
            getOptionValue={(option: Store) => option.url}
            components={{ Control, IndicatorSeparator }}
            value={harvester.store?.url}
            options={siteOptions}
            onChange={(e: Store | null) =>
              selectHandler(HARVESTER_FIELDS.STORE, e)
            }
          />
        </Col>
      </Row>*/}
      <Row m="8px 16px 16px 16px">
        <Input
          useLabel
          id="proxy"
          name="proxy"
          placeholder="Harvester Proxy (Optional)"
          onChange={(e: ChangeEvent<HTMLInputElement>) =>
            editHandler(e, HARVESTER_FIELDS.PROXY)
          }
          value={harvester.proxy}
        />
      </Row>
    </Container>
  );
};

const Container = styled.div`
  margin-top: auto;
  display: flex;
  flex-direction: column;
`;

const Col = styled.div<{ m: string; f?: number }>`
  display: flex;
  flex: 1;
  flex-direction: column;
  margin: ${({ m }) => m};
  ${({ f }) => (typeof f !== undefined ? `flex: ${f};` : "")}
`;

const Row = styled.div<{ m?: string }>`
  display: flex;
  justify-content: center;
  flex: 1;
  ${({ m }) =>
    typeof m !== "undefined"
      ? `
		margin: ${m};
	`
      : ""}
`;

export default Footer;
