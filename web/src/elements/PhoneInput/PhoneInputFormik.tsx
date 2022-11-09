import PhoneInput from "./index";
import { useField } from "formik";

// it should be similar to InputProps
export interface PhoneInputFormikProps {
  id: string;
  name: string;
  radius?: number;
  useLabel?: boolean;
  value?: string;
  country: string;
  error?: boolean;
  touched?: boolean;
  disabled?: boolean;
  autoFocus?: boolean;
  placeholder: string;
  maxLength?: number;
  onChange?: (e: any) => void;
  onBlur?: (e: any) => void;
  onFocus?: (e: any) => void;
}

const PhoneInputFormik = (props: PhoneInputFormikProps) => {
  const [field, meta] = useField(props.name);

  // only show error if there is an error and user touched the field
  const error = !!meta.error && !!meta.touched;

  return (
    <PhoneInput
      {...props}
      error={error}
      touched={meta.touched}
      onBlur={field.onBlur}
      onChange={field.onChange}
      value={field.value}
    />
  );
};

export default PhoneInputFormik;
