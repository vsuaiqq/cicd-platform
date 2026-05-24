import { useEffect, useRef, type ReactNode } from 'react'
import { useI18n } from '../../i18n'
import { IconEnv, IconGeneral, IconPipeline, IconSecrets, IconTeam } from './icons'
import {
  settingsLayout,
  settingsNav,
  settingsNavBtn,
  settingsNavIcon,
  settingsNavLabel,
  settingsPanel,
  settingsPanelBody,
  settingsPanelDesc,
  settingsPanelHead,
  settingsPanelTitle,
} from './styles'

export type SettingsSectionId = 'general' | 'env' | 'secrets' | 'team' | 'pipeline'

interface SectionDef {
  id: SettingsSectionId
  navLabel: string
  title: string
  desc: string
  icon: ReactNode
}

interface SettingsTabProps {
  active: SettingsSectionId
  onChange: (id: SettingsSectionId) => void
  children: ReactNode
}

export function SettingsTab({ active, onChange, children }: SettingsTabProps) {
  const { t } = useI18n()
  const navRef = useRef<HTMLElement>(null)

  useEffect(() => {
    const nav = navRef.current
    if (!nav || !window.matchMedia('(max-width: 900px)').matches) return
    const tab = nav.querySelector<HTMLElement>(`#settings-tab-${active}`)
    tab?.scrollIntoView({ behavior: 'smooth', inline: 'nearest', block: 'nearest' })
  }, [active])

  const sections: SectionDef[] = [
    {
      id: 'general',
      navLabel: t.project.settingsNavGeneral,
      title: t.project.settingsGeneral,
      desc: t.project.settingsGeneralDesc,
      icon: <IconGeneral />,
    },
    {
      id: 'env',
      navLabel: t.project.settingsNavEnv,
      title: t.project.panelEnvVars,
      desc: t.project.envVarsDesc,
      icon: <IconEnv />,
    },
    {
      id: 'secrets',
      navLabel: t.project.settingsNavSecrets,
      title: t.project.panelSecrets,
      desc: t.project.secretsDesc,
      icon: <IconSecrets />,
    },
    {
      id: 'team',
      navLabel: t.project.settingsNavTeam,
      title: t.members.panelTitle,
      desc: t.project.settingsTeamDesc,
      icon: <IconTeam />,
    },
    {
      id: 'pipeline',
      navLabel: t.project.settingsNavPipeline,
      title: t.project.panelPipelineConfig,
      desc: t.project.pipelineYamlDesc,
      icon: <IconPipeline />,
    },
  ]

  const current = sections.find((s) => s.id === active) ?? sections[0]

  return (
    <div css={settingsLayout}>
      <nav ref={navRef} css={settingsNav} role="tablist" aria-label={t.project.settingsNavLabel}>
        {sections.map((s) => {
          const isActive = active === s.id
          return (
            <button
              key={s.id}
              type="button"
              role="tab"
              id={`settings-tab-${s.id}`}
              aria-selected={isActive}
              aria-controls={`settings-panel-${s.id}`}
              title={s.title}
              css={settingsNavBtn(isActive)}
              onClick={() => onChange(s.id)}
            >
              <span css={settingsNavIcon}>{s.icon}</span>
              <span css={settingsNavLabel}>{s.navLabel}</span>
            </button>
          )
        })}
      </nav>
      <div
        css={settingsPanel}
        role="tabpanel"
        id={`settings-panel-${current.id}`}
        aria-labelledby={`settings-tab-${current.id}`}
      >
        <div css={settingsPanelHead}>
          <h2 css={settingsPanelTitle}>{current.title}</h2>
          <p css={settingsPanelDesc}>{current.desc}</p>
        </div>
        <div css={settingsPanelBody}>{children}</div>
      </div>
    </div>
  )
}
