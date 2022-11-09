import React, { PureComponent } from "react";

interface CurrencyConverterProps {
  from: string;
  to: string;
  language: string;
  value: string | number;
  date?: string;
  precision?: number;
}

export default class CurrencyConverter extends PureComponent<any, any> {
  constructor(props: CurrencyConverterProps) {
    super(props);

    this.state = {
      from: props.from.toUpperCase(),
      to: props.to.toUpperCase(),
      language: props.language,
      value: Number(props.value),
      date: props.date ? props.date : "latest",
      convertedValue: props.value,
      precision: props.precision
        ? props.precision > 0
          ? props.precision
          : 2
        : 2,
    };
  }

  componentDidMount() {
    const codes = [
      "CAD",
      "HKD",
      "ISK",
      "PHP",
      "DKK",
      "HUF",
      "CZK",
      "AUD",
      "RON",
      "SEK",
      "IDR",
      "INR",
      "BRL",
      "RUB",
      "HRK",
      "JPY",
      "THB",
      "CHF",
      "SGD",
      "PLN",
      "BGN",
      "TRY",
      "CNY",
      "NOK",
      "NZD",
      "ZAR",
      "USD",
      "MXN",
      "ILS",
      "GBP",
      "KRW",
      "MYR",
      "EUR",
    ];
    if (!(codes.includes(this.state.from) && codes.includes(this.state.to))) {
      throw new Error(
        `Country code is not supprted, supported country codes are: ${codes}`
      );
    } else if (typeof this.state.value !== "number") {
      throw new Error(`Input value of exchange is not of type number`);
    } else {
      fetch(
        `https://api.exchangerate.host/${this.state.date}?base=${this.state.from}`
      )
        .then((data) => data.json())
        .then((res) => {
          if (res.error) {
            return;
          } else {
            let value = this.state.value * res.rates[this.state.to];
            this.setState({
              convertedValue: value.toFixed(this.state.precision),
            });
          }
        })
        .catch(() => {});
    }
  }

  render() {
    return (
      <div>
        {new Intl.NumberFormat(this.state.language, {
          style: "currency",
          currency: this.state.to,
          minimumFractionDigits: this.state.precision,
          maximumFractionDigits: this.state.precision,
        }).format(this.state.convertedValue)}
      </div>
    );
  }
}
