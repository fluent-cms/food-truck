import React, { useState } from 'react';
import reactLogo from './assets/react.svg';
import { Button } from 'primereact/button';
import { InputText } from 'primereact/inputtext';

import 'primereact/resources/themes/lara-light-indigo/theme.css'; //theme
import 'primereact/resources/primereact.min.css'; //core css
import 'primeicons/primeicons.css'; //icons
import 'primeflex/primeflex.css'; // flex
import './App.css';
import Map from "./map/Map";


function App() {
  return (
      <div className="App">
          <Map />
          <div>Version 1.0</div>
      </div>
  );
}
export default App;