import * as React from "react";
import { useCallback, useEffect } from "react";
import { motion, useAnimation, transform } from "framer-motion";
import { Icon } from "react-feather";
import styled from "styled-components";
import NumberFormat from "react-number-format";

const mapRemainingToColor = transform([2, 6], ["#F26E86", "#ccc"]);
const mapRemainingToSpringVelocity = transform([0, 5], [50, 0]);

interface InputProps {
  id: string;
  name: string;
  radius?: number;
  format: string;
  mask?: string[];
  removeFormatting?: any;
  textTransform?: string;
  useLabel?: boolean;
  value: string | number;
  error?: boolean;
  touched?: boolean;
  disabled?: boolean;
  autoFocus?: boolean;
  placeholder: string;
  maxLength?: number;
  Icon?: Icon | null;
  onChange: (e: any) => void;
  onBlur?: (e: any) => void;
  onFocus?: (e: any) => void;
}

export const NumberFormatter = ({
  id,
  name,
  removeFormatting,
  format,
  mask,
  radius = 4,
  textTransform = "none",
  value,
  useLabel = false,
  error = false,
  touched = false,
  disabled = false,
  autoFocus = false,
  placeholder,
  maxLength = 0,
  Icon = null,
  onChange,
  onBlur,
  onFocus,
}: InputProps) => {
  const charactersRemaining = maxLength - `${value}`.length;
  const controls = useAnimation();

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
      onChange(e);
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
      {useLabel ? <label htmlFor={id}>{placeholder}</label> : null}
      <NumberFormatterComponent
        autoComplete="off"
        autoCorrect="off"
        removeFormatting={removeFormatting}
        showLabel={!!`${value}`.length}
        {...{
          id,
          name,
          value,
          radius,
          useLabel,
          disabled,
          format,
          mask,
          placeholder,
          onFocus,
          onBlur,
          autoFocus,
        }}
        onChange={handleChange}
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

NumberFormatter.defaultProps = {
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
    border-radius: ${({ radius }) => radius}px;
    font-size: 12px;
    font-family: Inter, sans-serif;
    font-weight: 400;
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
  span {
    display: block;
    font-weight: 700;
  }
`;

const NumberFormatterComponent = styled(NumberFormat)<{
  useLabel: boolean;
  radius: number;
  error: boolean;
  useMax: boolean;
  hasIcon: boolean;
  textTransform: string;
  showLabel: boolean;
}>`
  outline: none;
  border-radius: ${({ radius }) => radius}px;
  font-size: 12px;
  font-family: Inter, sans-serif;
  font-weight: 400;
  flex: 1;
  margin: 0;
  transition: all 0.2s ease-out;
  padding: ${({ useLabel }) => (useLabel ? "1.165em 8px" : "10px 0")};
  padding-right: ${({ useMax }) => (useMax ? "48" : "8")}px;
  padding-left: ${({ hasIcon }) => (hasIcon ? "32" : "8")}px;
  color: ${({ theme }) => theme.colors.h2};
  background-color: ${({ theme }) => theme.colors.sidebar};
  border: 1px solid
    ${({ theme, error }) => (error ? theme.colors.failed : theme.colors.border)};

  ${({ showLabel, useLabel }) =>
    showLabel && useLabel
      ? `
    padding-top: 1.875em;
    padding-bottom: .45em;
  `
      : ""}

  &:hover, &:focus {
    border-color: ${({ theme }) => theme.colors.h2};
  }

  &::placeholder {
    color: ${({ theme }) => theme.colors.placeholder};
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

export default NumberFormatter;
