import type { ReactNode } from "react";
import Link from "@docusaurus/Link";
import useDocusaurusContext from "@docusaurus/useDocusaurusContext";
import Layout from "@theme/Layout";
import Heading from "@theme/Heading";

import styles from "./index.module.scss";

export default function Home(): ReactNode {
  const { siteConfig } = useDocusaurusContext();
  return (
    <Layout
      title={`${siteConfig.title}`}
      description="Yiu Operations CLI 是一个运维命令行工具，旨在简化和自动化各种操作任务。"
    >
      <div className={styles.wrapper}>
        <div className={styles.logoWrapper}>
          <img
            className={styles.logo}
            src="/img/Yiu/icononly_transparent_nobuffer.png"
            alt={`${siteConfig.title} Logo`}
          />
        </div>
        <Heading as="h1" className="hero__title">
          {siteConfig.title}
        </Heading>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <div className={styles.buttons}>
          <Link
            className="button button--secondary button--lg"
            to="/docs/intro"
          >
            快速开始
          </Link>
        </div>
      </div>
    </Layout>
  );
}
