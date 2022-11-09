import { useCallback, useState } from "react";
import { motion } from "framer-motion";
import { useHistory } from "react-router-dom";
import { ArrowLeft, ArrowRight } from "react-feather";
import { useFormik, FormikProvider } from "formik";
import styled from "styled-components";

import { Typography, InputFormik, Buttons, Checkbox } from "../../elements";

import { Register as RegisterForm } from "../../forms";
import { METADATA } from "../../constants";

interface Form {
  key: string;
  email: string;
  password: string;
  terms: boolean;
}

const Register = () => {
  const history = useHistory();

  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const handleRegister = async () => {
    console.log(isValid);
    if (isValid && !isLoading) {
      setError("");
      setIsLoading(true);

      const start = Date.now();

      try {
        const res = await window.RPCAction("auth:register", [
          values.key,
          values.email,
          values.password,
          "abc",
        ]);

        if (!res.Success) {
          // NOTE: Wait at least a second before hearing back
          setTimeout(() => {
            setError(res.Message);
            setIsLoading(false);

            setTimeout(() => {
              setError("");
            }, 2500);
          }, 1000 - (Date.now() - start));
        } else {
          setSuccess(res.Message);
          setIsLoading(false);
        }
      } catch (err) {
        setError("Something went wrong. Please try again.");
        setIsLoading(false);

        setTimeout(() => {
          setError("");
        }, 2500);
      }
    }
  };

  const formikbag = useFormik<Form>({
    isInitialValid: false,
    validateOnMount: true,
    onSubmit: handleRegister,
    initialValues: RegisterForm.initialValues,
    validationSchema: RegisterForm.validationSchema,
  });

  const { values, isValid, submitForm, setFieldValue } = formikbag;

  const handleBack = useCallback(() => {
    history.push("/");
  }, [history]);

  const handler = useCallback((e: any, url: string) => {
    e.preventDefault();
    e.stopPropagation();

    // window.openUrlInBrowser(url)
  }, []);

  return (
    <Container>
      <Inner>
        <InnerSpacer>
          <Buttons.Tertiary
            variant="IconButton"
            width="auto"
            height="36"
            onClick={handleBack}
          >
            <BackArrow />
            Back to log in
          </Buttons.Tertiary>
        </InnerSpacer>
        {/* <InnerSpacer> */}
        <Title>Register</Title>
        {/* </InnerSpacer> */}
        <Spacer>
          <FormikProvider value={formikbag}>
            <InnerSpacer>
              <Label>Key</Label>
              <InputFormik
                id="key"
                name="key"
                maxLength={27}
                placeholder="Enter your key"
                type="key"
                autoFocus={true}
              />
            </InnerSpacer>
            <InnerSpacer>
              <Label>Email</Label>
              <InputFormik
                id="email"
                name="email"
                maxLength={60}
                placeholder="Enter your email address"
                type="email"
              />
            </InnerSpacer>
            <InnerSpacer>
              <Label>Password</Label>
              <InputFormik
                id="password"
                name="password"
                maxLength={60}
                placeholder="Enter your password"
                onSubmit={handleRegister}
                type="password"
              />
              <TermsRow onClick={() => setFieldValue("terms", !values.terms)}>
                <Col>
                  <Checkbox isChecked={values.terms} setIsChecked={() => {}} />
                </Col>
                <Col>
                  <TermsText>
                    I agree to the{" "}
                    <Emphasis
                      onClick={(e) => handler(e, METADATA.TERMS_OF_SERVICE)}
                    >
                      ToS{" "}
                    </Emphasis>{" "}
                    &{" "}
                    <Emphasis
                      onClick={(e) => handler(e, METADATA.PRIVACY_POLICY)}
                    >
                      Privacy Policy
                    </Emphasis>
                  </TermsText>
                </Col>
              </TermsRow>
            </InnerSpacer>
            <InnerSpacer>
              <FlexRow>
                <LoginButton
                  type="submit"
                  disabled={
                    !values.key ||
                    !values.email ||
                    !values.password ||
                    isLoading
                  }
                  whileHover={{ scale: 1.01 }}
                  whileTap={{ scale: 0.99 }}
                  onClick={submitForm}
                >
                  {isLoading ? (
                    <>Registering...</>
                  ) : (
                    <>
                      Continue
                      <Arrow />
                    </>
                  )}
                </LoginButton>
              </FlexRow>
            </InnerSpacer>
            {!!error && <ErrorMessage>{error}</ErrorMessage>}
            {!!success && <SuccessMessage>{success}</SuccessMessage>}
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
  padding: 32px;
`;

const Spacer = styled.div`
  margin-bottom: 48px;
`;

const InnerSpacer = styled.div`
  margin-bottom: 20px;
`;

const Title = styled(Typography.H2)`
  font-size: 24px;
  font-weight: 700;
  color: ${({ theme }) => theme.colors.h1};
  margin: 0 0 20px 0;
`;

const FlexRow = styled.div`
  display: flex;
`;

const TermsRow = styled.div`
  display: flex;
  justify-content: flex-start;
  align-items: center;
  margin: 16px 0 0 0;
`;

const Label = styled(Typography.Paragraph)`
  display: flex;
  margin: 0 0 8px 0;
  font-size: 14px;
  font-weight: 500;
  color: ${({ theme }) => theme.colors.h1};
`;

const Col = styled.div<{ fill?: boolean }>`
  display: flex;
  flex-direction: column;
  ${({ fill }) =>
    fill
      ? `
    flex: 1;
  `
      : ""}
`;

const TermsText = styled.div`
  display: flex;
  white-space: pre;
  color: ${({ theme }) => theme.colors.h1};
  font-size: 14px;
  margin: 0 0 0 8px;
  cursor: pointer;
  font-weight: 400;
`;

const Emphasis = styled.span`
  display: flex;
  color: ${({ theme }) => theme.colors.primary};
  font-size: 14px;
  cursor: pointer;
  font-weight: 500;
`;

const LoginButton = styled(motion.button)<{ disabled: boolean }>`
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
  cursor: ${({ disabled }) => (disabled ? "default" : "cursor")};
  opacity: ${({ disabled }) => (disabled ? 0.4 : 1)};
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

const SuccessMessage = styled(Typography.Paragraph)`
  color: ${({ theme }) => theme.colors.success};
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

export default Register;
