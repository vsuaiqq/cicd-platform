import { css } from '@emotion/react'
import Breadcrumb from '../Breadcrumb'
import { useI18n } from '../../../i18n'
import { runRow, runsTableHead, hero, tabList, tabShell } from '../../../pages/project/styles'
import {
  Skeleton,
  SkeletonButton,
  SkeletonCard,
  SkeletonCircle,
  SkeletonStack,
  SkeletonText,
} from './Skeleton'

const rowSk = css([
  runRow,
  {
    pointerEvents: 'none',
    '&:hover': { background: 'transparent' },
  },
])


export function SkeletonListRows({ count = 3 }: { count?: number }) {
  return (
    <>
      {Array.from({ length: count }, (_, i) => (
        <div
          key={i}
          css={css({
            display: 'flex',
            alignItems: 'center',
            gap: 12,
            margin: '1px 6px',
            padding: '8px 12px 8px 14px',
          })}
          aria-hidden
        >
          <Skeleton height={34} width={34} radius={9} />
          <SkeletonStack gap={6} css={css({ flex: 1 })}>
            <Skeleton height={10} width={`${55 + i * 12}%`} />
            <Skeleton height={10} width="38%" />
          </SkeletonStack>
        </div>
      ))}
    </>
  )
}

export function SkeletonBreadcrumb() {
  const { t } = useI18n()
  return <Breadcrumb items={[{ label: t.nav.dashboard, to: '/' }, { label: '…' }]} />
}

export function SkeletonPageHeader() {
  return (
    <SkeletonStack gap={10} css={css({ marginBottom: 24 })}>
      <Skeleton height={28} width="38%" style={{ maxWidth: 320 }} />
      <Skeleton height={14} width="52%" style={{ maxWidth: 420 }} />
    </SkeletonStack>
  )
}

export function SkeletonStatCards({ count = 3 }: { count?: number }) {
  return (
    <div
      css={css({
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fit, minmax(140px, 1fr))',
        gap: 14,
        marginBottom: 24,
      })}
      aria-hidden
    >
      {Array.from({ length: count }, (_, i) => (
        <div
          key={i}
          css={css({
            padding: '16px 18px',
            borderRadius: 'var(--radius-lg)',
            border: '1px solid var(--border-subtle)',
            background: 'var(--bg-card)',
          })}
        >
          <Skeleton height={12} width={60} />
          <Skeleton height={32} width={48} css={css({ marginTop: 8 })} />
        </div>
      ))}
    </div>
  )
}

export function SkeletonProjectGrid({ count = 3 }: { count?: number }) {
  return (
    <div
      css={css({
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fill, minmax(320px, 1fr))',
        gap: 14,
      })}
      aria-hidden
    >
      {Array.from({ length: count }, (_, i) => (
        <SkeletonCard key={i} height={150} />
      ))}
    </div>
  )
}

export function DashboardPageSkeleton() {
  return (
    <div aria-busy="true" aria-label="Loading">
      <SkeletonStatCards />
      <SkeletonProjectGrid />
    </div>
  )
}

export function ProjectPageSkeleton() {
  return (
    <div aria-busy="true" aria-label="Loading">
      <SkeletonBreadcrumb />
      <div css={hero}>
        <div css={css({ display: 'flex', gap: 16, alignItems: 'flex-start' })}>
          <Skeleton height={48} width={48} radius={12} />
          <div css={css({ flex: 1, minWidth: 0 })}>
            <Skeleton height={22} width="42%" style={{ maxWidth: 280, marginBottom: 12 }} />
            <Skeleton height={14} width="58%" style={{ maxWidth: 360 }} />
          </div>
        </div>
      </div>
      <div css={tabShell}>
        <div css={tabList}>
          {[72, 88, 96].map((w) => (
            <Skeleton key={w} height={34} width={w} radius={10} />
          ))}
        </div>
      </div>
      <SkeletonCard height={320} />
    </div>
  )
}

export function RunsTableSkeleton({ rows = 4 }: { rows?: number }) {
  return (
    <>
      <div css={runsTableHead} aria-hidden>
        <span />
        <span />
        <span />
        <span />
        <span />
      </div>
      {Array.from({ length: rows }, (_, i) => (
        <div key={i} css={rowSk} aria-hidden>
          <Skeleton height={24} width={72} />
          <div>
            <Skeleton height={14} width={`${45 + i * 8}%`} style={{ maxWidth: 200, marginBottom: 6 }} />
            <Skeleton height={10} width={88} />
          </div>
          <Skeleton height={12} width={40} css={css({ marginLeft: 'auto' })} />
          <Skeleton height={12} width={56} css={css({ marginLeft: 'auto' })} />
          <span />
        </div>
      ))}
    </>
  )
}

export function SettingsSectionSkeleton() {
  return (
    <SkeletonStack gap={16} css={css({ padding: '4px 0' })}>
      <Skeleton height={40} width="100%" radius="var(--radius-md)" />
      <SkeletonText lines={4} lineHeight={14} />
      <Skeleton height={36} width={120} radius="var(--radius-md)" />
    </SkeletonStack>
  )
}

export function AnalyticsPanelSkeleton() {
  return (
    <div aria-hidden>
      <SkeletonStack gap={20}>
        <div css={css({ display: 'flex', gap: 8 })}>
          <SkeletonButton width={48} />
          <SkeletonButton width={56} />
          <SkeletonButton width={48} />
        </div>
        <div
          css={css({
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fill, minmax(140px, 1fr))',
            gap: 12,
          })}
        >
          {Array.from({ length: 5 }, (_, i) => (
            <Skeleton key={i} height={72} width="100%" radius="var(--radius-lg)" />
          ))}
        </div>
        <Skeleton height={220} width="100%" radius="var(--radius-lg)" />
        <Skeleton height={160} width="100%" radius="var(--radius-lg)" />
      </SkeletonStack>
    </div>
  )
}

export function RunDetailSkeleton() {
  return (
    <div aria-busy="true" aria-label="Loading">
      <SkeletonBreadcrumb />
      <SkeletonStack gap={16} css={css({ marginTop: 8 })}>
        <div css={css({ display: 'flex', alignItems: 'center', gap: 12, flexWrap: 'wrap' })}>
          <Skeleton height={28} width={200} />
          <Skeleton height={24} width={80} radius="var(--radius-full)" />
        </div>
        <Skeleton height={14} width="45%" style={{ maxWidth: 360 }} />
        <div css={css({ display: 'flex', gap: 8, marginTop: 8 })}>
          <SkeletonButton width={100} />
          <SkeletonButton width={88} />
        </div>
        <div
          css={css({
            display: 'grid',
            gridTemplateColumns: '240px minmax(0, 1fr)',
            gap: 16,
            marginTop: 16,
            '@media (max-width: 900px)': { gridTemplateColumns: '1fr' },
          })}
        >
          <Skeleton height={320} width="100%" radius="var(--radius-lg)" />
          <Skeleton height={420} width="100%" radius="var(--radius-lg)" />
        </div>
      </SkeletonStack>
    </div>
  )
}

export function SetupPageSkeleton() {
  return (
    <div aria-busy="true" aria-label="Loading">
      <SkeletonBreadcrumb />
      <SkeletonPageHeader />
      <SkeletonStack gap={24}>
        <SkeletonCard height={140} />
        <SkeletonCard height={200} />
        <Skeleton height={40} width={160} radius="var(--radius-md)" />
      </SkeletonStack>
    </div>
  )
}

export function AppShellSkeleton() {
  return (
    <div
      css={css({
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: '60vh',
        padding: 24,
      })}
      aria-busy="true"
      aria-live="polite"
    >
      <SkeletonStack gap={12} css={css({ width: 'min(100%, 280px)', alignItems: 'center' })}>
        <SkeletonCircle size={48} />
        <Skeleton height={12} width="70%" />
        <Skeleton height={12} width="50%" />
      </SkeletonStack>
    </div>
  )
}
