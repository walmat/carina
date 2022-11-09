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
import { useDropzone } from "react-dropzone";
import styled from "styled-components";
import toast from "react-hot-toast";
import { useHotkeys } from "react-hotkeys-hook";

import TextAreaFormik from "../../../../elements/TextArea/TextAreaFormik";
import SelectFormik from "../../../../elements/Select/SelectFormik";

import {
  Buttons,
  IndicatorSeparator,
  Control,
  Loader,
} from "../../../../elements";
import { Modal } from "../../../../components";

import {
  makeStores,
  Store,
  StoreObject,
} from "../../../../stores/Main/reducers/stores";
import {
  makeAccounts,
  addAccounts,
  AccountGroup,
} from "../../../../stores/Main/reducers/accounts";
import { Account } from "../../../../forms";
import { notify } from "../../../../utils";

import { stores as defaultStores } from "../../../../constants";

interface AddAccountProps {
  group: AccountGroup;
  open: boolean;
  setOpen: Dispatch<SetStateAction<boolean>>;
}

const EOL = /\r?\n/;
let fileReader: any;

interface AddAccountForm {
  group: {
    id: string;
    name: string;
  };
  store: null | Store;
  accounts: string;
}

const extractAccountStores = (stores: StoreObject[]) => {
  const [shopify, footsites, ...rest] = stores;

  return [shopify, footsites, ...rest];
};

const AddAccounts = ({ group, open, setOpen }: AddAccountProps) => {
  const dispatch = useDispatch();
  const groups = useSelector(makeAccounts);
  const stores = useSelector(makeStores);

  const [loading, setLoading] = useState(false);
  const [buffer, setBuffer] = useState<string[]>([]);

  const [accountStores] = useState(
    extractAccountStores(
      stores?.length
        ? stores.filter((store) => store.usesAccounts)
        : (defaultStores as any)
    )
  );

  const onSubmit = (values: AddAccountForm) => {
    if (isValid) {
      const { group, store } = values;
      dispatch(addAccounts({ group, store, accounts: buffer }));
      setOpen(false);
      handleClear();
    }
  };

  const handleFileRead = (id: string) => {
    const { result } = fileReader;
    if (!result) {
      return;
    }

    toast.dismiss(id);

    const accounts = result.split(EOL).filter(Boolean);

    let buffer;
    if (values.accounts) {
      setFieldValue(
        "accounts",
        `${values.accounts}\n`.concat(accounts.join("\n")).trim()
      );
      buffer = `${values.accounts}`.trim().split("\n").concat(accounts);
    } else {
      setFieldValue("accounts", accounts.join("\n").trim());
      buffer = accounts.map((p: string) => p.trim());
    }

    setLoading(false);
    setBuffer(buffer);

    setTimeout(() => {
      notify(
        `${accounts.length} ${
          accounts.length > 1 ? "accounts" : "account"
        } loaded`,
        "success",
        { duration: 1500 }
      );
    }, 25);
  };

  const formikbag = useFormik<AddAccountForm>({
    isInitialValid: false,
    validateOnMount: true,
    onSubmit,
    ...Account.Create,
  });

  useEffect(() => {
    formikbag.setFieldValue("group", { id: group.id, name: group.name });
  }, [group.id]);

  const {
    resetForm,
    values,
    isValid,
    setFieldTouched,
    setFieldValue,
    handleSubmit,
  } = formikbag;

  const onDrop = useCallback(
    (acceptedFiles) => {
      if (!loading && acceptedFiles.length) {
        setLoading(true);
        const start = Date.now();
        let id = notify("Processing file", "loading");

        const timeout = setTimeout(() => {
          if (loading) {
            toast.dismiss(id);
            notify("Failed processing file", "error", { duration: 1500 });
          }
        }, 5000);

        for (const file of acceptedFiles) {
          try {
            handleFilePath(start, id, file);
          } catch (error) {
            // Noop...
          }
        }

        if (timeout) {
          clearTimeout(timeout);
        }
      }
    },
    [values.accounts, toast, loading, setLoading]
  );

  const { getRootProps, getInputProps } = useDropzone({
    onDrop,
    disabled: loading,
    accept: ".txt",
  });

  const handleFilePath = useCallback(
    (start: number, id: string, file: File) => {
      fileReader = new FileReader();
      fileReader.onloadend = () => {
        const end = Date.now();
        return setTimeout(() => {
          return handleFileRead(id);
        }, 1000 - (end - start));
      };
      fileReader.readAsText(file);
    },
    [fileReader]
  );

  const handleClear = useCallback(() => {
    setBuffer([]);
    setLoading(false);
    resetForm({});
  }, []);

  const handleClose = useCallback(() => {
    return setOpen(false);
  }, []);

  const formatAccounts = () => {
    setFieldTouched("accounts", true);

    const sanitized = values.accounts.trim();
    if (!sanitized?.length) {
      setBuffer([]);
      return;
    }

    const accounts = sanitized.split(EOL).filter(Boolean);

    setFieldValue("accounts", sanitized);
    setBuffer(accounts);
  };

  useHotkeys("command+c,ctrl+c", handleClear, []);
  useHotkeys("enter", () => handleSubmit(), []);

  console.log(accountStores, formikbag.values.store);

  return (
    <FormikProvider value={formikbag}>
      <Modal
        height="auto"
        width={525}
        amount={buffer.length}
        title="Add Accounts"
        show={open}
        setShow={handleClose}
      >
        <Fragment>
          <Selects>
            <SelectContainer m="0 8px 0 0">
              <SelectFormik
                required
                name="group"
                label="Account Group"
                placeholder="Account Group"
                components={{ IndicatorSeparator, Control }}
                getOptionLabel={(option: any) => option.name}
                getOptionValue={(option: any) => option.id}
                options={Object.values(groups).map(
                  // eslint-disable-next-line @typescript-eslint/no-unused-vars
                  ({ byId, ...group }) => group
                )}
              />
            </SelectContainer>
            <SelectContainer m="0 0 0 8px">
              <SelectFormik
                required
                autoFocus
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
            <TextAreaFormik
              id="accounts"
              name="accounts"
              placeholder="Paste accounts list"
              onBlur={formatAccounts}
            />
          </TextContainer>

          <Dropzone {...getRootProps()}>
            <input {...getInputProps()} />
            {loading ? (
              <Loader height={16} width={16} />
            ) : (
              <p>Or drag & drop</p>
            )}
          </Dropzone>

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

const Dropzone = styled.div`
  border: 1px dashed ${({ theme }) => theme.colors.primary};
  margin: 16px 0;
  height: 48px;
  border-radius: 4px;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 14px;
  background-color: ${({ theme }) => theme.colors.secondary};
  color: ${({ theme }) => theme.colors.primary};
  cursor: pointer;
`;

export default AddAccounts;
