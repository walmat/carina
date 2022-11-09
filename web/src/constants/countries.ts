const countries = [
  {
    name: "Afghanistan",
    code: "AF",
    provinces: [],
  },
  {
    name: "Aland Islands",
    code: "AX",
    provinces: [],
  },
  {
    name: "Albania",
    code: "AL",
    provinces: [],
  },
  {
    name: "Algeria",
    code: "DZ",
    provinces: [],
  },
  {
    name: "Andorra",
    code: "AD",
    provinces: [],
  },
  {
    name: "Angola",
    code: "AO",
    provinces: [],
  },
  {
    name: "Anguilla",
    code: "AI",
    provinces: [],
  },
  {
    name: "Antigua And Barbuda",
    code: "AG",
    provinces: [],
  },
  {
    name: "Argentina",
    code: "AR",
    provinces: [
      {
        name: "Buenos Aires",
        code: "B",
      },
      {
        name: "Catamarca",
        code: "K",
      },
      {
        name: "Chaco",
        code: "H",
      },
      {
        name: "Chubut",
        code: "U",
      },
      {
        name: "Ciudad AutÃ³noma de Buenos Aires",
        code: "C",
      },
      {
        name: "CÃ³rdoba",
        code: "X",
      },
      {
        name: "Corrientes",
        code: "W",
      },
      {
        name: "Entre RÃ­os",
        code: "E",
      },
      {
        name: "Formosa",
        code: "P",
      },
      {
        name: "Jujuy",
        code: "Y",
      },
      {
        name: "La Pampa",
        code: "L",
      },
      {
        name: "La Rioja",
        code: "F",
      },
      {
        name: "Mendoza",
        code: "M",
      },
      {
        name: "Misiones",
        code: "N",
      },
      {
        name: "NeuquÃ©n",
        code: "Q",
      },
      {
        name: "RÃ­o Negro",
        code: "R",
      },
      {
        name: "Salta",
        code: "A",
      },
      {
        name: "San Juan",
        code: "J",
      },
      {
        name: "San Luis",
        code: "D",
      },
      {
        name: "Santa Cruz",
        code: "Z",
      },
      {
        name: "Santa Fe",
        code: "S",
      },
      {
        name: "Santiago Del Estero",
        code: "G",
      },
      {
        name: "Tierra Del Fuego",
        code: "V",
      },
      {
        name: "TucumÃ¡n",
        code: "T",
      },
    ],
  },
  {
    name: "Armenia",
    code: "AM",
    provinces: [],
  },
  {
    name: "Aruba",
    code: "AW",
    provinces: [],
  },
  {
    name: "Australia",
    code: "AU",
    provinces: [
      {
        name: "Australian Capital Territory",
        code: "ACT",
      },
      {
        name: "New South Wales",
        code: "NSW",
      },
      {
        name: "Northern Territory",
        code: "NT",
      },
      {
        name: "Queensland",
        code: "QLD",
      },
      {
        name: "South Australia",
        code: "SA",
      },
      {
        name: "Tasmania",
        code: "TAS",
      },
      {
        name: "Victoria",
        code: "VIC",
      },
      {
        name: "Western Australia",
        code: "WA",
      },
    ],
  },
  {
    name: "Austria",
    code: "AT",
    provinces: [],
  },
  {
    name: "Azerbaijan",
    code: "AZ",
    provinces: [],
  },
  {
    name: "Bahamas",
    code: "BS",
    provinces: [],
  },
  {
    name: "Bahrain",
    code: "BH",
    provinces: [],
  },
  {
    name: "Bangladesh",
    code: "BD",
    provinces: [],
  },
  {
    name: "Barbados",
    code: "BB",
    provinces: [],
  },
  {
    name: "Belarus",
    code: "BY",
    provinces: [],
  },
  {
    name: "Belgium",
    code: "BE",
    provinces: [],
  },
  {
    name: "Belize",
    code: "BZ",
    provinces: [],
  },
  {
    name: "Benin",
    code: "BJ",
    provinces: [],
  },
  {
    name: "Bermuda",
    code: "BM",
    provinces: [],
  },
  {
    name: "Bhutan",
    code: "BT",
    provinces: [],
  },
  {
    name: "Bolivia",
    code: "BO",
    provinces: [],
  },
  {
    name: "Bonaire, Sint Eustatius and Saba",
    code: "BQ",
    provinces: [],
  },
  {
    name: "Bosnia And Herzegovina",
    code: "BA",
    provinces: [],
  },
  {
    name: "Botswana",
    code: "BW",
    provinces: [],
  },
  {
    name: "Bouvet Island",
    code: "BV",
    provinces: [],
  },
  {
    name: "Brazil",
    code: "BR",
    provinces: [
      {
        name: "Acre",
        code: "AC",
      },
      {
        name: "Alagoas",
        code: "AL",
      },
      {
        name: "AmapÃ¡",
        code: "AP",
      },
      {
        name: "Amazonas",
        code: "AM",
      },
      {
        name: "Bahia",
        code: "BA",
      },
      {
        name: "CearÃ¡",
        code: "CE",
      },
      {
        name: "Distrito Federal",
        code: "DF",
      },
      {
        name: "EspÃ­rito Santo",
        code: "ES",
      },
      {
        name: "GoiÃ¡s",
        code: "GO",
      },
      {
        name: "MaranhÃ£o",
        code: "MA",
      },
      {
        name: "Mato Grosso",
        code: "MT",
      },
      {
        name: "Mato Grosso do Sul",
        code: "MS",
      },
      {
        name: "Minas Gerais",
        code: "MG",
      },
      {
        name: "ParÃ¡",
        code: "PA",
      },
      {
        name: "ParaÃ­ba",
        code: "PB",
      },
      {
        name: "ParanÃ¡",
        code: "PR",
      },
      {
        name: "Pernambuco",
        code: "PE",
      },
      {
        name: "PiauÃ­",
        code: "PI",
      },
      {
        name: "Rio de Janeiro",
        code: "RJ",
      },
      {
        name: "Rio Grande do Norte",
        code: "RN",
      },
      {
        name: "Rio Grande do Sul",
        code: "RS",
      },
      {
        name: "RondÃ´nia",
        code: "RO",
      },
      {
        name: "Roraima",
        code: "RR",
      },
      {
        name: "Santa Catarina",
        code: "SC",
      },
      {
        name: "SÃ£o Paulo",
        code: "SP",
      },
      {
        name: "Sergipe",
        code: "SE",
      },
      {
        name: "Tocantins",
        code: "TO",
      },
    ],
  },
  {
    name: "British Indian Ocean Territory",
    code: "IO",
    provinces: [],
  },
  {
    name: "Brunei",
    code: "BN",
    provinces: [],
  },
  {
    name: "Bulgaria",
    code: "BG",
    provinces: [],
  },
  {
    name: "Burkina Faso",
    code: "BF",
    provinces: [],
  },
  {
    name: "Burundi",
    code: "BI",
    provinces: [],
  },
  {
    name: "Cambodia",
    code: "KH",
    provinces: [],
  },
  {
    name: "Canada",
    code: "CA",
    provinces: [
      {
        name: "Alberta",
        code: "AB",
      },
      {
        name: "British Columbia",
        code: "BC",
      },
      {
        name: "Manitoba",
        code: "MB",
      },
      {
        name: "New Brunswick",
        code: "NB",
      },
      {
        name: "Newfoundland",
        code: "NL",
      },
      {
        name: "Northwest Territories",
        code: "NT",
      },
      {
        name: "Nova Scotia",
        code: "NS",
      },
      {
        name: "Nunavut",
        code: "NU",
      },
      {
        name: "Ontario",
        code: "ON",
      },
      {
        name: "Prince Edward Island",
        code: "PE",
      },
      {
        name: "Quebec",
        code: "QC",
      },
      {
        name: "Saskatchewan",
        code: "SK",
      },
      {
        name: "Yukon",
        code: "YT",
      },
    ],
  },
  {
    name: "Cape Verde",
    code: "CV",
    provinces: [],
  },
  {
    name: "Cayman Islands",
    code: "KY",
    provinces: [],
  },
  {
    name: "Central African Republic",
    code: "CF",
    provinces: [],
  },
  {
    name: "Chad",
    code: "TD",
    provinces: [],
  },
  {
    name: "Chile",
    code: "CL",
    provinces: [],
  },
  {
    name: "China",
    code: "CN",
    provinces: [
      {
        name: "Anhui",
        code: "AH",
      },
      {
        name: "Beijing",
        code: "BJ",
      },
      {
        name: "Chongqing",
        code: "CQ",
      },
      {
        name: "Fujian",
        code: "FJ",
      },
      {
        name: "Gansu",
        code: "GS",
      },
      {
        name: "Guangdong",
        code: "GD",
      },
      {
        name: "Guangxi",
        code: "GX",
      },
      {
        name: "Guizhou",
        code: "GZ",
      },
      {
        name: "Hainan",
        code: "HI",
      },
      {
        name: "Hebei",
        code: "HE",
      },
      {
        name: "Heilongjiang",
        code: "HL",
      },
      {
        name: "Henan",
        code: "HA",
      },
      {
        name: "Hubei",
        code: "HB",
      },
      {
        name: "Hunan",
        code: "HN",
      },
      {
        name: "Inner Mongolia",
        code: "NM",
      },
      {
        name: "Jiangsu",
        code: "JS",
      },
      {
        name: "Jiangxi",
        code: "JX",
      },
      {
        name: "Jilin",
        code: "JL",
      },
      {
        name: "Liaoning",
        code: "LN",
      },
      {
        name: "Ningxia",
        code: "NX",
      },
      {
        name: "Qinghai",
        code: "QH",
      },
      {
        name: "Shaanxi",
        code: "SN",
      },
      {
        name: "Shandong",
        code: "SD",
      },
      {
        name: "Shanghai",
        code: "SH",
      },
      {
        name: "Shanxi",
        code: "SX",
      },
      {
        name: "Sichuan",
        code: "SC",
      },
      {
        name: "Tianjin",
        code: "TJ",
      },
      {
        name: "Xinjiang",
        code: "XJ",
      },
      {
        name: "Xizang",
        code: "YZ",
      },
      {
        name: "Yunnan",
        code: "YN",
      },
      {
        name: "Zhejiang",
        code: "ZJ",
      },
    ],
  },
  {
    name: "Christmas Island",
    code: "CX",
    provinces: [],
  },
  {
    name: "Cocos (Keeling) Islands",
    code: "CC",
    provinces: [],
  },
  {
    name: "Colombia",
    code: "CO",
    provinces: [
      {
        name: "Amazonas",
        code: "AMA",
      },
      {
        name: "Antioquia",
        code: "ANT",
      },
      {
        name: "Arauca",
        code: "ARA",
      },
      {
        name: "AtlÃ¡ntico",
        code: "ATL",
      },
      {
        name: "BogotÃ¡, D.C.",
        code: "DC",
      },
      {
        name: "BolÃ­var",
        code: "BOL",
      },
      {
        name: "BoyacÃ¡",
        code: "BOY",
      },
      {
        name: "Caldas",
        code: "CAL",
      },
      {
        name: "CaquetÃ¡",
        code: "CAQ",
      },
      {
        name: "Casanare",
        code: "CAS",
      },
      {
        name: "Cauca",
        code: "CAU",
      },
      {
        name: "Cesar",
        code: "CES",
      },
      {
        name: "ChocÃ³",
        code: "CHO",
      },
      {
        name: "CÃ³rdoba",
        code: "COR",
      },
      {
        name: "Cundinamarca",
        code: "CUN",
      },
      {
        name: "GuainÃ­a",
        code: "GUA",
      },
      {
        name: "Guaviare",
        code: "GUV",
      },
      {
        name: "Huila",
        code: "HUI",
      },
      {
        name: "La Guajira",
        code: "LAG",
      },
      {
        name: "Magdalena",
        code: "MAG",
      },
      {
        name: "Meta",
        code: "MET",
      },
      {
        name: "NariÃ±o",
        code: "NAR",
      },
      {
        name: "Norte de Santander",
        code: "NSA",
      },
      {
        name: "Putumayo",
        code: "PUT",
      },
      {
        name: "QuindÃ­o",
        code: "QUI",
      },
      {
        name: "Risaralda",
        code: "RIS",
      },
      {
        name: "San AndrÃ©s, Providencia y Santa Catalina",
        code: "SAP",
      },
      {
        name: "Santander",
        code: "SAN",
      },
      {
        name: "Sucre",
        code: "SUC",
      },
      {
        name: "Tolima",
        code: "TOL",
      },
      {
        name: "Valle del Cauca",
        code: "VAC",
      },
      {
        name: "VaupÃ©s",
        code: "VAU",
      },
      {
        name: "Vichada",
        code: "VID",
      },
    ],
  },
  {
    name: "Comoros",
    code: "KM",
    provinces: [],
  },
  {
    name: "Congo",
    code: "CG",
    provinces: [],
  },
  {
    name: "Congo, The Democratic Republic Of The",
    code: "CD",
    provinces: [],
  },
  {
    name: "Cook Islands",
    code: "CK",
    provinces: [],
  },
  {
    name: "Costa Rica",
    code: "CR",
    provinces: [],
  },
  {
    name: "CÃ´te d'Ivoire",
    code: "CI",
    provinces: [],
  },
  {
    name: "Croatia",
    code: "HR",
    provinces: [],
  },
  {
    name: "Cuba",
    code: "CU",
    provinces: [],
  },
  {
    name: "CuraÃ§ao",
    code: "CW",
    provinces: [],
  },
  {
    name: "Cyprus",
    code: "CY",
    provinces: [],
  },
  {
    name: "Czech Republic",
    code: "CZ",
    provinces: [],
  },
  {
    name: "Denmark",
    code: "DK",
    provinces: [],
  },
  {
    name: "Djibouti",
    code: "DJ",
    provinces: [],
  },
  {
    name: "Dominica",
    code: "DM",
    provinces: [],
  },
  {
    name: "Dominican Republic",
    code: "DO",
    provinces: [],
  },
  {
    name: "Ecuador",
    code: "EC",
    provinces: [],
  },
  {
    name: "Egypt",
    code: "EG",
    provinces: [
      {
        name: "6th of October",
        code: "SU",
      },
      {
        name: "Al Sharqia",
        code: "SHR",
      },
      {
        name: "Alexandria",
        code: "ALX",
      },
      {
        name: "Aswan",
        code: "ASN",
      },
      {
        name: "Asyut",
        code: "AST",
      },
      {
        name: "Beheira",
        code: "BH",
      },
      {
        name: "Beni Suef",
        code: "BNS",
      },
      {
        name: "Cairo",
        code: "C",
      },
      {
        name: "Dakahlia",
        code: "DK",
      },
      {
        name: "Damietta",
        code: "DT",
      },
      {
        name: "Faiyum",
        code: "FYM",
      },
      {
        name: "Gharbia",
        code: "GH",
      },
      {
        name: "Giza",
        code: "GZ",
      },
      {
        name: "Helwan",
        code: "HU",
      },
      {
        name: "Ismailia",
        code: "IS",
      },
      {
        name: "Kafr el-Sheikh",
        code: "KFS",
      },
      {
        name: "Luxor",
        code: "LX",
      },
      {
        name: "Matrouh",
        code: "MT",
      },
      {
        name: "Minya",
        code: "MN",
      },
      {
        name: "Monufia",
        code: "MNF",
      },
      {
        name: "New Valley",
        code: "WAD",
      },
      {
        name: "North Sinai",
        code: "SIN",
      },
      {
        name: "Port Said",
        code: "PTS",
      },
      {
        name: "Qalyubia",
        code: "KB",
      },
      {
        name: "Qena",
        code: "KN",
      },
      {
        name: "Red Sea",
        code: "BA",
      },
      {
        name: "Sohag",
        code: "SHG",
      },
      {
        name: "South Sinai",
        code: "JS",
      },
      {
        name: "Suez",
        code: "SUZ",
      },
    ],
  },
  {
    name: "El Salvador",
    code: "SV",
    provinces: [],
  },
  {
    name: "Equatorial Guinea",
    code: "GQ",
    provinces: [],
  },
  {
    name: "Eritrea",
    code: "ER",
    provinces: [],
  },
  {
    name: "Estonia",
    code: "EE",
    provinces: [],
  },
  {
    name: "Ethiopia",
    code: "ET",
    provinces: [],
  },
  {
    name: "Falkland Islands (Malvinas)",
    code: "FK",
    provinces: [],
  },
  {
    name: "Faroe Islands",
    code: "FO",
    provinces: [],
  },
  {
    name: "Fiji",
    code: "FJ",
    provinces: [],
  },
  {
    name: "Finland",
    code: "FI",
    provinces: [],
  },
  {
    name: "France",
    code: "FR",
    provinces: [],
  },
  {
    name: "French Guiana",
    code: "GF",
    provinces: [],
  },
  {
    name: "French Polynesia",
    code: "PF",
    provinces: [],
  },
  {
    name: "French Southern Territories",
    code: "TF",
    provinces: [],
  },
  {
    name: "Gabon",
    code: "GA",
    provinces: [],
  },
  {
    name: "Gambia",
    code: "GM",
    provinces: [],
  },
  {
    name: "Georgia",
    code: "GE",
    provinces: [],
  },
  {
    name: "Germany",
    code: "DE",
    provinces: [],
  },
  {
    name: "Ghana",
    code: "GH",
    provinces: [],
  },
  {
    name: "Gibraltar",
    code: "GI",
    provinces: [],
  },
  {
    name: "Greece",
    code: "GR",
    provinces: [],
  },
  {
    name: "Greenland",
    code: "GL",
    provinces: [],
  },
  {
    name: "Grenada",
    code: "GD",
    provinces: [],
  },
  {
    name: "Guadeloupe",
    code: "GP",
    provinces: [],
  },
  {
    name: "Guatemala",
    code: "GT",
    provinces: [
      {
        name: "Alta Verapaz",
        code: "AVE",
      },
      {
        name: "Baja Verapaz",
        code: "BVE",
      },
      {
        name: "Chimaltenango",
        code: "CMT",
      },
      {
        name: "Chiquimula",
        code: "CQM",
      },
      {
        name: "El Progreso",
        code: "EPR",
      },
      {
        name: "Escuintla",
        code: "ESC",
      },
      {
        name: "Guatemala",
        code: "GUA",
      },
      {
        name: "Huehuetenango",
        code: "HUE",
      },
      {
        name: "Izabal",
        code: "IZA",
      },
      {
        name: "Jalapa",
        code: "JAL",
      },
      {
        name: "Jutiapa",
        code: "JUT",
      },
      {
        name: "PetÃ©n",
        code: "PET",
      },
      {
        name: "Quetzaltenango",
        code: "QUE",
      },
      {
        name: "QuichÃ©",
        code: "QUI",
      },
      {
        name: "Retalhuleu",
        code: "RET",
      },
      {
        name: "SacatepÃ©quez",
        code: "SAC",
      },
      {
        name: "San Marcos",
        code: "SMA",
      },
      {
        name: "Santa Rosa",
        code: "SRO",
      },
      {
        name: "SololÃ¡",
        code: "SOL",
      },
      {
        name: "SuchitepÃ©quez",
        code: "SUC",
      },
      {
        name: "TotonicapÃ¡n",
        code: "TOT",
      },
      {
        name: "Zacapa",
        code: "ZAC",
      },
    ],
  },
  {
    name: "Guernsey",
    code: "GG",
    provinces: [],
  },
  {
    name: "Guinea",
    code: "GN",
    provinces: [],
  },
  {
    name: "Guinea Bissau",
    code: "GW",
    provinces: [],
  },
  {
    name: "Guyana",
    code: "GY",
    provinces: [],
  },
  {
    name: "Haiti",
    code: "HT",
    provinces: [],
  },
  {
    name: "Heard Island And Mcdonald Islands",
    code: "HM",
    provinces: [],
  },
  {
    name: "Holy See (Vatican City State)",
    code: "VA",
    provinces: [],
  },
  {
    name: "Honduras",
    code: "HN",
    provinces: [],
  },
  {
    name: "Hong Kong",
    code: "HK",
    provinces: [
      {
        name: "Hong Kong Island",
        code: "HK",
      },
      {
        name: "Kowloon",
        code: "KL",
      },
      {
        name: "New Territories",
        code: "NT",
      },
    ],
  },
  {
    name: "Hungary",
    code: "HU",
    provinces: [],
  },
  {
    name: "Iceland",
    code: "IS",
    provinces: [],
  },
  {
    name: "India",
    code: "IN",
    provinces: [
      {
        name: "Andaman and Nicobar",
        code: "AN",
      },
      {
        name: "Andhra Pradesh",
        code: "AP",
      },
      {
        name: "Arunachal Pradesh",
        code: "AR",
      },
      {
        name: "Assam",
        code: "AS",
      },
      {
        name: "Bihar",
        code: "BR",
      },
      {
        name: "Chandigarh",
        code: "CH",
      },
      {
        name: "Chattisgarh",
        code: "CG",
      },
      {
        name: "Dadra and Nagar Haveli",
        code: "DN",
      },
      {
        name: "Daman and Diu",
        code: "DD",
      },
      {
        name: "Delhi",
        code: "DL",
      },
      {
        name: "Goa",
        code: "GA",
      },
      {
        name: "Gujarat",
        code: "GJ",
      },
      {
        name: "Haryana",
        code: "HR",
      },
      {
        name: "Himachal Pradesh",
        code: "HP",
      },
      {
        name: "Jammu and Kashmir",
        code: "JK",
      },
      {
        name: "Jharkhand",
        code: "JH",
      },
      {
        name: "Karnataka",
        code: "KA",
      },
      {
        name: "Kerala",
        code: "KL",
      },
      {
        name: "Lakshadweep",
        code: "LD",
      },
      {
        name: "Madhya Pradesh",
        code: "MP",
      },
      {
        name: "Maharashtra",
        code: "MH",
      },
      {
        name: "Manipur",
        code: "MN",
      },
      {
        name: "Meghalaya",
        code: "ML",
      },
      {
        name: "Mizoram",
        code: "MZ",
      },
      {
        name: "Nagaland",
        code: "NL",
      },
      {
        name: "Orissa",
        code: "OR",
      },
      {
        name: "Puducherry",
        code: "PY",
      },
      {
        name: "Punjab",
        code: "PB",
      },
      {
        name: "Rajasthan",
        code: "RJ",
      },
      {
        name: "Sikkim",
        code: "SK",
      },
      {
        name: "Tamil Nadu",
        code: "TN",
      },
      {
        name: "Telangana",
        code: "TS",
      },
      {
        name: "Tripura",
        code: "TR",
      },
      {
        name: "Uttar Pradesh",
        code: "UP",
      },
      {
        name: "Uttarakhand",
        code: "UK",
      },
      {
        name: "West Bengal",
        code: "WB",
      },
    ],
  },
  {
    name: "Indonesia",
    code: "ID",
    provinces: [
      {
        name: "Aceh",
        code: "AC",
      },
      {
        name: "Bali",
        code: "BA",
      },
      {
        name: "Bangka Belitung",
        code: "BB",
      },
      {
        name: "Banten",
        code: "BT",
      },
      {
        name: "Bengkulu",
        code: "BE",
      },
      {
        name: "Gorontalo",
        code: "GO",
      },
      {
        name: "Jakarta",
        code: "JK",
      },
      {
        name: "Jambi",
        code: "JA",
      },
      {
        name: "Jawa Barat",
        code: "JB",
      },
      {
        name: "Jawa Tengah",
        code: "JT",
      },
      {
        name: "Jawa Timur",
        code: "JI",
      },
      {
        name: "Kalimantan Barat",
        code: "KB",
      },
      {
        name: "Kalimantan Selatan",
        code: "KS",
      },
      {
        name: "Kalimantan Tengah",
        code: "KT",
      },
      {
        name: "Kalimantan Timur",
        code: "KI",
      },
      {
        name: "Kalimantan Utara",
        code: "KU",
      },
      {
        name: "Kepulauan Riau",
        code: "KR",
      },
      {
        name: "Lampung",
        code: "LA",
      },
      {
        name: "Maluku",
        code: "MA",
      },
      {
        name: "Maluku Utara",
        code: "MU",
      },
      {
        name: "Nusa Tenggara Barat",
        code: "NB",
      },
      {
        name: "Nusa Tenggara Timur",
        code: "NT",
      },
      {
        name: "Papua",
        code: "PA",
      },
      {
        name: "Papua Barat",
        code: "PB",
      },
      {
        name: "Riau",
        code: "RI",
      },
      {
        name: "Sulawesi Barat",
        code: "SR",
      },
      {
        name: "Sulawesi Selatan",
        code: "SN",
      },
      {
        name: "Sulawesi Tengah",
        code: "ST",
      },
      {
        name: "Sulawesi Tenggara",
        code: "SG",
      },
      {
        name: "Sulawesi Utara",
        code: "SA",
      },
      {
        name: "Sumatra Barat",
        code: "SB",
      },
      {
        name: "Sumatra Selatan",
        code: "SS",
      },
      {
        name: "Sumatra Utara",
        code: "SU",
      },
      {
        name: "Yogyakarta",
        code: "YO",
      },
    ],
  },
  {
    name: "Iran, Islamic Republic Of",
    code: "IR",
    provinces: [],
  },
  {
    name: "Iraq",
    code: "IQ",
    provinces: [],
  },
  {
    name: "Ireland",
    code: "IE",
    provinces: [
      {
        name: "Carlow",
        code: "CW",
      },
      {
        name: "Cavan",
        code: "CN",
      },
      {
        name: "Clare",
        code: "CE",
      },
      {
        name: "Cork",
        code: "CO",
      },
      {
        name: "Donegal",
        code: "DL",
      },
      {
        name: "Dublin",
        code: "D",
      },
      {
        name: "Galway",
        code: "G",
      },
      {
        name: "Kerry",
        code: "KY",
      },
      {
        name: "Kildare",
        code: "KE",
      },
      {
        name: "Kilkenny",
        code: "KK",
      },
      {
        name: "Laois",
        code: "LS",
      },
      {
        name: "Leitrim",
        code: "LM",
      },
      {
        name: "Limerick",
        code: "LK",
      },
      {
        name: "Longford",
        code: "LD",
      },
      {
        name: "Louth",
        code: "LH",
      },
      {
        name: "Mayo",
        code: "MO",
      },
      {
        name: "Meath",
        code: "MH",
      },
      {
        name: "Monaghan",
        code: "MN",
      },
      {
        name: "Offaly",
        code: "OY",
      },
      {
        name: "Roscommon",
        code: "RN",
      },
      {
        name: "Sligo",
        code: "SO",
      },
      {
        name: "Tipperary",
        code: "TA",
      },
      {
        name: "Waterford",
        code: "WD",
      },
      {
        name: "Westmeath",
        code: "WH",
      },
      {
        name: "Wexford",
        code: "WX",
      },
      {
        name: "Wicklow",
        code: "WW",
      },
    ],
  },
  {
    name: "Isle Of Man",
    code: "IM",
    provinces: [],
  },
  {
    name: "Israel",
    code: "IL",
    provinces: [],
  },
  {
    name: "Italy",
    code: "IT",
    provinces: [
      {
        name: "Agrigento",
        code: "AG",
      },
      {
        name: "Alessandria",
        code: "AL",
      },
      {
        name: "Ancona",
        code: "AN",
      },
      {
        name: "Aosta",
        code: "AO",
      },
      {
        name: "Arezzo",
        code: "AR",
      },
      {
        name: "Ascoli Piceno",
        code: "AP",
      },
      {
        name: "Asti",
        code: "AT",
      },
      {
        name: "Avellino",
        code: "AV",
      },
      {
        name: "Bari",
        code: "BA",
      },
      {
        name: "Barletta-Andria-Trani",
        code: "BT",
      },
      {
        name: "Belluno",
        code: "BL",
      },
      {
        name: "Benevento",
        code: "BN",
      },
      {
        name: "Bergamo",
        code: "BG",
      },
      {
        name: "Biella",
        code: "BI",
      },
      {
        name: "Bologna",
        code: "BO",
      },
      {
        name: "Bolzano",
        code: "BZ",
      },
      {
        name: "Brescia",
        code: "BS",
      },
      {
        name: "Brindisi",
        code: "BR",
      },
      {
        name: "Cagliari",
        code: "CA",
      },
      {
        name: "Caltanissetta",
        code: "CL",
      },
      {
        name: "Campobasso",
        code: "CB",
      },
      {
        name: "Carbonia-Iglesias",
        code: "CI",
      },
      {
        name: "Caserta",
        code: "CE",
      },
      {
        name: "Catania",
        code: "CT",
      },
      {
        name: "Catanzaro",
        code: "CZ",
      },
      {
        name: "Chieti",
        code: "CH",
      },
      {
        name: "Como",
        code: "CO",
      },
      {
        name: "Cosenza",
        code: "CS",
      },
      {
        name: "Cremona",
        code: "CR",
      },
      {
        name: "Crotone",
        code: "KR",
      },
      {
        name: "Cuneo",
        code: "CN",
      },
      {
        name: "Enna",
        code: "EN",
      },
      {
        name: "Fermo",
        code: "FM",
      },
      {
        name: "Ferrara",
        code: "FE",
      },
      {
        name: "Firenze",
        code: "FI",
      },
      {
        name: "Foggia",
        code: "FG",
      },
      {
        name: "ForlÃ¬-Cesena",
        code: "FC",
      },
      {
        name: "Frosinone",
        code: "FR",
      },
      {
        name: "Genova",
        code: "GE",
      },
      {
        name: "Gorizia",
        code: "GO",
      },
      {
        name: "Grosseto",
        code: "GR",
      },
      {
        name: "Imperia",
        code: "IM",
      },
      {
        name: "Isernia",
        code: "IS",
      },
      {
        name: "L'Aquila",
        code: "AQ",
      },
      {
        name: "La Spezia",
        code: "SP",
      },
      {
        name: "Latina",
        code: "LT",
      },
      {
        name: "Lecce",
        code: "LE",
      },
      {
        name: "Lecco",
        code: "LC",
      },
      {
        name: "Livorno",
        code: "LI",
      },
      {
        name: "Lodi",
        code: "LO",
      },
      {
        name: "Lucca",
        code: "LU",
      },
      {
        name: "Macerata",
        code: "MC",
      },
      {
        name: "Mantova",
        code: "MN",
      },
      {
        name: "Massa-Carrara",
        code: "MS",
      },
      {
        name: "Matera",
        code: "MT",
      },
      {
        name: "Medio Campidano",
        code: "VS",
      },
      {
        name: "Messina",
        code: "ME",
      },
      {
        name: "Milano",
        code: "MI",
      },
      {
        name: "Modena",
        code: "MO",
      },
      {
        name: "Monza e Brianza",
        code: "MB",
      },
      {
        name: "Napoli",
        code: "NA",
      },
      {
        name: "Novara",
        code: "NO",
      },
      {
        name: "Nuoro",
        code: "NU",
      },
      {
        name: "Ogliastra",
        code: "OG",
      },
      {
        name: "Olbia-Tempio",
        code: "OT",
      },
      {
        name: "Oristano",
        code: "OR",
      },
      {
        name: "Padova",
        code: "PD",
      },
      {
        name: "Palermo",
        code: "PA",
      },
      {
        name: "Parma",
        code: "PR",
      },
      {
        name: "Pavia",
        code: "PV",
      },
      {
        name: "Perugia",
        code: "PG",
      },
      {
        name: "Pesaro e Urbino",
        code: "PU",
      },
      {
        name: "Pescara",
        code: "PE",
      },
      {
        name: "Piacenza",
        code: "PC",
      },
      {
        name: "Pisa",
        code: "PI",
      },
      {
        name: "Pistoia",
        code: "PT",
      },
      {
        name: "Pordenone",
        code: "PN",
      },
      {
        name: "Potenza",
        code: "PZ",
      },
      {
        name: "Prato",
        code: "PO",
      },
      {
        name: "Ragusa",
        code: "RG",
      },
      {
        name: "Ravenna",
        code: "RA",
      },
      {
        name: "Reggio Calabria",
        code: "RC",
      },
      {
        name: "Reggio Emilia",
        code: "RE",
      },
      {
        name: "Rieti",
        code: "RI",
      },
      {
        name: "Rimini",
        code: "RN",
      },
      {
        name: "Roma",
        code: "RM",
      },
      {
        name: "Rovigo",
        code: "RO",
      },
      {
        name: "Salerno",
        code: "SA",
      },
      {
        name: "Sassari",
        code: "SS",
      },
      {
        name: "Savona",
        code: "SV",
      },
      {
        name: "Siena",
        code: "SI",
      },
      {
        name: "Siracusa",
        code: "SR",
      },
      {
        name: "Sondrio",
        code: "SO",
      },
      {
        name: "Taranto",
        code: "TA",
      },
      {
        name: "Teramo",
        code: "TE",
      },
      {
        name: "Terni",
        code: "TR",
      },
      {
        name: "Torino",
        code: "TO",
      },
      {
        name: "Trapani",
        code: "TP",
      },
      {
        name: "Trento",
        code: "TN",
      },
      {
        name: "Treviso",
        code: "TV",
      },
      {
        name: "Trieste",
        code: "TS",
      },
      {
        name: "Udine",
        code: "UD",
      },
      {
        name: "Varese",
        code: "VA",
      },
      {
        name: "Venezia",
        code: "VE",
      },
      {
        name: "Verbano-Cusio-Ossola",
        code: "VB",
      },
      {
        name: "Vercelli",
        code: "VC",
      },
      {
        name: "Verona",
        code: "VR",
      },
      {
        name: "Vibo Valentia",
        code: "VV",
      },
      {
        name: "Vicenza",
        code: "VI",
      },
      {
        name: "Viterbo",
        code: "VT",
      },
    ],
  },
  {
    name: "Jamaica",
    code: "JM",
    provinces: [],
  },
  {
    name: "Japan",
    code: "JP",
    provinces: [
      {
        name: "Aichi",
        code: "JP-23",
      },
      {
        name: "Akita",
        code: "JP-05",
      },
      {
        name: "Aomori",
        code: "JP-02",
      },
      {
        name: "Chiba",
        code: "JP-12",
      },
      {
        name: "Ehime",
        code: "JP-38",
      },
      {
        name: "Fukui",
        code: "JP-18",
      },
      {
        name: "Fukuoka",
        code: "JP-40",
      },
      {
        name: "Fukushima",
        code: "JP-07",
      },
      {
        name: "Gifu",
        code: "JP-21",
      },
      {
        name: "Gunma",
        code: "JP-10",
      },
      {
        name: "Hiroshima",
        code: "JP-34",
      },
      {
        name: "Hokkaidō",
        code: "JP-01",
      },
      {
        name: "Hyōgo",
        code: "JP-28",
      },
      {
        name: "Ibaraki",
        code: "JP-08",
      },
      {
        name: "Ishikawa",
        code: "JP-17",
      },
      {
        name: "Iwate",
        code: "JP-03",
      },
      {
        name: "Kagawa",
        code: "JP-37",
      },
      {
        name: "Kagoshima",
        code: "JP-46",
      },
      {
        name: "Kanagawa",
        code: "JP-14",
      },
      {
        name: "Kōchi",
        code: "JP-39",
      },
      {
        name: "Kumamoto",
        code: "JP-43",
      },
      {
        name: "Kyōto",
        code: "JP-26",
      },
      {
        name: "Mie",
        code: "JP-24",
      },
      {
        name: "Miyagi",
        code: "JP-04",
      },
      {
        name: "Miyazaki",
        code: "JP-45",
      },
      {
        name: "Nagano",
        code: "JP-20",
      },
      {
        name: "Nagasaki",
        code: "JP-42",
      },
      {
        name: "Nara",
        code: "JP-29",
      },
      {
        name: "Niigata",
        code: "JP-15",
      },
      {
        name: "ÅŒita",
        code: "JP-44",
      },
      {
        name: "Okayama",
        code: "JP-33",
      },
      {
        name: "Okinawa",
        code: "JP-47",
      },
      {
        name: "ÅŒsaka",
        code: "JP-27",
      },
      {
        name: "Saga",
        code: "JP-41",
      },
      {
        name: "Saitama",
        code: "JP-11",
      },
      {
        name: "Shiga",
        code: "JP-25",
      },
      {
        name: "Shimane",
        code: "JP-32",
      },
      {
        name: "Shizuoka",
        code: "JP-22",
      },
      {
        name: "Tochigi",
        code: "JP-09",
      },
      {
        name: "Tokushima",
        code: "JP-36",
      },
      {
        name: "Tōkyō",
        code: "JP-13",
      },
      {
        name: "Tottori",
        code: "JP-31",
      },
      {
        name: "Toyama",
        code: "JP-16",
      },
      {
        name: "Wakayama",
        code: "JP-30",
      },
      {
        name: "Yamagata",
        code: "JP-06",
      },
      {
        name: "Yamaguchi",
        code: "JP-35",
      },
      {
        name: "Yamanashi",
        code: "JP-19",
      },
    ],
  },
  {
    name: "Jersey",
    code: "JE",
    provinces: [],
  },
  {
    name: "Jordan",
    code: "JO",
    provinces: [],
  },
  {
    name: "Kazakhstan",
    code: "KZ",
    provinces: [],
  },
  {
    name: "Kenya",
    code: "KE",
    provinces: [],
  },
  {
    name: "Kiribati",
    code: "KI",
    provinces: [],
  },
  {
    name: "Korea, Democratic People's Republic Of",
    code: "KP",
    provinces: [],
  },
  {
    name: "Kosovo",
    code: "XK",
    provinces: [],
  },
  {
    name: "Kuwait",
    code: "KW",
    provinces: [],
  },
  {
    name: "Kyrgyzstan",
    code: "KG",
    provinces: [],
  },
  {
    name: "Lao People's Democratic Republic",
    code: "LA",
    provinces: [],
  },
  {
    name: "Latvia",
    code: "LV",
    provinces: [],
  },
  {
    name: "Lebanon",
    code: "LB",
    provinces: [],
  },
  {
    name: "Lesotho",
    code: "LS",
    provinces: [],
  },
  {
    name: "Liberia",
    code: "LR",
    provinces: [],
  },
  {
    name: "Libyan Arab Jamahiriya",
    code: "LY",
    provinces: [],
  },
  {
    name: "Liechtenstein",
    code: "LI",
    provinces: [],
  },
  {
    name: "Lithuania",
    code: "LT",
    provinces: [],
  },
  {
    name: "Luxembourg",
    code: "LU",
    provinces: [],
  },
  {
    name: "Macao",
    code: "MO",
    provinces: [],
  },
  {
    name: "Macedonia, Republic Of",
    code: "MK",
    provinces: [],
  },
  {
    name: "Madagascar",
    code: "MG",
    provinces: [],
  },
  {
    name: "Malawi",
    code: "MW",
    provinces: [],
  },
  {
    name: "Malaysia",
    code: "MY",
    provinces: [
      {
        name: "Johor",
        code: "JHR",
      },
      {
        name: "Kedah",
        code: "KDH",
      },
      {
        name: "Kelantan",
        code: "KTN",
      },
      {
        name: "Kuala Lumpur",
        code: "KUL",
      },
      {
        name: "Labuan",
        code: "LBN",
      },
      {
        name: "Melaka",
        code: "MLK",
      },
      {
        name: "Negeri Sembilan",
        code: "NSN",
      },
      {
        name: "Pahang",
        code: "PHG",
      },
      {
        name: "Perak",
        code: "PRK",
      },
      {
        name: "Perlis",
        code: "PLS",
      },
      {
        name: "Pulau Pinang",
        code: "PNG",
      },
      {
        name: "Putrajaya",
        code: "PJY",
      },
      {
        name: "Sabah",
        code: "SBH",
      },
      {
        name: "Sarawak",
        code: "SWK",
      },
      {
        name: "Selangor",
        code: "SGR",
      },
      {
        name: "Terengganu",
        code: "TRG",
      },
    ],
  },
  {
    name: "Maldives",
    code: "MV",
    provinces: [],
  },
  {
    name: "Mali",
    code: "ML",
    provinces: [],
  },
  {
    name: "Malta",
    code: "MT",
    provinces: [],
  },
  {
    name: "Martinique",
    code: "MQ",
    provinces: [],
  },
  {
    name: "Mauritania",
    code: "MR",
    provinces: [],
  },
  {
    name: "Mauritius",
    code: "MU",
    provinces: [],
  },
  {
    name: "Mayotte",
    code: "YT",
    provinces: [],
  },
  {
    name: "Mexico",
    code: "MX",
    provinces: [
      {
        name: "Aguascalientes",
        code: "AGS",
      },
      {
        name: "Baja California",
        code: "BC",
      },
      {
        name: "Baja California Sur",
        code: "BCS",
      },
      {
        name: "Campeche",
        code: "CAMP",
      },
      {
        name: "Chiapas",
        code: "CHIS",
      },
      {
        name: "Chihuahua",
        code: "CHIH",
      },
      {
        name: "Ciudad de MÃ©xico",
        code: "DF",
      },
      {
        name: "Coahuila",
        code: "COAH",
      },
      {
        name: "Colima",
        code: "COL",
      },
      {
        name: "Durango",
        code: "DGO",
      },
      {
        name: "Guanajuato",
        code: "GTO",
      },
      {
        name: "Guerrero",
        code: "GRO",
      },
      {
        name: "Hidalgo",
        code: "HGO",
      },
      {
        name: "Jalisco",
        code: "JAL",
      },
      {
        name: "MÃ©xico",
        code: "MEX",
      },
      {
        name: "MichoacÃ¡n",
        code: "MICH",
      },
      {
        name: "Morelos",
        code: "MOR",
      },
      {
        name: "Nayarit",
        code: "NAY",
      },
      {
        name: "Nuevo LeÃ³n",
        code: "NL",
      },
      {
        name: "Oaxaca",
        code: "OAX",
      },
      {
        name: "Puebla",
        code: "PUE",
      },
      {
        name: "QuerÃ©taro",
        code: "QRO",
      },
      {
        name: "Quintana Roo",
        code: "Q ROO",
      },
      {
        name: "San Luis PotosÃ­",
        code: "SLP",
      },
      {
        name: "Sinaloa",
        code: "SIN",
      },
      {
        name: "Sonora",
        code: "SON",
      },
      {
        name: "Tabasco",
        code: "TAB",
      },
      {
        name: "Tamaulipas",
        code: "TAMPS",
      },
      {
        name: "Tlaxcala",
        code: "TLAX",
      },
      {
        name: "Veracruz",
        code: "VER",
      },
      {
        name: "YucatÃ¡n",
        code: "YUC",
      },
      {
        name: "Zacatecas",
        code: "ZAC",
      },
    ],
  },
  {
    name: "Moldova, Republic of",
    code: "MD",
    provinces: [],
  },
  {
    name: "Monaco",
    code: "MC",
    provinces: [],
  },
  {
    name: "Mongolia",
    code: "MN",
    provinces: [],
  },
  {
    name: "Montenegro",
    code: "ME",
    provinces: [],
  },
  {
    name: "Montserrat",
    code: "MS",
    provinces: [],
  },
  {
    name: "Morocco",
    code: "MA",
    provinces: [],
  },
  {
    name: "Mozambique",
    code: "MZ",
    provinces: [],
  },
  {
    name: "Myanmar",
    code: "MM",
    provinces: [],
  },
  {
    name: "Namibia",
    code: "NA",
    provinces: [],
  },
  {
    name: "Nauru",
    code: "NR",
    provinces: [],
  },
  {
    name: "Nepal",
    code: "NP",
    provinces: [],
  },
  {
    name: "Netherlands",
    code: "NL",
    provinces: [],
  },
  {
    name: "Netherlands Antilles",
    code: "AN",
    provinces: [],
  },
  {
    name: "New Caledonia",
    code: "NC",
    provinces: [],
  },
  {
    name: "New Zealand",
    code: "NZ",
    provinces: [
      {
        name: "Auckland",
        code: "AUK",
      },
      {
        name: "Bay of Plenty",
        code: "BOP",
      },
      {
        name: "Canterbury",
        code: "CAN",
      },
      {
        name: "Gisborne",
        code: "GIS",
      },
      {
        name: "Hawke's Bay",
        code: "HKB",
      },
      {
        name: "Manawatu-Wanganui",
        code: "MWT",
      },
      {
        name: "Marlborough",
        code: "MBH",
      },
      {
        name: "Nelson",
        code: "NSN",
      },
      {
        name: "Northland",
        code: "NTL",
      },
      {
        name: "Otago",
        code: "OTA",
      },
      {
        name: "Southland",
        code: "STL",
      },
      {
        name: "Taranaki",
        code: "TKI",
      },
      {
        name: "Tasman",
        code: "TAS",
      },
      {
        name: "Waikato",
        code: "WKO",
      },
      {
        name: "Wellington",
        code: "WGN",
      },
      {
        name: "West Coast",
        code: "WTC",
      },
    ],
  },
  {
    name: "Nicaragua",
    code: "NI",
    provinces: [],
  },
  {
    name: "Niger",
    code: "NE",
    provinces: [],
  },
  {
    name: "Nigeria",
    code: "NG",
    provinces: [
      {
        name: "Abia",
        code: "AB",
      },
      {
        name: "Abuja Federal Capital Territory",
        code: "FC",
      },
      {
        name: "Adamawa",
        code: "AD",
      },
      {
        name: "Akwa Ibom",
        code: "AK",
      },
      {
        name: "Anambra",
        code: "AN",
      },
      {
        name: "Bauchi",
        code: "BA",
      },
      {
        name: "Bayelsa",
        code: "BY",
      },
      {
        name: "Benue",
        code: "BE",
      },
      {
        name: "Borno",
        code: "BO",
      },
      {
        name: "Cross River",
        code: "CR",
      },
      {
        name: "Delta",
        code: "DE",
      },
      {
        name: "Ebonyi",
        code: "EB",
      },
      {
        name: "Edo",
        code: "ED",
      },
      {
        name: "Ekiti",
        code: "EK",
      },
      {
        name: "Enugu",
        code: "EN",
      },
      {
        name: "Gombe",
        code: "GO",
      },
      {
        name: "Imo",
        code: "IM",
      },
      {
        name: "Jigawa",
        code: "JI",
      },
      {
        name: "Kaduna",
        code: "KD",
      },
      {
        name: "Kano",
        code: "KN",
      },
      {
        name: "Katsina",
        code: "KT",
      },
      {
        name: "Kebbi",
        code: "KE",
      },
      {
        name: "Kogi",
        code: "KO",
      },
      {
        name: "Kwara",
        code: "KW",
      },
      {
        name: "Lagos",
        code: "LA",
      },
      {
        name: "Nasarawa",
        code: "NA",
      },
      {
        name: "Niger",
        code: "NI",
      },
      {
        name: "Ogun",
        code: "OG",
      },
      {
        name: "Ondo",
        code: "ON",
      },
      {
        name: "Osun",
        code: "OS",
      },
      {
        name: "Oyo",
        code: "OY",
      },
      {
        name: "Plateau",
        code: "PL",
      },
      {
        name: "Rivers",
        code: "RI",
      },
      {
        name: "Sokoto",
        code: "SO",
      },
      {
        name: "Taraba",
        code: "TA",
      },
      {
        name: "Yobe",
        code: "YO",
      },
      {
        name: "Zamfara",
        code: "ZA",
      },
    ],
  },
  {
    name: "Niue",
    code: "NU",
    provinces: [],
  },
  {
    name: "Norfolk Island",
    code: "NF",
    provinces: [],
  },
  {
    name: "Northern Ireland",
    code: "NB",
    provinces: [],
  },
  {
    name: "Norway",
    code: "NO",
    provinces: [],
  },
  {
    name: "Oman",
    code: "OM",
    provinces: [],
  },
  {
    name: "Pakistan",
    code: "PK",
    provinces: [],
  },
  {
    name: "Palestinian Territory, Occupied",
    code: "PS",
    provinces: [],
  },
  {
    name: "Panama",
    code: "PA",
    provinces: [
      {
        name: "Bocas del Toro",
        code: "PA-1",
      },
      {
        name: "ChiriquÃ­",
        code: "PA-4",
      },
      {
        name: "CoclÃ©",
        code: "PA-2",
      },
      {
        name: "ColÃ³n",
        code: "PA-3",
      },
      {
        name: "DariÃ©n",
        code: "PA-5",
      },
      {
        name: "EmberÃ¡",
        code: "PA-EM",
      },
      {
        name: "Herrera",
        code: "PA-6",
      },
      {
        name: "Kuna Yala",
        code: "PA-KY",
      },
      {
        name: "Los Santos",
        code: "PA-7",
      },
      {
        name: "NgÃ¶be-BuglÃ©",
        code: "PA-NB",
      },
      {
        name: "PanamÃ¡",
        code: "PA-8",
      },
      {
        name: "PanamÃ¡ Oeste",
        code: "PA-10",
      },
      {
        name: "Veraguas",
        code: "PA-9",
      },
    ],
  },
  {
    name: "Papua New Guinea",
    code: "PG",
    provinces: [],
  },
  {
    name: "Paraguay",
    code: "PY",
    provinces: [],
  },
  {
    name: "Peru",
    code: "PE",
    provinces: [],
  },
  {
    name: "Philippines",
    code: "PH",
    provinces: [],
  },
  {
    name: "Pitcairn",
    code: "PN",
    provinces: [],
  },
  {
    name: "Poland",
    code: "PL",
    provinces: [],
  },
  {
    name: "Portugal",
    code: "PT",
    provinces: [
      {
        name: "AÃ§ores",
        code: "PT-20",
      },
      {
        name: "Aveiro",
        code: "PT-01",
      },
      {
        name: "Beja",
        code: "PT-02",
      },
      {
        name: "Braga",
        code: "PT-03",
      },
      {
        name: "BraganÃ§a",
        code: "PT-04",
      },
      {
        name: "Castelo Branco",
        code: "PT-05",
      },
      {
        name: "Coimbra",
        code: "PT-06",
      },
      {
        name: "Ã‰vora",
        code: "PT-07",
      },
      {
        name: "Faro",
        code: "PT-08",
      },
      {
        name: "Guarda",
        code: "PT-09",
      },
      {
        name: "Leiria",
        code: "PT-10",
      },
      {
        name: "Lisboa",
        code: "PT-11",
      },
      {
        name: "Madeira",
        code: "PT-30",
      },
      {
        name: "Portalegre",
        code: "PT-12",
      },
      {
        name: "Porto",
        code: "PT-13",
      },
      {
        name: "SantarÃ©m",
        code: "PT-14",
      },
      {
        name: "SetÃºbal",
        code: "PT-15",
      },
      {
        name: "Viana do Castelo",
        code: "PT-16",
      },
      {
        name: "Vila Real",
        code: "PT-17",
      },
      {
        name: "Viseu",
        code: "PT-18",
      },
    ],
  },
  {
    name: "Qatar",
    code: "QA",
    provinces: [],
  },
  {
    name: "Republic of Cameroon",
    code: "CM",
    provinces: [],
  },
  {
    name: "Rest of World",
    code: "*",
    provinces: [],
  },
  {
    name: "Reunion",
    code: "RE",
    provinces: [],
  },
  {
    name: "Romania",
    code: "RO",
    provinces: [
      {
        name: "Alba",
        code: "AB",
      },
      {
        name: "Arad",
        code: "AR",
      },
      {
        name: "ArgeÈ™",
        code: "AG",
      },
      {
        name: "BacÄƒu",
        code: "BC",
      },
      {
        name: "Bihor",
        code: "BH",
      },
      {
        name: "BistriÈ›a-NÄƒsÄƒud",
        code: "BN",
      },
      {
        name: "BotoÈ™ani",
        code: "BT",
      },
      {
        name: "BrÄƒila",
        code: "BR",
      },
      {
        name: "BraÈ™ov",
        code: "BV",
      },
      {
        name: "BucureÈ™ti",
        code: "B",
      },
      {
        name: "BuzÄƒu",
        code: "BZ",
      },
      {
        name: "CÄƒlÄƒraÈ™i",
        code: "CL",
      },
      {
        name: "CaraÈ™-Severin",
        code: "CS",
      },
      {
        name: "Cluj",
        code: "CJ",
      },
      {
        name: "ConstanÈ›a",
        code: "CT",
      },
      {
        name: "Covasna",
        code: "CV",
      },
      {
        name: "DÃ¢mboviÈ›a",
        code: "DB",
      },
      {
        name: "Dolj",
        code: "DJ",
      },
      {
        name: "GalaÈ›i",
        code: "GL",
      },
      {
        name: "Giurgiu",
        code: "GR",
      },
      {
        name: "Gorj",
        code: "GJ",
      },
      {
        name: "Harghita",
        code: "HR",
      },
      {
        name: "Hunedoara",
        code: "HD",
      },
      {
        name: "IalomiÈ›a",
        code: "IL",
      },
      {
        name: "IaÈ™i",
        code: "IS",
      },
      {
        name: "Ilfov",
        code: "IF",
      },
      {
        name: "MaramureÈ™",
        code: "MM",
      },
      {
        name: "MehedinÈ›i",
        code: "MH",
      },
      {
        name: "MureÈ™",
        code: "MS",
      },
      {
        name: "NeamÈ›",
        code: "NT",
      },
      {
        name: "Olt",
        code: "OT",
      },
      {
        name: "Prahova",
        code: "PH",
      },
      {
        name: "SÄƒlaj",
        code: "SJ",
      },
      {
        name: "Satu Mare",
        code: "SM",
      },
      {
        name: "Sibiu",
        code: "SB",
      },
      {
        name: "Suceava",
        code: "SV",
      },
      {
        name: "Teleorman",
        code: "TR",
      },
      {
        name: "TimiÈ™",
        code: "TM",
      },
      {
        name: "Tulcea",
        code: "TL",
      },
      {
        name: "VÃ¢lcea",
        code: "VL",
      },
      {
        name: "Vaslui",
        code: "VS",
      },
      {
        name: "Vrancea",
        code: "VN",
      },
    ],
  },
  {
    name: "Russia",
    code: "RU",
    provinces: [
      {
        name: "Altai Krai",
        code: "ALT",
      },
      {
        name: "Altai Republic",
        code: "AL",
      },
      {
        name: "Amur Oblast",
        code: "AMU",
      },
      {
        name: "Arkhangelsk Oblast",
        code: "ARK",
      },
      {
        name: "Astrakhan Oblast",
        code: "AST",
      },
      {
        name: "Belgorod Oblast",
        code: "BEL",
      },
      {
        name: "Bryansk Oblast",
        code: "BRY",
      },
      {
        name: "Chechen Republic",
        code: "CE",
      },
      {
        name: "Chelyabinsk Oblast",
        code: "CHE",
      },
      {
        name: "Chukotka Autonomous Okrug",
        code: "CHU",
      },
      {
        name: "Chuvash Republic",
        code: "CU",
      },
      {
        name: "Irkutsk Oblast",
        code: "IRK",
      },
      {
        name: "Ivanovo Oblast",
        code: "IVA",
      },
      {
        name: "Jewish Autonomous Oblast",
        code: "YEV",
      },
      {
        name: "Kabardino-Balkarian Republic",
        code: "KB",
      },
      {
        name: "Kaliningrad Oblast",
        code: "KGD",
      },
      {
        name: "Kaluga Oblast",
        code: "KLU",
      },
      {
        name: "Kamchatka Krai",
        code: "KAM",
      },
      {
        name: "Karachayâ€“Cherkess Republic",
        code: "KC",
      },
      {
        name: "Kemerovo Oblast",
        code: "KEM",
      },
      {
        name: "Khabarovsk Krai",
        code: "KHA",
      },
      {
        name: "Khanty-Mansi Autonomous Okrug",
        code: "KHM",
      },
      {
        name: "Kirov Oblast",
        code: "KIR",
      },
      {
        name: "Komi Republic",
        code: "KO",
      },
      {
        name: "Kostroma Oblast",
        code: "KOS",
      },
      {
        name: "Krasnodar Krai",
        code: "KDA",
      },
      {
        name: "Krasnoyarsk Krai",
        code: "KYA",
      },
      {
        name: "Kurgan Oblast",
        code: "KGN",
      },
      {
        name: "Kursk Oblast",
        code: "KRS",
      },
      {
        name: "Leningrad Oblast",
        code: "LEN",
      },
      {
        name: "Lipetsk Oblast",
        code: "LIP",
      },
      {
        name: "Magadan Oblast",
        code: "MAG",
      },
      {
        name: "Mari El Republic",
        code: "ME",
      },
      {
        name: "Moscow",
        code: "MOW",
      },
      {
        name: "Moscow Oblast",
        code: "MOS",
      },
      {
        name: "Murmansk Oblast",
        code: "MUR",
      },
      {
        name: "Nizhny Novgorod Oblast",
        code: "NIZ",
      },
      {
        name: "Novgorod Oblast",
        code: "NGR",
      },
      {
        name: "Novosibirsk Oblast",
        code: "NVS",
      },
      {
        name: "Omsk Oblast",
        code: "OMS",
      },
      {
        name: "Orenburg Oblast",
        code: "ORE",
      },
      {
        name: "Oryol Oblast",
        code: "ORL",
      },
      {
        name: "Penza Oblast",
        code: "PNZ",
      },
      {
        name: "Perm Krai",
        code: "PER",
      },
      {
        name: "Primorsky Krai",
        code: "PRI",
      },
      {
        name: "Pskov Oblast",
        code: "PSK",
      },
      {
        name: "Republic of Adygeya",
        code: "AD",
      },
      {
        name: "Republic of Bashkortostan",
        code: "BA",
      },
      {
        name: "Republic of Buryatia",
        code: "BU",
      },
      {
        name: "Republic of Dagestan",
        code: "DA",
      },
      {
        name: "Republic of Ingushetia",
        code: "IN",
      },
      {
        name: "Republic of Kalmykia",
        code: "KL",
      },
      {
        name: "Republic of Karelia",
        code: "KR",
      },
      {
        name: "Republic of Khakassia",
        code: "KK",
      },
      {
        name: "Republic of Mordovia",
        code: "MO",
      },
      {
        name: "Republic of North Ossetiaâ€“Alania",
        code: "SE",
      },
      {
        name: "Republic of Tatarstan",
        code: "TA",
      },
      {
        name: "Rostov Oblast",
        code: "ROS",
      },
      {
        name: "Ryazan Oblast",
        code: "RYA",
      },
      {
        name: "Saint Petersburg",
        code: "SPE",
      },
      {
        name: "Sakha Republic (Yakutia)",
        code: "SA",
      },
      {
        name: "Sakhalin Oblast",
        code: "SAK",
      },
      {
        name: "Samara Oblast",
        code: "SAM",
      },
      {
        name: "Saratov Oblast",
        code: "SAR",
      },
      {
        name: "Smolensk Oblast",
        code: "SMO",
      },
      {
        name: "Stavropol Krai",
        code: "STA",
      },
      {
        name: "Sverdlovsk Oblast",
        code: "SVE",
      },
      {
        name: "Tambov Oblast",
        code: "TAM",
      },
      {
        name: "Tomsk Oblast",
        code: "TOM",
      },
      {
        name: "Tula Oblast",
        code: "TUL",
      },
      {
        name: "Tver Oblast",
        code: "TVE",
      },
      {
        name: "Tyumen Oblast",
        code: "TYU",
      },
      {
        name: "Tyva Republic",
        code: "TY",
      },
      {
        name: "Udmurtia",
        code: "UD",
      },
      {
        name: "Ulyanovsk Oblast",
        code: "ULY",
      },
      {
        name: "Vladimir Oblast",
        code: "VLA",
      },
      {
        name: "Volgograd Oblast",
        code: "VGG",
      },
      {
        name: "Vologda Oblast",
        code: "VLG",
      },
      {
        name: "Voronezh Oblast",
        code: "VOR",
      },
      {
        name: "Yamalo-Nenets Autonomous Okrug",
        code: "YAN",
      },
      {
        name: "Yaroslavl Oblast",
        code: "YAR",
      },
      {
        name: "Zabaykalsky Krai",
        code: "ZAB",
      },
    ],
  },
  {
    name: "Rwanda",
    code: "RW",
    provinces: [],
  },
  {
    name: "Saint BarthÃ©lemy",
    code: "BL",
    provinces: [],
  },
  {
    name: "Saint Helena",
    code: "SH",
    provinces: [],
  },
  {
    name: "Saint Kitts And Nevis",
    code: "KN",
    provinces: [],
  },
  {
    name: "Saint Lucia",
    code: "LC",
    provinces: [],
  },
  {
    name: "Saint Martin",
    code: "MF",
    provinces: [],
  },
  {
    name: "Saint Pierre And Miquelon",
    code: "PM",
    provinces: [],
  },
  {
    name: "Samoa",
    code: "WS",
    provinces: [],
  },
  {
    name: "San Marino",
    code: "SM",
    provinces: [],
  },
  {
    name: "Sao Tome And Principe",
    code: "ST",
    provinces: [],
  },
  {
    name: "Saudi Arabia",
    code: "SA",
    provinces: [],
  },
  {
    name: "Senegal",
    code: "SN",
    provinces: [],
  },
  {
    name: "Serbia",
    code: "RS",
    provinces: [],
  },
  {
    name: "Seychelles",
    code: "SC",
    provinces: [],
  },
  {
    name: "Sierra Leone",
    code: "SL",
    provinces: [],
  },
  {
    name: "Singapore",
    code: "SG",
    provinces: [],
  },
  {
    name: "Sint Maarten",
    code: "SX",
    provinces: [],
  },
  {
    name: "Slovakia",
    code: "SK",
    provinces: [],
  },
  {
    name: "Slovenia",
    code: "SI",
    provinces: [],
  },
  {
    name: "Solomon Islands",
    code: "SB",
    provinces: [],
  },
  {
    name: "Somalia",
    code: "SO",
    provinces: [],
  },
  {
    name: "South Africa",
    code: "ZA",
    provinces: [
      {
        name: "Eastern Cape",
        code: "EC",
      },
      {
        name: "Free State",
        code: "FS",
      },
      {
        name: "Gauteng",
        code: "GT",
      },
      {
        name: "KwaZulu-Natal",
        code: "NL",
      },
      {
        name: "Limpopo",
        code: "LP",
      },
      {
        name: "Mpumalanga",
        code: "MP",
      },
      {
        name: "North West",
        code: "NW",
      },
      {
        name: "Northern Cape",
        code: "NC",
      },
      {
        name: "Western Cape",
        code: "WC",
      },
    ],
  },
  {
    name: "South Georgia And The South Sandwich Islands",
    code: "GS",
    provinces: [],
  },
  {
    name: "South Korea",
    code: "KR",
    provinces: [
      {
        name: "Busan",
        code: "KR-26",
      },
      {
        name: "Chungbuk",
        code: "KR-43",
      },
      {
        name: "Chungnam",
        code: "KR-44",
      },
      {
        name: "Daegu",
        code: "KR-27",
      },
      {
        name: "Daejeon",
        code: "KR-30",
      },
      {
        name: "Gangwon",
        code: "KR-42",
      },
      {
        name: "Gwangju",
        code: "KR-29",
      },
      {
        name: "Gyeongbuk",
        code: "KR-47",
      },
      {
        name: "Gyeonggi",
        code: "KR-41",
      },
      {
        name: "Gyeongnam",
        code: "KR-48",
      },
      {
        name: "Incheon",
        code: "KR-28",
      },
      {
        name: "Jeju",
        code: "KR-49",
      },
      {
        name: "Jeonbuk",
        code: "KR-45",
      },
      {
        name: "Jeonnam",
        code: "KR-46",
      },
      {
        name: "Sejong",
        code: "KR-50",
      },
      {
        name: "Seoul",
        code: "KR-11",
      },
      {
        name: "Ulsan",
        code: "KR-31",
      },
    ],
  },
  {
    name: "South Sudan",
    code: "SS",
    provinces: [],
  },
  {
    name: "Spain",
    code: "ES",
    provinces: [
      {
        name: "A CoruÃ±a",
        code: "C",
      },
      {
        name: "Ãlava",
        code: "VI",
      },
      {
        name: "Albacete",
        code: "AB",
      },
      {
        name: "Alicante",
        code: "A",
      },
      {
        name: "AlmerÃ­a",
        code: "AL",
      },
      {
        name: "Asturias",
        code: "O",
      },
      {
        name: "Ãvila",
        code: "AV",
      },
      {
        name: "Badajoz",
        code: "BA",
      },
      {
        name: "Balears",
        code: "PM",
      },
      {
        name: "Barcelona",
        code: "B",
      },
      {
        name: "Burgos",
        code: "BU",
      },
      {
        name: "CÃ¡ceres",
        code: "CC",
      },
      {
        name: "CÃ¡diz",
        code: "CA",
      },
      {
        name: "Cantabria",
        code: "S",
      },
      {
        name: "CastellÃ³n",
        code: "CS",
      },
      {
        name: "Ceuta",
        code: "CE",
      },
      {
        name: "Ciudad Real",
        code: "CR",
      },
      {
        name: "CÃ³rdoba",
        code: "CO",
      },
      {
        name: "Cuenca",
        code: "CU",
      },
      {
        name: "Girona",
        code: "GI",
      },
      {
        name: "Granada",
        code: "GR",
      },
      {
        name: "Guadalajara",
        code: "GU",
      },
      {
        name: "GuipÃºzcoa",
        code: "SS",
      },
      {
        name: "Huelva",
        code: "H",
      },
      {
        name: "Huesca",
        code: "HU",
      },
      {
        name: "JaÃ©n",
        code: "J",
      },
      {
        name: "La Rioja",
        code: "LO",
      },
      {
        name: "Las Palmas",
        code: "GC",
      },
      {
        name: "LeÃ³n",
        code: "LE",
      },
      {
        name: "Lleida",
        code: "L",
      },
      {
        name: "Lugo",
        code: "LU",
      },
      {
        name: "Madrid",
        code: "M",
      },
      {
        name: "MÃ¡laga",
        code: "MA",
      },
      {
        name: "Melilla",
        code: "ML",
      },
      {
        name: "Murcia",
        code: "MU",
      },
      {
        name: "Navarra",
        code: "NA",
      },
      {
        name: "Ourense",
        code: "OR",
      },
      {
        name: "Palencia",
        code: "P",
      },
      {
        name: "Pontevedra",
        code: "PO",
      },
      {
        name: "Salamanca",
        code: "SA",
      },
      {
        name: "Santa Cruz de Tenerife",
        code: "TF",
      },
      {
        name: "Segovia",
        code: "SG",
      },
      {
        name: "Sevilla",
        code: "SE",
      },
      {
        name: "Soria",
        code: "SO",
      },
      {
        name: "Tarragona",
        code: "T",
      },
      {
        name: "Teruel",
        code: "TE",
      },
      {
        name: "Toledo",
        code: "TO",
      },
      {
        name: "Valencia",
        code: "V",
      },
      {
        name: "Valladolid",
        code: "VA",
      },
      {
        name: "Vizcaya",
        code: "BI",
      },
      {
        name: "Zamora",
        code: "ZA",
      },
      {
        name: "Zaragoza",
        code: "Z",
      },
    ],
  },
  {
    name: "Sri Lanka",
    code: "LK",
    provinces: [],
  },
  {
    name: "St. Vincent",
    code: "VC",
    provinces: [],
  },
  {
    name: "Sudan",
    code: "SD",
    provinces: [],
  },
  {
    name: "Suriname",
    code: "SR",
    provinces: [],
  },
  {
    name: "Svalbard And Jan Mayen",
    code: "SJ",
    provinces: [],
  },
  {
    name: "Swaziland",
    code: "SZ",
    provinces: [],
  },
  {
    name: "Sweden",
    code: "SE",
    provinces: [],
  },
  {
    name: "Switzerland",
    code: "CH",
    provinces: [],
  },
  {
    name: "Syria",
    code: "SY",
    provinces: [],
  },
  {
    name: "Taiwan",
    code: "TW",
    provinces: [],
  },
  {
    name: "Tajikistan",
    code: "TJ",
    provinces: [],
  },
  {
    name: "Tanzania, United Republic Of",
    code: "TZ",
    provinces: [],
  },
  {
    name: "Thailand",
    code: "TH",
    provinces: [
      {
        name: "Amnat Charoen",
        code: "TH-37",
      },
      {
        name: "Ang Thong",
        code: "TH-15",
      },
      {
        name: "Bangkok",
        code: "TH-10",
      },
      {
        name: "Bueng Kan",
        code: "TH-38",
      },
      {
        name: "Buriram",
        code: "TH-31",
      },
      {
        name: "Chachoengsao",
        code: "TH-24",
      },
      {
        name: "Chai Nat",
        code: "TH-18",
      },
      {
        name: "Chaiyaphum",
        code: "TH-36",
      },
      {
        name: "Chanthaburi",
        code: "TH-22",
      },
      {
        name: "Chiang Mai",
        code: "TH-50",
      },
      {
        name: "Chiang Rai",
        code: "TH-57",
      },
      {
        name: "Chon Buri",
        code: "TH-20",
      },
      {
        name: "Chumphon",
        code: "TH-86",
      },
      {
        name: "Kalasin",
        code: "TH-46",
      },
      {
        name: "Kamphaeng Phet",
        code: "TH-62",
      },
      {
        name: "Kanchanaburi",
        code: "TH-71",
      },
      {
        name: "Khon Kaen",
        code: "TH-40",
      },
      {
        name: "Krabi",
        code: "TH-81",
      },
      {
        name: "Lampang",
        code: "TH-52",
      },
      {
        name: "Lamphun",
        code: "TH-51",
      },
      {
        name: "Loei",
        code: "TH-42",
      },
      {
        name: "Lopburi",
        code: "TH-16",
      },
      {
        name: "Mae Hong Son",
        code: "TH-58",
      },
      {
        name: "Maha Sarakham",
        code: "TH-44",
      },
      {
        name: "Mukdahan",
        code: "TH-49",
      },
      {
        name: "Nakhon Nayok",
        code: "TH-26",
      },
      {
        name: "Nakhon Pathom",
        code: "TH-73",
      },
      {
        name: "Nakhon Phanom",
        code: "TH-48",
      },
      {
        name: "Nakhon Ratchasima",
        code: "TH-30",
      },
      {
        name: "Nakhon Sawan",
        code: "TH-60",
      },
      {
        name: "Nakhon Si Thammarat",
        code: "TH-80",
      },
      {
        name: "Nan",
        code: "TH-55",
      },
      {
        name: "Narathiwat",
        code: "TH-96",
      },
      {
        name: "Nong Bua Lam Phu",
        code: "TH-39",
      },
      {
        name: "Nong Khai",
        code: "TH-43",
      },
      {
        name: "Nonthaburi",
        code: "TH-12",
      },
      {
        name: "Pathum Thani",
        code: "TH-13",
      },
      {
        name: "Pattani",
        code: "TH-94",
      },
      {
        name: "Pattaya",
        code: "TH-S",
      },
      {
        name: "Phangnga",
        code: "TH-82",
      },
      {
        name: "Phatthalung",
        code: "TH-93",
      },
      {
        name: "Phayao",
        code: "TH-56",
      },
      {
        name: "Phetchabun",
        code: "TH-67",
      },
      {
        name: "Phetchaburi",
        code: "TH-76",
      },
      {
        name: "Phichit",
        code: "TH-66",
      },
      {
        name: "Phitsanulok",
        code: "TH-65",
      },
      {
        name: "Phra Nakhon Si Ayutthaya",
        code: "TH-14",
      },
      {
        name: "Phrae",
        code: "TH-54",
      },
      {
        name: "Phuket",
        code: "TH-83",
      },
      {
        name: "Prachin Buri",
        code: "TH-25",
      },
      {
        name: "Prachuap Khiri Khan",
        code: "TH-77",
      },
      {
        name: "Ranong",
        code: "TH-85",
      },
      {
        name: "Ratchaburi",
        code: "TH-70",
      },
      {
        name: "Rayong",
        code: "TH-21",
      },
      {
        name: "Roi Et",
        code: "TH-45",
      },
      {
        name: "Sa Kaeo",
        code: "TH-27",
      },
      {
        name: "Sakon Nakhon",
        code: "TH-47",
      },
      {
        name: "Samut Prakan",
        code: "TH-11",
      },
      {
        name: "Samut Sakhon",
        code: "TH-74",
      },
      {
        name: "Samut Songkhram",
        code: "TH-75",
      },
      {
        name: "Saraburi",
        code: "TH-19",
      },
      {
        name: "Satun",
        code: "TH-91",
      },
      {
        name: "Sing Buri",
        code: "TH-17",
      },
      {
        name: "Sisaket",
        code: "TH-33",
      },
      {
        name: "Songkhla",
        code: "TH-90",
      },
      {
        name: "Sukhothai",
        code: "TH-64",
      },
      {
        name: "Suphan Buri",
        code: "TH-72",
      },
      {
        name: "Surat Thani",
        code: "TH-84",
      },
      {
        name: "Surin",
        code: "TH-32",
      },
      {
        name: "Tak",
        code: "TH-63",
      },
      {
        name: "Trang",
        code: "TH-92",
      },
      {
        name: "Trat",
        code: "TH-23",
      },
      {
        name: "Ubon Ratchathani",
        code: "TH-34",
      },
      {
        name: "Udon Thani",
        code: "TH-41",
      },
      {
        name: "Uthai Thani",
        code: "TH-61",
      },
      {
        name: "Uttaradit",
        code: "TH-53",
      },
      {
        name: "Yala",
        code: "TH-95",
      },
      {
        name: "Yasothon",
        code: "TH-35",
      },
    ],
  },
  {
    name: "Timor Leste",
    code: "TL",
    provinces: [],
  },
  {
    name: "Togo",
    code: "TG",
    provinces: [],
  },
  {
    name: "Tokelau",
    code: "TK",
    provinces: [],
  },
  {
    name: "Tonga",
    code: "TO",
    provinces: [],
  },
  {
    name: "Trinidad and Tobago",
    code: "TT",
    provinces: [],
  },
  {
    name: "Tunisia",
    code: "TN",
    provinces: [],
  },
  {
    name: "Turkey",
    code: "TR",
    provinces: [],
  },
  {
    name: "Turkmenistan",
    code: "TM",
    provinces: [],
  },
  {
    name: "Turks and Caicos Islands",
    code: "TC",
    provinces: [],
  },
  {
    name: "Tuvalu",
    code: "TV",
    provinces: [],
  },
  {
    name: "Uganda",
    code: "UG",
    provinces: [],
  },
  {
    name: "Ukraine",
    code: "UA",
    provinces: [],
  },
  {
    name: "United Arab Emirates",
    code: "AE",
    provinces: [
      {
        name: "Abu Dhabi",
        code: "AZ",
      },
      {
        name: "Ajman",
        code: "AJ",
      },
      {
        name: "Dubai",
        code: "DU",
      },
      {
        name: "Fujairah",
        code: "FU",
      },
      {
        name: "Ras al-Khaimah",
        code: "RK",
      },
      {
        name: "Sharjah",
        code: "SH",
      },
      {
        name: "Umm al-Quwain",
        code: "UQ",
      },
    ],
  },
  {
    name: "United Kingdom",
    code: "GB",
    provinces: [],
  },
  {
    name: "United States",
    code: "US",
    provinces: [
      {
        name: "Alabama",
        code: "AL",
      },
      {
        name: "Alaska",
        code: "AK",
      },
      {
        name: "American Samoa",
        code: "AS",
      },
      {
        name: "Arizona",
        code: "AZ",
      },
      {
        name: "Arkansas",
        code: "AR",
      },
      {
        name: "Armed Forces Americas",
        code: "AA",
      },
      {
        name: "Armed Forces Europe",
        code: "AE",
      },
      {
        name: "Armed Forces Pacific",
        code: "AP",
      },
      {
        name: "California",
        code: "CA",
      },
      {
        name: "Colorado",
        code: "CO",
      },
      {
        name: "Connecticut",
        code: "CT",
      },
      {
        name: "Delaware",
        code: "DE",
      },
      {
        name: "District of Columbia",
        code: "DC",
      },
      {
        name: "Federated States of Micronesia",
        code: "FM",
      },
      {
        name: "Florida",
        code: "FL",
      },
      {
        name: "Georgia",
        code: "GA",
      },
      {
        name: "Guam",
        code: "GU",
      },
      {
        name: "Hawaii",
        code: "HI",
      },
      {
        name: "Idaho",
        code: "ID",
      },
      {
        name: "Illinois",
        code: "IL",
      },
      {
        name: "Indiana",
        code: "IN",
      },
      {
        name: "Iowa",
        code: "IA",
      },
      {
        name: "Kansas",
        code: "KS",
      },
      {
        name: "Kentucky",
        code: "KY",
      },
      {
        name: "Louisiana",
        code: "LA",
      },
      {
        name: "Maine",
        code: "ME",
      },
      {
        name: "Marshall Islands",
        code: "MH",
      },
      {
        name: "Maryland",
        code: "MD",
      },
      {
        name: "Massachusetts",
        code: "MA",
      },
      {
        name: "Michigan",
        code: "MI",
      },
      {
        name: "Minnesota",
        code: "MN",
      },
      {
        name: "Mississippi",
        code: "MS",
      },
      {
        name: "Missouri",
        code: "MO",
      },
      {
        name: "Montana",
        code: "MT",
      },
      {
        name: "Nebraska",
        code: "NE",
      },
      {
        name: "Nevada",
        code: "NV",
      },
      {
        name: "New Hampshire",
        code: "NH",
      },
      {
        name: "New Jersey",
        code: "NJ",
      },
      {
        name: "New Mexico",
        code: "NM",
      },
      {
        name: "New York",
        code: "NY",
      },
      {
        name: "North Carolina",
        code: "NC",
      },
      {
        name: "North Dakota",
        code: "ND",
      },
      {
        name: "Northern Mariana Islands",
        code: "MP",
      },
      {
        name: "Ohio",
        code: "OH",
      },
      {
        name: "Oklahoma",
        code: "OK",
      },
      {
        name: "Oregon",
        code: "OR",
      },
      {
        name: "Palau",
        code: "PW",
      },
      {
        name: "Pennsylvania",
        code: "PA",
      },
      {
        name: "Puerto Rico",
        code: "PR",
      },
      {
        name: "Rhode Island",
        code: "RI",
      },
      {
        name: "South Carolina",
        code: "SC",
      },
      {
        name: "South Dakota",
        code: "SD",
      },
      {
        name: "Tennessee",
        code: "TN",
      },
      {
        name: "Texas",
        code: "TX",
      },
      {
        name: "Utah",
        code: "UT",
      },
      {
        name: "Vermont",
        code: "VT",
      },
      {
        name: "Virgin Islands",
        code: "VI",
      },
      {
        name: "Virginia",
        code: "VA",
      },
      {
        name: "Washington",
        code: "WA",
      },
      {
        name: "West Virginia",
        code: "WV",
      },
      {
        name: "Wisconsin",
        code: "WI",
      },
      {
        name: "Wyoming",
        code: "WY",
      },
    ],
  },
  {
    name: "United States Minor Outlying Islands",
    code: "UM",
    provinces: [],
  },
  {
    name: "Uruguay",
    code: "UY",
    provinces: [],
  },
  {
    name: "Uzbekistan",
    code: "UZ",
    provinces: [],
  },
  {
    name: "Vanuatu",
    code: "VU",
    provinces: [],
  },
  {
    name: "Venezuela",
    code: "VE",
    provinces: [],
  },
  {
    name: "Vietnam",
    code: "VN",
    provinces: [],
  },
  {
    name: "Virgin Islands, British",
    code: "VG",
    provinces: [],
  },
  {
    name: "Wallis And Futuna",
    code: "WF",
    provinces: [],
  },
  {
    name: "Western Sahara",
    code: "EH",
    provinces: [],
  },
  {
    name: "Yemen",
    code: "YE",
    provinces: [],
  },
  {
    name: "Zambia",
    code: "ZM",
    provinces: [],
  },
  {
    name: "Zimbabwe",
    code: "ZW",
    provinces: [],
  },
];

export { countries };
