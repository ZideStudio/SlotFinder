import React, { useEffect, useState } from "react";
import "./LoaderPage.scss";
import LogoColorUrl from "@Front/assets/svg/logo/colored_logo_no_bg.svg?url";
import { useTranslation } from "react-i18next";

const LoaderPage = () => {
  const [visible, setVisible] = useState(false);
  const { t } = useTranslation("loader");

  useEffect(() => {
    const timer = setTimeout(() => {
      setVisible(true);
    }, 100);

    return () => clearTimeout(timer);
  }, []);

  if (!visible) {
    return null;
  }

  return (
    <div className="loader-page">
      <div
        className="loader-logo-wrapper"
        role="status" // oxlint-disable-line jsx-a11y/prefer-tag-over-role
        aria-label={t("loading")}
      >
        <div
          className="loader-page__liquid-container"
          style={
            {
              "--logo-url": `url(${LogoColorUrl})`,
            } as React.CSSProperties
          }
        >
          <div className="loader-page__liquid-fill">
            <div className="wave wave-1" />
            <div className="wave wave-2" />
          </div>
        </div>
      </div>
      {/* oxlint-disable-next-line jsx-a11y/prefer-tag-over-role*/}
    </div>
  );
};

export default LoaderPage;
