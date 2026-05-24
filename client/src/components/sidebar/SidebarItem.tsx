import { NavLink, type NavLinkProps } from 'react-router-dom'
import type { ReactNode } from 'react'
import * as s from './styles'

type Base = {
  icon: ReactNode
  label: string
  title?: string
  suffix?: ReactNode
  className?: string
  danger?: boolean
}

type LinkItem = Base & {
  to: NavLinkProps['to']
  end?: boolean
  onClick?: never
}

type ButtonItem = Base & {
  to?: never
  end?: never
  onClick: () => void
  active?: boolean
}

export type SidebarItemProps = LinkItem | ButtonItem

export function SidebarItem(props: SidebarItemProps) {
  const { icon, label, title, suffix, className } = props
  const body = (
    <>
      <span data-slot="icon" css={s.itemIcon}>{icon}</span>
      <span css={s.itemLabel}>{label}</span>
      {suffix ? <span css={s.itemSuffix}>{suffix}</span> : null}
    </>
  )

  if ('to' in props && props.to !== undefined) {
    return (
      <NavLink
        to={props.to}
        end={props.end}
        title={title ?? label}
        css={props.danger ? s.logoutItem : s.itemLink}
        className={({ isActive }) => [isActive ? 'active' : '', className].filter(Boolean).join(' ')}
      >
        {body}
      </NavLink>
    )
  }

  return (
    <button
      type="button"
      title={title ?? label}
      css={props.danger ? s.logoutItem : s.itemButton}
      className={[props.active ? 'active' : '', className].filter(Boolean).join(' ') || undefined}
      onClick={props.onClick}
      aria-expanded={props.active}
    >
      {body}
    </button>
  )
}
