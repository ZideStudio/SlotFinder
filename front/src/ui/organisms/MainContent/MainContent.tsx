import { type PropsWithChildren } from 'react';
import classes from './MainContent.module.css';

export const MainContent = ({ children }: PropsWithChildren) => <main className={classes.main}>{children}</main>;
