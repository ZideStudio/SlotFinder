import { Button } from '@Front/ui/molecules/Button/Button';
import AddCalendarIcon from '@material-symbols/svg-400/outlined/calendar_add_on.svg?react';
import Person from '@material-symbols/svg-400/outlined/person.svg?react';
import logo from '../../../../public/assets/logo.png';

import './Header.scss';

export const Header = () => (
  <header className="header">
    <img src={logo} alt="Slot Finder logo" className="header__logo" />
    <div className="header__buttons">
      <Button icon={AddCalendarIcon} variant="secondary" aria-label="add event" className="header__button" />
      <Button icon={Person} variant="secondary" aria-label="user profile" className="header__button" />
    </div>
  </header>
);
