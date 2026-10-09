import { getRequestConfig } from "next-intl/server";

import { type Locale, resolveLocale } from "./config";

// 显式导入消息文件，避免打包器将整个 messages 目录作为上下文模块加载。
const loaders: Record<Locale, () => Promise<{ default: Record<string, unknown> }>> = {
  zh: () => import("../../messages/zh.json"),
};

// 不使用 i18n 路径路由，因此在此处确定语言，不读取 requestLocale，以兼容静态导出。
export default getRequestConfig(async () => {
  const locale = resolveLocale();
  const messages = (await loaders[locale]()).default;
  return { locale, messages };
});
