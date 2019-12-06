import React from "react";
import axios from "axios";
import { Grid, Button } from "semantic-ui-react";

const googleLogin = async e => {
  e.preventDefault();
  await axios
    .get("/api/v1/auth/google/login")
    .then(res => (window.location.href = res.data.redirect));
};

const LandingPage = () => {
  return (
    <Grid>
      <Grid.Row centered>
        <h2 className="ui huge header">
          {"Service Deployment Project Website"}
        </h2>
      </Grid.Row>
      <Grid.Row centered>
        <Button primary onClick={googleLogin}>
          Login with Google
        </Button>
      </Grid.Row>
    </Grid>
  );
};

export default LandingPage;
