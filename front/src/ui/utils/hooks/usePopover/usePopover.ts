import { useId, type CSSProperties } from "react";

type TriggerProps = {
  "data-popover-trigger": "";
  popoverTarget: string;
  style: CSSProperties;
};

type PopoverProps = {
  id: string;
  style: CSSProperties;
};

type UsePopoverReturnType = {
  triggerProps: TriggerProps;
  popoverProps: PopoverProps;
};

export const usePopover = (): UsePopoverReturnType => {
  const id = useId();
  const popoverId = `popover-${id}`;
  const anchorName = `--popover-${id.replaceAll(":", "")}`;
  const anchorStyle = { "--popover-anchor-name": anchorName } as CSSProperties;

  return {
    triggerProps: {
      "data-popover-trigger": "",
      popoverTarget: popoverId,
      style: anchorStyle,
    },
    popoverProps: {
      id: popoverId,
      style: anchorStyle,
    },
  };
};
