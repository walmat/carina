import NumberFormatter from "./index";
import { Icon } from "react-feather";
import { useField } from "formik";

// it should be similar to NumberFormatterProps
export interface NumberFormatterFormikProps {
  id: string;
  name: string;
  radius?: number;
  format: string;
  mask?: string[];
  removeFormatting?: any;
  textTransform?: string;
  useLabel?: boolean;
  value?: string | number;
  error?: boolean;
  touched?: boolean;
  disabled?: boolean;
  autoFocus?: boolean;
  placeholder: string;
  maxLength?: number;
  Icon?: Icon | null;
  onChange?: (e: any) => void;
  onBlur?: (e: any) => void;
  onFocus?: (e: any) => void;
}

const NumberFormatterFormik = (props: NumberFormatterFormikProps) => {
  const [field, meta] = useField(props.name);

  // only show error if there is an error and user touched the field
  const error = !!meta.error && !!meta.touched;

  return (
    <NumberFormatter
      {...props}
      error={error}
      touched={meta.touched}
      onBlur={field.onBlur}
      onChange={field.onChange}
      value={field.value}
    />
  );
};

export default NumberFormatterFormik;
