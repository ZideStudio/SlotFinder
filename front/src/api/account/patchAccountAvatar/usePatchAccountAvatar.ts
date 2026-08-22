import { patchV1AccountAvatar } from "@Front/api/generated/account/account";
import type { PatchV1AccountAvatarBody } from "@Front/api/generated/slotFinderAPI.schemas";
import type { ErrorResponse } from "@Front/types/ErrorResponse";
import { useMutation } from "@tanstack/react-query";

type UsePatchAccountAvatarProps = {
  onError?: (error: ErrorResponse) => void;
};

type UsePatchAccountAvatarReturn = {
  patchAccountAvatar: (data: PatchV1AccountAvatarBody) => void;
  isLoading: boolean;
};

export const usePatchAccountAvatar = ({
  onError,
}: UsePatchAccountAvatarProps = {}): UsePatchAccountAvatarReturn => {
  const mutation = useMutation<void, ErrorResponse, PatchV1AccountAvatarBody>({
    mutationKey: ["accountAvatarPatch"],
    mutationFn: (data) => patchV1AccountAvatar(data),
    onError,
  });

  return {
    patchAccountAvatar: mutation.mutate,
    isLoading: mutation.isPending,
  };
};
