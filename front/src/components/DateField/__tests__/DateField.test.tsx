import { render, screen } from '@testing-library/react';
import { ReactNode } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { DateField } from '../DateField';
import userEvent from '@testing-library/user-event';

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

describe('DateField', () => {
  it('renders without crashing', () => {
    render(
      <FormWrapper>
        <DateField name="birthDate" label="Birth Date" placeholder="Select your birth date" />
      </FormWrapper>,
    );

    expect(screen.getByLabelText('Birth Date')).toBeInTheDocument();
    expect(screen.getByLabelText('Birth Date')).toHaveAttribute('name', 'birthDate');
    expect(screen.getByPlaceholderText('Select your birth date')).toBeInTheDocument();
  });

  it('displays the error message from form state when validation fails', async () => {
    const WrapperWithError = () => {
      const methods = useForm({ defaultValues: { birthDate: '' } });

      const { setError } = methods;

      return (
        <FormProvider {...methods}>
          <DateField name="birthDate" label="Birth Date" />
          <button onClick={() => setError('birthDate', { message: 'This field is required' })}>Trigger error</button>
        </FormProvider>
      );
    };

    render(<WrapperWithError />);

    await screen.findByRole('button', { name: 'Trigger error' }).then(button => button.click());

    await expect(screen.findByText('This field is required')).resolves.toBeInTheDocument();
  });

  it('updates the input value on user typing', async () => {
    render(
      <FormWrapper>
        <DateField name="eventDate" label="Event Date" />
      </FormWrapper>,
    );

    const input = screen.getByLabelText('Event Date');
    await userEvent.type(input, '2024-01-01');
    expect(input).toHaveValue('2024-01-01');
  });
});
