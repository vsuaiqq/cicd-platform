

const S = { width: 20, height: 20, fill: 'none' as const, 'aria-hidden': true }

export function IconDashboard() {
  return (
    <svg {...S} viewBox="0 0 20 20">
      <rect x="2.5" y="2.5" width="6" height="6" rx="1.5" stroke="currentColor" strokeWidth="1.25" />
      <rect x="11.5" y="2.5" width="6" height="6" rx="1.5" stroke="currentColor" strokeWidth="1.25" />
      <rect x="2.5" y="11.5" width="6" height="6" rx="1.5" stroke="currentColor" strokeWidth="1.25" />
      <rect x="11.5" y="11.5" width="6" height="6" rx="1.5" stroke="currentColor" strokeWidth="1.25" />
    </svg>
  )
}

export function IconPlus() {
  return (
    <svg {...S} viewBox="0 0 20 20">
      <path d="M10 4.5v11M4.5 10h11" stroke="currentColor" strokeWidth="1.25" strokeLinecap="round" />
    </svg>
  )
}

export function IconDocs() {
  return (
    <svg {...S} viewBox="0 0 20 20">
      <path d="M6 3.5h6.5L15.5 7v9.5a1 1 0 01-1 1h-8a1 1 0 01-1-1v-11a1 1 0 011-1z" stroke="currentColor" strokeWidth="1.25" strokeLinejoin="round" />
      <path d="M12.5 3.5V7H16M7 10.5h6M7 13.5h4" stroke="currentColor" strokeWidth="1.25" strokeLinecap="round" />
    </svg>
  )
}

export function IconBell() {
  return (
    <svg {...S} viewBox="0 0 20 20">
      <path d="M10 3.5a4.5 4.5 0 00-4.5 4.5v2.2L3.5 13h13l-2-2.8V8a4.5 4.5 0 00-4.5-4.5z" stroke="currentColor" strokeWidth="1.25" strokeLinejoin="round" />
      <path d="M8.2 14.5a1.8 1.8 0 003.6 0" stroke="currentColor" strokeWidth="1.25" strokeLinecap="round" />
    </svg>
  )
}

export function IconSettings() {
  return (
    <svg {...S} viewBox="0 0 20 20">
      <circle cx="10" cy="10" r="2.5" stroke="currentColor" strokeWidth="1.25" />
      <path d="M10 4v1.2M10 14.8V16M4 10h1.2M14.8 10H16M5.8 5.8l.85.85M13.35 13.35l.85.85M5.8 14.2l.85-.85M13.35 6.65l.85-.85" stroke="currentColor" strokeWidth="1.25" strokeLinecap="round" />
    </svg>
  )
}

export function IconSearch() {
  return (
    <svg {...S} viewBox="0 0 20 20">
      <circle cx="9" cy="9" r="5" stroke="currentColor" strokeWidth="1.25" />
      <path d="M13 13l3.5 3.5" stroke="currentColor" strokeWidth="1.25" strokeLinecap="round" />
    </svg>
  )
}

export function IconLogout() {
  return (
    <svg {...S} viewBox="0 0 20 20">
      <path d="M12.5 4h2.5a1 1 0 011 1v10a1 1 0 01-1 1h-2.5" stroke="currentColor" strokeWidth="1.25" strokeLinecap="round" />
      <path d="M8.5 13l2.5-3-2.5-3M11 10H4.5" stroke="currentColor" strokeWidth="1.25" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

export function IconChevronRight() {
  return (
    <svg width={16} height={16} viewBox="0 0 16 16" fill="none" aria-hidden>
      <path d="M6 4l4 4-4 4" stroke="currentColor" strokeWidth="1.25" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

export function IconPanelClose({ collapsed }: { collapsed: boolean }) {
  return (
    <svg width={18} height={18} viewBox="0 0 20 20" fill="none" aria-hidden style={{
      transform: collapsed ? 'rotate(180deg)' : 'none',
      transition: 'transform var(--duration-sidebar) var(--ease-sidebar)',
    }}>
      <path d="M12.5 5.5L8 10l4.5 4.5" stroke="currentColor" strokeWidth="1.25" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

export function IconSun() {
  return (
    <svg width={18} height={18} viewBox="0 0 20 20" fill="none" aria-hidden>
      <circle cx="10" cy="10" r="3.5" stroke="currentColor" strokeWidth="1.25" />
      <path d="M10 3v1.5M10 15.5V17M3 10h1.5M15.5 10H17M5.2 5.2l1 1M13.8 13.8l1 1M5.2 14.8l1-1M13.8 6.2l1-1" stroke="currentColor" strokeWidth="1.25" strokeLinecap="round" />
    </svg>
  )
}

export function IconMoon() {
  return (
    <svg width={18} height={18} viewBox="0 0 20 20" fill="none" aria-hidden>
      <path d="M16 12.5A6.5 6.5 0 017.5 4 6.5 6.5 0 1016 12.5z" stroke="currentColor" strokeWidth="1.25" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}
