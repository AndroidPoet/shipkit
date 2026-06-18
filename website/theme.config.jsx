import { useConfig } from 'nextra-theme-docs'
import { useRouter } from 'next/router'

const Logo = () => (
  <span style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', fontWeight: 700 }}>
    <svg width="22" height="22" viewBox="0 0 24 24" fill="none" aria-hidden="true">
      <path
        d="M12 2 3 6.5v5.8c0 4.9 3.6 8.3 9 9.7 5.4-1.4 9-4.8 9-9.7V6.5L12 2Zm-1 13.4-3.5-3.5L9 10.4l2 2 4-4 1.5 1.5L11 15.4Z"
        fill="#0EA5A4"
      />
    </svg>
    <span>Shipkit</span>
  </span>
)

export default {
  logo: <Logo />,
  project: {
    link: 'https://github.com/AndroidPoet/shipkit',
  },
  docsRepositoryBase: 'https://github.com/AndroidPoet/shipkit/tree/main/website',
  color: {
    hue: 178,
    saturation: 85,
  },
  footer: {
    content: (
      <span>
        MIT © {new Date().getFullYear()}{' '}
        <a href="https://github.com/AndroidPoet/shipkit" target="_blank" rel="noreferrer">
          Shipkit
        </a>
        . The release cockpit for mobile apps.
      </span>
    ),
  },
  head: function useHead() {
    const { frontMatter } = useConfig()
    const { asPath } = useRouter()
    const pageTitle = frontMatter?.title
    const title = pageTitle ? `${pageTitle} – Shipkit` : 'Shipkit'
    const description =
      frontMatter?.description ??
      'Shipkit — one AI-agent-friendly command surface for Google Play, App Store Connect, RevenueCat, and CI release automation for mobile apps.'
    const base = 'https://androidpoet.github.io/shipkit'
    const path = asPath === '/' ? '' : asPath.split('?')[0].split('#')[0]
    const canonical = `${base}${path}`
    const ogImage = `${base}/favicon.svg`
    return (
      <>
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <title>{title}</title>
        <meta name="description" content={description} />
        <link rel="canonical" href={canonical} />
        <link rel="icon" href={`${base}/favicon.svg`} type="image/svg+xml" />
        <meta name="theme-color" content="#0EA5A4" />
        <meta property="og:type" content="website" />
        <meta property="og:site_name" content="Shipkit" />
        <meta property="og:url" content={canonical} />
        <meta property="og:title" content={pageTitle ?? 'Shipkit'} />
        <meta property="og:description" content={description} />
        <meta property="og:image" content={ogImage} />
        <meta name="twitter:card" content="summary_large_image" />
        <meta name="twitter:title" content={pageTitle ?? 'Shipkit'} />
        <meta name="twitter:description" content={description} />
        <meta name="twitter:image" content={ogImage} />
      </>
    )
  },
  sidebar: {
    defaultMenuCollapseLevel: 1,
  },
  toc: {
    backToTop: true,
  },
  navigation: {
    prev: true,
    next: true,
  },
  darkMode: true,
}
