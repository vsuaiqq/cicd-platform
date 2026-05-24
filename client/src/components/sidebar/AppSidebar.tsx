import { NavLink } from 'react-router-dom'
import { FlowBrandLogo } from '../FlowBrand'
import { useI18n } from '../../i18n'
import { SidebarItem } from './SidebarItem'
import { NotificationBell, type RunNotification } from './NotificationBell'
import {
  IconDashboard,
  IconDocs,
  IconLogout,
  IconPanelClose,
  IconSearch,
  IconSettings,
  IconPlus,
  IconSun,
  IconMoon,
  IconChevronRight,
} from './icons'
import * as s from './styles'

export interface AppSidebarProps {
  collapsed: boolean
  onToggleCollapse: () => void
  profileDisplayName: string
  profileSubline: string
  profileInitials: string
  connectedTitle: string
  theme: 'light' | 'dark'
  onToggleTheme: () => void
  lang: 'en' | 'ru'
  onToggleLang: () => void
  onOpenSearch: () => void
  onOpenSettings: () => void
  onOpenHelp: () => void
  onLogout: () => void
  notifications: RunNotification[]
  notifPanelOpen: boolean
  onNotifToggle: () => void
  onNotifMarkAllRead: (e: React.MouseEvent) => void
  onNotifUpdate: (fn: (prev: RunNotification[]) => RunNotification[]) => void
}

export function AppSidebar({
  collapsed,
  onToggleCollapse,
  profileDisplayName,
  profileSubline,
  profileInitials,
  connectedTitle,
  theme,
  onToggleTheme,
  lang,
  onToggleLang,
  onOpenSearch,
  onOpenSettings,
  onOpenHelp,
  onLogout,
  notifications,
  notifPanelOpen,
  onNotifToggle,
  onNotifMarkAllRead,
  onNotifUpdate,
}: AppSidebarProps) {
  const { t } = useI18n()

  return (
    <aside
      className="app-sidebar"
      css={s.sidebarRoot}
      data-collapsed={collapsed || undefined}
      aria-label={t.aria.primaryNav}
    >
      <header css={s.sidebarHeader}>
        <NavLink to="/" css={s.brandLink} title="Flow">
          <FlowBrandLogo collapsed={collapsed} />
        </NavLink>
        <button
          type="button"
          css={s.collapseButton}
          onClick={onToggleCollapse}
          aria-label={collapsed ? t.nav.expandSidebar : t.nav.collapseSidebar}
          title={collapsed ? t.nav.expandSidebar : t.nav.collapseSidebar}
        >
          <IconPanelClose collapsed={collapsed} />
        </button>
      </header>

      <div css={s.sidebarScroll}>
        <p css={s.sectionTitle}>{t.nav.overview}</p>
        <SidebarItem to="/" end icon={<IconDashboard />} label={t.nav.dashboard} />

        <div css={s.divider} />

        <p css={s.sectionTitle}>{t.nav.projects}</p>
        <SidebarItem to="/projects/new" icon={<IconPlus />} label={t.nav.newProject} />

        <div css={[s.divider, { marginTop: 'auto' }]} />

        <p css={s.sectionTitle}>{t.nav.help}</p>
        <SidebarItem icon={<IconDocs />} label={t.nav.docs} onClick={onOpenHelp} />
      </div>

      <footer css={s.sidebarFooter}>
        <NotificationBell
          notifications={notifications}
          panelOpen={notifPanelOpen}
          onToggle={onNotifToggle}
          onMarkAllRead={onNotifMarkAllRead}
          onUpdate={onNotifUpdate}
        />

        <SidebarItem icon={<IconSettings />} label={t.nav.settings} onClick={onOpenSettings} />

        <SidebarItem
          icon={<IconSearch />}
          label={t.nav.quickSearch}
          title={`${t.nav.quickSearch} (⌘K)`}
          onClick={onOpenSearch}
          suffix={<span css={s.kbHint}>⌘K</span>}
        />

        <div css={s.prefsRow}>
          <button
            type="button"
            css={s.prefChip}
            onClick={onToggleTheme}
            title={theme === 'dark' ? t.prefs.light : t.prefs.dark}
            aria-label={theme === 'dark' ? t.prefs.light : t.prefs.dark}
          >
            {theme === 'dark' ? <IconMoon /> : <IconSun />}
            <span data-pref-label>{theme === 'dark' ? t.prefs.dark : t.prefs.light}</span>
          </button>
          <button
            type="button"
            css={s.prefChip}
            onClick={onToggleLang}
            title={t.prefs.language}
            aria-label={t.prefs.language}
          >
            {lang === 'en' ? 'EN' : 'RU'}
          </button>
        </div>

        <NavLink
          to="/profile"
          title={profileDisplayName}
          css={s.itemLink}
          className={({ isActive }) => (isActive ? 'active' : '')}
        >
          <span data-slot="icon" css={s.itemIcon}>
            <span css={s.avatar}>
              {profileInitials}
              <span css={s.avatarStatus} title={connectedTitle} />
            </span>
          </span>
          <span css={s.profileMeta}>
            <span css={s.profileName}>{profileDisplayName}</span>
            <span css={s.profileSub}>{profileSubline}</span>
          </span>
          <span css={s.itemSuffix}>
            <IconChevronRight />
          </span>
        </NavLink>

        <SidebarItem
          icon={<IconLogout />}
          label={t.nav.signOut}
          onClick={onLogout}
          danger
        />
      </footer>
    </aside>
  )
}
