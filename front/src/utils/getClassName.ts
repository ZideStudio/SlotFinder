export const clsx = (...classNames: (string | boolean | undefined)[]) => classNames.filter(Boolean).join(' ');

type getClassNameParams = {
  defaultClassName: string;
  modifiers?: (string | boolean | undefined)[] | string;
  className?: string;
};

export const getClassName = ({ defaultClassName, modifiers = [], className }: getClassNameParams) => {
  const formattedModifiers = Array.isArray(modifiers) ? modifiers : modifiers.trim().split(' ');
  const parsedModifiers = formattedModifiers.filter(Boolean).map(modifier => `${defaultClassName}--${modifier}`);

  return clsx(defaultClassName, ...parsedModifiers, className);
};
