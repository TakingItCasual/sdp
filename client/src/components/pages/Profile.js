import React, { useState, useEffect } from "react";
import { withRouter } from "react-router-dom";
import axios from "axios";
import { Grid, Form } from "semantic-ui-react";

import Header from "../Header";

const Profile = props => {
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");

  useEffect(() => {
    axios.get("/api/v1/priv/user").then(res => {
      setFirstName(res.data.first_name);
      setLastName(res.data.last_name);
      setEmail(res.data.school_email);
    });
  }, []);

  const handleSubmit = async e => {
    e.preventDefault();
    if (firstName && lastName && email) {
      await axios.put("/api/v1/priv/user", {
        first_name: firstName,
        last_name: lastName,
        school_email: email
      });
      props.history.push("/users");
    } else {
      alert("Please fill in all forms.");
    }
  };

  return (
    <>
      <Header />
      <Form method="PUT" onSubmit={handleSubmit}>
        <Form.Input
          label="First Name: "
          value={firstName}
          onChange={e => setFirstName(e.target.value)}
        />
        <Form.Input
          label="Last Name: "
          value={lastName}
          onChange={e => setLastName(e.target.value)}
        />
        <Form.Input
          label="School Email: "
          value={email}
          onChange={e => setEmail(e.target.value)}
        />
        <button className="ui primary submit button">Submit</button>
      </Form>
    </>
  );
};

export default withRouter(Profile);
