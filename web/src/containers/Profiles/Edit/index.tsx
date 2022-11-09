import {
  Fragment,
  Dispatch,
  SetStateAction,
  useCallback,
  useState,
  useMemo,
  useEffect,
} from "react";
import {useDispatch, useSelector} from "react-redux";
import { useFormik, FormikProvider } from "formik";
import styled from "styled-components";
import { useHotkeys } from "react-hotkeys-hook";

import { Modal } from "../../../components";

import { editProfiles, makeProfiles, ProfileGroup} from "../../../stores/Main/reducers/profiles";

import { Profiles } from "../../../forms";
import { Buttons, Checkbox, Typography } from "../../../elements";

import Shipping from "./shipping";
import Billing from "./billing";
import Payment from "./payment";

import { Steps } from "./types";

interface CreateProfileProps {
  group: ProfileGroup;
  ids: string[];
  open: boolean;
  setOpen: Dispatch<SetStateAction<boolean>>;
}

const stepMap = {
  shipping: {
    done: ["billing", "payment"],
    active: "shipping",
  },
  billing: {
    done: ["payment"],
    active: "billing",
  },
  payment: {
    done: [],
    active: "payment",
  },
};

const renderFields = (
  step: Steps,
  values: any,
  errors: any,
  touched: any,
  matches: boolean,
  handleChange: any,
  setFieldTouched: any,
  setFieldValue: any
) => {
  switch (step) {
    default:
    case "shipping": {
      return (
        <Shipping
          values={values}
          errors={errors}
          touched={touched}
          matches={matches}
          handleChange={handleChange}
          setFieldTouched={setFieldTouched}
          setFieldValue={setFieldValue}
        />
      );
    }

    case "billing": {
      return (
        <Billing
          values={values}
          errors={errors}
          touched={touched}
          handleChange={handleChange}
          setFieldTouched={setFieldTouched}
          setFieldValue={setFieldValue}
        />
      );
    }

    case "payment": {
      return (
        <Payment
          values={values}
          errors={errors}
          touched={touched}
          handleChange={handleChange}
          setFieldTouched={setFieldTouched}
          setFieldValue={setFieldValue}
        />
      );
    }
  }
};

const EditProfile = ({ group, ids, open, setOpen }: CreateProfileProps) => {
  const dispatch = useDispatch();

  const profiles = useSelector(makeProfiles);
  const [matches, setMatches] = useState(false);
  const [step, setStep] = useState<Steps>("shipping");

  const handleSubmitForm = () => {
    dispatch(editProfiles({
      ...values,
      group,
      matches,
      ids
    }));

    setOpen(false);
    resetForm({});
  };

  const progressStepper: any = () => {
    if (step === "billing") {
      setStep("payment");
    }

    if (step === "shipping") {
      if (matches) {
        setStep("payment");
      } else {
        setStep("billing");
      }
    }

    if (step === "payment") {
      handleSubmitForm();
    }
  };

  const formikbag = useFormik<{}>({
    isInitialValid: true,
    validateOnMount: false,
    onSubmit: handleSubmitForm,
    initialValues: Profiles.Edit.initialValues,
  });

  const {
    resetForm,
    touched,
    values,
    errors,
    handleChange,
    setTouched,
    setFieldTouched,
    setFieldValue,
  } = formikbag;

  const handleClear = useCallback(() => {
    resetForm({ values: { ...values, ...Profiles.Edit.sections[step] } });
  }, [step, values]);

  const toggleMatches = useCallback(() => {
    setTouched({});
    setMatches((prev) => !prev);
  }, [matches]);

  const label = useMemo(() => step === "payment" ? "Save" : "Next", [step]);

  useHotkeys("command+c,ctrl+c", handleClear, []);
  useHotkeys("return", progressStepper, []);

  const handleClose = useCallback(() => {
    return setOpen(false);
  }, []);

  useEffect(() => {
    const profile = ids.length === 1 ? profiles[group.id].byId[ids[0]] : null;
    if (profile) {
      setFieldValue('shipping', profile.shipping);
      setFieldValue('billing', profile.billing);
      setFieldValue('payment', profile.payment);
      setFieldValue('payment.exp', `${profile.payment.expMonth}/${profile.payment.expYear.slice(2)}`);
      setFieldValue('name', profile.name);
    }
  }, [setFieldValue]);

  return (
    <Modal
      height="auto"
      width={515}
      title="Edit Profiles"
      crumbs={{
        step,
        steps: Object.keys(stepMap),
        stepMap,
        setStep,
        matches,
        progressStepper,
      }}
      amount={ids.length > 1 ? ids.length : ids[0].slice(0, 5)}
      show={open}
      setShow={handleClose}
    >
      <FormikProvider value={formikbag}>
        <Fragment>
          {renderFields(
            step,
            values,
            errors,
            touched,
            matches,
            handleChange,
            setFieldTouched,
            setFieldValue
          )}
          <Row m="16px 0 0 0">
            <Buttons.Tertiary
              variant="IconButton"
              command="⌘ C"
              text="Clear"
              width={84}
              height={39}
              onClick={handleClear}
            />
            {step === "shipping" ? (
              <Matches onClick={toggleMatches}>
                <MatchesText>Same Billing Information</MatchesText>
                <Checkbox isChecked={matches} setIsChecked={() => {}} />
              </Matches>
            ) : null}
            <Buttons.Primary
              variant="IconButton"
              command="↩︎"
              text={label}
              width={84}
              height={39}
              onClick={progressStepper}
            />
          </Row>
        </Fragment>
      </FormikProvider>
    </Modal>
  );
};

const Row = styled.div<{ m: string }>`
  display: flex;
  justify-content: space-between;
  margin: ${({ m }) => m};
`;

const Matches = styled.div`
  display: flex;
  cursor: pointer;
  margin-right: 16px;
  margin-left: auto;
  justify-content: center;
  align-items: center;
`;

const MatchesText = styled(Typography.H2)`
  font-size: 12px;
  margin: 0 8px;
  font-weight: 400;
  display: flex;
  cursor: pointer;
  color: ${({ theme }) => theme.colors.paragraph};
`;

export default EditProfile;
