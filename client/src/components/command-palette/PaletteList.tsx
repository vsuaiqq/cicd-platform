import type { ListEntry, PaletteItem } from './buildItems'
import { IconEnter, IconSearch } from './icons'
import { SkeletonListRows } from '../ui'
import {
  emptyHint,
  emptyIcon,
  emptyTitle,
  emptyWrap,
  groupHead,
  list,
  row,
  rowBody,
  rowEnter,
  rowIcon,
  rowMark,
  rowSub,
  rowTitle,
} from './styles'

function Highlight({ text, query }: { text: string; query: string }) {
  if (!query.trim()) return <>{text}</>
  const q = query.trim().toLowerCase()
  const lower = text.toLowerCase()
  const idx = lower.indexOf(q)
  if (idx === -1) return <>{text}</>
  return (
    <>
      {text.slice(0, idx)}
      <mark css={rowMark}>{text.slice(idx, idx + q.length)}</mark>
      {text.slice(idx + q.length)}
    </>
  )
}

interface PaletteListProps {
  listRef: React.RefObject<HTMLDivElement | null>
  entries: ListEntry[]
  cursor: number
  query: string
  loadingProjects: boolean
  emptyTitleText: string
  emptyHintText: string
  noResultsTitle: string
  noResultsHint: string
  onCursor: (index: number) => void
  onSelect: (item: PaletteItem) => void
}

export function PaletteList({
  listRef,
  entries,
  cursor,
  query,
  loadingProjects,
  emptyTitleText,
  emptyHintText,
  noResultsTitle,
  noResultsHint,
  onCursor,
  onSelect,
}: PaletteListProps) {
  const q = query.trim()

  if (entries.length === 0 && !loadingProjects) {
    return (
      <div css={list} ref={listRef} id="command-palette-list" role="listbox">
        <div css={emptyWrap}>
          <div css={emptyIcon}>
            <IconSearch />
          </div>
          <div css={emptyTitle}>{q ? noResultsTitle : emptyTitleText}</div>
          <div css={emptyHint}>{q ? `"${query}" — ${noResultsHint}` : emptyHintText}</div>
        </div>
      </div>
    )
  }

  let itemIndex = -1

  return (
    <div css={list} ref={listRef} id="command-palette-list" role="listbox">
      {entries.map((entry) => {
        if (entry.type === 'group') {
          return (
            <div key={entry.id} css={groupHead} role="presentation">
              {entry.label}
            </div>
          )
        }

        itemIndex += 1
        const idx = itemIndex
        const { item } = entry
        const active = cursor === idx

        return (
          <div
            key={item.id}
            id={`cmd-item-${idx}`}
            data-index={idx}
            data-active={active}
            css={row}
            role="option"
            aria-selected={active}
            onMouseEnter={() => onCursor(idx)}
            onMouseDown={(e) => {
              e.preventDefault()
              onSelect(item)
            }}
          >
            <span css={rowIcon(item.tone)}>{item.icon}</span>
            <span css={rowBody}>
              <span css={rowTitle}>
                <Highlight text={item.label} query={q} />
              </span>
              {item.sublabel && (
                <span css={rowSub}>
                  <Highlight text={item.sublabel} query={q} />
                </span>
              )}
            </span>
            <span css={rowEnter}>
              <IconEnter />
            </span>
          </div>
        )
      })}
      {loadingProjects && <SkeletonListRows />}
    </div>
  )
}
