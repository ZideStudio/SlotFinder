#!/usr/bin/env bash
echo "Starting post-create script..."

if ! command -v rtk >/dev/null 2>&1; then
  echo "Warning: RTK is not installed; continuing without automatic Copilot hook setup."
  exit 0
fi

rtk init --copilot --auto-patch >/tmp/rtk-init.log 2>&1
status=$?

if [ "$status" -ne 0 ]; then
  echo "Warning: RTK initialization failed; continuing without automatic Copilot hook setup."
  if [ -s /tmp/rtk-init.log ]; then
    cat /tmp/rtk-init.log
  fi
fi

echo "Post-create script successfully completed."

exit 0
