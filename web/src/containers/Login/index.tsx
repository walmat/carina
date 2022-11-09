import React, { useCallback, useState } from "react";
import { motion } from 'framer-motion';
import { useHistory } from 'react-router-dom';
import { ArrowRight } from 'react-feather';
import { useFormik, FormikProvider } from "formik";
import styled from "styled-components";

import { Typography, InputFormik } from "../../elements";

import { Login as LoginForm } from '../../forms';

interface Form {
  email: string;
  password: string;
}

const AppleIcon = () => {
  return (
    <svg width="14" height="16" viewBox="0 0 14 16" fill="none" xmlns="http://www.w3.org/2000/svg" xmlnsXlink="http://www.w3.org/1999/xlink">
      <rect x="0.0994873" width="13.5385" height="16" fill="url(#pattern0)"/>
      <defs>
        <pattern id="pattern0" patternContentUnits="objectBoundingBox" width="1" height="1">
          <use xlinkHref="#image0" transform="translate(0 -0.00266569) scale(0.0019802 0.00167555)"/>
        </pattern>
        <image id="image0" width="505" height="600" xlinkHref="data:image/png;base64,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"/>
      </defs>
    </svg>
  )
}

const GoogleIcon = () => {
  return (
    <svg width="15" height="14" viewBox="0 0 17 16" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path fillRule="evenodd" clipRule="evenodd" d="M8.86077 9.52V6.544H16.3488C16.4608 7.048 16.5488 7.52 16.5488 8.184C16.5488 12.752 13.4848 16 8.86877 16C4.45277 16 0.868774 12.416 0.868774 8C0.868774 3.584 4.45277 0 8.86877 0C11.0288 0 12.8368 0.792001 14.2208 2.088L11.9488 4.296C11.3728 3.752 10.3648 3.112 8.86877 3.112C6.22077 3.112 4.06077 5.312 4.06077 8.008C4.06077 10.704 6.22077 12.904 8.86877 12.904C11.9328 12.904 13.0608 10.784 13.2688 9.528H8.86077V9.52Z" fill="#F4F4F4"/>
    </svg>
  )
}

const renderSubForm = (email: string, password: string, error: string, loading: boolean, history: any, submitForm: any) => {
  const handleNavigate = useCallback(() => {
    history.push('/forgot')
  }, []);

  // NOTE: Render alternative methods
  if (!email) {
    return (
      <Col>
        <InnerSpacer>
          <FlexRow>
            <Line />
            <LineText>or</LineText>
            <Line />
          </FlexRow>
        </InnerSpacer>
        <FlexRow>
          <Col fill>
            {/* TODO: Apple OAuth Flow */}
            <AlternativeButton
              m="0 8px 0 0"
            >
              <AlternativeIcon component={<AppleIcon />} />
              <AlternativeText>
                Login with Apple
              </AlternativeText>
            </AlternativeButton>
          </Col>
          <Col fill>
            {/* TODO: Google OAuth Flow */}
            <AlternativeButton
              m="0 0 0 8px"
            >
              <AlternativeIcon component={<GoogleIcon />} />
              <AlternativeText>
                Login with Google
              </AlternativeText>
            </AlternativeButton>
          </Col>
        </FlexRow>
      </Col>
    )
  }

  // NOTE: Render password and log in button
  return (
    <Col>
      <InnerSpacer>
        <Label>Password</Label>
        <InputFormik
          id="password"
          name="password"
          type="password"
          maxLength={60}
          onSubmit={submitForm}
          placeholder="Enter your password"
        />
        <ForgotPassword onClick={handleNavigate}>Forgot password?</ForgotPassword>
      </InnerSpacer>
      <InnerSpacer>
        <FlexRow>
          <LoginButton
            type="submit"
            disabled={!password || !email || loading}
            whileHover={{ scale: 1.01 }}
            whileTap={{ scale: 0.99 }}
            onClick={submitForm}
          >
            {loading ? <>Authenticating...</> : (
              <>
                Log in
                <Arrow />
              </>
            )}
          </LoginButton>
        </FlexRow>
      </InnerSpacer>
      <ErrorMessage>
        {error}
      </ErrorMessage>
    </Col>
  )
}

const Login = () => {
  const history = useHistory();

  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const handleLogin = async () => {
    if (isValid && !isLoading) {
      setError('');
      setIsLoading(true);

      const start = Date.now();

      try {
        const { Success, Message } = await window.RPCAction('auth:login', [values.email, values.password]);
        if (!Success) {
          // NOTE: Wait at least a second before hearing back
          setTimeout(() => {
            setError(Message);
            setIsLoading(false);

            setTimeout(() => {
              setError('');
            }, 2500);
          }, 1000 - (Date.now() - start))
        } else {
          setTimeout(() => {
            setIsLoading(false);
            }, 1000 - (Date.now() - start))
        }
      } catch (err) {
        setError("Something went wrong. Please try again.");
        setIsLoading(false);

        setTimeout(() => {
          setError('');
        }, 2500);
      }
    }
  };

  const formikbag = useFormik<Form>({
    isInitialValid: false,
    validateOnMount: true,
    onSubmit: handleLogin,
    initialValues: LoginForm.initialValues,
    validationSchema: LoginForm.validationSchema,
  });

  const { values, isValid } = formikbag;

  const handleRegister = useCallback(() => {
    history.push('/register');
  }, [history]);

  return (
    <Container>
      <Inner>
        <Spacer>
          <Title>Welcome to Nebula</Title>
          <Helper>Don’t have an account?{' '}<Emphasis onClick={handleRegister}>Create one here</Emphasis></Helper>
        </Spacer>
        <Spacer>
          <FormikProvider value={formikbag}>
            <InnerSpacer>
              <Label>Email</Label>
              <InputFormik
                id="email"
                name="email"
                maxLength={60}
                placeholder="Enter your email address"
                type="email"
                autoFocus={true}
              />
            </InnerSpacer>
            {renderSubForm(values.email, values.password, error, isLoading, history, handleLogin)}
          </FormikProvider>
        </Spacer>
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
  padding: 0 32px;
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

const Emphasis = styled.span`
  color: ${({ theme }) => theme.colors.primary};
  font-weight: 500;
  cursor: pointer;
`;

const FlexRow = styled.div`
  display: flex;
`;

const Label = styled(Typography.Paragraph)`
  display: flex;
  margin: 0 0 8px 0;
  font-size: 14px;
  font-weight: 500;
  color: ${({ theme }) => theme.colors.h1};
`;

const Line = styled.div`
  display: flex;
  margin: auto 0;
  height: 1px;
  flex: 1;
  background-color: ${({ theme }) => theme.colors.tertiaryHue};
`;

const LineText = styled.div`
  display: flex;
  color: ${({ theme }) => theme.colors.h2};
  margin: auto 8px;
`;

const Col = styled.div<{ fill?: boolean }>`
  display: flex;
  flex-direction: column;
  ${({ fill }) => fill ? `
    flex: 1;
  ` : ''}
`;

const AlternativeButton = styled(motion.button)<{ m?: string }>`
  padding: 10px 0;
  border: none;
  display: flex;
  color: #fff;
  align-items: center;
  justify-content: center;
  background-color: ${({ theme }) => theme.colors.altLogin};
  border-radius: 4px;
  ${({ m }) => m ? `margin: ${m}` : ''};
`;

const AlternativeIcon = styled(({ component, ...props }) => React.cloneElement(component, props))`
  display: flex;
`;

const AlternativeText = styled(Typography.Paragraph)`
  margin: 2px 0 0 8px;
  cursor: pointer;
  font-size: 14px;
  color: #fff;
`;

const ForgotPassword = styled.div`
  display: flex;
  justify-content: flex-end;
  color: ${({ theme }) => theme.colors.primary};
  font-size: 12px;
  margin: 8px 0 0 0;
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

const Arrow = styled(ArrowRight)`
  color: #fff;
  height: 14px;
  width: auto;
  stroke-width: 3px;
  margin-left: 4px;
  cursor: pointer;
`;

export default Login;
