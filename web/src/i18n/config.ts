// 支持的语言与默认值。简体中文是默认界面语言，韩语作为可选语言保留。
export const LOCALES = ["zh", "ko"] as const;
export type Locale = (typeof LOCALES)[number];
export const DEFAULT_LOCALE: Locale = "zh";

// 在构建时确定语言。为兼容静态导出，不使用 cookies()、headers() 等动态 API。
// NEXT_PUBLIC_LOCALE 为空或不受支持时回退到简体中文。
export function resolveLocale(): Locale {
  const raw = process.env.NEXT_PUBLIC_LOCALE;
  return LOCALES.includes(raw as Locale) ? (raw as Locale) : DEFAULT_LOCALE;
}
