import { userEvent } from "@vitest/browser/context";
import { type ReactNode } from "react";
import { FormProvider, useForm } from "react-hook-form";
import { render } from "vitest-browser-react";
import { screen } from "@testing-library/dom";
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

    expect(screen.getByLabelText("Profile Picture")).toBeInTheDocument();
    expect(screen.getByLabelText("Profile Picture")).toHaveAttribute(
      "name",
      "profilePicture",
    );
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

    await userEvent.click(
      screen.getByRole("button", { name: "Trigger error" }),
    );

    await expect(
      screen.findByText("This field is required"),
    ).resolves.toBeInTheDocument();
  });

  it("updates the input value on user uploading a file", async () => {
    await render(
      <FormWrapper>
        <PictureUploadField name="profilePicture" label="Profile Picture" />
      </FormWrapper>,
    );

    const input = screen.getByLabelText("Profile Picture") as HTMLInputElement;
    const file = new File(["dummy content"], "example.png", {
      type: "image/png",
    });

    await userEvent.upload(input, file);

    expect(input.files).not.toBeNull();
    expect(input.files).toHaveLength(1);
  });
});
