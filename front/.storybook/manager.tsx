import React, { useCallback } from 'react';
import { ShareAltIcon } from '@storybook/icons';
import { IconButton } from 'storybook/internal/components';
import { addons, types, useStorybookApi } from 'storybook/manager-api';

const ADDON_ID = 'open-canvas-new-tab';
const TOOL_ID = `${ADDON_ID}/tool`;

const OpenCanvasInNewTabTool = () => {
  const api = useStorybookApi();

  const handleClick = useCallback(() => {
    const story = api.getCurrentStoryData();
    if (story) {
      window.open(`./iframe.html?id=${story.id}&viewMode=story`, '_blank');
    }
  }, [api]);

  return (
    <IconButton title="Open canvas in new tab" onClick={handleClick}>
      <ShareAltIcon />
    </IconButton>
  );
};

addons.register(ADDON_ID, () => {
  addons.add(TOOL_ID, {
    type: types.TOOL,
    title: 'Open canvas in new tab',
    render: OpenCanvasInNewTabTool,
  });
});
