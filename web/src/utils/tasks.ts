import { Platform } from "../stores/Main/reducers/stores";
import { Modes } from "../constants";

export const modesForPlatform = (platform: Platform) => {
  switch (platform) {
    default:
      return Modes.DefaultModes;
    case "Shopify":
      return Modes.ShopifyModes;
    case "Supreme":
      return Modes.SupremeModes;
  }
};

export const modesForStoreUrl = (url: string) => {
  switch (url) {
    case "https://www.footlocker.com":
    case "https://www.footlocker.ca":
    case "https://www.kidsfootlocker.com":
    case "https://www.eastbay.com":
    case "https://www.champssports.com":
    case "https://www.footaction.com":
    case "https://www.walmart.com":
      return Modes.DefaultModes;
    case "https://www.supremenewyork.com":
      return Modes.SupremeModes;
    case "Shopify":
    default:
      return Modes.ShopifyModes;
  }
};
