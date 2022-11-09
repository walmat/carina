import React from "react";
import { useSelector } from "react-redux";
// @ts-ignore
import Select, { createFilter } from "react-select";

import { makeTheme } from '../../stores/Main/reducers/theme';

export const IndicatorSeparator = () => null;

export const primaryStyles = (
  theme: number | null,
  startIcon = false,
  type: string
) => {
  let borderColor = theme === 0 ? "#979797" : "#616161";

  return {
    control: (styles: any, state: any) => {
      const { isFocused } = state;

      const overrides: any = {};
      if (isFocused) {
        overrides.borderColor = theme === 0 ? "#202126" : "#f4f4f4";
      }

      return {
        ...styles,
        display: "flex",
        flex: 1,
        borderWidth: `1px`,
        borderStyle: "solid",
        borderColor,
        minWidth: 80,
        color: theme === 0 ? "#616161" : "#D8D8D8",
        backgroundColor: theme === 0 ? "#fff" : "#2E2F34",
        height: "28px",
        fontSize: "12px",
        minHeight: "28px",
        borderRadius:
          type === "full"
            ? "4px"
            : type === "first"
            ? "4px 0 0 4px"
            : "0 4px 4px 0",
        outline: "none",
        cursor: "pointer",
        boxShadow: "none",
        whiteSpace: "nowrap",
        overflow: "hidden",
        textOverflow: "ellipsis",
        transition: "background-color, color, border .15 ease-in-out",
        ":hover": {
          cursor: "pointer",
          ...overrides,
        },
        ...overrides,
      };
    },
    input: (styles: any) => ({
      ...styles,
      color: theme === 0 ? "#616161" : "#d8d8d8",
    }),
    indicator: (styles: any) => ({
      ...styles,
      padding: "0 8px",
    }),
    clearIndicator: (styles: any) => ({
      ...styles,
      zIndex: 9999,
      padding: 0,
      ":hover": {
        ...styles[":hover"],
        color: "hsl(0,0%,70%)",
      },
    }),
    noOptionsMessage: (styles: any) => ({
      ...styles,
      padding: "4px 8px",
    }),
    indicatorsContainer: (styles: any) => ({
      ...styles,
      height: "100%",
      margin: "0 8px 0 0",
    }),
    dropdownIndicator: (styles: any) => ({
      ...styles,
      color: theme === 0 ? "#616161" : "#d8d8d8",
      padding: 0,
      display: "flex",
      justifyContent: "center",
      alignItems: "center",
      width: 12,
      ":hover": {
        ...styles[":hover"],
        color: theme === 0 ? "#616161" : "#d8d8d8",
      },
    }),
    multiValue: (styles: any) => ({
      ...styles,
      minWidth: "unset",
      margin: "0 8px 0 0",
      color: "#fff",
      background: "linear-gradient(97.53deg, #776DF2 0%, #978DFD 100%)",
    }),
    multiValueLabel: (styles: any) => ({
      ...styles,
      color: "#fff",
      fontWeight: 500,
      textTransform: "capitalize",
      padding: "4px",
      paddingLeft: 8,
    }),
    multiValueRemove: (styles: any) => ({
      ...styles,
      color: "#fff",
      paddingLeft: 0,
      ":hover": {
        opacity: 0.6,
      },
    }),
    menuPortal: (styles: any) => ({
      ...styles,
      zIndex: 9999,
    }),
    placeholder: (styles: any) => ({
      ...styles,
      color: theme === 0 ? "hsl(0,0%,70%)" : "hsl(0,0%,40%)",
    }),
    option: (styles: any, { isFocused, isSelected }: any) => {
      const color = theme === 0 ? "#616161" : "#d8d8d8";

      const retVal = {
        ...styles,
        fontSize: "12px",
        fontWeight: 400,
        padding: "4px 8px",
        textOverflow: "clip",
        whiteSpace: "nowrap",
        cursor: "pointer",
        outline: "none",
        boxShadow: "none",
        borderRadius: 0,
        color,
        overflow: "hidden",
        ":active": {
          backgroundColor: "#786EF2",
          color: "#fff",
        },
        ":selected": {
          color: "#fff",
        },
      };
      if (isSelected) {
        return {
          ...retVal,
          color: "#fff",
          fontWeight: 500,
          backgroundColor: "#776DF2",
        };
      }
      if (isFocused) {
        return {
          ...retVal,
          color: "#786EF2",
          backgroundColor: "rgba(120, 110, 242, 0.2)",
        };
      }

      return retVal;
    },
    valueContainer: (styles: any, { isMulti }: any) => {
      const ret = {
        ...styles,
        maxHeight: "29px",
        padding: "0px 8px",
        margin: 0,
        cursor: "pointer",
        position: "static",
      };

      if (startIcon) {
        ret.padding = 0;
        ret.margin = "0 0 0 36px";
      }

      if (isMulti) {
        return {
          ...ret,
          overflowX: "scroll",
          overflowY: "hidden",
          flexWrap: "nowrap",
          "::-webkit-scrollbar": {
            width: "0px",
            height: "0px",
            background: "transparent",
          },
        };
      }
      return ret;
    },
    singleValue: ({ maxWidth, position, top, transform, ...styles }: any) => ({
      ...styles,
      margin: 0,
      flex: 1,
      display: "flex",
      color: theme === 0 ? "#616161" : "#d8d8d8",
      cursor: "pointer",
    }),
    menu: ({ ...styles }: any) => ({
      ...styles,
      width: "unset",
      minWidth: "100%",
      maxHeight: 190,
      padding: 0,
      margin: "4px 0",
      overflow: "hidden",
      borderRadius: 2,
      backgroundColor: theme === 0 ? "#fff" : "#2E2F34",
    }),
    menuList: (styles: any) => ({
      ...styles,
      maxHeight: 190,
      fontSize: "12px",
      padding: 0,
      borderRadius: 2,
      backgroundColor: theme === 0 ? "#fff" : "#2E2F34",
    }),
  };
};

interface SelectProps {
  name: string;
  type?: "full" | "first" | "second";
  required?: boolean;
  isMulti?: boolean;
  isClearable?: boolean;
  startIcon?: boolean;
  placeholder: string;
  value: any;
  onChange: any;
  onBlur?: any;
  options: any[];
  getOptionLabel?: any;
  getOptionValue?: any;
}

const DateSelect = ({
  required = true,
  type = "full",
  isClearable = false,
  startIcon = false,
  placeholder,
  name,
  value,
  onChange,
  onBlur,
  options,
  getOptionLabel,
  getOptionValue,
  ...rest
}: SelectProps) => {
  const theme = useSelector(makeTheme);

  const base: any = { ...rest };
  if (getOptionLabel) {
    base.getOptionLabel = getOptionLabel;
  }

  if (getOptionValue) {
    base.getOptionValue = getOptionValue;
  }

  return (
    <Select
      name={name}
      required={required}
      isClearable={isClearable}
      ignoreAccents={false}
      backspaceRemovesValue={false}
      filterOption={createFilter({ ignoreAccents: false })}
      menuPortalTarget={document.body}
      menuPlacement="auto"
      menuPosition="fixed"
      classNamePrefix="select"
      placeholder={placeholder}
      components={{ IndicatorSeparator }}
      value={value}
      onChange={onChange}
      onBlur={onBlur}
      options={options}
      styles={primaryStyles(theme, startIcon, type)}
      {...base}
    />
  );
};

DateSelect.defaultProps = {
  required: true,
  type: "full",
  startIcon: false,
  isCreatable: false,
  isMulti: false,
  isClearable: false,
};

export default DateSelect;
