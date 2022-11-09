import { useRef, useEffect } from 'react';

const useScroll = (callBack: any) => {
	const ref = useRef(null);

	useEffect(() => {
		// @ts-ignore
		if (ref?.current?.addEventListener && callBack) {
			// @ts-ignore
			ref.current.addEventListener('scroll', callBack);
		}

		return () => {
			// @ts-ignore
			if (ref.current?.removeEventListener && callBack) {
				// @ts-ignore
				ref.current.removeEventListener('scroll', callBack);
			}
		};
	}, [ref, callBack]);

	return ref;
};

export { useScroll };