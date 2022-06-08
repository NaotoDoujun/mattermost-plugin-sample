import React from 'react';
import manifest from './manifest';
import { openRootModal } from './actions';
import reducer from './reducer';
import Root from './components/root';

// Courtesy of https://feathericons.com/
const Icon = () => <i className='icon fa fa-plug' />;

class Plugin {

  initialize(registry, store) {

    registry.registerRootComponent(Root);

    registry.registerChannelHeaderButtonAction(
      <Icon />,
      () => store.dispatch(openRootModal()),
      "Similar Words Search"
    );

    registry.registerReducer(reducer);
  }
}

window.registerPlugin(manifest.id, new Plugin());