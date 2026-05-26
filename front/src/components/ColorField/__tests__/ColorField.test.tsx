import { fireEvent, render, screen } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';
import { type ReactNode } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { ColorField } from '../ColorField';

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

describe('ColorField', () => {
  it('should render color input with correct label and name attribute', () => {
    render(
      <FormWrapper>
        <ColorField name="color" label="Color" description="Choose a color" />
      </FormWrapper>,
    );

    expect(screen.getByLabelText('Color')).toBeInTheDocument();
    expect(screen.getByLabelText('Color')).toHaveAttribute('name', 'color');
  });

  it('displays the error message from form state when validation fails', async () => {
    const user = userEvent.setup();

    const WrapperWithError = () => {
      const methods = useForm({ defaultValues: { color: '' } });

      const { setError } = methods;

      return (
        <FormProvider {...methods}>
          <ColorField name="color" label="Color" description="Choose a color" />
          <button onClick={() => setError('color', { message: 'This field is required' })}>Trigger error</button>
        </FormProvider>
      );
    };

    render(<WrapperWithError />);

    await user.click(screen.getByRole('button', { name: 'Trigger error' }));

    const errorMessage = await screen.findByText('This field is required');
    const input = screen.getByLabelText('Color');

    expect(errorMessage).toBeInTheDocument();
    expect(input).toHaveAttribute('aria-invalid', 'true');
    expect(input).toHaveAttribute('aria-describedby', errorMessage.id);
  });

  it('updates the color value on selection', () => {
    render(
      <FormWrapper>
        <ColorField name="color" label="Color" description="Choose a color" />
      </FormWrapper>,
    );

    const input = screen.getByLabelText('Color');
    fireEvent.change(input, { target: { value: '#ff0000' } });

    expect(input).toHaveValue('#ff0000');
  });

  it('should reflect defaultValues in the color input', () => {
    render(
      <FormWrapper defaultValues={{ color: '#ff0000' }}>
        <ColorField name="color" label="Color" description="Choose a color" />
      </FormWrapper>,
    );

    expect(screen.getByLabelText('Color')).toHaveValue('#ff0000');
  });

  it('should update the color input after reset', async () => {
    const user = userEvent.setup();

    const WrapperWithReset = () => {
      const methods = useForm({ defaultValues: { color: '#ff0000' } });
      return (
        <FormProvider {...methods}>
          <ColorField name="color" label="Color" description="Choose a color" />
          <button onClick={() => methods.reset({ color: '#00ff00' })}>Reset</button>
        </FormProvider>
      );
    };

    render(<WrapperWithReset />);

    await user.click(screen.getByRole('button', { name: 'Reset' }));

    expect(screen.getByLabelText('Color')).toHaveValue('#00ff00');
  });
});
