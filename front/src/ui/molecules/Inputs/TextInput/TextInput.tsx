import { TextInputAtom } from '@Front/ui/atoms/Inputs/TextInputAtom/TextInputAtom';
import { Field } from '@Front/ui/utils/components/Field';

export const TextInput = ({ ...textInputAtomProps }) => {
  return <Field input={TextInputAtom} {...textInputAtomProps} />;
};
