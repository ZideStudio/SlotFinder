import { userEvent } from "@vitest/browser/context";
import { type ReactNode } from "react";
import { FormProvider, useForm } from "react-hook-form";
import { render } from "vitest-browser-react";
import { page } from "vitest/browser";
import { PictureUploadField } from "../PictureUploadField";

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

describe("PictureUploadField", () => {
  it("renders without crashing", async () => {
    await render(
      <FormWrapper>
        <PictureUploadField name="profilePicture" label="Profile Picture" />
      </FormWrapper>,
    );

    await expect
      .element(page.getByLabelText("Profile Picture"))
      .toBeInTheDocument();
    await expect
      .element(page.getByLabelText("Profile Picture"))
      .toHaveAttribute("name", "profilePicture");
  });

  it("displays the error message from form state when validation fails", async () => {
    const WrapperWithError = () => {
      const methods = useForm({ defaultValues: { profilePicture: undefined } });

      const { setError } = methods;

      return (
        <FormProvider {...methods}>
          <PictureUploadField name="profilePicture" label="Profile Picture" />
          <button
            onClick={() =>
              setError("profilePicture", { message: "This field is required" })
            }
          >
            Trigger error
          </button>
        </FormProvider>
      );
    };

    await render(<WrapperWithError />);

    await userEvent.click(page.getByRole("button", { name: "Trigger error" }));

    await expect
      .element(page.getByText("This field is required"))
      .toBeInTheDocument();
  });

  it("updates the input value on user uploading a file", async () => {
    await render(
      <FormWrapper>
        <PictureUploadField name="profilePicture" label="Profile Picture" />
      </FormWrapper>,
    );

    const inputLocator = page.getByLabelText("Profile Picture");
    const file = new File(["dummy content"], "example.png", {
      type: "image/png",
    });

    await userEvent.upload(inputLocator, file);

    const inputEl = (await inputLocator.element()) as HTMLInputElement;
    expect(inputEl.files).not.toBeNull();
    expect(inputEl.files).toHaveLength(1);
  });
});
