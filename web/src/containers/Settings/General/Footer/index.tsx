import React, {useCallback} from 'react';
import styled from "styled-components";
import { motion } from 'framer-motion';
import { Twitter, Instagram } from 'react-feather';

import { Typography } from '../../../../elements';

import { METADATA } from '../../../../constants';

const Footer = () => {
	const handler = useCallback((url: string) => astilectron.sendMessage({ type: 'openExternalUrl', data: url }), [])

	return (
		<Container>
			{/*@ts-ignore*/}
			<FooterLink mr={8} ml={16} onClick={() => handler(METADATA.PRIVACY_POLICY)}>
				Privacy Policy
			</FooterLink>
			{/*@ts-ignore*/}
			<FooterLink mr={8} ml={8} onClick={() => handler(METADATA.TERMS_OF_SERVICE)}>
				Terms of Service
			</FooterLink>
			{/*@ts-ignore*/}
			<FooterLink mr={8} ml={8} onClick={() => handler(METADATA.APP_DOCUMENTATION)}>
				Documentation
			</FooterLink>
			<IconContainer
				whileHover={{ scale: 1.15 }}
				whileTap={{ scale: 0.95 }}
				onClick={() => handler(METADATA.APP_TWITTER)}
			>
				<Icon component={<Twitter />} />
			</IconContainer>
			<IconContainer
				whileHover={{ scale: 1.15 }}
				whileTap={{ scale: 0.95 }}
				onClick={() => handler(METADATA.APP_INSTAGRAM)}
			>
				<Icon component={<Instagram />} />
			</IconContainer>
		</Container>
	);
};

const Container = styled.div`
  display: flex;
  position: absolute;
  bottom: 0;
  margin-bottom: -16px;
`;

const FooterLink = styled(Typography.H2)<{ mr: number; ml: number }>`
	margin-right: ${({ mr }) => mr}px;
	margin-left: ${({ ml }) => ml}px;
	cursor: pointer;
	color: ${({ theme }) => theme.colors.primary};
	font-size: 12px;
`;

const IconContainer = styled(motion.div)`
	display: flex;
	justify-content: center;
	align-items: center;
  width: 24px;
  height: 24px;
  cursor: pointer !important;
  border-radius: 4px;
  margin: auto 8px;
`

const Icon = styled(
	({ component, ...props }) => React.cloneElement(component, props)
)`
  width: 14px;
  height: 14px;
  cursor: pointer !important;
  color: ${({ theme }) => theme.colors.h2};
`

export default Footer;
