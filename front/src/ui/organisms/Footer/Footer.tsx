import { type ReactNode } from 'react';
import classes from './Footer.module.css';

type FooterProps = {
  footerContent?: ReactNode;
};

export const Footer = ({ footerContent }: FooterProps) => <footer className={classes.footer}>{footerContent}</footer>;
