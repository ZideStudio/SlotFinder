import { Button } from '@Front/ui/molecules/Button/Button';
import AddCalendarIcon from '@material-symbols/svg-400/outlined/calendar_add_on.svg?react';
import Person from '@material-symbols/svg-400/outlined/person.svg?react';
import logo from '../../../assets/images/SLOT_FINDER-V24x.png';

import './Header.scss';

export const Header = () => (
  <div className="header">
    <img src={logo} alt="Slot Finder" className="logo" />
    <div className="buttons">
      <Button icon={AddCalendarIcon} variant="secondary" aria-label="add event" className='button' />
      <Button icon={Person} variant="secondary" aria-label="user profile" className='button' />
    </div>
  </div>
);
