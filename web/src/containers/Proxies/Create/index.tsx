import React, {
  Fragment,
  Dispatch,
  SetStateAction,
  useCallback,
  useState, useEffect,
} from "react";
import { useDispatch, useSelector } from "react-redux";
import { useFormik, FormikProvider } from "formik";
import { useDropzone } from "react-dropzone";
import styled from "styled-components";
import toast from "react-hot-toast";
import { useHotkeys } from "react-hotkeys-hook";

import {
  Buttons,
  IndicatorSeparator,
  Control,
  Loader,
} from "../../../elements";
import { Modal } from "../../../components";

import SelectFormik from '../../../elements/Select/SelectFormik';
import TextAreaFormik from '../../../elements/TextArea/TextAreaFormik';

import { makeProxies, addProxies, ProxyGroup } from '../../../stores/Main/reducers/proxies';
import { Proxies } from '../../../forms';
import { notify } from "../../../utils";

type Props = {
  group: ProxyGroup;
  open: boolean;
  setOpen: Dispatch<SetStateAction<boolean>>;
};

const EOL = /\r?\n/;
let fileReader: any;

type Form = {
  group:  {
    id: string;
    name: string;
  };
  proxies: string;
};

const AddProxies = ({ group, open, setOpen }: Props) => {
  const dispatch = useDispatch();
  const groups = useSelector(makeProxies);

  const [loading, setLoading] = useState(false);
  const [buffer, setBuffer] = useState<string[]>([]);

  const onSubmit = (values: Form) => {
    if (isValid) {
      const { group } = values;
      dispatch(addProxies({ group, proxies: buffer }));
      setBuffer([]);
      handleClear();
    }
  };

  const handleFileRead = (id: string) => {
    const { result } = fileReader;
    if (!result) {
      return;
    }

    toast.dismiss(id);

    const proxies = result.split(EOL).filter(Boolean);

    let buffer;
    if (values.proxies) {
      setFieldValue(
        "proxies",
        `${values.proxies}\n`.concat(proxies.join("\n")).trim()
      );
      buffer = `${values.proxies}`.trim().split("\n").concat(proxies);
    } else {
      setFieldValue(
        "proxies",
        proxies.join("\n").trim()
      );
      buffer = proxies.map((p: string) => p.trim());
    }

    setLoading(false);
    setBuffer(buffer);

    setTimeout(() => {
      notify(`${proxies.length} proxies loaded`, "success", { duration: 1500 });
    }, 25);
  };

  const formikbag = useFormik<Form>({
    isInitialValid: false,
    validateOnMount: true,
    onSubmit,
    ...Proxies.Create,
  });

  useEffect(() => {
    formikbag.setFieldValue('group', { id: group.id, name: group.name })
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
            notify( "Failed processing file", "error", { duration: 1500 });
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
    [values.proxies, toast, loading, setLoading]
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
    [values.proxies, fileReader]
  );

  const handleClear = useCallback(() => {
    setBuffer([]);
    setLoading(false);
    resetForm({});
  }, [buffer]);

  const handleClose = useCallback(() => {
    return setOpen(false);
  }, []);

  const formatProxies = () => {
    setFieldTouched("proxies", true);

    const sanitized = values.proxies.trim();
    if (!sanitized?.length) {
      setBuffer([]);
      return;
    }

    const proxies = sanitized.split(EOL).filter(Boolean);

    setFieldValue("proxies", sanitized);
    setBuffer(proxies);
  };

  useHotkeys("command+c,ctrl+c", handleClear, []);
  useHotkeys("enter", () => handleSubmit(), []);

  return (
    <FormikProvider value={formikbag}>
      <Modal
        height="auto"
        width={325}
        amount={buffer.length}
        title="Add Proxies"
        show={open}
        setShow={handleClose}
      >
        <Fragment>
          <SelectContainer>
            <SelectFormik
              name="group"
              required
              label="Proxy Group"
              placeholder="Proxy Group"
              components={{ IndicatorSeparator, Control }}
              getOptionLabel={(option: any) => option.name}
              getOptionValue={(option: any) => option.id}
              options={Object.values(groups).map(
                // eslint-disable-next-line @typescript-eslint/no-unused-vars
                ({ byId, ...group }) => group
              )}
            />
          </SelectContainer>

          <TextContainer>
            <TextAreaFormik
              autoFocus
              id="proxies"
              name="proxies"
              placeholder="Paste proxy list"
              onBlur={formatProxies}
            />
          </TextContainer>

          <Dropzone {...getRootProps()}>
            <input {...getInputProps()} />
            {loading ? <Loader height={16} width={16} /> : <p>Or drag & drop</p>}
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

const SelectContainer = styled.div`
  width: 325px;
  min-width: 325px;
  max-width: 325px;
`;

const TextContainer = styled.div`
  display: flex;
  flex: 1;
  margin-top: 16px;
  width: 325px;
  min-width: 325px;
  max-width: 325px;
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

export default AddProxies;
