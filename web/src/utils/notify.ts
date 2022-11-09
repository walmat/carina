import toast from "react-hot-toast";

type Type = "success" | "loading" | "error";

const defaultOptions = {
	duration: 5000,
	className: 'toasty-mctoastboy',
	position: 'bottom-right'
}

export const notify = (message: string, type: Type, options = {}) => {
	return toast[type](message, { ...defaultOptions, ...options });
};
