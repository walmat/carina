import { join } from "path";
import { readFileSync } from "fs";

let _pkginfo = {};

// // eslint-disable-next-line no-undef
// if (typeof _pkginfo !== 'undefined' && _pkginfo !== null) {
//   // eslint-disable-next-line no-undef
//   _pkginfo = PKG_INFO;
// } else {
//   // TODO: Not sure how to get cwd/__dirname in golang
//   _pkginfo = JSON.parse(
//     readFileSync(join('../', 'package.json'), { encoding: 'utf8' })
//   );
// }

export const pkginfo = _pkginfo;
