import React, { Dispatch, SetStateAction, useCallback, useState } from "react";
import { Popover } from "react-tiny-popover";
import styled from "styled-components";
import { motion } from "framer-motion";
import { ExternalLink } from "react-feather";

import { Tooltip } from "../../components";
import { Select, IndicatorSeparator } from "../../elements";
import { Action } from "../../types";

type SortBy = {
  id: string;
  desc: boolean;
};

type Props = {
  actions: Action[];
  setSortBy?: Dispatch<SetStateAction<SortBy>>;
  sortBy?: SortBy;
  setHiddenColumns?: Dispatch<SetStateAction<string[]>>;
  hiddenColumns?: string[];
  groupBy?: string;
  setGroupBy?: Dispatch<SetStateAction<string>>;
};

type ActionProps = {
  type: "Primary" | "First" | "Second" | "Action" | "Button";
  title: string;
  Icon: any;
  popover?: boolean;
  onClick?: any;
  isOpen: boolean;
  setIsOpen: Dispatch<SetStateAction<boolean>>;
  headers?: any;
  setSortBy?: Dispatch<SetStateAction<SortBy>>;
  sortBy?: SortBy;
  setHiddenColumns?: Dispatch<SetStateAction<string[]>>;
  hiddenColumns?: string[];
  groupBy?: string;
  setGroupBy?: Dispatch<SetStateAction<string>>;
  disabled?: boolean;
};

interface Tip {
  children: any;
  disabled?: boolean;
  title: string;
}

interface CustomComponentProps extends React.ComponentPropsWithoutRef<"div"> {}

const CustomComponent = React.forwardRef<HTMLDivElement, CustomComponentProps>(
  (props, ref) => <div ref={ref}>{props.children}</div>
);

const Tip = ({ children, disabled, title }: Tip) => {
  if (disabled) {
    return null;
  }

  return <Tooltip text={title}>{children}</Tooltip>;
};

const ActionComponent = ({
  type,
  onClick,
  title,
  Icon,
  popover,
  isOpen,
  setIsOpen,
  setSortBy = () => {},
  sortBy,
  hiddenColumns,
  setHiddenColumns = () => {},
  groupBy,
  setGroupBy = () => {},
  headers,
  disabled,
}: ActionProps) => {
  let children;

  const toggleOpen = () => setIsOpen((prev) => !prev);

  switch (type) {
    default:
    case "Button":
      children = (
        <Tip title={title} disabled={disabled}>
          <Button
            key={title}
            onClick={onClick}
            whileHover={{ scale: disabled ? 1 : 1.05 }}
            whileTap={{ scale: disabled ? 1 : 0.95 }}
          >
            <IconContainer component={<Icon />} />
          </Button>
        </Tip>
      );
      break;
    case "Action":
      children = (
        <Tip title={title} disabled={disabled}>
          <ActionContainer
            key={title}
            onClick={onClick}
            whileHover={{ scale: disabled ? 1 : 1.05 }}
            whileTap={{ scale: disabled ? 1 : 0.95 }}
          >
            <IconContainer component={<Icon />} />
          </ActionContainer>
        </Tip>
      );
      break;
    case "Primary":
      children = (
        <Tip title={title} disabled={disabled}>
          <Primary
            key={title}
            disabled={disabled}
            onClick={popover ? toggleOpen : onClick}
            whileHover={{ scale: disabled ? 1 : 1.05 }}
            whileTap={{ scale: disabled ? 1 : 0.95 }}
          >
            <IconContainer component={<Icon />} />
          </Primary>
        </Tip>
      );
      break;
    case "First":
      children = (
        <Tip title={title} disabled={disabled}>
          <First
            key={title}
            onClick={onClick}
            whileHover={{ scale: disabled ? 1 : 1.05 }}
            whileTap={{ scale: disabled ? 1 : 0.95 }}
          >
            <IconContainer component={<Icon />} />
          </First>
        </Tip>
      );
      break;
    case "Second":
      children = (
        <Tip title={title} disabled={disabled}>
          <Second
            key={title}
            onClick={onClick}
            whileHover={{ scale: disabled ? 1 : 1.05 }}
            whileTap={{ scale: disabled ? 1 : 0.95 }}
          >
            <IconContainer component={<Icon />} />
          </Second>
        </Tip>
      );
      break;
  }

  if (popover) {
    const handleClick = useCallback(() => {
      setIsOpen((prev) => !prev);
    }, []);

    // TODO: open delay window
    const handleOpenDelays = useCallback(() => {}, []);

    return (
      <Popover
        isOpen={isOpen}
        padding={4}
        reposition={false}
        content={
          <PopoverContainer>
            <Row useBorder>
              <PopoverHeading>Settings</PopoverHeading>
            </Row>
            <Row useBorder>
              <SelectRow>
                <SelectHeading>Ordering</SelectHeading>
                <Select
                  size="small"
                  name="ordering"
                  isClearable={false}
                  placeholder="None"
                  value={{
                    key: sortBy?.id,
                    label: headers?.find((h: any) => h.key === sortBy?.id)
                      ?.label,
                  }}
                  components={{ IndicatorSeparator }}
                  getOptionLabel={(option: any) => option.label}
                  getOptionValue={(option: any) => option.key}
                  onChange={(value: any) => {
                    if (!value) {
                      return setSortBy({ id: "id", desc: false });
                    }
                    return setSortBy({ id: value.key, desc: false });
                  }}
                  options={headers.filter(
                    (h: any) => !hiddenColumns?.includes(h.key)
                  )}
                />
              </SelectRow>
              <SelectRow>
                <SelectHeading>Grouping</SelectHeading>
                <Select
                  size="small"
                  name="grouping"
                  isClearable
                  placeholder="None"
                  value={
                    groupBy
                      ? {
                          key: groupBy,
                          label: headers?.find((h: any) => h.key === groupBy)
                            ?.label,
                        }
                      : null
                  }
                  components={{ IndicatorSeparator }}
                  getOptionLabel={(option: any) => option.label}
                  getOptionValue={(option: any) => option.key}
                  onChange={(value: any) => {
                    setGroupBy(value.key);
                  }}
                  options={headers.filter(
                    (h: any) => !hiddenColumns?.includes(h.key)
                  )}
                />
              </SelectRow>
            </Row>
            <Row useBorder>
              <SectionHeading m="0 0 12px 0">Table Columns</SectionHeading>
              <ButtonsContainer>
                {headers.map((h: any) => (
                  <ToggleButton
                    isSelected={!hiddenColumns?.includes(h.key)}
                    onClick={() => {
                      setHiddenColumns((prev) => {
                        if (prev?.includes(h.key)) {
                          return prev.filter((c: any) => c !== h.key);
                        } else {
                          if (prev.length === 5) {
                            return prev;
                          }
                          return [...prev, h.key];
                        }
                      });
                    }}
                  >
                    {h.label}
                  </ToggleButton>
                ))}
              </ButtonsContainer>
            </Row>
            <Row useCursor useBorder={false} onClick={handleOpenDelays}>
              <SectionHeading m="0 0 0 0">
                Delays & Scheduling
                <ExternalLink />
              </SectionHeading>
            </Row>
          </PopoverContainer>
        }
        positions={["bottom"]}
        align="end"
        onClickOutside={handleClick}
      >
        <CustomComponent>{children}</CustomComponent>
      </Popover>
    );
  }

  return children;
};

const Actions = ({
  actions,
  setSortBy,
  sortBy,
  hiddenColumns,
  setHiddenColumns,
  groupBy,
  setGroupBy,
}: Props) => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <Container>
      {actions.map((action) =>
        ActionComponent({
          ...action,
          isOpen,
          setIsOpen,
          setSortBy,
          sortBy,
          hiddenColumns,
          setHiddenColumns,
          groupBy,
          setGroupBy,
        })
      )}
    </Container>
  );
};

const Container = styled.div`
  display: flex;
  margin-right: -8px;
  cursor: pointer;
  margin-left: auto;
  align-items: center;
`;

const PopoverContainer = styled.div`
  height: fit-content;
  width: 200px;
  border: 1px solid ${({ theme }) => theme.colors.border};
  box-shadow: rgba(0, 0, 0, 0.15) 0px 5px 15px 0px;
  border-radius: 4px;
  display: flex;
  flex-direction: column;
  background-color: ${({ theme }) => theme.colors.sidebar};
  z-index: 100;
`;

const Row = styled.div<{ useCursor?: boolean; useBorder?: boolean }>`
  font-weight: 400;
  padding: 10px 12px;
  color: ${({ theme }) => theme.colors.h2};
  ${({ theme, useBorder }) =>
    useBorder ? `border-bottom: 1px solid ${theme.colors.separator};` : ""}
  ${({ useCursor }) =>
    useCursor
      ? `
    cursor: pointer;
    & > * {
      cursor: pointer;
    }
  `
      : ""}
`;

const PopoverHeading = styled.div`
  font-weight: 500;
  font-size: 12px;
  line-height: 12px;
  color: ${({ theme }) => theme.colors.subHeading};
`;

const SelectHeading = styled.span`
  font-weight: 400;
  font-size: 12px;
  line-height: 12px;
  color: ${({ theme }) => theme.colors.paragraph};
`;

const SelectRow = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 400;
  font-size: 12px;
  color: ${({ theme }) => theme.colors.muted};

  &:first-child {
    margin-bottom: 8px;
  }

  & .select__control {
    width: 70px;
  }
`;

const ButtonsContainer = styled.div`
  display: flex;
  align-items: center;
  flex-wrap: wrap;
`;

const ToggleButton = styled(motion.div)<{ isSelected?: boolean }>`
  height: 18px;
  padding: 2px 8px;
  display: flex;
  border-radius: 4px;
  font-size: 10px;
  line-height: 12px;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  background-color: ${({ theme, isSelected }) =>
    isSelected ? theme.colors.primary : "transparent"};
  color: ${({ theme, isSelected }) =>
    isSelected ? "#FFFFFF" : theme.colors.primary};
  border-width: 1px;
  border-style: solid;
  border-color: ${({ theme }) => theme.colors.primary};
  margin-right: 8px;
  margin-bottom: 6px;
`;

const SectionHeading = styled.div<{ m: string }>`
  font-weight: 400;
  font-size: 12px;
  line-height: 12px;
  color: ${({ theme }) => theme.colors.subHeading};
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: ${({ m }) => m};

  & > svg {
    width: 14px;
    height: 14px;
  }
`;

const Button = styled(motion.div)`
  width: 165px;
  height: 32px;
  margin: 0 8px;
  color: #ffffff;
  display: flex;
  border-radius: 4px;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  background-color: ${({ theme }) => theme.colors.primary};
`;

const Primary = styled(motion.div)<{ disabled?: boolean }>`
  width: 32px;
  height: 32px;
  border-radius: 4px;
  margin: 0 8px;
  color: #fff;
  opacity: ${({ disabled }) => (disabled ? 0.4 : 1)};
  color: ${({ theme }) => theme.colors.paragraph};
  background-color: ${({ theme }) => theme.colors.view};
  display: flex;
  justify-content: center;
  align-items: center;

  cursor: ${({ disabled }) => (disabled ? "default" : "pointer")};

  & > svg,
  line,
  path,
  polyline {
    cursor: ${({ disabled }) => (disabled ? "default" : "pointer")};
  }
`;

const First = styled(motion.div)`
  width: 31px;
  height: 31px;
  margin: 0 0 0 8px;
  cursor: pointer;
  display: flex;
  justify-content: center;
  align-items: center;
  color: ${({ theme }) => theme.colors.paragraph};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 4px;
  border-right-width: 0.5px;
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
`;

const Second = styled(motion.div)`
  width: 31px;
  height: 31px;
  margin: 0 8px 0 0;
  cursor: pointer;
  display: flex;
  justify-content: center;
  align-items: center;
  color: ${({ theme }) => theme.colors.paragraph};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 4px;
  border-left-width: 0.5px;
  border-top-left-radius: 0;
  border-bottom-left-radius: 0;

  & > a {
    font-size: inherit;
    color: inherit;
    width: 14px;
    display: flex;
    justify-content: center;
    align-items: center;
  }
`;

const ActionContainer = styled(motion.div)`
  width: 32px;
  height: 32px;
  display: flex;
  border-radius: 4px;
  margin: 0 8px;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  color: #fff;
  background-color: ${({ theme }) => theme.colors.primary};
`;

const IconContainer = styled(({ component, ...props }) =>
  React.cloneElement(component, props)
)`
  display: flex;
  justify-content: center;
  align-items: center;
  margin: auto 0;
  width: 14px;
  height: 14px;
`;

export default Actions;
