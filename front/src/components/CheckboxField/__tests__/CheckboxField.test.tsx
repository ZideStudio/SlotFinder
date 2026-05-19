import { render, screen } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';
import { useForm, FormProvider } from 'react-hook-form';
import { CheckboxField } from '../CheckboxField';
import { type ReactNode } from 'react';

const FormWrapper = ({
  children,
  defaultValues = {},
}: {
  children: ReactNode;
  defaultValues?: Record<string, unknown>;
}) => {
  const methods = useForm({ defaultValues });
  return <FormProvider {...methods}>{children}</FormProvider>;
};

describe('CheckboxField', () => {
  it('should render checkbox input with correct label and name attribute', () => {
    render(
      <FormWrapper>
        <CheckboxField name="acceptTerms" label="Accept Terms" />
      </FormWrapper>,
    );

    expect(screen.getByLabelText('Accept Terms')).toBeInTheDocument();
    expect(screen.getByLabelText('Accept Terms')).toHaveAttribute('name', 'acceptTerms');
  });

  it('displays the error message from form state when validation fails', async () => {
    const user = userEvent.setup();

    const WrapperWithError = () => {
      const methods = useForm({ defaultValues: { acceptTerms: false } });

      const { setError } = methods;

      return (
        <FormProvider {...methods}>
          <CheckboxField name="acceptTerms" label="Accept Terms" />
          <button onClick={() => setError('acceptTerms', { message: 'This field is required' })}>Trigger error</button>
        </FormProvider>
      );
    };

    render(<WrapperWithError />);

    await user.click(screen.getByRole('button', { name: 'Trigger error' }));

    const errorMessage = await screen.findByText('This field is required');
    const input = screen.getByRole('checkbox', { name: 'Accept Terms' });

    expect(errorMessage).toBeInTheDocument();
    expect(input).toHaveAttribute('aria-invalid', 'true');
    expect(input).toHaveAttribute('aria-describedby', errorMessage.id);
  });

  it('updates the input checked state on user click', async () => {
    const user = userEvent.setup();

    render(
      <FormWrapper>
        <CheckboxField name="acceptTerms" label="Accept Terms" />
      </FormWrapper>,
    );

    const checkbox = screen.getByLabelText('Accept Terms');
    await user.click(checkbox);

    expect(checkbox).toBeChecked();
  });
});
