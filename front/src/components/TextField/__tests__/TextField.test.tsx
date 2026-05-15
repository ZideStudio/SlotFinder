import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { useForm, FormProvider } from 'react-hook-form';
import { TextField } from '../TextField';
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

describe('TextField', () => {
  it('renders without crashing', () => {
    render(
      <FormWrapper>
        <TextField name="username" label="Username" />
      </FormWrapper>,
    );

    expect(screen.getByLabelText('Username')).toBeInTheDocument();
  });

  it('renders with the correct name attribute', () => {
    render(
      <FormWrapper>
        <TextField name="email" label="Email" />
      </FormWrapper>,
    );

    expect(screen.getByLabelText('Email')).toHaveAttribute('name', 'email');
  });

  it('passes extra props down to TextInput', () => {
    render(
      <FormWrapper>
        <TextField name="search" label="Search" placeholder="Search…" />
      </FormWrapper>,
    );

    expect(screen.getByPlaceholderText('Search…')).toBeInTheDocument();
  });

  it('displays the error message from form state when validation fails', async () => {
    const WrapperWithError = () => {
      const methods = useForm({ defaultValues: { username: '' } });

      const { setError } = methods;

      return (
        <FormProvider {...methods}>
          <TextField name="username" label="Username" />
          <button onClick={() => setError('username', { message: 'This field is required' })}>Trigger error</button>
        </FormProvider>
      );
    };

    render(<WrapperWithError />);

    await userEvent.click(screen.getByRole('button', { name: 'Trigger error' }));

    expect(await screen.findByText('This field is required')).toBeInTheDocument();
  });

  it('updates the input value on user typing', async () => {
    render(
      <FormWrapper>
        <TextField name="firstName" label="First Name" />
      </FormWrapper>,
    );

    const input = screen.getByLabelText('First Name');
    await userEvent.type(input, 'Test');

    expect(input).toHaveValue('Test');
  });
});
