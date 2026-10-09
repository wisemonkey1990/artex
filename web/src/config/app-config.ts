import { resolveLocale } from "@/i18n/config";

import packageJson from "../../package.json";

const currentYear = new Date().getFullYear();

// 浏览器标题与搜索引擎摘要会写入所有页面的静态 metadata，并根据当前语言选择。
const META_BY_LOCALE = {
  zh: {
    title: "ARTEX — 自主渗透测试控制台",
    description: "LLM 驱动的自主渗透测试系统控制台",
  },
} as const;

export const APP_CONFIG = {
  name: "ARTEX",
  version: packageJson.version,
  copyright: `© ${currentYear}, ARTEX.`,
  meta: META_BY_LOCALE[resolveLocale()],
};
