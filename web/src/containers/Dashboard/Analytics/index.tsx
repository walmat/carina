import React, { useState } from 'react';
import { useSelector } from "react-redux";
import styled from 'styled-components';
import { Folder } from 'react-feather';
import moment from 'moment';
import { ResponsiveContainer, AreaChart, CartesianGrid, XAxis, Tooltip, Area } from 'recharts';

import { DatePicker } from '../../../components';
import { Typography, DateSelect } from '../../../elements';
import { makeTheme } from "../../../stores/Main/reducers/theme";

import Pill from './pill';

const buildXAxis = (timeframe: any, checkouts: any[]) => {
	const { id, value } = timeframe;

	switch (id) {
		default:
		case 'day':
			return [moment(value.start).format('MM/DD'), '4 am', '8 am', 'Noon', '4 pm', '8 pm', moment(value.end).add(1, 'day').format('MM/DD')];
		case 'week':
			return [...Array(7)].map((_, i) => moment(value.start).add(i, 'days').format('MM/DD'));
		case 'month':
			return [...Array(7)].map((_, i) => moment(value.start).add(i > 0 ? (i * 5) - 1 : i, 'days').format('MM/DD'));
		case '3mo':
			return [...Array(7)].map((_, i) => moment(value.start).add(i > 0 ? (i * 15) - 1 : i, 'days').format('MM/DD'));
		case 'year':
			return [...Array(12)].map((_, i) => {
				const next = moment(value.start).add(i, 'months');
				return next.format("MMM 'YY");
			});
		case 'all':
			// TODO: How do I do this?
			// Maybe grab moment() of first checkout and split it into 7 equal parts?
			return [];
		case 'custom':
			// TODO: Custom start/end dates can be chosen using the datepicker.
			// We will have to determine a best fit x-axis
			// NOTE: Possible fall back to the above ranges depending on how long the chosen range is
			// For example, if they choose April 7th - April 14th, use `week`
			// We will just need to account for edge cases (e.g. - 2 day range - April 7th - April 9th)

			const period = value.end.diff(value.start);
			const formatter = () => period < 518400000 ? 'MM/DD h:mm A' : 'MM/DD';
			const step = (i: number) => (i * period) / 6;

			const middle = [...Array(5)].map((_, i) => moment(value.start).add(step(i + 1), 'ms').format(formatter()));

			return [
				value.start.format('MM/DD'),
				...middle,
				value.end.format('MM/DD')
			];
	}
};

const Analytics = () => {
	const theme = useSelector(makeTheme);
	// @ts-ignore
	const [timeframe, setTimeframe] = useState<any>(DatePicker.DATE_PICKER_OPTIONS[0]);
	// @ts-ignore
	const [comparison, setComparison] = useState<any>(DatePicker.COMPARISON_OPTIONS[0]);

	//these values need to be fetched from gql
	const [total, setTotal] = useState(0);
	const [previous, setPrevious] = useState(0);

	const axis = buildXAxis(timeframe, []);

	const data = axis.map(x => ({ name: x, amount: 0 }))

	return (
		<Container>
			<Row useBorder m="" p="0 16px" height="48px">
				<Title>Analytics</Title>
				<Pill fetch={() => { }} />

				<Selects>
					<DatePicker timeframe={timeframe} setTimeframe={setTimeframe} />
					<ComparisonText>compared to</ComparisonText>
					<Shadow>
						<DateSelect
							required
							type="full"
							isClearable={false}
							isMulti={false}
							placeholder="None"
							name="comparison"
							// @ts-ignore
							options={DatePicker.COMPARISON_OPTIONS}
							onChange={setComparison}
							value={comparison}
						/>
					</Shadow>
					<Save />
				</Selects>
			</Row>
			<Col>
				<AmountSpent>Amount Spent</AmountSpent>
				<AmountContainer>
					<Amount>${total}</Amount>
					{comparison.value.start ? <ComparisonAmount>${previous}</ComparisonAmount> : null}
				</AmountContainer>
				<RowHeight>
					<ResponsiveContainer width="100%" height="100%">
						<AreaChart
							data={data}
							throttleDelay={100}
							margin={{
								top: 16,
								right: 32,
								left: 32,
								bottom: 0,
							}}
						>
							<CartesianGrid stroke={theme !== 0 ? '#616161' : '#d8d8d8'} strokeDasharray="8 8" />
							<XAxis dataKey="name" fontSize={10} stroke={theme === 0 ? '#616161' : '#d8d8d8'} axisLine={{ stroke: 'none' }} tickLine={{ stroke: theme !== 0 ? '#616161' : '#d8d8d8' }} />
							<Tooltip />
							<Area type="monotone" dataKey="amount" stroke="#786EF2" fill="rgba(120, 110, 242, 0.2)" activeDot={{ r: 4 }} />
						</AreaChart>
					</ResponsiveContainer>
				</RowHeight>
			</Col>
		</Container>
	);
};

const Container = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  border-radius: ${({ theme }) => theme.space.M}px;
  padding: ${({ theme }) => theme.space.NONE}px;
  background-color: ${({ theme }) => theme.colors.sidebar};
`;

const Col = styled.div`
	display: flex;
	flex: 1;
	height: 48px;
	flex-direction: column;
`;

const Row = styled.div<{ m: string; p: string; height?: string; useBorder?: boolean }>`
	display: flex;
	align-items: center;
	margin: ${({ m }) => m};
	padding: ${({ p }) => p};
	${({ height }) => height ? `height: ${height}` : ''};
	${({ theme, useBorder }) => useBorder ? `border-bottom: 1px solid ${theme.colors.separator}` : ''};
`;

const Shadow = styled.div`
	display: flex;
	box-shadow: 0px 1px 2px rgba(0, 0, 0, 0.1);
`;

const RowHeight = styled.div`
	display: flex;
	height: calc(100% - 52px);
	width: 100%;
`;

const Title = styled(Typography.H4)`
  font-size: 1.35vw;
  font-weight: 500;
  color: ${({ theme }) => theme.colors.h2};
  margin: 0 8px 0 0;
`;

const Selects = styled.div`
	margin-left: auto;
	display: flex;
`;

const Save = styled(Folder)`
	margin: auto 0 auto 16px;
	width: 16px;
	color: ${({ theme }) => theme.colors.paragraph};
	
	&:hover {
		color: hsl(0, 0%, 70%);
	}
	
	&:active {
    color: ${({ theme }) => theme.colors.h2};
	}
`;

const ComparisonText = styled(Typography.Paragraph)`
	display: flex;
	align-items: center;
	color: ${({ theme }) => theme.colors.paragraph};
	font-size: 12px;
	margin: 0 8px;
	font-weight: 400;
`;

const AmountSpent = styled(Typography.H4)`
	display: flex;
	margin: 16px 16px 4px 16px;
	color: ${({ theme }) => theme.colors.h2};
	font-size: 12px;
	font-weight: 400;
`;

const AmountContainer = styled.div`
	display: flex;
`;

const Amount = styled(Typography.H4)`
	display: flex;
	flex-direction: column;
	margin: 2px 0 0 16px;
	color: ${({ theme }) => theme.colors.primary};
	font-size: 12px;
	font-weight: 500;
`;

const ComparisonAmount = styled(Typography.H4)`
	display: flex;
	flex-direction: column;
	margin: 2px 16px 0 auto;
	color: ${({ theme }) => theme.colors.paragraph};
	opacity: 0.6;
	font-size: 12px;
	font-weight: 400;
`;


export default Analytics;
