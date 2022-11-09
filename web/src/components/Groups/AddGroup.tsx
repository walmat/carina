import React, { Fragment } from "react";
import { useDispatch } from "react-redux";
import { isEmpty } from "lodash";
import styled from "styled-components";
import { useFormik, FormikProvider } from "formik";

import { Buttons } from "../../elements";
import { Create } from "../../forms";
import { useKeyPress } from "../../hooks";
import InputFormik from "../../elements/Input/InputFormik";
import { FormikHelpers } from "formik/dist/types";

export interface AddGroupProps {
  groups: GroupMap;
  create: any;
}

export interface AddGroupForm {
  name: string;
}

export type GroupMetadata = {
  id: string;
  name: string;
};

/*
	This maps a UUID (generated by the Go backend)
	into an object containing the group's metadata
	as well as the items in the container.

	e.g: Task groups, Proxy groups.

	The group items are usually referenced by the byId
	accessor.
*/
export type GroupMap = {
  [id: string]: GroupMetadata;
};

const AddGroup = ({ groups, create }: AddGroupProps) => {
  const dispatch = useDispatch();

  const onSubmit = async (
    values: AddGroupForm,
    formikHelpers: FormikHelpers<AddGroupForm>
  ) => {
    const { setFieldError, resetForm, validateForm } = formikHelpers;

    const errors = await validateForm(values);

    if (!isEmpty(errors)) {
      return Object.keys(errors).map((key) =>
        setFieldError(key, `${key} is required`)
      );
    }

    const { name } = values;

    if (
      Object.values(groups).find((group: GroupMetadata) => group.name === name)
    ) {
      return [setFieldError("name", `Group name is already in use`)];
    }

    dispatch(create(name.trim()));
    resetForm({});
  };

  const formikbag = useFormik<AddGroupForm>({
    isInitialValid: false,
    validateOnMount: false,
    onSubmit,
    ...Create,
  });
  const { setFieldTouched, handleSubmit } = formikbag;

  useKeyPress("Enter", handleSubmit, []);

  return (
    <FormikProvider value={formikbag}>
      <Fragment key="Add-Group">
        <Row bottom={16} top={0}>
          <InputFormik
            autoFocus
            id="name"
            name="name"
            useLabel
            maxLength={25}
            onFocus={() => setFieldTouched("name", true)}
            placeholder="Group name"
          />
        </Row>

        <RowEnd>
          <Buttons.Primary
            variant="IconButton"
            command="↩︎"
            text="Create"
            width={94}
            height={36}
            onClick={handleSubmit}
          />
        </RowEnd>
      </Fragment>
    </FormikProvider>
  );
};

const Row = styled.div<{ top: number; bottom: number }>`
  display: flex;
  flex: 1;
  margin-top: ${({ top }) => `${top}`}px;
  margin-bottom: ${({ bottom }) => `${bottom}`}px;
`;

const RowEnd = styled.div`
  display: flex;
  flex: 1;
  align-items: flex-end;
  justify-content: flex-end;
`;

export default AddGroup;
