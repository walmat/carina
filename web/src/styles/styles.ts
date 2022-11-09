import { createGlobalStyle } from "styled-components";
import { StyledSystemProps } from "./dynamic";

const fontFamilies: { heading: string; body: string } = {
  heading: "Inter, serif",
  body: "Roboto, sans-serif",
};

export const GlobalStyles = createGlobalStyle`
  *:not(svg) {
    -webkit-touch-callout: none;
    -webkit-user-select: none;
    -khtml-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    user-select: none;
    text-rendering: auto;
    outline: none;
    scroll-behavior: smooth;
    font-family: Roboto !important;
    cursor: default;
  }
  
  // NOTE: Possible fix for the tooltip flicker
  .tooltip { pointer-events: none; }
  
  ::-webkit-scrollbar {
    background: transparent;
    display: block;
    cursor: grab;
    width: 8px;
    max-height: 50px;
  }
  
  ::-webkit-scrollbar-track {
    width: 8px;
    margin: 8px 0 4px 0;
    border-radius: 4px;
    background: transparent;
  }
  
  ::-webkit-scrollbar-button, ::-webkit-scrollbar-corner, ::-webkit-resizer {
    display: none;
  }
  
  ::-webkit-scrollbar-thumb {
    min-height: 64px;
    max-height: 64px;
    border-radius: 4px;
    background: ${({ theme }: any) => theme.colors.unscrollbar};
    
    &:hover {
      background: ${({ theme }: any) => theme.colors.scrollbar};
    }
  }
  
  h1, h2, h3, h4, h5, h6, a {
    cursor: default;
  }
  
  input, select, textarea {
    cursor: text;
    -webkit-touch-callout: text;
    -webkit-user-select: text;
    -khtml-user-select: text;
    -moz-user-select: text;
    -ms-user-select: text;
    user-select: text;
  }

  *:focus {
    outline: none;
  }
  
  textarea {
    overflow-x: hidden;
    resize: none;
  }
  
  webview {
    width: 100%;
    height: 100%;
  }
  
  g,
  path,
  circle,
  rect,
  button,
  img,
  a,
  svg,
  svg > * {
    text-decoration: none;
    cursor: pointer;
    -webkit-transition-duration: 0s;
    -moz-transition-duration: 0s;
    -o-transition-duration: 0s;
    transition-duration: 0s;
  }

  html, body {
    margin: 0;
    flex: 1;
    display: flex;
    border-radius: 8px;
    background-color: ${({ theme }) => theme.colors.background};
    flex-direction: column;
  }

  body,
  #root {
    border-radius: 8px;
    height: 100vh;
    width: 100vw;
    display: flex;
    background-color: ${({ theme }) => theme.colors.background};
    border-radius: 8px;
    overflow: hidden;
  }
  
  input[type="date"]::-webkit-inner-spin-button,
	input[type="date"]::-webkit-calendar-picker-indicator {
	    display: none;
	    -webkit-appearance: none;
	}

	// MARK: 2FA Code Styles
	.react-2fa-code {
	  display: flex;
	  flex-direction: column;
	  flex: 1;
	  
	  input {
	    background: transparent !important;
      border: 1px solid ${({ theme }) => theme.colors.altLogin} !important;
	    border-left: 0.5px solid ${({ theme }) => theme.colors.altLogin} !important;
	    border-right: 0.5px solid ${({ theme }) =>
        theme.colors.altLogin} !important;
	    color: ${({ theme }) => theme.colors.paragraph} !important;
	  }
	  
	  input:focus {
	    border: 1px solid ${({ theme }) => theme.colors.primary} !important;
	    caret-color: ${({ theme }) => theme.colors.primary} !important;
	  }
	  
	  input:focus + input {
      border: 1px solid ${({ theme }) => theme.colors.altLogin} !important;
	  }
	  
	  div:first-of-type {
	    display: flex !important;
	    align-items: center !important;
	    justify-content: center !important;
	  }
	  
	  input:first-child {
	    border-bottom-left-radius: 8px !important;
	    border-top-left-radius: 8px !important;
	  }
	  
	  input:last-child {
	    border-bottom-right-radius: 8px !important;
	    border-top-right-radius: 8px !important;
	  }
	  
	  input:nth-of-type(3) {
	    margin-right: 12px;
	  }
	  
	  input:nth-of-type(4) {
	    margin-left: 12px;
    }
  }

  div.rc-menu-container > ul > div:nth-child(2) > li > ul {
    max-width: 160px;
    padding: 8px;
  }
  
  div.rc-menu-container > ul > div:nth-child(2) > li > ul > div > li.rc-menu__item.rc-menu__item {
    padding-left: 0;
    border-radius: 2px;
  }
	
	.rc-menu {
	  width: 224px;
	  color: ${({ theme }) => theme.colors.paragraph};
    background: ${({ theme }) => theme.colors.sidebar};
    border: 1px solid ${({ theme }) => theme.colors.border};
    box-sizing: border-box;
    box-shadow: rgba(0,0,0,0.15) 4px 4px 4px -4px
        
    &__divider {
      background-color: ${({ theme }) => theme.colors.separator};
    }
	  
	  &__group {
	    margin-bottom: 16px;
	  }
	  
	  &__group:last-child {
	    margin-bottom: 0;
	  }
	  
	  &__item {
	    &--focusable {
	      padding: 0 0 8px 0;
	      
	      &:hover {
	        background-color: transparent;
	      }
	    }
	  
      &:hover > p {
        color: ${({ theme }) => theme.colors.h2};
      }
      
      &:hover > svg {
        color: ${({ theme }) => theme.colors.h2};
      }

	    &--hover {
        background-color: ${({ theme }) => theme.colors.secondary};
	    }
	  }
	  
	  .rc-menu__submenu > .rc-menu__item::after {
	    color: ${({ theme }) => theme.colors.primary};
	    transform: scale(0.75);
	    right: 1.5rem;
	  }
	  
	  .rc-menu__submenu > .rc-menu__item {
      padding: 0.585rem 1.5rem;
	  }
	  
	  &__header {
	    text-transform: capitalize;
	    color: ${({ theme }) => theme.colors.h2};
	    padding: 0.2rem 16px;
	  }
	}
	
	.toasty-mctoastboy {
	  border: 1px solid ${({ theme }) => theme.colors.border};
    color: ${({ theme }) => theme.colors.h2};
	  background-color: ${({ theme }) => theme.colors.sidebar};
	}
`;

interface Typography {
  H1: StyledSystemProps;
  H2: StyledSystemProps;
  H3: StyledSystemProps;
  H4: StyledSystemProps;
  H5: StyledSystemProps;
  Paragraph: StyledSystemProps;
  Link: StyledSystemProps;
}
export const typography: Typography = {
  H1: {
    fontSize: [50, 51, 52, 57],
    fontWeight: 700,
    fontFamily: fontFamilies.heading,
    as: "h1",
  },
  H2: {
    fontSize: [37, 39, 41, 43],
    fontWeight: 700,
    fontFamily: fontFamilies.heading,
    as: "h2",
  },
  H3: {
    fontSize: [27, 28, 30, 32],
    fontWeight: 700,
    fontFamily: fontFamilies.heading,
    as: "h3",
  },
  H4: {
    fontSize: [18, 20, 22, 28],
    fontWeight: 700,
    fontFamily: fontFamilies.heading,
    as: "h4",
  },
  H5: {
    fontWeight: 700,
    fontSize: [16, 17, 19, 21],
    fontFamily: fontFamilies.heading,
    as: "h5",
  },
  Paragraph: {
    fontSize: [14, 15, 15, 16],
    fontWeight: 400,
    fontFamily: fontFamilies.body,
    as: "p",
  },
  Link: {
    fontWeight: 700,
    fontSize: [14, 15, 15, 16],
    fontFamily: fontFamilies.body,
  },
};
