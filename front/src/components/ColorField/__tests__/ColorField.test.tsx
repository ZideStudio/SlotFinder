import { render, screen, fireEvent } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';
import { useForm, FormProvider } from 'react-hook-form';
import { type ReactNode } from 'react';
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
});
