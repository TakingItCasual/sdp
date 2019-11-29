import React from 'react';
import { Link } from 'react-router-dom';

const LandingPage = () => {
  return (
    <div>
      <h2>Service Deployment Project Website</h2>
      <Link to="/profile">Profile</Link>
    </div>
  );
};

export default LandingPage;
