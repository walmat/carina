import {
  Home,
  Calendar,
  Trello,
  GitPullRequest,
  Users,
  Lock,
  Tool,
  Settings,
  Mail,
  Package,
  Globe,
  Star,
  Activity,
  Icon,
} from "react-feather";

export const getIcon = (name: string): Icon => {
  return ICONS[name];
};

type Icons = {
  [name: string]: Icon;
};

export const ICONS: Icons = {
  Dashboard: Home,
  Calendar,
  Tasks: Trello,
  Workflows: GitPullRequest,
  Profiles: Users,
  Proxies: Lock,
  General: Settings,
  Defaults: Tool,
  Accounts: Mail,
  "Shipping Rates": Package,
  Webhooks: Globe,
  "Quick Tasks": Star,
  Integrations: Activity,
};
