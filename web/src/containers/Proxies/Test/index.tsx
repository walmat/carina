import React, { Dispatch, SetStateAction, useCallback, useMemo } from "react";
import {useDispatch} from "react-redux";
import { useFormik, FormikProvider } from "formik";
import { motion } from "framer-motion";
import styled from "styled-components";
import { useHotkeys } from "react-hotkeys-hook";

import { Buttons, Checkbox, Typography } from "../../../elements";
import { Modal } from "../../../components";
import { testProxies } from "../../../stores/Main/reducers/proxies";

type Props = {
  open: boolean;
  setOpen: Dispatch<SetStateAction<boolean>>;
  selected: string[];
  groupId: string;
};

type Platform = {
  label: string;
  selected: boolean;
  value: string;
};

type Form = {
  platforms: Platform[];
};

const platforms = [
  {
    label: "All Stores",
    value: "All",
  },
  {
    label: "Shopify",
    value: "https://kith.com",
  },
  {
    label: "Supreme",
    value: "https://www.supremenewyork.com",
  },
  {
    label: "Yeezy Supply",
    value: "https://www.yeezysupply.com",
  },
  {
    label: "Footsites",
    value: "https://www.footlocker.com",
  },
  {
    label: "Walmart",
    value: "https://www.walmart.com",
  },
];

const TestProxies = ({ open, setOpen, selected, groupId }: Props) => {
  const dispatch = useDispatch();

  const onSubmit = () => {
    if (isValid) {
      if (values.platforms.some(p => p.value === 'All')) {
        const [, ...plats] = platforms;
        dispatch(testProxies({ ids: selected, platforms: plats, groupId }));
      } else {
        dispatch(testProxies({ ids: selected, platforms: values.platforms, groupId }));
      }

      handleClear();
      handleClose();
    }
  };

  const initialValues = useMemo(() => ({
    platforms: platforms.map(p => ({ ...p, selected: false })),
  }), []);

  const formikbag = useFormik<Form>({
    isInitialValid: true,
    validateOnMount: true,
    onSubmit,
    initialValues,
  });

  const { resetForm, values, isValid, setFieldValue, handleSubmit } = formikbag;

  const handleClear = useCallback(() => {
    resetForm({});
  }, []);

  const handleClose = useCallback(() => {
    return setOpen(false);
  }, []);

  const handleSelect = useCallback(
    (index: number, state: boolean) => {
      const numSelected = values.platforms.filter(
        (platform) => platform.selected
      ).length;

      if (numSelected === 1) {
        const current = values.platforms.findIndex(
          (platform) => platform.selected
        );
        if (current === index) {
          return;
        }
      }

      // We're possibly selecting all options
      if (numSelected === values.platforms.length - 2) {
        const lastOption = values.platforms.findIndex((platform, at) => {
          if (at === 0) {
            return false;
          }

          if (!platform.selected) {
            return true;
          }
        });

        // let's unselect all options and just tick the `all stores` option
        if (lastOption === index) {
          setFieldValue(`platforms.${0}.selected`, true);
          values.platforms.forEach((_, index) => {
            if (index > 0) {
              setFieldValue(`platforms.${index}.selected`, false);
            }
          });

          return;
        }
      }

      if (index === 0) {
        const isSelected = values.platforms[0].selected;
        // NOTE: uncheck all other boxes
        if (!isSelected) {
          setFieldValue(`platforms.${0}.selected`, true);
          values.platforms.forEach((_, index) => {
            if (index > 0) {
              setFieldValue(`platforms.${index}.selected`, false);
            }
          });
        }
        return;
      }

      if (values.platforms[0].selected) {
        setFieldValue(`platforms.${0}.selected`, false);
      }

      setFieldValue(`platforms.${index}.selected`, !state);
    },
    [values]
  );

  useHotkeys("command+c,ctrl+c", handleClear, []);
  useHotkeys("enter", () => handleSubmit(), []);

  return (
    <FormikProvider value={formikbag}>
      <Modal
        height={384}
        width={325}
        title="Choose Stores"
        show={open}
        setShow={handleClose}
      >
        <Container>
          {values.platforms.map((platform, index) => (
            <CheckboxContainer
              onClick={() => handleSelect(index, platform.selected)}
            >
              <Checkbox isChecked={platform.selected} setIsChecked={() => {}} />
              <TextContainer>
                <Label selected={platform.selected}>{platform.label}</Label>
              </TextContainer>
            </CheckboxContainer>
          ))}

          <Row m="16px 0 0 0">
            <Buttons.Tertiary
              variant="IconButton"
              command="⌘ C"
              text="Clear"
              width={84}
              height={39}
              onClick={handleClear}
            />
            <Buttons.Primary
              variant="IconButton"
              command="↩︎"
              text="Test"
              width={84}
              height={39}
              onClick={handleSubmit}
            />
          </Row>
        </Container>
      </Modal>
    </FormikProvider>
  );
};

const Container = styled.div`
  display: flex;
  flex-direction: column;
`;

const Row = styled.div<{ m: string }>`
  display: flex;
  justify-content: space-between;
  margin: ${({ m }) => m};
  cursor: pointer;
`;

const CheckboxContainer = styled(motion.div)`
  display: flex;
  cursor: pointer;
  padding: 4px 12px;
  border-radius: 4px;
  margin-bottom: 8px;

  &:hover {
    background-color: ${({ theme }) => theme.colors.secondary};
  }

  & > * {
    cursor: pointer;
  }
`;

const TextContainer = styled.div`
  display: flex;
  flex: 1;
  height: 32px;
  margin: auto 0;
  border-radius: 4px;
  cursor: pointer;
`;

const Label = styled(Typography.H2)<{ selected: boolean }>`
  font-size: 12px;
  font-weight: 500;
  color: ${({ theme }) => theme.colors.h2};
  cursor: pointer;
  margin: 0 8px;
  display: flex;
  align-items: center;
`;

export default TestProxies;
