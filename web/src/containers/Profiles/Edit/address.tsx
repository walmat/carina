import PlacesAutocomplete, {
  geocodeByPlaceId,
  Suggestion,
} from "react-places-autocomplete";
import { Home } from "react-feather";
import styled from "styled-components";

import { Input } from "../../../elements";

import { countries } from "../../../constants";

interface AddressProps {
  values: any;
  errors: any;
  touched: any;
  type: "billing" | "shipping";
  value: string;
  setFieldValue: any;
  setFieldTouched: any;
}

const Address = ({
  values,
  errors,
  touched,
  type,
  value,
  setFieldValue,
  setFieldTouched,
}: AddressProps) => {
  const onChange = (value: string) => {
    setFieldValue(`${type}.line1`, value);
  };

  const handleSelect = async (_: any, placeId: string) => {
    const [place] = await geocodeByPlaceId(placeId);

    const { address_components: addressComponents } = place;

    const { long_name: streetNumber = "" } =
      addressComponents.find((comp: any) =>
        comp.types.includes("street_number")
      ) || {};

    const { long_name: route = "" } =
      addressComponents.find((comp: any) => comp.types.includes("route")) || {};

    setFieldValue(`${type}.line1`, `${streetNumber} ${route}`);

    const { short_name: countryCode } =
      addressComponents.find((comp: any) => comp.types.includes("country")) ||
      {};

    if (countries) {
      const country = countries.find(({ code }: any) => code === countryCode);

      if (country) {
        const { provinces, ...countryValue } = country;
        setFieldValue(`${type}.country`, countryValue);

        if (
          values[type].country &&
          values[type].country.value !== country.code
        ) {
          setFieldValue(`${type}.state`, null);
        }

        const { short_name: provinceCode } =
          addressComponents.find((comp: any) =>
            comp.types.includes("administrative_area_level_1")
          ) || {};

        const province = country.provinces.find(
          ({ code }: any) => code === provinceCode
        );

        if (province) {
          setFieldValue(`${type}.state`, province);
        }
      }
    }

    const { long_name: city = "" } =
      addressComponents.find((comp: any) => comp.types.includes("locality")) ||
      {};

    setFieldValue(`${type}.city`, city);

    const { long_name: postalCode = "" } =
      place.address_components.find((comp: any) =>
        comp.types.includes("postal_code")
      ) || {};

    setFieldValue(`${type}.postCode`, postalCode);
  };

  return (
    <PlacesAutocomplete
      value={value}
      debounce={250}
      highlightFirstSuggestion
      shouldFetchSuggestions={`${type}.line1`.length > 3}
      onChange={onChange}
      onSelect={handleSelect}
    >
      {({ getInputProps, suggestions, getSuggestionItemProps }: any) => (
        <>
          <Input
            {...getInputProps({
              id: `${type}.line1`,
              name: `${type}.line1`,
              required: true,
              useLabel: true,
              textTransform: "capitalize",
              error: !!errors?.[type]?.line1,
              touched: !!touched?.[type]?.line1,
              onBlur: () => setFieldTouched(`${type}.line1`, true),
              placeholder: "Address",
              key: `${type}.line1`,
              Icon: Home,
            })}
          />
          {suggestions.length ? (
            <MenuList>
              {suggestions.map((suggestion: Suggestion) => {
                return (
                  <MenuItem
                    key={suggestion.placeId}
                    active={suggestion.active}
                    {...getSuggestionItemProps(suggestion)}
                  >
                    {suggestion.description}
                  </MenuItem>
                );
              })}
            </MenuList>
          ) : null}
        </>
      )}
    </PlacesAutocomplete>
  );
};

const MenuList = styled.div`
  display: flex;
  flex-direction: column;
  position: absolute;
  margin-top: 40px;
  z-index: 999;
  width: 384px;
  border-radius: 4px;
  max-height: 152px;
  overflow-x: hidden;
  overflow-y: scroll;
  background-color: ${({ theme }) => theme.colors.sidebar};
  box-shadow: 0 0 0 1px hsla(0, 0%, 0%, 0.1), 0 4px 11px hsla(0, 0%, 0%, 0.1);
`;

const MenuItem = styled.div<{ active: boolean }>`
  display: flex;
  padding: 12px 8px;
  color: ${({ theme }) => theme.colors.paragraph};
  font-size: 12px;
  font-weight: 400;
  cursor: pointer;

  ${({ active, theme }) =>
    active
      ? `
		background-color: ${theme.colors.primary};
		color: #fff;
	`
      : ""}
`;

export default Address;
