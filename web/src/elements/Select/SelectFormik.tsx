import { useField } from "formik";
import Select from "./index";

interface SelectFormikProps {
  useDefault?: boolean;
  autoFocus?: boolean;
  name: string;
  disabled?: boolean;
  required?: boolean;
  isCreatable?: boolean;
  isMulti?: boolean;
  isClearable?: boolean;
  Icon?: any;
  label?: string;
  closeMenuOnSelect?: boolean;
  isOptionDisabled?: (option: any) => boolean;
  noOptionsMessage?: () => undefined;
  error?: boolean;
  touched?: boolean;
  placeholder: string;
  components?: any;
  value?: any;
  onChange?: any;
  onFocus?: any;
  onBlur?: any;
  options: any[];
  getOptionLabel?: any;
  getOptionValue?: any;
  size?: "small" | "default";
}
const SelectFormik = (props: SelectFormikProps) => {
  const [field, meta, helpers] = useField(props.name);

  return (
    <Select
      value={field.value}
      onBlur={field.onBlur}
      onChange={(event: any) => {
        helpers.setValue(event);
      }}
      {...props}
    />
  );
};

export default SelectFormik;
