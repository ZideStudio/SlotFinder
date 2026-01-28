import { TextInputAtom, TextInputAtomProps } from '@Front/ui/atoms/Inputs/TextInputAtom/TextInputAtom';
import { Field } from '@Front/ui/utils/components/Field';

export const TextInput = (props: TextInputAtomProps) => {
  return <Field input={TextInputAtom} {...props} />;
};
