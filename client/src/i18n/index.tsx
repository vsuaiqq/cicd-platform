import React, { createContext, useCallback, useContext, useState } from 'react'
import type { en } from './en'
import { en as enT } from './en'
import { ru } from './ru'

export type Lang = 'en' | 'ru'
export type Translations = typeof en

const LANG_KEY = 'cicd_lang'

function getStoredLang(): Lang {
  try {
    const l = localStorage.getItem(LANG_KEY)
    if (l === 'en' || l === 'ru') return l
  } catch {  }

  const browser = (typeof navigator !== 'undefined' ? navigator.language : 'en').slice(0, 2).toLowerCase()
  return browser === 'ru' ? 'ru' : 'en'
}

const translations: Record<Lang, Translations> = { en: enT, ru }

interface I18nContextValue {
  lang: Lang
  t: Translations
  setLang: (lang: Lang) => void
}

const I18nContext = createContext<I18nContextValue>({
  lang: 'en',
  t: enT,
  setLang: () => {},
})

export function I18nProvider({ children }: { children: React.ReactNode }) {
  const [lang, setLangState] = useState<Lang>(getStoredLang)

  const setLang = useCallback((l: Lang) => {
    setLangState(l)
    try { localStorage.setItem(LANG_KEY, l) } catch {  }
  }, [])

  return (
    <I18nContext.Provider value={{ lang, t: translations[lang], setLang }}>
      {children}
    </I18nContext.Provider>
  )
}

export function useI18n() {
  return useContext(I18nContext)
}
