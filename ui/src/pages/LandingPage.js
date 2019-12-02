import React from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom';

const googleLogin = async e => {
  e.preventDefault();
  await axios.get('/api/v1/auth/google/login').then(res => window.location.href = res.data.redirect)
};

const LandingPage = () => {
  return (
    <div>
      <h2>Service Deployment Project Website</h2>
      <button onClick={googleLogin}>Login with Google</button>
    </div>
  );
};

export default LandingPage;
