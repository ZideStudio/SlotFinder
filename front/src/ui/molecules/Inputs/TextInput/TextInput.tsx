import { TextInputAtom } from '@Front/ui/atoms/Inputs/TextInputAtom/TextInputAtom';
import './TextInput.scss';
import { Field } from '@Front/ui/utils/components/Field';

export const TextInput = ({ ...textInputAtomProps }) => {
  return <Field defaultClassName="ds-text-input" input={TextInputAtom} {...textInputAtomProps} />;
};
