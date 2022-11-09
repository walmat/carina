import React, {
  Fragment,
  Dispatch,
  SetStateAction,
  useCallback,
  useEffect,
} from "react";
import { useDispatch } from "react-redux";
import { useFormik, FormikProvider } from "formik";
import styled from "styled-components";
import { useHotkeys } from "react-hotkeys-hook";

import { Buttons, InputFormik } from "../../../elements";
import { Modal } from "../../../components";

import { editProxies, ProxyGroup } from '../../../stores/Main/reducers/proxies';
import { Proxies } from '../../../forms';

type Props = {
  group: ProxyGroup;
  ids: string[];
  open: boolean;
  setOpen: Dispatch<SetStateAction<boolean>>;
};

type Form = {
  group:  {
    id: string;
    name: string;
  };
  proxy: string;
};

const EditProxies = ({ group, ids, open, setOpen }: Props) => {
  const dispatch = useDispatch();

  const onSubmit = (values: Form) => {
    if (isValid) {
      const { group, proxy } = values;
      dispatch(editProxies({ group, id: ids[0], proxy }));
      handleClear();
      setOpen(false);
    }
  };

  const formikbag = useFormik<Form>({
    isInitialValid: false,
    validateOnMount: true,
    onSubmit,
    ...Proxies.Edit,
  });

  const {
    resetForm,
    isValid,
    setFieldValue,
    handleSubmit,
  } = formikbag;

  const handleClear = useCallback(() => {
    resetForm({});
  }, []);

  const handleClose = useCallback(() => {
    return setOpen(false);
  }, []);

  useHotkeys("command+c,ctrl+c", handleClear, []);
  useHotkeys("enter", () => handleSubmit(), []);

  useEffect(() => {
    if (ids.length) {
      const { host, port, username, password } = group.byId[ids[0]];
      
      if (username && password) {
        setFieldValue('proxy', `${host}:${port}:${username}:${password}`);
      } else {
        setFieldValue('proxy', `${host}:${port}`);
      }
    }
  }, [ids]);

  return (
    <FormikProvider value={formikbag}>
      <Modal
        height="auto"
        width={325}
        amount={ids.length ? ids[0].slice(0, 5) : undefined}
        title="Editing Proxy"
        show={open}
        setShow={handleClose}
      >
        <Fragment>
          <TextContainer>
            <InputFormik
              autoFocus
              id="proxy"
              name="proxy"
              label="Proxy address"
              placeholder="Proxy address"
            />
          </TextContainer>

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
              text="Save"
              width={84}
              height={39}
              onClick={handleSubmit}
            />
          </Row>
        </Fragment>
      </Modal>
    </FormikProvider>
  );
};

const Row = styled.div<{ m: string }>`
  display: flex;
  justify-content: space-between;
  margin: ${({ m }) => m};
`;

const TextContainer = styled.div`
  display: flex;
  flex: 1;
  margin-top: 16px;
  width: 325px;
  min-width: 325px;
  max-width: 325px;
`;

export default EditProxies;
