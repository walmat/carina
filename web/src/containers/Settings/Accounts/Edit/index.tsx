import React, {
  Fragment,
  Dispatch,
  SetStateAction,
  useCallback,
  useState,
  useEffect,
} from "react";
import { useDispatch, useSelector } from "react-redux";
import { useFormik, FormikProvider } from "formik";
import styled from "styled-components";
import { useHotkeys } from "react-hotkeys-hook";

import SelectFormik from "../../../../elements/Select/SelectFormik";

import {
  Buttons,
  IndicatorSeparator,
  Control,
  InputFormik,
} from "../../../../elements";
import { Modal } from "../../../../components";

import {
  makeStores,
  Store,
  StoreObject,
} from "../../../../stores/Main/reducers/stores";
import {
  editAccounts,
  AccountGroup,
} from "../../../../stores/Main/reducers/accounts";
import { Account } from "../../../../forms";

import { stores as defaultStores } from "../../../../constants";

interface AddAccountProps {
  group: AccountGroup;
  ids: string[];
  open: boolean;
  setOpen: Dispatch<SetStateAction<boolean>>;
}

interface AddAccountForm {
  store: null | Store;
  account: string;
}

const extractAccountStores = (stores: StoreObject[]) => {
  const [shopify, footsites, ...rest] = stores;

  return [shopify, footsites, ...rest];
};

const AddAccounts = ({ group, ids, open, setOpen }: AddAccountProps) => {
  const dispatch = useDispatch();
  const stores = useSelector(makeStores);

  const [accountStores] = useState(
    extractAccountStores(
      stores?.length
        ? stores.filter((store) => store.usesAccounts)
        : (defaultStores as any)
    )
  );

  const onSubmit = (values: AddAccountForm) => {
    if (isValid) {
      const { store, account } = values;
      dispatch(editAccounts({ group, id: ids[0], store, account }));
      handleClear();
      setOpen(false);
    }
  };

  const formikbag = useFormik<AddAccountForm>({
    isInitialValid: false,
    validateOnMount: true,
    onSubmit,
    ...Account.Edit,
  });

  useEffect(() => {
    formikbag.setFieldValue("group", { id: group.id, name: group.name });

    if (ids.length) {
      const acc = group.byId[ids[0]];

      formikbag.setFieldValue("store", acc.store);
      formikbag.setFieldValue("account", `${acc.username}:${acc.password}`);
    }
  }, [group.id]);

  const { resetForm, isValid, handleSubmit } = formikbag;

  const handleClear = useCallback(() => {
    resetForm({});
  }, []);

  const handleClose = useCallback(() => {
    return setOpen(false);
  }, []);

  useHotkeys("command+c,ctrl+c", handleClear, []);
  useHotkeys("enter", () => handleSubmit(), []);

  console.log(accountStores, formikbag.values.store);

  return (
    <FormikProvider value={formikbag}>
      <Modal
        height="auto"
        width={325}
        amount={ids.length ? ids[0].slice(0, 5) : undefined}
        title="Edit Account"
        show={open}
        setShow={handleClose}
      >
        <Fragment>
          <Selects>
            <SelectContainer m="0">
              <SelectFormik
                required
                isClearable
                name="store"
                label="Account Store"
                placeholder="Account Store"
                components={{ IndicatorSeparator, Control }}
                getOptionLabel={(option: any) => option.name}
                getOptionValue={(option: any) => option.url}
                options={accountStores}
              />
            </SelectContainer>
          </Selects>

          <TextContainer>
            <InputFormik
              id="account"
              name="account"
              label="Username / Password"
              placeholder="Username:Password"
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
              text="Create"
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

const Selects = styled.div`
  display: flex;
`;

const SelectContainer = styled.div<{ m: string }>`
  display: flex;
  flex: 1;
  margin: ${({ m }) => m};
  flex-direction: column;
`;

const TextContainer = styled.div`
  display: flex;
  flex: 1;
  margin-top: 16px;
`;

export default AddAccounts;
