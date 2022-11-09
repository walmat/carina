import { useCallback, useState } from "react";
import { motion } from "framer-motion";
import { useHistory } from "react-router-dom";
import { ArrowLeft, ArrowRight } from "react-feather";
import { useFormik, FormikProvider } from "formik";
import styled from "styled-components";

import { Typography, InputFormik, Buttons } from "../../elements";

import { ForgotPassword as ForgotPasswordForm } from "../../forms";

interface Form {
  email: string;
}

const ForgotPassword = () => {
  const history = useHistory();

  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");
  const [msg, setMsg] = useState("");

  const handleSubmit = async () => {
    if (isValid && !isLoading) {
      setError("");
      setIsLoading(true);

      const start = Date.now();

      try {
        // TODO: Solve ReCAPTCHA and send proper response
        const res = await window.RPCAction("auth:forgot", [
          values.email,
          "abc",
        ]);
        if (!res) {
          // NOTE: Wait at least a second before hearing back
          setTimeout(() => {
            setError(res.Message);
            setIsLoading(false);

            setTimeout(() => {
              setError("");
            }, 2500);
          }, 1000 - (Date.now() - start));
        } else {
          setMsg("Request sent. Please check your email.");
          setIsLoading(false);
          handleForward();
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
    onSubmit: handleSubmit,
    ...ForgotPassword,
    initialValues: ForgotPasswordForm.initialValues,
    validationSchema: ForgotPasswordForm.validationSchema,
  });

  const { values, isValid, submitForm } = formikbag;

  const handleBack = useCallback(() => {
    history.push("/");
  }, [history]);

  const handleForward = useCallback(() => {
    history.push("/change");
  }, [history]);

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
        <Spacer>
          <Title>Forgot Password</Title>
          <Helper>Enter your email address below</Helper>
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
                onSubmit={handleSubmit}
                type="email"
              />
            </InnerSpacer>
            <InnerSpacer>
              <FlexRow>
                <LoginButton
                  type="submit"
                  disabled={!values.email || isLoading}
                  whileHover={{ scale: 1.01 }}
                  whileTap={{ scale: 0.99 }}
                  onClick={submitForm}
                >
                  {isLoading ? (
                    <>Checking...</>
                  ) : (
                    <>
                      Continue
                      <Arrow />
                    </>
                  )}
                </LoginButton>
              </FlexRow>
            </InnerSpacer>
            <Message error={!!error}>{msg || error}</Message>
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

const Label = styled(Typography.Paragraph)`
  display: flex;
  margin: 0 0 8px 0;
  font-size: 14px;
  font-weight: 500;
  color: ${({ theme }) => theme.colors.h1};
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

const Message = styled(Typography.Paragraph)<{ error: boolean }>`
  color: ${({ theme, error }) =>
    error ? theme.colors.failed : theme.colors.success};
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

export default ForgotPassword;
