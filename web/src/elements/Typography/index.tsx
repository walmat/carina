import React from "react";
import { Link as RouterLink, LinkProps } from "react-router-dom";
import DynamicStyledSystemComponent, {
  StyledSystemProps,
} from "../../styles/dynamic";
import { typography } from "../../styles";

type AnchorProps = StyledSystemProps &
  Pick<LinkProps, "to"> & {
    onClick?: (event: React.MouseEvent<HTMLAnchorElement>) => void;
  };

const Link: React.FC<AnchorProps> = ({ to, onClick, children, ...props }) => (
  <RouterLink to={to} onClick={onClick}>
    <DynamicStyledSystemComponent {...typography.Link} {...props}>
      {children}
    </DynamicStyledSystemComponent>
  </RouterLink>
);

interface TypographyComponentProps {
  H1: React.FC<StyledSystemProps>;
  H2: React.FC<StyledSystemProps>;
  H3: React.FC<StyledSystemProps>;
  H4: React.FC<StyledSystemProps>;
  H5: React.FC<StyledSystemProps>;
  Paragraph: React.FC<StyledSystemProps>;
  Link: React.FC<AnchorProps>;
}

const createComponent: (
  textStyle: StyledSystemProps,
  displayName: string
) => React.FC<StyledSystemProps> = (textStyle, displayName) => {
  const component: React.FC<StyledSystemProps> = (props) => (
    <DynamicStyledSystemComponent {...textStyle} {...props}>
      {props.children}
    </DynamicStyledSystemComponent>
  );
  component.displayName = displayName;
  return component;
};

const Typography: TypographyComponentProps = {
  H1: createComponent(typography.H1, "H1"),
  H2: createComponent(typography.H2, "H2"),
  H3: createComponent(typography.H3, "H3"),
  H4: createComponent(typography.H4, "H4"),
  H5: createComponent(typography.H5, "H5"),
  Paragraph: createComponent(typography.Paragraph, "Paragraph"),
  Link: Link,
};

export { Typography };
