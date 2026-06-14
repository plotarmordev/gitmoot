import React, {type ReactNode} from 'react';
import Link from '@docusaurus/Link';
import useBaseUrl from '@docusaurus/useBaseUrl';
import {useThemeConfig, type FooterLinkItem} from '@docusaurus/theme-common';

const footerTagline = 'Local-first coordination for AI agents working through GitHub.';

function opensInNewTab(href: string): boolean {
  return /^https?:\/\//.test(href) && !href.startsWith('https://gitmoot.io');
}

function FooterLink({item}: {item: FooterLinkItem}): ReactNode {
  if (item.href) {
    return (
      <Link
        className="gitmoot-footer__link"
        href={item.href}
        target={opensInNewTab(item.href) ? '_blank' : undefined}
        rel={opensInNewTab(item.href) ? 'noreferrer' : undefined}>
        {item.label}
      </Link>
    );
  }

  if (item.to) {
    return (
      <Link className="gitmoot-footer__link" to={item.to}>
        {item.label}
      </Link>
    );
  }

  return null;
}

export default function Footer(): ReactNode {
  const {footer} = useThemeConfig();
  const logoSrc = useBaseUrl(footer?.logo?.src ?? '/img/gitmoot-logo.svg');

  if (!footer) {
    return null;
  }

  const links = footer.links as FooterLinkItem[];

  return (
    <footer className="footer gitmoot-footer">
      <div className="container container-fluid gitmoot-footer__inner">
        <div className="gitmoot-footer__brand">
          <Link className="gitmoot-footer__wordmark" to={footer.logo?.href ?? '/intro'}>
            <img
              src={logoSrc}
              alt=""
              width={footer.logo?.width ?? 28}
              height={footer.logo?.height ?? 28}
            />
            <span>Gitmoot</span>
          </Link>
          <p>{footerTagline}</p>
        </div>

        <div className="gitmoot-footer__meta">
          <nav className="gitmoot-footer__links" aria-label="Footer navigation">
            {links.map((item) => (
              <FooterLink key={`${item.label}-${item.href ?? item.to}`} item={item} />
            ))}
          </nav>
          {footer.copyright && <p>{footer.copyright}</p>}
        </div>
      </div>
    </footer>
  );
}
