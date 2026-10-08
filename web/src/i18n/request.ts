import { getRequestConfig } from "next-intl/server";

import { type Locale, resolveLocale } from "./config";

// locale 별 메시지 로더. 명시적 import 맵으로 둬서 번들러가 messages/ 디렉터리 전체를
// context 모듈로 끌어들이지 않게 한다(로컬에만 존재하는 zh.sources.json 제외).
const loaders: Record<Locale, () => Promise<{ default: Record<string, unknown> }>> = {
  ko: () => import("../../messages/ko.json"),
  zh: () => import("../../messages/zh.json"),
};

// next-intl 요청 설정. i18n 경로 라우팅(세그먼트·미들웨어)을 쓰지 않는 구성이라
// locale 在此处直接确定，不读取 requestLocale，以兼容静态导出。
// 简体中文和韩语文案分别保存在 messages/zh.json 与 messages/ko.json。
export default getRequestConfig(async () => {
  const locale = resolveLocale();
  const messages = (await loaders[locale]()).default;
  return { locale, messages };
});
