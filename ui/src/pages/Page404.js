import React from 'react';
import { Link } from 'react-router-dom';

const Page404 = () => {
  return (
    <div>
      <h2>404 page not found</h2>
      <Link to="/">Back to Landing Page</Link>
    </div>
  );
};

export default Page404;
