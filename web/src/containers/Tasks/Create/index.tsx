import React, {
  Dispatch,
  SetStateAction,
  useCallback,
  useState,
  useMemo,
  useEffect,
} from "react";
import { useDispatch, useSelector } from "react-redux";
import styled from "styled-components";
import { motion } from "framer-motion";
import { useFormik, FormikProvider } from "formik";
import { useHotkeys } from "react-hotkeys-hook";

import { makeStores } from "../../../stores/Main/reducers/stores";
import { makeProfiles } from "../../../stores/Main/reducers/profiles";
import { makeProxies } from "../../../stores/Main/reducers/proxies";
import {
  Tasks as TasksType,
  TaskGroup,
  addTasks,
} from "../../../stores/Main/reducers/tasks";

import { sizes, stores as defaultStores } from "../../../constants";
import { modesForPlatform } from "../../../utils/tasks";
import { Modal, Counter } from "../../../components";
import { Buttons, Control } from "../../../elements";
import { Tasks, CreateTaskForm } from "../../../forms";
import SelectFormik from "../../../elements/Select/SelectFormik";
import InputFormik from "../../../elements/Input/InputFormik";

interface CreateTaskProps {
  groups: TasksType;
  group: TaskGroup;
  open: boolean;
  setOpen: Dispatch<SetStateAction<boolean>>;
}

const CreateTask = ({ groups, group, open, setOpen }: CreateTaskProps) => {
  const dispatch = useDispatch();

  const stores = useSelector(makeStores);
  const profiles = useSelector(makeProfiles);
  const proxies = useSelector(makeProxies);

  const [amount, setAmount] = useState<number>(1);

  const onSubmit = (values: CreateTaskForm) => {
    if (isValid) {
      dispatch(addTasks({ ...values, amount }));
    }
  };

  const formikbag = useFormik<CreateTaskForm>({
    isInitialValid: false,
    validateOnMount: false,
    onSubmit,
    initialValues: Tasks.Create.initialValues(group),
    validationSchema: Tasks.Create.validationSchema,
  });

  useEffect(() => {
    formikbag.setFieldValue("groups", [{ id: group.id, name: group.name }]);
  }, [group.id]);

  const { resetForm, values, isValid, setFieldValue, handleSubmit } = formikbag;

  const handleClear = useCallback(() => {
    resetForm({});
    setAmount(1);
  }, []);

  const handleClose = useCallback(() => {
    return setOpen(false);
  }, []);

  useHotkeys("command+c,ctrl+c", handleClear, []);
  useHotkeys("return", () => handleSubmit(), []);

  return useMemo(
    () => (
      <FormikProvider value={formikbag}>
        <Modal
          height="auto"
          width={515}
          title="Add Tasks"
          show={open}
          setShow={handleClose}
        >
          <Col>
            <Row>
              <Col m="0 8px 8px 0">
                <SelectFormik
                  autoFocus
                  required
                  isMulti
                  isClearable
                  label="Task Groups"
                  closeMenuOnSelect={false}
                  name="groups"
                  placeholder="Task Groups"
                  components={{ Control }}
                  getOptionLabel={(option: any) => option.name}
                  getOptionValue={(option: any) => option.id}
                  onChange={(event: any) => {
                    if (!event) {
                      return setFieldValue("groups", []);
                    }

                    return setFieldValue(
                      "groups",
                      // @ts-ignore
                      event.map(({ byId, ...group }) => group)
                    );
                  }}
                  options={Object.values(groups)}
                />
              </Col>
              <Col m="0 0 8px 8px">
                <SelectFormik
                  required
                  isCreatable
                  isClearable
                  name="store"
                  label="Store / Platform"
                  placeholder="Search stores"
                  getOptionLabel={(option: any) => option.name}
                  getOptionValue={(option: any) => option.url}
                  components={{ Control }}
                  onChange={(event: any) => {
                    if (!event) {
                      setFieldValue("mode", null);
                      return setFieldValue("store", null);
                    }

                    setFieldValue("mode", null);

                    return setFieldValue("store", event);
                  }}
                  options={
                    !stores?.length
                      ? []
                      : stores.map((platform) => ({
                          ...platform,
                          options: platform.options.map((site) => ({
                            ...site,
                            type: platform.type,
                          })),
                        }))
                  }
                />
              </Col>
            </Row>
            <Row>
              <Col m="8px 8px 8px 0">
                <SelectFormik
                  required
                  name="mode"
                  placeholder="Mode"
                  label="Task Mode"
                  components={{ Control }}
                  value={values.mode ? values.mode : null}
                  onChange={(event: any) => {
                    if (!event) {
                      return setFieldValue("mode", null);
                    }

                    return setFieldValue("mode", event);
                  }}
                  options={!values.store ? [] : values.store.modes}
                />
              </Col>
              <Col m="8px 0 8px 8px">
                <SelectFormik
                  required
                  isMulti
                  isClearable
                  closeMenuOnSelect={false}
                  name="sizes"
                  placeholder="Sizes"
                  label="Desired Sizes"
                  components={{ Control }}
                  value={values.sizes.map((size: string) => ({
                    label: size,
                    value: size,
                  }))}
                  onChange={(event: any) => {
                    if (!event) {
                      return setFieldValue("sizes", []);
                    }

                    return setFieldValue(
                      "sizes",
                      event.map(({ value }: any) => value)
                    );
                  }}
                  options={sizes}
                />
              </Col>
            </Row>
            <Row m="8px 0">
              <InputFormik
                id="product"
                name="product"
                label="Keywords / Variants / Early Links"
                useLabel
                placeholder="Enter product link, sku, or variants"
              />
            </Row>
            <Row>
              <Col m="8px 8px 8px 0">
                <SelectFormik
                  required
                  isMulti
                  isClearable
                  closeMenuOnSelect={false}
                  name="profiles"
                  placeholder="Profiles"
                  label="Billing Profiles"
                  components={{ Control }}
                  getOptionLabel={(option: any) => option.name}
                  getOptionValue={(option: any) => option.id}
                  onChange={(event: any) => {
                    if (!event) {
                      return setFieldValue("profiles", []);
                    }

                    return setFieldValue("profiles", event);
                  }}
                  options={Object.values(profiles).map(
                    ({ name: label, byId }, index) => ({
                      label,
                      index,
                      options: Object.values(byId).map(({ id, name }) => ({
                        id,
                        group: label,
                        name,
                      })),
                    })
                  )}
                />
              </Col>
              <Col m="8px 0 8px 8px">
                <SelectFormik
                  required
                  isClearable
                  name="proxies"
                  placeholder="Proxy Group"
                  label="Proxy Group"
                  components={{ Control }}
                  getOptionLabel={(option: any) => option.name}
                  getOptionValue={(option: any) => option.id}
                  onChange={(event: any) => {
                    if (!event) {
                      return setFieldValue("proxies", null);
                    }

                    const { byId, ...proxy } = event;

                    return setFieldValue("proxies", proxy);
                  }}
                  options={Object.values(proxies)}
                />
              </Col>
            </Row>
            <Row m="16px 0 0 0">
              <Buttons.Tertiary
                variant="IconButton"
                command="⌘ C"
                text="Clear"
                width={84}
                height={39}
                onClick={handleClear}
              />
              <CounterContainer>
                <Counter count={amount} setCount={setAmount} />
              </CounterContainer>
              <Buttons.Primary
                variant="IconButton"
                command="↩︎"
                text="Create"
                width={84}
                height={39}
                onClick={handleSubmit}
              />
            </Row>
          </Col>
        </Modal>
      </FormikProvider>
    ),
    [open, values, amount]
  );
};

const Col = styled(motion.div)<{ m?: string }>`
  display: flex;
  flex-direction: column;
  flex: 1;
  ${({ m }) => (typeof m !== "undefined" ? `margin: ${m};` : "")}
`;

const Row = styled(motion.div)<{ m?: string }>`
  display: flex;
  justify-content: space-between;
  ${({ m }) => (typeof m !== "undefined" ? `margin: ${m};` : "")}
`;

const CounterContainer = styled.div`
  display: flex;
  margin-right: 16px;
  margin-left: auto;
`;

export default CreateTask;
