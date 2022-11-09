import { countries } from "../constants";

const buildOptions = (list: any[], name: string, code: string) =>
  list.map((datum) => ({ name: datum[name], code: datum[code] }));

export const isProvinceDisabled = (country: any, disabled?: boolean) => {
  if (country && country.value) {
    const { provinces } = getCountry(country.value);
    if (!provinces || !provinces.length) {
      return true;
    }
  }
  return disabled;
};

export const getProvinces = (code: string) => {
  const country = getCountry(code);
  return country && country.provinces;
};

const getCountry = (code: string) =>
  Object.assign(
    {},
    countries.find((country) => country.code === code)
  );

const buildCountryOptions = () => buildOptions(countries, "name", "code");

const buildProvinceOptions = (country: any) => {
  if (country && country.code) {
    return buildOptions(getProvinces(country.code), "name", "code");
  }

  return [];
};

export const ProfileUtils = {
  countries,
  isProvinceDisabled,
  getProvinces,
  buildCountryOptions,
  buildProvinceOptions,
};
