import React, {
	Fragment,
	Dispatch,
	SetStateAction,
	useCallback,
	useState
} from 'react';
import { useDispatch } from 'react-redux';
import { FormikProvider, useFormik } from 'formik';
import styled from 'styled-components';
import { useHotkeys } from 'react-hotkeys-hook';

import { Modal } from '../../../../components';

import { addDefault } from '../../../../stores/Main/reducers/defaults';
import { Defaults as DefaultsForm } from '../../../../forms';
import { Buttons } from "../../../../elements";

import Generic from "./generic";
import Profiles from "./profiles";
import Tasks from "./tasks";

import { Steps } from "./types";

interface CreateDefaultProps {
  open: boolean;
  setOpen: Dispatch<SetStateAction<boolean>>;
}

const stepMap = {
  generic: {
    done: ["profiles", "tasks"],
    active: "generic",
  },
  profiles: {
    done: ["tasks"],
    active: "profiles",
  },
  tasks: {
    done: [],
    active: "tasks",
  },
};

const renderFields = (
  step: Steps,
  values: any,
  errors: any,
  touched: any,
  handleChange: any,
  setFieldTouched: any,
  setFieldValue: any
) => {
  switch (step) {
    default:
    case "generic": {
      return (
        <Generic
          values={values}
          errors={errors}
          touched={touched}
          handleChange={handleChange}
          setFieldTouched={setFieldTouched}
          setFieldValue={setFieldValue}
        />
      );
    }

    case "profiles": {
      return (
        <Profiles
          values={values}
          errors={errors}
          touched={touched}
          handleChange={handleChange}
          setFieldTouched={setFieldTouched}
          setFieldValue={setFieldValue}
        />
      );
    }

    case "tasks": {
      return (
        <Tasks
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

const CreateDefault = ({ open, setOpen }: CreateDefaultProps) => {
  const dispatch = useDispatch();
  const [step, setStep] = useState<Steps>("generic");

  const onSubmit = (values: {}) => {
    if (step === "tasks") {
      dispatch(addDefault({ ...values }));

      resetForm({});
      setStep("generic");
    }
  };

  const progressStepper: any = () => {
    return submitForm().then(() => {
      if (isValid) {
        setTouched({});

        if (step === "generic") {
          setStep("profiles");
        }

        if (step === "profiles") {
          setStep("tasks");
        }
      }

      return null;
    });
  };

  const formikbag = useFormik<{}>({
    isInitialValid: false,
    validateOnMount: true,
    onSubmit,
    initialValues: DefaultsForm.initialValues,
    validationSchema: DefaultsForm.validationSchema[step],
  });

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
    resetForm({ values: { ...values, ...DefaultsForm.sections[step] } });
  }, [step, values]);

  const label = step === "tasks" ? "Save" : "Next";

  useHotkeys("command+c,ctrl+c", handleClear, []);
  useHotkeys("enter", progressStepper, []);

  const handleClose = useCallback(() => {
    return setOpen(false);
  }, []);

  return (
    <FormikProvider value={formikbag}>
      <Modal
        height="auto"
        width={515}
        title="Add Default"
        crumbs={{
          step,
          steps: Object.keys(stepMap),
          stepMap,
          matches: false,
          setStep,
          progressStepper,
        }}
        show={open}
        setShow={handleClose}
      >
        <Fragment>
          {renderFields(
            step,
            values,
            errors,
            touched,
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
      </Modal>
    </FormikProvider>
  );
};

const Row = styled.div<{ m: string }>`
  display: flex;
  justify-content: space-between;
  margin: ${({ m }) => m};
`;

export default CreateDefault;
