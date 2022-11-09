const isValidDiscord =
  /(discord|discordapp).com\/api\/webhooks\/([^\/]+)\/([^\/]+)/;
// TODO: Strengthen this check
const isValidWebhook = (url: string) =>
  (/[discord|discordapp]\.com/.test(url) && isValidDiscord.test(url)) ||
  /hooks\.slack\.com/.test(url) ||
  /webhooks\.aycd\.io/.test(url);

const getType = (url: string) =>
  /slack/i.test(url) ? "slack" : /aycd/i.test(url) ? "aycd" : "discord";

export { isValidWebhook, getType };
