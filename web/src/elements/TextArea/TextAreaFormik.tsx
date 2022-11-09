import TextArea from "./index";
import { useField } from "formik";

// it should be similar to TextAreaProps
export interface TextAreaFormikProps {
  id: string;
  name: string;
  placeholder: string;
  autoFocus?: boolean;
  onChange?: (e: any) => void;
  onFocus?: (e: any) => void;
  onBlur?: (e: any) => void;
}

const TextAreaFormik = (props: TextAreaFormikProps) => {
  const [field, meta] = useField(props.name);

  return (
    <TextArea
      error={!!meta.error}
      touched={meta.touched}
      onBlur={field.onBlur}
      onChange={field.onChange}
      value={field.value}
      {...props}
    />
  );
};

export default TextAreaFormik;
