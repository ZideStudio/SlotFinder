import { getV1AccountMe } from "@Front/api/generated/account/account";
import type { AccountAccountResponseDto } from "@Front/api/generated/slotFinderAPI.schemas";
import { useQuery } from "@tanstack/react-query";

const ONE_HOUR = 60 * 60 * 1000;
export const getAccountMeQueryKey = ["getAccountMe"];

export const useGetAccountMe = () =>
  useQuery<AccountAccountResponseDto>({
    queryKey: getAccountMeQueryKey,
    queryFn: getV1AccountMe,
    staleTime: ONE_HOUR,
    gcTime: ONE_HOUR,
    refetchOnWindowFocus: false,
    retry: 1,
  });
