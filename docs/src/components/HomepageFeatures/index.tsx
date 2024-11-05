import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  Svg: React.ComponentType<React.ComponentProps<'svg'>>;
  description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Simple Management Interface',
    Svg: require('@site/static/img/dashboard.svg').default,
    description: (
      <>
        Mantrae provides an intuitive UI for managing Traefik routers, middlewares,
        and DNS entries, making setup and configuration straightforward.
      </>
    ),
  },
  {
    title: 'Flexible DNS Automation',
    Svg: require('@site/static/img/dns.svg').default,
    description: (
      <>
        Automate DNS entries with Cloudflare, PowerDNS, or Technitium integration,
        allowing seamless domain configuration for your routers. Check out the
        <a href="/docs/usage/dns"> dns provider documentation</a>
      </>
    ),
  },
  {
    title: 'Distributed Container Management',
    Svg: require('@site/static/img/container.svg').default,
    description: (
      <>
        Use Mantrae Agents to gather container information from multiple hosts,
        centralizing Traefik label management even without Traefik installed.
        Check out the <a href="/docs/usage/agents">agent documentation</a>
      </>
    ),
  },
];

function Feature({ title, Svg, description }: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <Svg className={styles.featureSvg} role="img" />
      </div>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): JSX.Element {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
