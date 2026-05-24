import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { syncRecentWithProjects } from '../../lib/recentProjects'
import { useAppSelector } from '../../store'
import { useListProjectsQuery } from '../../store/api/apiSlice'
import { useI18n } from '../../i18n'
import { buildPaletteItems } from './buildItems'
import { PaletteList } from './PaletteList'
import { IconSearch } from './icons'
import {
  backdrop,
  dialog,
  footer,
  footerHint,
  footerMeta,
  kbd,
  searchIcon,
  searchInput,
  searchRow,
  spinner,
} from './styles'

interface CommandPaletteProps {
  open: boolean
  onClose: () => void
  onOpenSettings?: () => void
}

export default function CommandPalette({ open, onClose, onOpenSettings }: CommandPaletteProps) {
  const navigate = useNavigate()
  const { t } = useI18n()
  const { sessionValidated } = useAppSelector((s) => s.auth)
  const [query, setQuery] = useState('')
  const [cursor, setCursor] = useState(0)
  const inputRef = useRef<HTMLInputElement>(null)
  const listRef = useRef<HTMLDivElement>(null)

  const {
    data: projects = [],
    isLoading,
    isFetching,
    isError,
  } = useListProjectsQuery(undefined, {
    skip: !open || !sessionValidated,
  })

  const loading = isLoading || (isFetching && projects.length === 0)
  const recentIds = useMemo(() => syncRecentWithProjects(projects), [projects])

  useEffect(() => {
    if (!open) return
    setQuery('')
    setCursor(0)
    const id = window.setTimeout(() => inputRef.current?.focus(), 40)
    return () => window.clearTimeout(id)
  }, [open])

  const go = useCallback(
    (path: string) => {
      navigate(path)
      onClose()
    },
    [navigate, onClose],
  )

  const copy = useMemo(
    () => ({
      navigate: t.palette.navigate,
      actions: t.palette.actions,
      recent: t.palette.recent,
      allProjects: t.palette.allProjects,
      projects: t.palette.projects,
      dashboardLabel: t.palette.dashboardLabel,
      dashboardSub: t.palette.dashboardSub,
      newProjectLabel: t.palette.newProjectLabel,
      newProjectSub: t.palette.newProjectSub,
      profileLabel: t.palette.profileLabel,
      profileSub: t.palette.profileSub,
      settingsLabel: t.palette.settingsLabel,
      settingsSub: t.palette.settingsSub,
    }),
    [t],
  )

  const { entries } = useMemo(
    () =>
      buildPaletteItems({
        projects: loading ? [] : projects,
        recentIds: loading ? [] : recentIds,
        query,
        copy,
        go,
        onClose,
        onOpenSettings,
      }),
    [loading, projects, recentIds, query, copy, go, onClose, onOpenSettings],
  )

  const flatSelectable = useMemo(
    () =>
      entries
        .filter((e): e is Extract<(typeof entries)[number], { type: 'item' }> => e.type === 'item')
        .map((e) => e.item),
    [entries],
  )

  useEffect(() => {
    setCursor(0)
  }, [query])

  useEffect(() => {
    if (cursor >= flatSelectable.length) setCursor(Math.max(0, flatSelectable.length - 1))
  }, [cursor, flatSelectable.length])

  useEffect(() => {
    if (!open) return
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        e.preventDefault()
        onClose()
        return
      }
      if (e.key === 'ArrowDown') {
        e.preventDefault()
        setCursor((c) => (flatSelectable.length ? (c + 1) % flatSelectable.length : 0))
      }
      if (e.key === 'ArrowUp') {
        e.preventDefault()
        setCursor((c) => (flatSelectable.length ? (c - 1 + flatSelectable.length) % flatSelectable.length : 0))
      }
      if (e.key === 'Enter' && flatSelectable[cursor]) {
        e.preventDefault()
        flatSelectable[cursor].onSelect()
      }
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [open, flatSelectable, cursor, onClose])

  useEffect(() => {
    const el = listRef.current?.querySelector<HTMLElement>(`[data-index="${cursor}"]`)
    el?.scrollIntoView({ block: 'nearest' })
  }, [cursor])

  if (!open) return null

  const emptyHintText = isError ? t.palette.loadErrorHint : t.palette.emptyQueryHint
  const emptyTitleText = isError ? t.palette.loadErrorTitle : t.palette.emptyQueryTitle

  return (
    <div
      css={backdrop}
      role="presentation"
      onMouseDown={(e) => {
        if (e.target === e.currentTarget) onClose()
      }}
    >
      <div
        css={dialog}
        role="dialog"
        aria-modal="true"
        aria-label={t.nav.quickSearch}
        onMouseDown={(e) => e.stopPropagation()}
      >
        <div css={searchRow}>
          <span css={searchIcon}>
            <IconSearch />
          </span>
          <input
            ref={inputRef}
            css={searchInput}
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder={t.palette.placeholder}
            autoComplete="off"
            autoCorrect="off"
            spellCheck={false}
            aria-autocomplete="list"
            aria-controls="command-palette-list"
            aria-activedescendant={flatSelectable[cursor] ? `cmd-item-${cursor}` : undefined}
          />
          {loading && <span css={spinner} aria-hidden />}
          <kbd css={kbd}>esc</kbd>
        </div>

        <PaletteList
          listRef={listRef}
          entries={entries}
          cursor={cursor}
          query={query}
          loadingProjects={loading}
          emptyTitleText={emptyTitleText}
          emptyHintText={emptyHintText}
          noResultsTitle={t.palette.noResults}
          noResultsHint={t.palette.noResultsHint}
          onCursor={setCursor}
          onSelect={(item) => item.onSelect()}
        />

        <footer css={footer}>
          <span css={footerHint}>
            <kbd css={kbd}>↑</kbd>
            <kbd css={kbd}>↓</kbd>
            {t.palette.hintNavigate}
          </span>
          <span css={footerHint}>
            <kbd css={kbd}>↵</kbd>
            {t.palette.hintOpen}
          </span>
          <span css={footerHint}>
            <kbd css={kbd}>esc</kbd>
            {t.palette.hintClose}
          </span>
          {!loading && projects.length > 0 && (
            <span css={footerMeta}>
              {t.palette.projectCount.replace('{n}', String(projects.length))}
            </span>
          )}
        </footer>
      </div>
    </div>
  )
}
