import React, { useCallback, useMemo } from "react";
import { useDispatch, useSelector } from "react-redux";
import styled from "styled-components";
import { Moon, Sun } from "react-feather";

import Routes from "./Routes";

import { makeTheme, setTheme } from "../../stores/Login/reducers/theme";
import Dragbar from "../../components/Dragbar";
import Actions from "../../components/Actions";
import { METADATA } from "../../constants";
import { Typography } from "../../elements";
import Slideshow from './Slideshow';

import "../../i18n/config";

const Logo = () => {
	const theme = useSelector(makeTheme);

	const color = useMemo(() => theme === 0 ? '#202126' : '#FFFFFF', [theme]);

	return (
		<svg width="127" height="25" viewBox="0 0 127 25" fill="none" xmlns="http://www.w3.org/2000/svg">
			<path d="M4.97181 16.7022L5.47108 14.8915L4.17299 13.9263L3.23438 15.8035L3.68038 16.8886L4.97181 16.7022Z" fill={color} />
			<path d="M23.6377 2.70264L25.8145 3.12202L25.7279 4.68638L25.3884 5.04585L23.2915 3.89421L23.6377 2.70264Z" fill={color} />
			<path d="M16.0356 24.0312C22.6717 24.0312 28.0513 18.6517 28.0513 12.0156C28.0513 5.37958 22.6717 0 16.0356 0C9.3996 0 4.02002 5.37958 4.02002 12.0156C4.02002 18.6517 9.3996 24.0312 16.0356 24.0312Z" fill="url(#paint0_linear)"/>
			<path d="M23.6378 2.70251C24.3901 3.00207 24.9626 3.46805 25.3287 4.08714C26.9197 6.84308 23.5513 11.7758 17.8064 15.0909C12.0616 18.406 6.11035 18.8587 4.51936 16.1028C4.15323 15.4704 4.05338 14.7315 4.17321 13.9193C0.971261 16.9149 -0.633041 19.6309 0.232349 21.1354C1.73014 23.7315 10.0046 21.7611 18.7118 16.7285C27.4189 11.6959 33.2637 5.51837 31.7659 2.92885C30.9005 1.43771 27.7917 1.45103 23.6378 2.70251Z" fill="white" />
			<path d="M43.3864 12.9545C43.392 11.8295 44.0625 11.1705 45.0398 11.1705C46.0114 11.1705 46.5966 11.8068 46.5909 12.875V18H49.0114V12.4432C49.0114 10.4091 47.8182 9.15909 46 9.15909C44.7045 9.15909 43.767 9.79545 43.375 10.8125H43.2727V9.27273H40.9659V18H43.3864V12.9545ZM54.8963 18.1705C57.0554 18.1705 58.5099 17.1193 58.8509 15.5L56.6122 15.3523C56.3679 16.017 55.7429 16.3636 54.9361 16.3636C53.7259 16.3636 52.9588 15.5625 52.9588 14.2614V14.2557H58.902V13.5909C58.902 10.625 57.1065 9.15909 54.7997 9.15909C52.2315 9.15909 50.5668 10.983 50.5668 13.6761C50.5668 16.4432 52.2088 18.1705 54.8963 18.1705ZM52.9588 12.7557C53.0099 11.7614 53.7656 10.9659 54.8395 10.9659C55.8906 10.9659 56.6179 11.7159 56.6236 12.7557H52.9588ZM60.527 18H62.9134V16.6023H63.0213C63.3565 17.3295 64.0895 18.142 65.4986 18.142C67.4872 18.142 69.0384 16.5682 69.0384 13.6477C69.0384 10.6477 67.419 9.15909 65.5043 9.15909C64.044 9.15909 63.3452 10.0284 63.0213 10.7386H62.9474V6.36364H60.527V18ZM62.8963 13.6364C62.8963 12.0795 63.5554 11.0852 64.7315 11.0852C65.9304 11.0852 66.5668 12.125 66.5668 13.6364C66.5668 15.1591 65.919 16.2159 64.7315 16.2159C63.5668 16.2159 62.8963 15.1932 62.8963 13.6364ZM76.25 14.2841C76.2557 15.4545 75.4545 16.0682 74.5852 16.0682C73.6705 16.0682 73.0795 15.4261 73.0739 14.3977V9.27273H70.6534V14.8295C70.6591 16.8693 71.8523 18.1136 73.608 18.1136C74.9205 18.1136 75.8636 17.4375 76.2557 16.4148H76.3466V18H78.6705V9.27273H76.25V14.2841ZM83.027 6.36364H80.6065V18H83.027V6.36364ZM87.4545 18.1648C88.7443 18.1648 89.5795 17.6023 90.0057 16.7898H90.0739V18H92.3693V12.1136C92.3693 10.0341 90.608 9.15909 88.6648 9.15909C86.5739 9.15909 85.1989 10.1591 84.8636 11.75L87.1023 11.9318C87.267 11.3523 87.7841 10.9261 88.6534 10.9261C89.4773 10.9261 89.9489 11.3409 89.9489 12.0568V12.0909C89.9489 12.6534 89.3523 12.7273 87.8352 12.875C86.108 13.0341 84.5568 13.6136 84.5568 15.5625C84.5568 17.2898 85.7898 18.1648 87.4545 18.1648ZM88.1477 16.4943C87.4034 16.4943 86.8693 16.1477 86.8693 15.483C86.8693 14.8011 87.4318 14.4659 88.2841 14.3466C88.8125 14.2727 89.6761 14.1477 89.9659 13.9545V14.8807C89.9659 15.7955 89.2102 16.4943 88.1477 16.4943ZM94.6903 18H95.9858V16.6591H96.1449C96.4403 17.1364 97.0085 18.1818 98.6903 18.1818C100.872 18.1818 102.395 16.4318 102.395 13.6591C102.395 10.9091 100.872 9.15909 98.6676 9.15909C96.9631 9.15909 96.4403 10.2045 96.1449 10.6591H96.0312V6.36364H94.6903V18ZM96.0085 13.6364C96.0085 11.6818 96.8722 10.3636 98.5085 10.3636C100.213 10.3636 101.054 11.7955 101.054 13.6364C101.054 15.5 100.19 16.9773 98.5085 16.9773C96.8949 16.9773 96.0085 15.6136 96.0085 13.6364ZM107.991 18.1818C110.355 18.1818 111.946 16.3864 111.946 13.6818C111.946 10.9545 110.355 9.15909 107.991 9.15909C105.628 9.15909 104.037 10.9545 104.037 13.6818C104.037 16.3864 105.628 18.1818 107.991 18.1818ZM107.991 16.9773C106.196 16.9773 105.378 15.4318 105.378 13.6818C105.378 11.9318 106.196 10.3636 107.991 10.3636C109.787 10.3636 110.605 11.9318 110.605 13.6818C110.605 15.4318 109.787 16.9773 107.991 16.9773ZM117.788 9.27273H115.925V7.18182H114.584V9.27273H113.266V10.4091H114.584V15.8636C114.584 17.3864 115.811 18.1136 116.947 18.1136C117.447 18.1136 117.766 18.0227 117.947 17.9545L117.675 16.75C117.561 16.7727 117.379 16.8182 117.084 16.8182C116.493 16.8182 115.925 16.6364 115.925 15.5V10.4091H117.788V9.27273ZM125.831 11.2273C125.411 9.98864 124.467 9.15909 122.74 9.15909C120.899 9.15909 119.536 10.2045 119.536 11.6818C119.536 12.8864 120.251 13.6932 121.854 14.0682L123.308 14.4091C124.189 14.6136 124.604 15.0341 124.604 15.6364C124.604 16.3864 123.808 17 122.558 17C121.462 17 120.774 16.5284 120.536 15.5909L119.263 15.9091C119.575 17.392 120.797 18.1818 122.581 18.1818C124.609 18.1818 125.99 17.0739 125.99 15.5682C125.99 14.3523 125.229 13.5852 123.672 13.2045L122.376 12.8864C121.342 12.6307 120.876 12.2841 120.876 11.6136C120.876 10.8636 121.672 10.3182 122.74 10.3182C123.911 10.3182 124.393 10.9659 124.626 11.5682L125.831 11.2273Z" fill={color} />
			<defs>
				<linearGradient id="paint0_linear" x1="10.0272" y1="22.4209" x2="22.0443" y2="1.60609" gradientUnits="userSpaceOnUse">
					<stop stopColor="#776AF2"/>
					<stop offset="1" stopColor="#A097FD"/>
				</linearGradient>
			</defs>
		</svg>
	);
}

const TopBar = () => {
	return (
		<Float>
			<Flex margin="0">
				<Logo />
			</Flex>
			<Flex margin="0 0 0 auto">
				<Actions useMargin={false} simple />
			</Flex>
		</Float>
	);
};

const renderThemeIcon = (theme: number, dispatch: any) => {
	const handleTheme = ({ value }: any) => {
		if (theme === value) {
			return;
		}

		return dispatch(setTheme({ theme }));
	};

	switch (theme) {
		default:
		case 0:
			return <NightIcon onClick={handleTheme} />
		case 1:
			return <DayIcon onClick={handleTheme} />
	}
}

const Footer = () => {
	const dispatch = useDispatch();
	const theme = useSelector(makeTheme);

	const handler = useCallback((url: string) => astilectron.sendMessage({ type: 'openExternalUrl', data: url }), [])

	return (
		<Bottom>
			{/* @ts-ignore */}
			<Link onClick={() => handler(METADATA.TERMS_OF_SERVICE)}>Terms of Service</Link>
			{/* @ts-ignore */}
			<Link onClick={() => handler(METADATA.PRIVACY_POLICY)}>Privacy Policy</Link>
			{renderThemeIcon(theme, dispatch)}
		</Bottom>
	)
};

const App = () => {

	return useMemo(
		() => (
			<Container>
				<Dragbar />
				<TopBar />
				<Row>
					<Routes />
					<Slideshow />
				</Row>
				<Footer />
			</Container>
		),
		[]
	);
};

const Container = styled.div`
  display: flex;
  flex-direction: column;
  flex: 1 0 0;
  min-width: 0;
  background-color: ${({ theme }) => theme.colors.background};
`;

const Float = styled.div`
  position: absolute;
  display: flex;
  top: 32px;
  left: 32px;
  width: calc(100% - 64px);
  z-index: 999;
  -webkit-app-region: drag;
`;

const Bottom = styled.div`
  position: absolute;
  display: flex;
  align-items: center;
  bottom: 32px;
  left: 32px;
  width: calc(100% - 64px);
  z-index: 999;
`;

const Row = styled.div`
  display: flex;
  flex: 1;
`;

const Col = styled.div`
  display: flex;
  flex-direction: column;
  background-color: ${({ theme }) => theme.colors.sidebar};
  flex: 1;
`;

const Flex = styled.div<{ margin: string; }>`
  display: flex;
  flex-direction: column;
  margin: ${({ margin }) => margin};
`;

const Link = styled(Typography.Link)`
  font-size: 12px;
  cursor: pointer;
  color: ${({ theme }) => theme.colors.primary};
  margin-right: 32px;
`;

const DayIcon = styled(Sun)`
  width: auto;
  height: 18px;
  fill: #fff;
  stroke: #fff;
`;

const NightIcon = styled(Moon)`
  width: auto;
  height: 18px;
  fill: #616161;
  stroke: #616161;
`;

export default App;