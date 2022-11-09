import React, { Fragment } from "react";
import { useDispatch } from "react-redux";
import { isEmpty } from "lodash";
import styled from "styled-components";
import { useFormik, FormikProvider } from "formik";

import { Buttons, Control } from "../../elements";
import { Remove } from "../../forms";
import { useKeyPress } from "../../hooks";
import { FormikHelpers } from "formik/dist/types";
import SelectFormik from "../../elements/Select/SelectFormik";

export interface RemoveGroupProps {
  remove: any;
  groups: any;
}

export interface RemoveGroupForm {
  groups: any[];
}

const RemoveGroup = ({ remove, groups }: RemoveGroupProps) => {
  const dispatch = useDispatch();

  const onSubmit = async (
    values: RemoveGroupForm,
    formikHelpers: FormikHelpers<RemoveGroupForm>
  ) => {
    const { setFieldError, resetForm, validateForm } = formikHelpers;

    const errors = await validateForm(values);
    if (!isEmpty(errors)) {
      return Object.keys(errors).map((key) =>
        setFieldError(key, `${key} is required`)
      );
    }

    const { groups } = values;
    dispatch(remove(groups.map(({ id }) => id)));
    resetForm({});
  };

  const formikbag = useFormik<RemoveGroupForm>({
    isInitialValid: false,
    validateOnMount: false,
    onSubmit,
    ...Remove,
  });

  const { handleSubmit, setFieldValue } = formikbag;

  useKeyPress("Enter", handleSubmit, []);

  return (
    <FormikProvider value={formikbag}>
      <Fragment key="Remove-Group">
        <Row bottom={16} top={0}>
          <SelectFormik
            useDefault={false}
            required
            isMulti
            autoFocus
            isClearable
            closeMenuOnSelect={false}
            name="groups"
            label="Remove Group(s)"
            placeholder="Select groups"
            components={{ Control }}
            getOptionLabel={(option: any) => option.name}
            getOptionValue={(option: any) => option.id}
            isOptionDisabled={(option: any) => option.id === "default"}
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
        </Row>

        <RowEnd>
          <Buttons.Primary
            variant="IconButton"
            command="↩︎"
            text="Remove"
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
  justify-content: flex-end;
`;

export default RemoveGroup;
