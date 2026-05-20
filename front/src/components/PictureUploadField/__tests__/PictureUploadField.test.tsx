import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { type ReactNode } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { PictureUploadField } from '../PictureUploadField';
import { assert } from 'vitest';

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

describe('PictureUploadField', () => {
  it('renders without crashing', () => {
    render(
      <FormWrapper>
        <PictureUploadField name="profilePicture" label="Profile Picture" />
      </FormWrapper>,
    );

    expect(screen.getByLabelText('Profile Picture')).toBeInTheDocument();
    expect(screen.getByLabelText('Profile Picture')).toHaveAttribute('name', 'profilePicture');
  });

  it('displays the error message from form state when validation fails', async () => {
    const WrapperWithError = () => {
      const methods = useForm({ defaultValues: { profilePicture: undefined } });

      const { setError } = methods;

      return (
        <FormProvider {...methods}>
          <PictureUploadField name="profilePicture" label="Profile Picture" />
          <button onClick={() => setError('profilePicture', { message: 'This field is required' })}>
            Trigger error
          </button>
        </FormProvider>
      );
    };

    render(<WrapperWithError />);

    await userEvent.click(screen.getByRole('button', { name: 'Trigger error' }));

    await expect(screen.findByText('This field is required')).resolves.toBeInTheDocument();
  });

  it('updates the input value on user uploading a file', async () => {
    const user = userEvent.setup();
    render(
      <FormWrapper>
        <PictureUploadField name="profilePicture" label="Profile Picture" />
      </FormWrapper>,
    );

    const input = screen.getByLabelText('Profile Picture') as HTMLInputElement;
    const file = new File(['dummy content'], 'example.png', { type: 'image/png' });

    await user.upload(input, file);
    // Assert that the file input's files property has been updated with the uploaded file
    assert(input.files !== null);

    expect(input.files).toHaveLength(1);
    expect(input.files[0]).toBe(file);
  });
});
