import React, {
  Fragment,
  Dispatch,
  SetStateAction,
  useCallback,
  useState,
  useMemo, useEffect,
} from "react";
import { useDispatch } from "react-redux";
import { useFormik, FormikProvider } from "formik";
import styled from "styled-components";
import { useHotkeys } from "react-hotkeys-hook";

import { Modal } from "../../../components";

import { addProfiles, ProfileGroup } from "../../../stores/Main/reducers/profiles";

import { Profiles } from "../../../forms";
import { Buttons, Checkbox, Typography } from "../../../elements";

import Shipping from "./shipping";
import Billing from "./billing";
import Payment from "./payment";

import { Steps } from "./types";

interface CreateProfileProps {
  group: ProfileGroup;
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

const CreateProfile = ({ group, open, setOpen }: CreateProfileProps) => {
  const dispatch = useDispatch();
  const [matches, setMatches] = useState(false);
  const [step, setStep] = useState<Steps>("shipping");

  const handleSubmitForm = () => {
    if (step === "payment") {
      dispatch(addProfiles({ ...values, matches }));

      setMatches(false);
      resetForm({});
      setStep("shipping");
    }
  };

  const progressStepper: any = () => {
    return submitForm().then(() => {
      if (isValid) {
        setTouched({});

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
      }

      return null;
    });
  };

  const formikbag = useFormik<{}>({
    isInitialValid: false,
    validateOnMount: true,
    onSubmit: handleSubmitForm,
    initialValues: Profiles.Create.initialValues,
    validationSchema: Profiles.Create.validationSchema[step],
  });

  useEffect(() => {
    formikbag.setFieldValue('group', { id: group.id, name: group.name })
  }, [group.id]);

  const {
    resetForm,
    touched,
    values,
    errors,
    isValid,
    submitForm,
    handleChange,
    setTouched,
    setFieldTouched,
    setFieldValue,
  } = formikbag;

  const handleClear = useCallback(() => {
    resetForm({ values: { ...values, ...Profiles.Create.sections[step] } });
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

  return (
    <Modal
      height="auto"
      width={515}
      title="Add Profile"
      crumbs={{
        step,
        steps: Object.keys(stepMap),
        stepMap,
        setStep,
        matches,
        progressStepper,
      }}
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

export default CreateProfile;
