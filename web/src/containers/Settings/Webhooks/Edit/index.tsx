import React, {Fragment, Dispatch, SetStateAction, useCallback, useEffect} from "react";
import { useDispatch, useSelector } from "react-redux";
import { useFormik, FormikProvider } from "formik";
import { Repeat } from 'react-feather';
import styled from "styled-components";
import { useHotkeys } from "react-hotkeys-hook";
import toast from "react-hot-toast";

import { Buttons, Toggle, Typography, Control } from "../../../../elements";
import { Modal } from "../../../../components";

import { editWebhook, Field, Webhook } from '../../../../stores/Main/reducers/webhooks';
import { makeProfileGroupOptions, ProfileGroup } from '../../../../stores/Main/reducers/profiles';

import { fields } from '../../../../constants';
import { Webhooks } from "../../../../forms";

import { getType, isValidWebhook } from "../utils";
import InputFormik from '../../../../elements/Input/InputFormik';
import SelectFormik from "../../../../elements/Select/SelectFormik";
import { noop } from "../../../../utils";

interface EditWebhookProps {
  open: boolean;
  webhook: Webhook;
  setOpen: Dispatch<SetStateAction<boolean>>;
}

interface AddWebhookForm {
  profiles: { id: string; name: string; }[];
  url: string;
  fields: Field[];
  declines: boolean;
  sensitivity: boolean;
}

const defaultFields = fields.filter(f => f.default);

const cheapComparer = (a1: any[], a2: any[]) => JSON.stringify(a1) === JSON.stringify(a2);

const EditWebhook = ({ open, webhook, setOpen }: EditWebhookProps) => {
  const dispatch = useDispatch();
  const profiles = useSelector(makeProfileGroupOptions);

  const onSubmit = (values: AddWebhookForm) => {
    Object.keys(values).map((k) => setFieldTouched(k, true));

    if (isValid) {
      setTouched({});
      handleEditWebhook(values).then(() => {});
    }
  };

  const handleEditWebhook = async (values: AddWebhookForm) => {
    try {
      const { url } = values;

      if (!isValidWebhook(url)) {
        toast.error("Invalid webhook");
        return;
      }

      console.log(values);

      const type = getType(url);
      dispatch(editWebhook({ ...values, type }));

      setOpen(false);
      handleClear();
    } catch (e) {
      // noop...
    }
  };

  const formikbag = useFormik<AddWebhookForm>({
    isInitialValid: false,
    validateOnMount: true,
    ...Webhooks,
    onSubmit,
  });

  const {
    values,
    resetForm,
    isValid,
    setTouched,
    setFieldValue,
    setFieldTouched,
    handleSubmit,
  } = formikbag;

  const handleClear = useCallback(() => {
    resetForm({});
  }, []);

  useHotkeys("command+c,ctrl+c", handleClear, []);
  // @ts-ignore
  useHotkeys("enter", handleSubmit, []);

  const diffLength = cheapComparer(defaultFields, values.fields.filter(f => f.enabled))
  const resetToDefaults = useCallback(() => {
    setFieldValue('fields', Webhooks.initialValues.fields);
  }, []);

  useEffect(() => {
    if (webhook) {
      Object.entries(webhook).forEach(([key, value]: [any, any]) => {
        // NOTE: We gotta do a little extra sauce here
        if (key === 'fields') {
          const f = fields.map(f => ({ ...f, enabled: value.some((v: string) => v === f.name)} ));
          return setFieldValue('fields', f);
        }

        return setFieldValue(key, value);
      });
    }
  }, [webhook]);

  return (
    <FormikProvider value={formikbag}>
      <Modal
        show={open}
        setShow={setOpen}
        amount={webhook.id.slice(0, 5)}
        title="Edit Webhook"
        width={368}
        height="auto"
      >
        <Fragment>
          <Row m="0 0 16px 0">
            <SelectFormik
              isClearable
              isMulti
              name="profiles"
              label="Profile Groups"
              placeholder="All Profiles"
              components={{ Control }}
              getOptionValue={(option: ProfileGroup) => option.id}
              getOptionLabel={(option: ProfileGroup) => option.name}
              options={profiles}
            />
          </Row>

          <Row m="0 0 8px 0">
            <ColFill m="0 8px 0 0">
              <InputFormik
                id="name"
                name="name"
                useLabel
                label="Webhook Name"
                placeholder="Webhook Name"
              />
            </ColFill>

            <ColFill m="0 0 0 8px">
              <InputFormik
                id="url"
                name="url"
                useLabel
                label="Webhook Link"
                placeholder="Webhook Link"
              />
            </ColFill>
          </Row>

          <BottomBorder />

          <RowStart m="0 0 12px 0">
            <FieldsTitle>
              Webhook Fields
            </FieldsTitle>
            {!diffLength && (
              <DefaultButton onClick={resetToDefaults}>
                <RepeatIcon />
                Defaults
              </DefaultButton>
            )}
          </RowStart>

          <RowWrap m="0">
            {values.fields.map((field, index) => {
              let margin = '0 8px 8px 0';
              if (index > 4) {
                margin = '0 8px 0 0';
              }

              if (index === values.fields.length - 1) {
                margin = '0'
              }

              return (
                <Buttons.Toggle
                  onClick={() => setFieldValue(`fields[${index}].enabled`, !field.enabled)}
                  active={field.enabled}
                  margin={margin}
                  height={24}
                >
                  {field.name}
                </Buttons.Toggle>
              )
            })}
          </RowWrap>

          <BottomBorder />

          <Row m="0 0 16px 0" onClick={() => setFieldValue('declines', !values.declines)}>
            <ColFill>
              <Title>Enable Declines</Title>
              <SubText>Attempts to send all decline payloads.</SubText>
            </ColFill>
            <Col>
              <Toggle on={values.declines} check={noop} />
            </Col>
          </Row>

          <Row m="0 0 8px 0" onClick={() => setFieldValue('sensitivity', !values.sensitivity)}>
            <ColFill>
              <Title>Enable Sensitivity</Title>
              <SubText>Hides sensitive fields when applicable.</SubText>
            </ColFill>
            <Col>
              <Toggle on={values.sensitivity} check={noop} />
            </Col>
          </Row>

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

const RowWrap = styled.div<{ m: string }>`
  display: flex;
  flex-wrap: wrap;
  margin: ${({ m }) => m};
`;

const RowStart = styled.div<{ m: string }>`
  display: flex;
  min-height: 24px;
  justify-content: flex-start;
  margin: ${({ m }) => m};
`;

const FieldsTitle = styled(Typography.H5)`
  color: ${({ theme }) => theme.colors.h2};
  font-size: 16px;
  margin: 0;
`;

const DefaultButton = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 4px 8px;
  margin-left: 8px;
  margin-top: -4px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  border: 1px solid ${({ theme }) => theme.colors.stop};
  color: ${({ theme }) => theme.colors.stop};
  background-color: rgba(255, 177, 94, 0.2);
`;

const RepeatIcon = styled(Repeat)`
  height: 12px;
  width: auto;
  margin-right: 6px;
  color: ${({ theme }) => theme.colors.stop};
`;

const Col = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
`;

const ColFill = styled.div<{ m?: string }>`
  display: flex;
  flex: 1;
  flex-direction: column;
  ${({ m }) => m ? `margin:${m};` : ''}
`;

const Title = styled(Typography.H2)`
  font-size: 14px;
  margin: 0 0 4px 0;
  font-weight: 500;
  color: ${({ theme }) => theme.colors.subHeading};
`;

const SubText = styled(Typography.Paragraph)`
  font-size: 12px;
  margin: 0;
  font-weight: 400;
  color: ${({ theme }) => theme.colors.paragraph};
`;

const BottomBorder = styled.div`
  display: flex;
  margin: 16px 0;
  padding: 0.5px; 0;
  background-color: ${({ theme }) => theme.colors.separator};
`;

export default EditWebhook;
