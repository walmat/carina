import * as React from "react";
import styled from "styled-components";

interface TextAreaProps {
  id: string;
  name: string;
  value: string;
  autoFocus?: boolean;
  error?: boolean;
  touched?: boolean;
  placeholder: string;
  onChange: (e: any) => void;
  onFocus?: (e: any) => void;
  onBlur?: (e: any) => void;
}

export const TextArea = ({
  id,
  name,
  value,
  autoFocus = false,
  error = false,
  touched = false,
  placeholder,
  onChange,
  onFocus,
  onBlur,
}: TextAreaProps) => {
  return (
    <Container error={error && touched}>
      <textarea
        rows={5}
        wrap="off"
        spellCheck={false}
        autoComplete="off"
        autoCapitalize="off"
        autoCorrect="off"
        {...{ autoFocus, id, value, name, placeholder, onFocus, onChange, onBlur }}
      />
    </Container>
  );
};

TextArea.defaultProps = {
  Icon: null,
};

const Container = styled.div<{ error: boolean }>`
  display: flex;
  flex: 1 1 auto;

  div {
    color: #ccc;
    background: linear-gradient(
      to right,
      ${({ theme }) => theme.colors.fade}, 0) 0%,
      ${({ theme }) => theme.colors.fade}, 1) 20%
    );
    position: absolute;
    top: 50%;
    right: 8px;
    transform: translateY(-50%);
    padding-right: 8px;
    padding-left: 8px;
  }

  div,
  textarea {
    border-radius: 4px;
    font-size: 12px;
    font-family: Inter, sans-serif;
    font-weight: 400;
  }

  textarea {
    outline: none;
    flex: 1;
    padding: 8px;
    margin: 0;
    color: ${({ theme }) => theme.colors.h2};
    background-color: ${({ theme }) => theme.colors.sidebar};
    border: 1px solid
      ${({ theme, error }) =>
        error ? theme.colors.failed : theme.colors.border};

    &::placeholder {
      color: ${({ theme }) => theme.colors.placeholder};
    }

    &:hover,
    &:focus {
      border-color: ${({ theme }) => theme.colors.h2};
    }
  }
`;

export default TextArea;
