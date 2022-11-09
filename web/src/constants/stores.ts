const stores = [
  {
    index: 0,
    type: "shopify",
    label: "Shopify",
    options: [
      {
        name: "A-Ma-Maniere",
        url: "https://www.a-ma-maniere.com",
      },
      {
        name: "Anti Social Social Club",
        url: "https://www.antisocialsocialclub.com",
      },
      {
        name: "Bape",
        url: "https://us.bape.com",
      },
      {
        name: "Bodega",
        url: "https://bdgastore.com",
      },
      {
        name: "CDG US",
        url: "https://us.cdgcdgcdg.com",
      },
      {
        name: "CDG UK",
        url: "https://us.cdgcdgcdg.com",
      },
      {
        name: "CDG JP",
        url: "https://jp.cdgcdgcdg.com",
      },
      {
        name: "Concepts",
        url: "https://cncpts.com",
      },
      {
        name: "Culture Kings",
        url: "https://www.culturekings.com.au",
      },
      {
        name: "DSM JP",
        url: "https://eflash-jp.doverstreetmarket.com",
      },
      {
        name: "DSM SG",
        url: "https://eflash-sg.doverstreetmarket.com",
      },
      {
        name: "DSM UK",
        url: "https://eflash.doverstreetmarket.com",
      },
      {
        name: "DSM US",
        url: "https://eflash-us.doverstreetmarket.com",
      },
      {
        name: "DTLR",
        url: "https://www.dtlr.com",
      },
      {
        name: "Fear Of God",
        url: "https://fearofgod.com",
      },
      {
        name: "Feature",
        url: "https://featuresneakerboutique.com",
      },
      {
        name: "Funko",
        url: "https://checkout.funko.com",
      },
      {
        name: "Hanon Shop",
        url: "https://www.hanon-shop.com",
      },
      {
        name: "Haven",
        url: "https://shop.havenshop.com",
      },
      {
        name: "H. Lorenzo",
        url: "https://www.hlorenzo.com",
      },
      {
        name: "Jimmy Jazz",
        supported: "true",
        url: "https://www.jimmyjazz.com",
      },
      {
        name: "KAWS",
        url: "https://kawsone.com",
      },
      {
        name: "Kith",
        url: "https://kith.com",
      },
      {
        name: "Lapstone & Hammer",
        url: "https://www.lapstoneandhammer.com",
      },
      {
        name: "Livestock",
        url: "https://www.deadstock.ca",
      },
      {
        name: "Nebula",
        url: "https://shop.nebulabots.com",
      },
      {
        name: "NRML",
        url: "https://nrml.ca",
      },
      {
        name: "Nohble",
        url: "https://nohble.com",
      },
      {
        name: "Notre",
        url: "https://www.notre-shop.com",
      },
      {
        name: "Obey Giant",
        url: "https://store.obeygiant.com",
      },
      {
        name: "Off The Hook",
        url: "https://offthehook.ca",
      },
      {
        name: "Oneness 287",
        url: "https://www.onenessboutique.com",
      },
      {
        name: "Packer Shoes",
        url: "https://packershoes.com",
      },
      {
        name: "Palace JP",
        url: "https://shop-jp.palaceskateboards.com",
      },
      {
        name: "Palace UK",
        url: "https://shop.palaceskateboards.com",
      },
      {
        name: "Palace US",
        url: "https://shop-usa.palaceskateboards.com",
      },
      {
        name: "Saint Alfred",
        url: "https://www.saintalfred.com",
      },
      {
        name: "Shoe Palace",
        url: "https://www.shoepalace.com",
      },
      {
        name: "Shop Nice Kicks",
        url: "https://shopnicekicks.com",
      },
      {
        name: "Social Status PGH",
        url: "https://www.socialstatuspgh.com",
      },
      {
        name: "Stussy US",
        url: "https://www.stussy.com",
      },
      {
        name: "Stussy UK",
        url: "https://www.stussy.co.uk",
      },
      {
        name: "Travis Scott Shop",
        url: "https://shop.travisscott.com",
      },
      {
        name: "Trophyroom",
        url: "https://www.trophyroomstore.com",
      },
      {
        name: "Undefeated",
        url: "https://undefeated.com",
      },
      {
        name: "Union Jordan LA",
        url: "https://www.unionjordanla.com",
      },
      {
        name: "Unknown",
        url: "https://www.unknwn.com",
      },
      {
        name: "Wish ATL",
        url: "https://wishatl.com",
      },
      {
        name: "Xhibition",
        url: "https://www.xhibition.co",
      },
    ],
    usesRates: false,
    usesAccounts: false,
  },
  {
    index: 1,
    type: "sup",
    label: "Supreme",
    options: [
      {
        name: "Supreme",
        url: "https://www.supremenewyork.com",
      },
    ],
    usesRates: false,
    usesAccounts: false,
  },
  {
    index: 2,
    type: "foots",
    label: "Footsites",
    options: [
      {
        name: "Foot Locker",
        url: "https://www.footlocker.com",
      },
      {
        name: "Foot Locker",
        url: "https://www.footlocker.co.uk",
      },
      {
        name: "Foot Action",
        url: "https://www.footaction.com",
      },
      {
        name: "Champs Sports",
        url: "https://www.champssports.com",
      },
      {
        name: "Eastbay",
        url: "https://www.eastbay.com",
      },
    ],
    usesRates: false,
    usesAccounts: false,
  },
  {
    index: 3,
    type: "demandware",
    label: "Demandware",
    options: [
      {
        name: "Yeezy Supply",
        url: "https://www.yeezysupply.com",
      },
      {
        name: "New Balance",
        url: "https://www.newbalance.com",
      },
    ],
    usesRates: false,
    usesAccounts: false,
  },
  {
    index: 4,
    type: "others",
    label: "Others",
    options: [
      {
        name: "Walmart",
        url: "https://www.walmart.com",
      },
      {
        name: "Pokémon Center",
        url: "https://www.pokemoncenter.com",
      },
      {
        name: "Nike",
        url: "https://www.nike.com",
      },
    ],
    usesRates: false,
    usesAccounts: false,
  },
];

export { stores };
