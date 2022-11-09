import React, { ReactElement } from "react";
import { useSelector } from "react-redux";
import Select, { components, createFilter } from 'react-windowed-select';
// @ts-ignore
import Creatable from "react-select/creatable";
import { Icon } from "react-feather";

import { makeTheme } from '../../stores/Main/reducers/theme';

export const IndicatorSeparator = () => null;

interface ControlProps {
  Icon: Icon;
  children: ReactElement;
}

const mvStyle = {
  margin: "0 0 0 4px",
};

const spanStyle = { cursor: "pointer", display: "flex", alignItems: "center" };
const iconStyle = {
  zIndex: 999,
  position: "absolute",
  left: 9.5,
  height: 12,
  width: 12,
  top: 15,
};

const labelStyle: any = (
  showLabel: boolean,
  theme: number,
  hasIcon: boolean
) => ({
  transition: "all .2s ease-out",
  fontSize: 11,
  fontWeight: 400,
  position: "absolute",
  top: 0,
  zIndex: 1,
  width: "100%",
  userSelect: "none",
  transform: showLabel ? "none" : "translateY(2.5px)",
  pointerEvents: "none",
  overflow: "hidden",
  whiteSpace: "nowrap",
  textOverflow: "ellipsis",
  boxSizing: "border-box",
  color: theme === 0 ? "rgba(0, 0, 0, 0.4)" : "rgba(255, 255, 255, 0.4)",
  opacity: showLabel ? 1 : 0,
  margin: "0.5em 0",
  padding: "0 8px",
  paddingLeft: hasIcon ? "33px" : "8px",
});

export const MultiValue = ({ children, ...props }: ControlProps) => {
  // @ts-ignore
  const { value } = props.selectProps;
  if (value.length === 1) {
    return (
      // @ts-ignore
      <components.MultiValue {...props}>{children}</components.MultiValue>
    );
  }

  if (value[0].name === children) {
    return (
      // @ts-ignore
      <components.MultiValue {...props}>
        {children}
        <span style={mvStyle}> • </span>
      </components.MultiValue>
    );
  }

  if (value[value.length - 1].name === children) {
    return (
      // @ts-ignore
      <components.MultiValue {...props}>{children}</components.MultiValue>
    );
  }

  return (
    // @ts-ignore
    <components.MultiValue {...props}>
      {children}
      <span style={mvStyle}> • </span>
    </components.MultiValue>
  );
};

export const Control = ({ children, ...props }: ControlProps) => {
  // @ts-ignore
  const { Icon, label, theme } = props.selectProps;
  // @ts-ignore
  const { hasValue } = props;

  if (label) {
    return (
      // @ts-ignore
      <components.Control {...props}>
        <span style={spanStyle}>{Icon && <Icon style={iconStyle} />}</span>
        {hasValue ? (
          <label style={labelStyle(hasValue, theme, !!Icon)}>{label}</label>
        ) : null}
        {children}
      </components.Control>
    );
  }

  return (
    // @ts-ignore
    <components.Control {...props}>
      <span style={spanStyle}>{Icon && <Icon style={iconStyle} />}</span>
      {children}
    </components.Control>
  );
};

export const primaryStyles = (
  theme: number | null,
  startIcon = false,
  size?: string,
  label?: string,
  showLabel?: boolean,
  error?: boolean,
  touched?: boolean
) => {
  let borderColor = theme === 0 ? "#979797" : "#616161";
  if (error && touched) {
    borderColor = "#F26E86";
  }

  const height = size === "small" ? 24 : 36;
  const fontSize = size === "small" ? 10 : 12;
  const padding = size === "small" ? "6px 8px" : "12px 8px";

  const container =
    size === "small"
      ? {}
      : {
          display: "flex",
          flex: 1,
        };

  return {
    container: (styles: any) => ({
      ...styles,
      ...container,
    }),
    control: (styles: any, state: any) => {
      const { isFocused } = state;

      const overrides: any = {};
      if (isFocused) {
        overrides.borderColor = theme === 0 ? "#202126" : "#f4f4f4";
      }

      return {
        ...styles,
        borderWidth: `1px`,
        borderStyle: "solid",
        borderColor,
        flex: 1,
        minWidth: 88,
        color: theme === 0 ? "#616161" : "#d8d8d8",
        backgroundColor: theme === 0 ? "#fff" : "#2E2F34",
        height,
        fontSize,
        minHeight: height,
        borderRadius: 4,
        outline: "none",
        cursor: "pointer",
        boxShadow: "none",
        whiteSpace: "nowrap",
        overflow: "hidden",
        textOverflow: "ellipsis",
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
      padding: padding,
    }),
    indicatorsContainer: (styles: any) => ({
      ...styles,
      height: "100%",
      margin: "0 8px",
    }),
    dropdownIndicator: (styles: any) => ({
      ...styles,
      color: theme === 0 ? "#616161" : "#d8d8d8",
      padding: 0,
      display: "flex",
      justifyContent: "center",
      alignItems: "center",
      width: size === "small" ? 12 : 16,
      ":hover": {
        ...styles[":hover"],
        color: theme === 0 ? "#616161" : "#d8d8d8",
      },
    }),
    multiValue: (styles: any) => ({
      ...styles,
      minWidth: "unset",
      margin: "0 8px 0 0 !important",
      color: "#fff",
      background: "linear-gradient(97.53deg, #776DF2 0%, #978DFD 100%)",
    }),
    multiValueLabel: (styles: any) => ({
      ...styles,
      color: "#fff",
      fontWeight: 400,
      textTransform: "capitalize",
      padding: "2px 6px",
      marginTop: -2,
      paddingLeft: 0,
      fontSize: "100%",
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
      margin: 0,
      color: theme === 0 ? "hsl(0,0%,70%)" : "hsl(0,0%,40%)",
    }),
    option: (styles: any, { isFocused, isSelected }: any) => {
      const color = theme === 0 ? "#616161" : "#d8d8d8";

      const retVal = {
        ...styles,
        fontSize: fontSize,
        fontWeight: 400,
        padding: padding,
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
        padding: "2px 8px",
        margin: 0,
        cursor: "pointer",
        position: "static",
        "& > div": {
          margin:  0,
        },
      };

      if (startIcon) {
        ret.padding = 0;
        ret.margin = "0 0 0 32px";
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
    singleValue: (styles: any) => ({
      ...styles,
      margin: 0,
      maxWidth: "65%",
      color: theme === 0 ? "#616161" : "#d8d8d8",
      textTransform: "capitalize",
      cursor: "pointer",
    }),
    groupHeading: (styles: any) => ({
      ...styles,
      padding: "2px 8px 0 8px",
    }),
    menu: (styles: any) => ({
      ...styles,
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
      fontSize: fontSize,
      padding: 0,
      borderRadius: 2,
      backgroundColor: theme === 0 ? "#fff" : "#2E2F34",
    }),
  };
};

export const primaryStylesWithLabel = (
  theme: number | null,
  startIcon = false,
  size?: string,
  label?: string,
  showLabel?: boolean,
  error?: boolean,
  touched?: boolean
) => {
  let borderColor = theme === 0 ? "#979797" : "#616161";
  if (error && touched) {
    borderColor = "#F26E86";
  }

  const fontSize = size === "small" ? 10 : 12;
  const padding = size === "small" ? "6px 8px" : "12px 8px";

  const container =
    size === "small"
      ? {}
      : {
          display: "flex",
          flex: 1,
          padding: 0,
        };

  return {
    container: (styles: any) => ({
      ...styles,
      ...container,
    }),
    control: (styles: any, state: any) => {
      const { isFocused } = state;

      const overrides: any = {};
      if (isFocused) {
        overrides.borderColor = theme === 0 ? "#202126" : "#f4f4f4";
      }

      return {
        ...styles,
        borderWidth: `1px`,
        borderStyle: "solid",
        maxHeight: 44,
        padding: label ? "1.165em 0px" : "10px 0",
        borderColor,
        flex: 1,
        minWidth: 88,
        color: theme === 0 ? "#616161" : "#d8d8d8",
        backgroundColor: theme === 0 ? "#fff" : "#2E2F34",
        fontSize: fontSize,
        borderRadius: 4,
        outline: "none",
        cursor: "pointer",
        boxShadow: "none",
        whiteSpace: "nowrap",
        overflow: "hidden",
        textOverflow: "ellipsis",
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
      padding: padding,
    }),
    indicatorsContainer: (styles: any) => ({
      ...styles,
      height: "100%",
      margin: "0 8px",
    }),
    dropdownIndicator: (styles: any) => ({
      ...styles,
      color: theme === 0 ? "#616161" : "#d8d8d8",
      padding: 0,
      display: "flex",
      justifyContent: "center",
      alignItems: "center",
      width: size === "small" ? 12 : 16,
      ":hover": {
        ...styles[":hover"],
        color: theme === 0 ? "#616161" : "#d8d8d8",
      },
    }),
    multiValue: (styles: any) => ({
      ...styles,
      minWidth: "unset",
      fontSize: "100%",
      color: theme === 0 ? "#616161" : "#d8d8d8",
      background: "transparent",
    }),
    multiValueLabel: (styles: any) => ({
      ...styles,
      color: theme === 0 ? "#616161" : "#d8d8d8",
      fontWeight: 400,
      fontSize: "100%",
      textTransform: "capitalize",
      padding: "0px 6px 4px 6px",
      paddingLeft: 0,
    }),
    multiValueRemove: (styles: any) => ({
      ...styles,
      display: "none",
    }),
    menuPortal: (styles: any) => ({
      ...styles,
      zIndex: 9999,
    }),
    placeholder: (styles: any) => ({
      ...styles,
      margin: 0,
      color: theme === 0 ? "hsl(0,0%,70%)" : "hsl(0,0%,40%)",
    }),
    option: (styles: any, { isFocused, isSelected }: any) => {
      const color = theme === 0 ? "#616161" : "#d8d8d8";

      const retVal = {
        ...styles,
        fontSize: fontSize,
        fontWeight: 400,
        padding: padding,
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
      let ret = {
        ...styles,
        maxHeight: "29px",
        padding: "2px 8px",
        margin: 0,
        height: "100%",
        position: "static",
        "& > div": {
          margin: 0,
        },
      };

      if (showLabel) {
        ret = {
          ...ret,
          position: "absolute",
          height: "unset",
          bottom: 0,
          padding: "4px 0px 0px 8px",
          width: "100%",
          "& > div": {
            margin: 'auto 0 4px 0',
          },
        };
      }

      if (startIcon) {
        ret.padding = 0;
        ret.margin = "0 0 0 32px";
      }

      if (isMulti) {
        return {
          ...ret,
          overflowX: "scroll",
          overflowY: "hidden",
          flexWrap: "nowrap",
          padding: "3px 8px 0 8px",
          "::-webkit-scrollbar": {
            width: "0px",
            height: "0px",
            background: "transparent",
          },
        };
      }
      return ret;
    },
    singleValue: (styles: any) => ({
      ...styles,
      margin: 0,
      maxWidth: "65%",
      color: theme === 0 ? "#202126" : "#f4f4f4",
      textTransform: "capitalize",
      cursor: "pointer",
    }),
    groupHeading: (styles: any) => ({
      ...styles,
      padding: "2px 8px 0 8px",
    }),
    menu: (styles: any) => ({
      ...styles,
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
      fontSize,
      padding: 0,
      borderRadius: 2,
      backgroundColor: theme === 0 ? "#fff" : "#2E2F34",
    }),
  };
};

interface SelectProps {
  useDefault?: boolean;
  autoFocus?: boolean;
  disabled?: boolean;
  name: string;
  required?: boolean;
  isCreatable?: boolean;
  isMulti?: boolean;
  isClearable?: boolean;
  Icon?: any;
  label?: string;
  closeMenuOnSelect?: boolean;
  isOptionDisabled?: (option: any) => boolean;
  noOptionsMessage?: () => undefined;
  error?: boolean;
  touched?: boolean;
  placeholder: string;
  components?: any;
  value: any;
  onChange: any;
  onFocus?: any;
  onBlur?: any;
  options: any[];
  getOptionLabel?: any;
  getOptionValue?: any;
  size?: "small" | "default";
}

const SelectPrimitive = ({
  useDefault = true,
  autoFocus = false,
  required = true,
  disabled = false,
  isCreatable = false,
  isMulti = false,
  isClearable = false,
  Icon = null,
  label = "",
  closeMenuOnSelect = true,
  isOptionDisabled = () => false,
  noOptionsMessage = () => undefined,
  error = false,
  touched = false,
  placeholder,
  components,
  name,
  value,
  onChange,
  onFocus,
  onBlur,
  options,
  getOptionLabel,
  getOptionValue,
  size,
}: SelectProps) => {
  const theme = useSelector(makeTheme);

  const base: any = {};
  if (getOptionLabel) {
    base.getOptionLabel = getOptionLabel;
  }

  if (getOptionValue) {
    base.getOptionValue = getOptionValue;
  }

  const styles = label ? primaryStylesWithLabel : primaryStyles;
  const showLabel = isMulti ? value.length : !!value;

  if (isCreatable) {
    // TODO: Creatable select
    return (
      <Creatable
        createOptionPosition="first"
        closeMenuOnScroll
        isDisabled={disabled}
        autoFocus={autoFocus}
        name={name}
        Icon={Icon}
        required={required}
        isMulti={isMulti}
        theme={theme}
        label={label}
        isClearable={isClearable}
        closeMenuOnSelect={closeMenuOnSelect}
        noOptionsMessage={noOptionsMessage}
        isOptionDisabled={isOptionDisabled}
        ignoreAccents={false}
        onSelectResetsInput={false}
        filterOption={createFilter({ ignoreAccents: false })}
        menuPortalTarget={document.body}
        menuPlacement="auto"
        classNamePrefix="select"
        placeholder={placeholder}
        components={{
          IndicatorSeparator,
          ...components,
        }}
        value={value}
        onFocus={onFocus}
        onChange={onChange}
        onBlur={onBlur}
        options={
          useDefault ? options : options.filter(({ id }) => id !== "default")
        }
        styles={styles(
          theme,
          !!Icon,
          size,
          label,
          showLabel,
          Boolean(error),
          Boolean(touched)
        )}
        {...base}
      />
    );
  }

  return (
    <Select
      closeMenuOnScroll
      autoFocus={autoFocus}
      name={name}
      Icon={Icon}
      isDisabled={disabled}
      required={required}
      isMulti={isMulti}
      theme={theme}
      label={label}
      isClearable={isClearable}
      closeMenuOnSelect={closeMenuOnSelect}
      noOptionsMessage={noOptionsMessage}
      isOptionDisabled={isOptionDisabled}
      ignoreAccents={false}
      filterOption={createFilter({ ignoreAccents: false })}
      menuPortalTarget={document.body}
      menuPlacement="auto"
      classNamePrefix="select"
      placeholder={placeholder}
      components={{
        IndicatorSeparator,
        MultiValue,
        ...components,
      }}
      value={value}
      onFocus={onFocus}
      onChange={onChange}
      onBlur={onBlur}
      options={
        useDefault ? options : options.filter(({ id }) => id !== "default")
      }
      styles={styles(
        theme,
        !!Icon,
        size,
        label,
        showLabel,
        Boolean(error),
        Boolean(touched)
      )}
      {...base}
    />
  );
};

SelectPrimitive.defaultProps = {
  required: true,
  startIcon: false,
  isCreatable: false,
  isMulti: false,
  isClearable: false,
  label: null,
  size: "default",
};

SelectPrimitive.IndicatorSeparator = IndicatorSeparator;
SelectPrimitive.Control = Control;

export default SelectPrimitive;
