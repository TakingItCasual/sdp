import React from 'react';
import { BrowserRouter, Route, Switch } from 'react-router-dom';

import LandingPage from './pages/LandingPage';
import Profile from './pages/Profile';
import Page404 from './pages/Page404';

import './App.css';

function App() {
  return (
    <BrowserRouter>
      <Switch>
        <Route exact path='/' component={LandingPage} />
        <Route exact path='/profile' component={Profile} />
        <Route component={Page404} />
      </Switch>
    </BrowserRouter>
  );
}

export default App;
