import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { type ReactNode } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { SelectField } from '../SelectField';

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

describe('SelectField', () => {
  const options = [
    { value: 'male', label: 'Male' },
    { value: 'female', label: 'Female' },
  ];
  it('renders without crashing', () => {
    render(
      <FormWrapper>
        <SelectField name="gender" label="Gender" options={options} />
      </FormWrapper>,
    );

    expect(screen.getByLabelText('Gender')).toBeInTheDocument();
    expect(screen.getByLabelText('Gender')).toHaveAttribute('name', 'gender');
  });

  it('displays the error message from form state when validation fails', async () => {
    const WrapperWithError = () => {
      const methods = useForm({ defaultValues: { gender: '' } });

      const { setError } = methods;

      return (
        <FormProvider {...methods}>
          <SelectField name="gender" label="Gender" options={options} />
          <button onClick={() => setError('gender', { message: 'This field is required' })}>Trigger error</button>
        </FormProvider>
      );
    };

    render(<WrapperWithError />);

    await userEvent.click(screen.getByRole('button', { name: 'Trigger error' }));

    await expect(screen.findByText('This field is required')).resolves.toBeInTheDocument();
  });

  it('updates the input value on user typing', async () => {
    render(
      <FormWrapper>
        <SelectField name="gender" label="Gender" options={options} />
      </FormWrapper>,
    );

    const input = screen.getByLabelText('Gender');
    await userEvent.selectOptions(input, 'female');

    expect(input).toHaveValue('female');
  });
});
