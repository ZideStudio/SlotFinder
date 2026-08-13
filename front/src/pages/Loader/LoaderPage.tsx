import React from "react";
import "./LoaderPage.scss";
import LogoColor from "../../../public/assets/colored_logo_no_bg.svg";

const LoaderPage: React.FC = () => (
  <div className="loader-page">
    <div className="loader-logo-wrapper" aria-label="loading">
      <div
        className="loader-page__liquid-container"
        style={{ 
          "--logo-url": `url(${LogoColor})`,
         } as React.CSSProperties}
      >
        <div className="loader-page__liquid-fill">
          <div className="wave wave-1" />
          <div className="wave wave-2" />
        </div>
      </div>
    </div>
  </div>
);

export default LoaderPage;
