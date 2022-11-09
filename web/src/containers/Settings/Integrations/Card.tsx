import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import styled from 'styled-components';
import { useFormik, FormikProvider } from 'formik';
import * as yup from 'yup';

import { Card } from '../../../components';
import { Typography, Toggle, Input, Buttons } from "../../../elements";
import { AYCD, TwoCaptcha, CapMonster, Scout } from '../../../icons';

import { makeTheme } from '../../../stores/Main/reducers/theme';
import { Integration, saveIntegration, disableIntegration, stash, Credentials } from '../../../stores/Main/reducers/integrations';

type Props = {
	integration: Integration;
};

type Form = {
	credentials: Credentials[];
};

const getIcon = (type: string) => {
	switch (type) {
		default:
		case 'aycd':
			return <AYCD />;
		case '2captcha':
			return <TwoCaptcha />;
		case 'capmonster':
			return <CapMonster />;
		case 'scout':
			return <Scout />;
	}
};

const CardComponent = ({ integration }: Props) => {
	const theme = useSelector(makeTheme);
	const dispatch = useDispatch();

	const { id, name, active, credentials } = integration;

	const handleOAuth = () => {};

	const onSubmit = async (values: Form) => {
		if (values.credentials.every(cred => cred.value)) {
			const { credentials: _credentials } = values;
			dispatch(saveIntegration({ id, credentials: _credentials }));
			dispatch(stash())
		}

		return null;
	};

	const handleDisable = () => {
		dispatch(disableIntegration(id));
		dispatch(stash())
	}

	const formikbag = useFormik<Form>({
		isInitialValid: false,
		validateOnMount: true,
		initialValues: {
			credentials: credentials || []
		},
		validationSchema: yup.object().shape({
			credentials: yup.array().of(
				yup.object().shape({
					label: yup.string(),
					value: yup.string()
				})
			)
		}),
		onSubmit
	});

	const {
		values,
		errors,
		isValid,
		setFieldValue,
		handleSubmit,
	} = formikbag;

	return (
		<FormikProvider value={formikbag}>
			<Card key={id}>
				<Flex height={32} ai="center">
					<IconContainer>
						{getIcon(id)}
					</IconContainer>
					<Title>{name}</Title>
					<ToggleContainer>
						<Toggle on={active} check={active ? handleDisable : () => handleSubmit()} theme={theme} />
					</ToggleContainer>
				</Flex>
				<Flex fd="column" mt="auto" ai="center">
					{credentials?.length ? credentials.map(({ label, value }, index) => {
						return (
							<InputContainer>
								<Input
									id="url"
									name="url"
									useLabel
									masked={active}
									disabled={active}
									touched={false}
									// @ts-ignore
									error={!!(errors.credentials || []).find(cred => cred.label === label)}
									placeholder={`${label}`}
									onChange={({ target: { value }}) => setFieldValue(`credentials.${index}.value`, value)}
									// @ts-ignore
									value={values.credentials.find(cred => cred.label === label).value}
								/>
							</InputContainer>
						);
					}) : (
						<InputContainer>
							<Buttons.Primary width="100%" height={36} text="Login with Scout" onClick={handleOAuth} />
						</InputContainer>
					)}

				</Flex>
			</Card>
		</FormikProvider>
	);
};

const Flex = styled.div<{ ai: string; height?: number; fd?: string; mt?: string; mb?: string }>`
	display: flex;
	${({ height }) => height ? `height: ${height}px` : ''};
	${({ fd }) => fd ? `flex-direction: ${fd}` : ''};
	${({ mt }) => mt ? `margin-top: ${mt}` : ''};
	${({ mb }) => mb ? `margin-bottom: ${mb}` : ''};
`;

const IconContainer = styled.div`
	margin: auto 8px auto 0;
`;

const Title = styled(Typography.H2)`
	display: flex;
	color: ${({ theme }) => theme.colors.h2};
	text-transform: capitalize;
	font-size: 1.35vw;
	font-weight: 700;
	margin: auto 16px auto 8px;
`;

const InputContainer = styled.div`
	display: flex;
	flex-direction: column;
	margin-top: 16px;
	flex: 1 1 auto;
`;

const ToggleContainer = styled.div`
	margin: auto 0 auto auto;
`;

export default CardComponent;
