import {useCallback, useEffect, useState} from "react";
import { motion } from 'framer-motion';
import { useHistory } from 'react-router-dom';
import {ArrowLeft, ArrowRight} from 'react-feather';
import { useFormik, FormikProvider } from "formik";
import styled from "styled-components";
import ReactCodeInput from 'react-verification-code-input';

import { Typography, Buttons, Checkbox } from "../../elements";
import { TwoFactor as TwoFactorForm } from '../../forms';

interface Form {
  code: string;
  remember: boolean;
}

const FormComponent = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async () => {
    if (!isLoading) {
      setError('');
      setIsLoading(true);

      const start = Date.now();

      try {
        const res = await window.RPCAction('auth:2fa', [values.code]);
        if (!res) {
          // NOTE: Wait at least a second before hearing back
          setTimeout(() => {
            setError("Invalid code. Please try again.");
            setIsLoading(false);

            setTimeout(() => {
              setError('');
            }, 2500);
          }, 1000 - (Date.now() - start))
        }
      } catch (err) {
        setTimeout(() => {
          setError("Something went wrong. Please try again.");
          setIsLoading(false);

          setTimeout(() => {
            setError('');
          }, 2500);
        },1000 - (Date.now() - start));
      }
    }
  };

  const formikbag = useFormik<Form>({
    isInitialValid: false,
    validateOnMount: true,
    // @ts-ignore
    onSubmit: handleSubmit,
    ...TwoFactorForm
  });

  const { values, submitForm, setFieldValue } = formikbag;

  return (
    <Spacer>
      <FormikProvider value={formikbag}>
        <InnerSpacer>
          <FlexRow>
            <ReactCodeInput
              disabled={isLoading}
              type="number"
              className="react-2fa-code"
              fields={6}
              onChange={value => setFieldValue('code', value)}
              fieldWidth={62.416}
              fieldHeight={63}
              placeholder={['', '', '', '', '', '']}
              values={values.code.split(' ')}
            />
          </FlexRow>
          <RememberRow onClick={() => setFieldValue('remember', !values.remember)}>
            <Col>
              <Checkbox isChecked={values.remember} setIsChecked={() => {}} />
            </Col>
            <Col>
              <RememberText>
                Remember this device for 30-days.
              </RememberText>
            </Col>
          </RememberRow>
        </InnerSpacer>
        <InnerSpacer>
          <FlexRow>
            <LoginButton
              type="submit"
              disabled={values.code.split("").length < 6 || isLoading}
              whileHover={{ scale: 1.01 }}
              whileTap={{ scale: 0.99 }}
              onClick={submitForm}
            >
              {isLoading ? <>Checking...</> : (
                <>
                  Continue
                  <Arrow />
                </>
              )}
            </LoginButton>
          </FlexRow>
        </InnerSpacer>
        <ErrorMessage>
          {error}
        </ErrorMessage>
      </FormikProvider>
    </Spacer>
  );
};

const TwoFactorAuthentication = () => {
  const history = useHistory();

  // NOTE: Hack to focus the first input on page load
  useEffect(() => {
    const input = document.querySelector('.react-2fa-code > div > input:first-child')
    if (input) {
      // @ts-ignore
      input.focus();
    }

  }, []);

  const handleBack = useCallback(() => {
    history.push('/');
  }, [history]);

  return (
    <Container>
      <Inner>
        <InnerSpacer>
          <Buttons.Tertiary  variant="IconButton" width="auto" height="36" onClick={handleBack}>
            <BackArrow />
            Back to log in
          </Buttons.Tertiary>
        </InnerSpacer>
        <Spacer>
          <Title>2-Factor Authentication</Title>
          <Helper>Enter the 6-digit authentication code</Helper>
        </Spacer>
        <FormComponent />
      </Inner>
    </Container>
  );
};

const Container = styled.div`
  display: flex;
  padding: 64px 0 32px 0;
  flex: 1;
`;

const Inner = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  flex: 1;
  padding: 32px;
`;

const Spacer = styled.div`
  margin-bottom: 48px;
`;

const InnerSpacer = styled.div`
  margin-bottom: 32px;
`;

const Title = styled(Typography.H2)`
  font-size: 24px;
  font-weight: 700;
  color: ${({ theme }) => theme.colors.h1};
  margin: 0 0 8px 0;
`;

const Helper = styled(Typography.Paragraph)`
  font-size: 14px;
  color: ${({ theme }) => theme.colors.h2};
  font-weight: 400;
  white-space: pre;
  margin: 0;
`;

const FlexRow = styled.div`
  display: flex;
`;

const RememberRow = styled.div`
  display: flex;
  justify-content: flex-start;
  align-items: center;
  margin: 16px 0 0 0;
`;

const Col = styled.div<{ fill?: boolean }>`
  display: flex;
  flex-direction: column;
  ${({ fill }) => fill ? `
    flex: 1;
  ` : ''}
`;

const RememberText = styled.div`
  display: flex;
  white-space: pre;
  color: ${({ theme }) => theme.colors.h1};
  font-size: 14px;
  margin: 0 0 0 8px;
  cursor: pointer;
  font-weight: 400;
`;

const LoginButton = styled(motion.button)<{ disabled: boolean; }>`
  display: flex;
  flex: 1;
  align-items: center;
  justify-content: center;
  height: 42px;
  box-sizing: border-box;
  border: 1px solid ${({ theme }) => theme.colors.primary};
  color: #fff;
  font-size: 14px;
  font-weight: 500;
  cursor: ${({ disabled }) => disabled ? 'default' : 'cursor'};
  opacity: ${({ disabled }) => disabled ? 0.4 : 1};
  background-color: ${({ theme }) => theme.colors.primary};
  border-radius: 4px;
  margin: 0;
`;

const ErrorMessage = styled(Typography.Paragraph)`
  color: ${({ theme }) => theme.colors.failed};
  font-size: 12px;
  font-weight: 400;
  height: 12px; // hack to keep the div always at 12px height even when there is no error message
  margin: 0;
  display: flex;
  justify-content: center;
`;

const BackArrow = styled(ArrowLeft)`
  color: #fff;
  height: 14px;
  width: auto;
  stroke-width: 3px;
  margin-right: 4px;
  cursor: pointer;
`;

const Arrow = styled(ArrowRight)`
  color: #fff;
  height: 14px;
  width: auto;
  stroke-width: 3px;
  margin-left: 4px;
  cursor: pointer;
`;

export default TwoFactorAuthentication;
