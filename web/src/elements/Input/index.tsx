import * as React from "react";
import { useCallback, useEffect } from "react";
import { motion, useAnimation, transform } from "framer-motion";
import { Icon } from "react-feather";
import styled from "styled-components";

const mapRemainingToColor = transform([2, 6], ["#F26E86", "#ccc"]);
const mapRemainingToSpringVelocity = transform([0, 5], [50, 0]);

export interface InputProps {
  id: string;
  name: string;
  type?: string;
  radius?: number;
  textTransform?: string;
  label?: string;
  useLabel?: boolean;
  value: string | number;
  restriction?: string;
  error?: boolean;
  touched?: boolean;
  masked?: boolean;
  disabled?: boolean;
  autoFocus?: boolean;
  placeholder: string;
  maxLength?: number;
  onSubmit?: any;
  Icon?: Icon | null;
  onChange: (e: any) => void;
  onBlur?: (e: any) => void;
  onFocus?: (e: any) => void;
}

const makeValue = (value: string, masked: boolean) => {
  if (masked) {
    return "•".repeat(value.length);
  }

  return value;
};

export const Input = ({
  id,
  name,
  type = "text",
  radius = 4,
  textTransform = "none",
  restriction = "",
  value,
  useLabel = false,
  label = "",
  error = false,
  touched = false,
  masked = false,
  disabled = false,
  autoFocus = false,
  placeholder,
  maxLength = 0,
  onSubmit = null,
  Icon = null,
  onChange,
  onBlur,
  onFocus,
  ...props
}: InputProps) => {
  const charactersRemaining = maxLength - `${value}`.length;
  const controls = useAnimation();

  const onKeyPress = useCallback((e) => {
    if ((e.key === 'Enter' || e.keyCode === 13) && typeof onSubmit === 'function') {
      onSubmit();
    }
  }, [onSubmit]);

  const handleChange = useCallback(
    (e) => {
      if (!maxLength) {
        onChange(e);
        return;
      }

      if (
        !charactersRemaining &&
        e.target.value.length - 1 === `${value}`.length
      ) {
        return;
      }

      if (
        restriction === "alpha" &&
        (e.target.value === "" || Number.isNaN(Number(e.target.value)))
      ) {
        onChange(e);
      }

      if (
        restriction === "numerical" &&
        (e.target.value === "" || !Number.isNaN(Number(e.target.value)))
      ) {
        onChange(e);
      }

      if (restriction === "") {
        onChange(e);
      }
    },
    [onChange, value, charactersRemaining]
  );

  useEffect(() => {
    if (!maxLength || charactersRemaining > 6) return;

    controls.start({
      scale: 1,
      transition: {
        type: "spring",
        velocity: mapRemainingToSpringVelocity(charactersRemaining),
        stiffness: 700,
        damping: 80,
      },
    });
  }, [controls, maxLength, charactersRemaining, `${value}`.length]);

  return (
    <Container
      useLabel={useLabel}
      radius={radius}
      textTransform={textTransform}
      error={error && touched}
      showLabel={!!`${value}`.length}
      useMax={maxLength > 0}
      hasIcon={!!Icon}
    >
      {Icon ? (
        <StartIcon>
          <Icon height={14} />
        </StartIcon>
      ) : null}
      {useLabel ? <label htmlFor={id}>{label || placeholder}</label> : null}
      <input
        spellCheck={false}
        autoComplete="off"
        autoCorrect="off"
        {...{
          id,
          name,
          type,
          placeholder,
          disabled,
          onFocus,
          onBlur,
          autoFocus,
          onKeyPress
        }}
        value={makeValue(`${value}`, masked)}
        onChange={handleChange}
        {...props}
      />
      {maxLength ? (
        <div>
          <motion.span
            animate={controls}
            style={{ color: mapRemainingToColor(charactersRemaining) }}
          >
            {charactersRemaining}
          </motion.span>
        </div>
      ) : null}
    </Container>
  );
};

Input.defaultProps = {
  maxLength: 0,
  Icon: null,
};

const Container = styled.div<{
  useLabel: boolean;
  radius: number;
  error: boolean;
  useMax: boolean;
  hasIcon: boolean;
  textTransform: string;
  showLabel: boolean;
}>`
  position: relative;
  display: flex;
  flex: 1 1 auto;
  max-height: 44px;

  label {
    transition: all 0.2s ease-out;
    font-size: 11px;
    font-weight: 400;
    position: absolute;
    top: 0;
    width: 100%;
    padding: 0 0.9166666667em;
    z-index: 1;
    -webkit-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    user-select: none;
    -webkit-transform: translateY(3px);
    transform: ${({ showLabel }) => (showLabel ? "none" : "translateY(2.5px)")};
    pointer-events: none;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
    -webkit-box-sizing: border-box;
    box-sizing: border-box;
    color: ${({ theme }) => theme.colors.tooltip};
    opacity: ${({ showLabel }) => (showLabel ? 1 : 0)};
    margin: 0.5em 0;
    padding: 0 11px;
    padding-left: ${({ hasIcon }) => (hasIcon ? "33" : "9")}px;
  }

  div {
    color: #ccc;
    background: linear-gradient(
      to right,
      ${({ theme }) => theme.colors.fade}, 0) 0%,
      ${({ theme }) => theme.colors.fade}, 1) 20%
    );
    position: absolute;
    top: 50%;
    right: 4px;
    transform: translateY(-50%);
    padding-right: 8px;
    padding-left: 8px;
  }

  div,
  input {
    border-radius: ${({ radius }) => radius}px;
    font-size: 12px;
    font-family: Inter, sans-serif;
    font-weight: 400;
  }

  div,
  span {
    display: block;
    font-weight: 700;
  }

  input {
    outline: none;
    flex: 1;
    margin: 0;
    transition: all 0.2s ease-out;
    padding: ${({ useLabel }) => (useLabel ? "1.165em 8px" : "10px 0")};
    padding-right: ${({ useMax }) => (useMax ? "48" : "8")}px;
    padding-left: ${({ hasIcon }) => (hasIcon ? "32" : "8")}px;
    color: ${({ theme }) => theme.colors.h2};
    text-transform: ${({ textTransform }) => textTransform};
    background-color: ${({ theme }) => theme.colors.sidebar};
    border: 1px solid
      ${({ theme, error }) =>
        error ? theme.colors.failed : theme.colors.border};

    ${({ showLabel, useLabel }) =>
      showLabel && useLabel
        ? `
      padding-top: 1.845em;
      padding-bottom: .45em;
    `
        : ""}

    &:hover, &:focus {
      border-color: ${({ theme }) => theme.colors.h2};
    }

    &::placeholder {
      color: ${({ theme }) => theme.colors.placeholder};
    }
  }
`;

const StartIcon = styled.span`
  position: absolute;
  top: 52.5%;
  left: 0;
  cursor: default;
  color: ${({ theme }) => theme.colors.h2};
  transform: translateY(-50%);
  padding-right: 8px;
  padding-left: 12px;

  & > svg {
    height: 12px;
    width: 12px;
  }
`;

export default Input;
