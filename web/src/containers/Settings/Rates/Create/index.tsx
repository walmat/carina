import React, {
  Fragment,
  Dispatch,
  SetStateAction,
  useCallback,
  useState,
} from "react";
import { useDispatch, useSelector } from "react-redux";
import { useFormik, FormikProvider } from "formik";
import styled from "styled-components";
import { Link } from "react-feather";
import toast from "react-hot-toast";
import { useHotkeys } from "react-hotkeys-hook";

import { Buttons, IndicatorSeparator, Control } from '../../../../elements';
import { Modal } from '../../../../components';

import { addRates } from '../../../../stores/Main/reducers/rates';
import { makeStores, Store } from '../../../../stores/Main/reducers/stores';
import { makeProfiles, Profile } from "../../../../stores/Main/reducers/profiles";
import { Rates } from '../../../../forms';
import { stores as defaultStores } from '../../../../constants';
import SelectFormik from '../../../../elements/Select/SelectFormik';
import InputFormik from '../../../../elements/Input/InputFormik';

type Props = {
  open: boolean;
  setOpen: Dispatch<SetStateAction<boolean>>;
};

type Form = {
  store: null | Store;
  profile: null | Profile;
  value: string;
};

const getPlaceholder = (value: string) => {
  if (!value.startsWith("http")) {
    return "Save";
  }
  return "Fetch";
};

const AddRates = ({ open, setOpen }: Props) => {
  const dispatch = useDispatch();
  const profiles = useSelector(makeProfiles);
  const stores = useSelector(makeStores);

  const [loading, setLoading] = useState(false);

  const onSubmit = (values: Form) => {
    if (isValid) {
      const { store, profile, value } = values;

      // NOTE: Direct rate provided
      if (/^.+-.+-.+/i.test(value) && !/^http/i.test(value)) {
        dispatch(addRates({ store, profile, rates: [value] }));
        return;
      }

      // TODO: If given a direct rate, let's just add it right away
      // TODO: Fetch rates w/ golang backend
      const id = notify(null, "Retrieving rates", toast.loading);
      handleFetch(id);
    }
  };

  const notify = (
    id: string | null,
    message: string,
    notifier = toast.success
  ) => {
    if (!id) {
      return notifier(message);
    }

    return notifier(message, { id });
  };

  const handleFetch = async (id: string) => {
    const { store, profile, value } = values;
    try {
      setLoading(true);

      notify(id, "Fetching rates", toast.loading);

      // @ts-ignore
      const rateInfo = await window.RPCAction("rates:fetch", ['shop', value, profile.id]);
      console.log(rateInfo);

      notify(id, `Found ${rateInfo.rates.length} rates`, toast.success);
      dispatch(addRates({ store: { name: store }, profile, rates: rateInfo.rates }));

      handleClear();
    } catch (e) {
      notify(id, `Failed fetching rates`, toast.error);
      console.log(e);
    }
  };

  const formikbag = useFormik<Form>({
    isInitialValid: false,
    validateOnMount: true,
    onSubmit: onSubmit,
    ...Rates,
  });

  const {
    resetForm,
    touched,
    values,
    errors,
    isValid,
    setFieldTouched,
    setFieldValue,
    handleSubmit,
  } = formikbag;

  const handleClear = useCallback(() => {
    setLoading(false);
    resetForm({});
  }, []);

  const handleClose = useCallback(() => {
    return setOpen(false);
  }, []);

  useHotkeys("command+c,ctrl+c", handleClear, []);
  useHotkeys("enter", () => handleSubmit(), []);

  return (
    <FormikProvider value={formikbag}>
      <Modal
        height="auto"
        width={525}
        title="Add Rates"
        show={open}
        setShow={handleClose}
      >
        <Fragment>
          <Selects>
            <SelectContainer m="0 8px 0 0">
              <SelectFormik
                autoFocus
                name="store"
                required
                label="Store / Platform"
                placeholder="Choose Store"
                components={{ IndicatorSeparator, Control }}
                getOptionLabel={(option: any) => option.name}
                getOptionValue={(option: any) => option.url}
                options={
                  stores.length
                    ? stores.filter((store) => store.usesRates)
                    : defaultStores
                }
              />
            </SelectContainer>
            <SelectContainer m="0 0 0 8px">
              <SelectFormik
                name="profile"
                required
                label="Billing Profile"
                placeholder="Choose Profile"
                components={{ IndicatorSeparator, Control }}
                getOptionLabel={(option: any) => option.name}
                getOptionValue={(option: any) => option.id}
                onChange={(event: any) => {
                  if (!event) {
                    return setFieldValue("profile", []);
                  }

                  return setFieldValue("profile", event);
                }}
                options={Object.values(profiles).map(
                  ({ name: label, byId }, index) => ({
                    label,
                    index,
                    options: Object.values(byId).map(({ id, name }) => ({
                      id,
                      name,
                    })),
                  })
                )}
              />
            </SelectContainer>
          </Selects>

          <TextContainer>
            <InputFormik
              id="value"
              name="value"
              useLabel
              label="Product / Rate"
              placeholder="Paste product information or pre-determined rate"
              disabled={loading}
              Icon={Link}
            />
          </TextContainer>

          <Row m="16px 0 0 0">
            <Buttons.Tertiary
              variant="IconButton"
              text="Clear"
              command="⌘ C"
              width={84}
              height={39}
              onClick={handleClear}
            />
            <Buttons.Primary
              variant="IconButton"
              disabled={loading}
              loading={loading}
              text={getPlaceholder(values.value)}
              command="↩︎"
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
  margin-top: 16px;
`;

export default AddRates;
