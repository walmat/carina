import Input from "./index";
import { Icon } from "react-feather";
import { useField } from "formik";

// it should be similar to InputProps
export interface InputFormikProps {
  id: string;
  name: string;
  type?: string;
  radius?: number;
  textTransform?: string;
  label?: string;
  useLabel?: boolean;
  onSubmit?: any;
  masked?: boolean;
  disabled?: boolean;
  autoFocus?: boolean;
  restriction?: string;
  placeholder: string;
  maxLength?: number;
  Icon?: Icon | null;
  onFocus?: (e: any) => void;
}

const InputFormik = (props: InputFormikProps) => {
  const [field, meta] = useField(props.name);

  return (
    <Input
      error={!!meta.error}
      touched={meta.touched}
      onBlur={field.onBlur}
      onChange={field.onChange}
      value={field.value}
      {...props}
    />
  );
};

export default InputFormik;
