import type { CSSProperties, ReactNode } from 'react'
import { css, type SerializedStyles } from '@emotion/react'

const block = css({
  display: 'block',
  flexShrink: 0,
  minWidth: 0,
})

export interface SkeletonProps {
  width?: number | string
  height?: number | string
  radius?: number | string
  circle?: boolean
  className?: string
  css?: SerializedStyles
  style?: CSSProperties
}


export function Skeleton({
  width = '100%',
  height = 12,
  radius,
  circle,
  className,
  css: extra,
  style,
}: SkeletonProps) {
  const h = height ?? 12
  return (
    <span
      className={['skeleton', className].filter(Boolean).join(' ')}
      css={[block, extra]}
      style={{
        width: circle ? h : width,
        height: h,
        borderRadius: circle ? '50%' : (radius ?? 'var(--radius-md)'),
        ...style,
      }}
      aria-hidden
    />
  )
}

const stackBase = css({
  display: 'flex',
  flexDirection: 'column',
})

export function SkeletonStack({
  gap = 8,
  children,
  css: extra,
}: {
  gap?: number
  children: ReactNode
  css?: SerializedStyles
}) {
  return (
    <div css={[stackBase, { gap }, extra]} aria-hidden>
      {children}
    </div>
  )
}


export function SkeletonText({
  lines = 3,
  lineHeight = 12,
  gap = 8,
  lastWidth = '62%',
}: {
  lines?: number
  lineHeight?: number
  gap?: number
  lastWidth?: number | string
}) {
  return (
    <SkeletonStack gap={gap}>
      {Array.from({ length: lines }, (_, i) => (
        <Skeleton
          key={i}
          height={lineHeight}
          width={i === lines - 1 ? lastWidth : '100%'}
        />
      ))}
    </SkeletonStack>
  )
}

export function SkeletonCircle({ size = 40 }: { size?: number }) {
  return <Skeleton circle height={size} width={size} />
}

export function SkeletonButton({ width = 96 }: { width?: number | string }) {
  return <Skeleton height={34} width={width} radius="var(--radius-md)" />
}

export function SkeletonCard({
  height = 150,
  css: extra,
}: {
  height?: number | string
  css?: SerializedStyles
}) {
  return <Skeleton height={height} width="100%" radius="var(--radius-lg)" css={extra} />
}
