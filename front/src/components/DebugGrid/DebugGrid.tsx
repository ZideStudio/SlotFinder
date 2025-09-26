import './DebugGrid.scss';

type DebugGridProps = {
  cols?: number;
  isCheckedByDefault?: boolean;
};

// oxlint-disable-next-line no-magic-numbers
export const DebugGrid = ({ cols = 12, isCheckedByDefault = false }: DebugGridProps) => {
  if (!import.meta.env.DEV) {
    return null;
  }

  return (
    <>
      <label className="debug-grid-checkbox">
        Grid <input type="checkbox" name="debuggrid" defaultChecked={isCheckedByDefault} />
      </label>

      <div className="debug-grid" role="presentation">
        <div className="grid">
          {[...Array(cols).keys()].map((col: number) => (
            <div key={col} className="cols" />
          ))}
        </div>
      </div>
    </>
  );
};
